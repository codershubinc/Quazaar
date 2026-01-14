package usage

import (
	"encoding/json"
	"log"
	"math"
	"net/http"

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

func GetStorageUsage() (*DiskInfo, error) {
	// We only care about root "/"
	usage, err := disk.Usage("/")
	if err != nil {
		return nil, err
	}

	// We might optionally want partition info to get the device name and fstype
	// but usage.Path usually returns the path asked for.
	// To get device/fstype accurately for root, we can search partitions matching "/"
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

// HandleGetStorageUsage handles GET /api/v0.1/system/storage
func HandleGetStorageUsage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	storage, err := GetStorageUsage()
	if err != nil {
		log.Printf("❌ Failed to get storage usage: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Failed to get storage usage",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"storage": storage,
	})
}
