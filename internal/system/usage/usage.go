package usage

import (
	"Quazaar/internal/system/info/cpu"
	"Quazaar/internal/system/info/gpu"
	"Quazaar/internal/system/info/host"
	"Quazaar/internal/system/info/ram"
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

type HostInfo = host.HostInfo

func GetSystemUsage() (*SystemUsage, error) {
	usage := &SystemUsage{}

	// CPU
	cpuInfo := cpu.GetCPUInfo()
	usage.CPUUsage = cpuInfo.CPUUsage
	usage.CPUTemp = cpuInfo.CPUTemp
	usage.CPUWattage = cpuInfo.CPUWattage

	// Memory
	if ramInfo, err := ram.GetRAMInfo(); err == nil {
		usage.MemoryUsage = ramInfo.MemoryUsage
		usage.MemoryUsed = ramInfo.MemoryUsed
		usage.MemoryTotal = ramInfo.MemoryTotal
	}

	// GPU
	if gpuInfo, err := gpu.GetGPUInfo(); err == nil {
		usage.GPUUsage = gpuInfo.GPUUsage
		usage.GPUMemory = gpuInfo.GPUMemory
		usage.GPUMemoryUsed = gpuInfo.GPUMemoryUsed
		usage.GPUMemoryTotal = gpuInfo.GPUMemoryTotal
		usage.GPUWattage = gpuInfo.GPUWattage
	}

	// Host
	if hostInfo, err := host.GetHostInfo(); err == nil {
		usage.Uptime = hostInfo.Uptime
		usage.Distro = hostInfo.Distro
		usage.Kernel = hostInfo.Kernel
		usage.Username = hostInfo.Username
	}

	// Total Wattage
	usage.Wattage = usage.CPUWattage + usage.GPUWattage

	return usage, nil
}

func GetHostInfo() (*host.HostInfo, error) {
	return host.GetHostInfo()
}
