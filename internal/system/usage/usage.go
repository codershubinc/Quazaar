package usage

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
	"github.com/godbus/dbus/v5"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

type SystemUsage struct {
	CPUUsage       float64 `json:"cpu_usage"`
	CPUTemp        float64 `json:"cpu_temp"`
	MemoryUsage    float64 `json:"memory_usage"`
	MemoryUsed     uint64  `json:"memory_used"`
	MemoryTotal    uint64  `json:"memory_total"`
	GPUUsage       float64 `json:"gpu_usage"`
	GPUMemory      float64 `json:"gpu_memory"`
	GPUMemoryUsed  uint64  `json:"gpu_memory_used"`
	GPUMemoryTotal uint64  `json:"gpu_memory_total"`
	CPUWattage     float64 `json:"cpu_wattage"`
	GPUWattage     float64 `json:"gpu_wattage"`
	Wattage        float64 `json:"wattage"`
	Uptime         uint64  `json:"uptime"`
	Distro         string  `json:"distro"`
	Kernel         string  `json:"kernel"`
	Username       string  `json:"username"`
}

type HostInfo struct {
	Hostname        string `json:"hostname"`
	OS              string `json:"os"`
	Platform        string `json:"platform"`
	PlatformVersion string `json:"platform_version"`
	KernelVersion   string `json:"kernel_version"`
	Arch            string `json:"arch"`
	CPUModel        string `json:"cpu_model"`
	CPUCores        int    `json:"cpu_cores"`
	Username        string `json:"username"`
	Uptime          uint64 `json:"uptime"`
}

// Globals for performance (caching & persistence)
var (
	dbusConn    *dbus.Conn
	dbusOnce    sync.Once
	cachedPaths = make(map[string]string)
	pathMutex   sync.RWMutex
)

// InitDBus establishes a persistent connection to the system bus
func InitDBus() {
	dbusOnce.Do(func() {
		var err error
		dbusConn, err = dbus.SystemBus()
		if err != nil {
			log.Printf("Warning: Failed to connect to DBus: %v", err)
		}
	})
}

func GetSystemUsage() (*SystemUsage, error) {
	// Ensure DBus is initialized
	InitDBus()

	usage := &SystemUsage{}

	// 1. CPU Usage
	cpuPercent, err := cpu.Percent(0, false)
	if err == nil && len(cpuPercent) > 0 {
		usage.CPUUsage = math.Round(cpuPercent[0]*100) / 100
	}

	// 2. Memory Usage (Total - Available logic)
	vMem, err := mem.VirtualMemory()
	if err == nil {
		usedBytes := vMem.Total - vMem.Available
		usage.MemoryUsed = usedBytes
		usage.MemoryTotal = vMem.Total
		if usage.MemoryTotal > 0 {
			usage.MemoryUsage = math.Round((float64(usedBytes)/float64(vMem.Total)*100)*100) / 100
		}
	}

	// 3. CPU Stats (Temp & Power)
	if t, err := getCPUTemperature(); err == nil {
		usage.CPUTemp = t
	}
	if w, err := getCPUWattage(); err == nil {
		usage.CPUWattage = w
	}

	// 4. GPU Stats (Hybrid: NVML + Sysfs)
	var totalGpuUsage, totalPower float64
	var maxMemUsed, maxMemTotal uint64
	var gpuCount int

	// Check NVIDIA (NVML first, fallback to exec if needed)
	nvUsage, _, nvMemUsed, nvMemTotal, nvPower, errNv := getNvidiaStats()
	if errNv == nil {
		totalGpuUsage += nvUsage
		totalPower += nvPower
		if nvMemTotal > maxMemTotal {
			maxMemTotal = nvMemTotal
			maxMemUsed = nvMemUsed
		}
		gpuCount++
	}

	// Check Sysfs (Integrated/AMD)
	sysUsage, _, sysMemUsed, sysMemTotal, sysPower, errSys := getSysfsGPUStats()
	if errSys == nil {
		totalGpuUsage += sysUsage
		totalPower += sysPower
		// Only override memory if it's the main large pool
		if maxMemTotal == 0 || sysMemTotal > maxMemTotal {
			maxMemTotal = sysMemTotal
			maxMemUsed = sysMemUsed
		}
		gpuCount++
	}

	// Consolidate GPU Data
	if gpuCount > 0 {
		usage.GPUUsage = math.Round((totalGpuUsage/float64(gpuCount))*100) / 100
		usage.GPUMemoryUsed = maxMemUsed
		usage.GPUMemoryTotal = maxMemTotal
		usage.GPUWattage = totalPower

		if usage.GPUMemoryTotal > 0 {
			usage.GPUMemory = math.Round((float64(usage.GPUMemoryUsed)/float64(usage.GPUMemoryTotal))*100*100) / 100
		}
	}

	// Total System Wattage
	usage.Wattage = usage.CPUWattage + usage.GPUWattage

	// Host Info (Uptime & OS)
	hostInfo, err := host.Info()
	if err == nil {
		usage.Uptime = hostInfo.Uptime
		usage.Distro = fmt.Sprintf("%s %s", casesTitle(hostInfo.Platform), hostInfo.PlatformVersion)
		usage.Kernel = hostInfo.KernelVersion
	}

	// Current User
	currentUser, err := user.Current()
	if err == nil {
		usage.Username = currentUser.Username
	}

	return usage, nil
}

// --- Helper Functions ---

func getCachedPath(key string) (string, bool) {
	pathMutex.RLock()
	defer pathMutex.RUnlock()
	val, ok := cachedPaths[key]
	return val, ok
}

func setCachedPath(key, value string) {
	pathMutex.Lock()
	defer pathMutex.Unlock()
	cachedPaths[key] = value
}

func readFileAsFloat(path string) (float64, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(strings.TrimSpace(string(content)), 64)
}

// getCPUTemperature with Caching
func getCPUTemperature() (float64, error) {
	if path, ok := getCachedPath("cpu_temp"); ok {
		if val, err := readFileAsFloat(path); err == nil {
			return math.Round((val/1000.0)*100) / 100, nil
		}
	}

	matches, _ := filepath.Glob("/sys/class/hwmon/hwmon*")
	for _, hwmonPath := range matches {
		nameBytes, _ := os.ReadFile(filepath.Join(hwmonPath, "name"))
		name := strings.TrimSpace(string(nameBytes))

		if strings.Contains(name, "k10temp") || strings.Contains(name, "coretemp") || strings.Contains(name, "zenpower") {
			target := filepath.Join(hwmonPath, "temp1_input")
			if val, err := readFileAsFloat(target); err == nil && val > 0 {
				setCachedPath("cpu_temp", target)
				return math.Round((val/1000.0)*100) / 100, nil
			}
		}
	}
	return 0, os.ErrNotExist
}

// getCPUWattage with Persistent DBus
func getCPUWattage() (float64, error) {
	if dbusConn != nil {
		obj := dbusConn.Object("org.freedesktop.UPower", "/org/freedesktop/UPower/devices/DisplayDevice")
		var energyRate dbus.Variant
		err := obj.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.UPower.Device", "EnergyRate").Store(&energyRate)
		if err == nil {
			if watts, ok := energyRate.Value().(float64); ok && watts > 0 {
				return watts, nil
			}
		}
	}

	if path, ok := getCachedPath("cpu_power"); ok {
		if val, err := readFileAsFloat(path); err == nil {
			return val / 1_000_000.0, nil
		}
	}

	matches, _ := filepath.Glob("/sys/class/hwmon/hwmon*")
	for _, hwmonPath := range matches {
		nameBytes, _ := os.ReadFile(filepath.Join(hwmonPath, "name"))
		name := strings.TrimSpace(string(nameBytes))
		if strings.Contains(name, "amdgpu") || strings.Contains(name, "nvidia") {
			continue
		}
		target := filepath.Join(hwmonPath, "power1_average")
		if val, err := readFileAsFloat(target); err == nil && val > 0 {
			setCachedPath("cpu_power", target)
			return val / 1_000_000.0, nil
		}
	}
	return 0, os.ErrNotExist
}

// getNvidiaStats using NVML (High Performance) with Fallback
func getNvidiaStats() (float64, float64, uint64, uint64, float64, error) {
	// Try NVML first (Native CGO)
	if ret := nvml.Init(); ret == nvml.SUCCESS {
		defer nvml.Shutdown()
		device, ret := nvml.DeviceGetHandleByIndex(0)
		if ret == nvml.SUCCESS {
			util, _ := device.GetUtilizationRates()
			mem, _ := device.GetMemoryInfo()
			power, _ := device.GetPowerUsage() // mW

			memPct := 0.0
			if mem.Total > 0 {
				memPct = (float64(mem.Used) / float64(mem.Total)) * 100
			}
			return float64(util.Gpu), memPct, mem.Used, mem.Total, float64(power) / 1000.0, nil
		}
	}

	// Fallback to nvidia-smi exec (Slower but works without CGO setup)
	return getNvidiaStatsExec()
}

func getNvidiaStatsExec() (float64, float64, uint64, uint64, float64, error) {
	cmd := exec.Command("nvidia-smi", "--query-gpu=utilization.gpu,memory.used,memory.total,power.draw", "--format=csv,noheader,nounits")
	out, err := cmd.Output()
	if err != nil {
		return 0, 0, 0, 0, 0, err
	}

	output := strings.TrimSpace(string(out))
	lines := strings.Split(output, "\n")
	if len(lines) == 0 || lines[0] == "" {
		return 0, 0, 0, 0, 0, nil
	}

	parts := strings.Split(lines[0], ",")
	if len(parts) != 4 {
		return 0, 0, 0, 0, 0, nil
	}

	parseFloat := func(s string) float64 {
		val, _ := strconv.ParseFloat(strings.TrimSpace(s), 64)
		return val
	}

	gpuUtil := parseFloat(parts[0])
	memUsedMb := parseFloat(parts[1])
	memTotalMb := parseFloat(parts[2])
	powerW := parseFloat(parts[3])

	memUtil := 0.0
	if memTotalMb > 0 {
		memUtil = math.Round((memUsedMb/memTotalMb)*100*100) / 100
	}

	return gpuUtil, memUtil, uint64(memUsedMb * 1024 * 1024), uint64(memTotalMb * 1024 * 1024), powerW, nil
}

// getSysfsGPUStats (AMD/Intel)
func getSysfsGPUStats() (float64, float64, uint64, uint64, float64, error) {
	matches, _ := filepath.Glob("/sys/class/drm/card*/device/gpu_busy_percent")
	for _, path := range matches {
		content, err := os.ReadFile(path)
		if err == nil {
			val, err := strconv.ParseFloat(strings.TrimSpace(string(content)), 64)
			if err == nil {
				dir := filepath.Dir(path)
				used, _ := readFileAsFloat(filepath.Join(dir, "mem_info_vram_used"))
				total, _ := readFileAsFloat(filepath.Join(dir, "mem_info_vram_total"))
				memPct := 0.0
				if total > 0 {
					memPct = (used / total) * 100
				}
				power := 0.0
				if pPath, ok := getCachedPath("gpu_sysfs_power"); ok {
					if pVal, err := readFileAsFloat(pPath); err == nil {
						power = pVal / 1_000_000.0
					}
				} else {
					hwmonMatches, _ := filepath.Glob(filepath.Join(dir, "hwmon", "hwmon*", "power1_average"))
					if len(hwmonMatches) > 0 {
						setCachedPath("gpu_sysfs_power", hwmonMatches[0])
						if pVal, err := readFileAsFloat(hwmonMatches[0]); err == nil {
							power = pVal / 1_000_000.0
						}
					}
				}
				return val, memPct, uint64(used), uint64(total), power, nil
			}
		}
	}
	return 0, 0, 0, 0, 0, os.ErrNotExist
}

// casesTitle converts first letter to uppercase (strings.Title is deprecated)
func casesTitle(str string) string {
	if len(str) == 0 {
		return ""
	}
	return strings.ToUpper(str[:1]) + str[1:]
}

// GetHostInfo retrieves detailed host information
func GetHostInfo() (*HostInfo, error) {
	// 1. Host Stats
	hInfo, err := host.Info()
	if err != nil {
		return nil, err
	}

	// 2. User Info
	currentUser, err := user.Current()
	username := "unknown"
	if err == nil {
		username = currentUser.Username
	}

	// 3. CPU Info (Get first core to find model name)
	cInfos, err := cpu.Info()
	cpuModel := "Unknown CPU"
	cores := 0
	if err == nil && len(cInfos) > 0 {
		cpuModel = cInfos[0].ModelName
		cores = len(cInfos) // roughly logical core count
	}

	return &HostInfo{
		Hostname:        hInfo.Hostname,
		OS:              hInfo.OS,
		Platform:        hInfo.Platform,
		PlatformVersion: hInfo.PlatformVersion,
		KernelVersion:   hInfo.KernelVersion,
		Arch:            hInfo.KernelArch,
		CPUModel:        cpuModel,
		CPUCores:        cores,
		Username:        username,
		Uptime:          hInfo.Uptime,
	}, nil
}
