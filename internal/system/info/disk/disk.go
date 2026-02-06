package disk

import (
"log"
"math"

"github.com/shirou/gopsutil/v3/disk"
)

type DiskInfo struct {
	Path        string  `json:"path"`
	Device      string  `json:"device"`
	Fstype      string  `json:"fstype"`
	Total       uint64  `json:"total"`
	Free        uint64  `json:"free"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"used_percent"`
}

func GetDiskUsage() (*DiskInfo, error) {
	// We only care about root "/"
	usage, err := disk.Usage("/")
	if err != nil {
		log.Printf("Failed to get disk usage: %v", err)
		return nil, err
	}

	var device, fstype string
	partitions, err := disk.Partitions(false)
	if err == nil {
		for _, p := range partitions {
			if p.Mountpoint == "/" {
				device = p.Device
				fstype = p.Fstype
				break
			}
		}
	}

	info := &DiskInfo{
		Path:        "/",
		Device:      device,
		Fstype:      fstype,
		Total:       usage.Total,
		Free:        usage.Free,
		Used:        usage.Used,
		UsedPercent: math.Round(usage.UsedPercent*100) / 100,
	}

	return info, nil
}
