package cpu

import (
"log"
"math"
"os"
"path/filepath"
"strconv"
"strings"
"sync"

"github.com/godbus/dbus/v5"
"github.com/shirou/gopsutil/v3/cpu"
)

type CPUInfo struct {
	CPUUsage   float64 `json:"cpu_usage"`
	CPUTemp    float64 `json:"cpu_temp"`
	CPUWattage float64 `json:"cpu_wattage"`
}

var (
dbusConn    *dbus.Conn
dbusOnce    sync.Once
pathCache   = make(map[string]string)
pathMutex   sync.RWMutex
)

func InitDBus() {
	dbusOnce.Do(func() {
		var err error
		dbusConn, err = dbus.SystemBus()
		if err != nil {
			log.Printf("Warning: Failed to connect to DBus: %v", err)
		}
	})
}

func GetCPUInfo() CPUInfo {
	InitDBus()
	return CPUInfo{
		CPUUsage:   getCPUUsage(),
		CPUTemp:    getCPUTemp(),
		CPUWattage: getCPUWattage(),
	}

}

func getCPUUsage() float64 {
	cpuPercent, err := cpu.Percent(0, false)
	if err == nil && len(cpuPercent) > 0 {
		return math.Round(cpuPercent[0]*100) / 100
	}
	return -1
}

func getCPUTemp() float64 {
	if t, err := getCPUTemperature(); err == nil {
		return t
	}
	return -1
}

func getCPUWattage() float64 {
	if w, err := getCPUWattageValue(); err == nil {
		return w
	}
	return 0
}

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

func getCPUWattageValue() (float64, error) {
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

// helper functions for caching file paths

func readFileAsFloat(path string) (float64, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(strings.TrimSpace(string(content)), 64)
}

func getCachedPath(key string) (string, bool) {
	pathMutex.RLock()
	defer pathMutex.RUnlock()
	path, exists := pathCache[key]
	return path, exists
}

func setCachedPath(key, path string) {
	pathMutex.Lock()
	defer pathMutex.Unlock()
	pathCache[key] = path
}
