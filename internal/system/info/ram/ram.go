package ram

import (
"math"

"github.com/shirou/gopsutil/v3/mem"
)

type RAMInfo struct {
	MemoryUsage float64 `json:"memory_usage"`
	MemoryUsed  uint64  `json:"memory_used"`
	MemoryTotal uint64  `json:"memory_total"`
}

func GetRAMInfo() (*RAMInfo, error) {
	info := &RAMInfo{}
	vMem, err := mem.VirtualMemory()
	if err == nil {
		usedBytes := vMem.Total - vMem.Available
		info.MemoryUsed = usedBytes
		info.MemoryTotal = vMem.Total
		if info.MemoryTotal > 0 {
			info.MemoryUsage = math.Round((float64(usedBytes)/float64(vMem.Total)*100)*100) / 100
		}
	}
	return info, err
}
