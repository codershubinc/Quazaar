package system

import (
	"Quazaar/internal/system/usage"
	"encoding/json"
	"log"
	"net/http"
)

// SystemInfo represents system information response
type SystemInfo struct {
	OS       string `json:"os"`
	Kernel   string `json:"kernel"`
	CPUModel string `json:"cpu_model"`
}

// HandleGetSystemInfo handles GET /api/v0.1/system/info - returns system information
func HandleGetSystemInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get host information
	hostInfo, err := usage.GetHostInfo()
	if err != nil {
		log.Printf("❌ Failed to get system info: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Failed to get system information",
			"message": err.Error(),
		})
		return
	}

	// Format OS string (Platform + PlatformVersion)
	osInfo := hostInfo.Platform
	if hostInfo.PlatformVersion != "" {
		osInfo += " " + hostInfo.PlatformVersion
	}

	systemInfo := SystemInfo{
		OS:       osInfo,
		Kernel:   hostInfo.KernelVersion,
		CPUModel: hostInfo.CPUModel,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"info":    systemInfo,
	})

	log.Printf("✅ System info retrieved: %s, %s", systemInfo.OS, systemInfo.CPUModel)
}
