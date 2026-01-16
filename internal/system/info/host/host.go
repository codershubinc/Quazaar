package host

import (
	"fmt"
	"os/user"
	"strings"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
)

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
	// Additional fields for SystemUsage compatibility
	Distro string `json:"distro"`
	Kernel string `json:"kernel"`
}

// casesTitle converts first letter to uppercase
func casesTitle(str string) string {
	if len(str) == 0 {
		return ""
	}
	return strings.ToUpper(str[:1]) + str[1:]
}

func GetHostInfo() (*HostInfo, error) {
	hInfo, err := host.Info()
	if err != nil {
		return nil, err
	}

	currentUser, err := user.Current()
	username := "unknown"
	if err == nil {
		username = currentUser.Username
	}

	cInfos, err := cpu.Info()
	cpuModel := "Unknown CPU"
	cores := 0
	if err == nil && len(cInfos) > 0 {
		cpuModel = cInfos[0].ModelName
		cores = len(cInfos)
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
		Distro:          fmt.Sprintf("%s %s", casesTitle(hInfo.Platform), hInfo.PlatformVersion),
		Kernel:          hInfo.KernelVersion,
	}, nil
}
