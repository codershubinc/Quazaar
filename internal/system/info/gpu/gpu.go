package gpu

import (
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
)

type GPUInfo struct {
	GPUUsage       float64 `json:"gpu_usage"`
	GPUMemory      float64 `json:"gpu_memory"`
	GPUMemoryUsed  uint64  `json:"gpu_memory_used"`
	GPUMemoryTotal uint64  `json:"gpu_memory_total"`
	GPUWattage     float64 `json:"gpu_wattage"`
}

var (
	cachedPaths = make(map[string]string)
	pathMutex   sync.RWMutex
)

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

func GetGPUInfo() (*GPUInfo, error) {
	info := &GPUInfo{}
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
		info.GPUUsage = math.Round((totalGpuUsage/float64(gpuCount))*100) / 100
		info.GPUMemoryUsed = maxMemUsed
		info.GPUMemoryTotal = maxMemTotal
		info.GPUWattage = totalPower

		if info.GPUMemoryTotal > 0 {
			info.GPUMemory = math.Round((float64(info.GPUMemoryUsed)/float64(info.GPUMemoryTotal))*100*100) / 100
		}
	}

	return info, nil
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
