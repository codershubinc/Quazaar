package system

import (
	"encoding/json"
	"log"
	"net/http"
)

// HandleGetWiFiInfo handles GET /api/v0.1/system/wifi - returns WiFi information
func HandleGetWiFiInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	wifiInfo, err := GetWiFiInfo()
	if err != nil {
		log.Printf("❌ Failed to get WiFi info: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Failed to get WiFi information",
			"message": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"wifi":    wifiInfo,
	})

	log.Printf("✅ WiFi info retrieved: %s", wifiInfo.SSID)
}

// HandleGetBluetoothDevices handles GET /api/v0.1/system/bluetooth - returns connected Bluetooth devices
func HandleGetBluetoothDevices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	devices, err := GetBluetoothDevices()
	if err != nil {
		log.Printf("❌ Failed to get Bluetooth devices: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Failed to get Bluetooth devices",
			"message": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"count":   len(devices),
		"devices": devices,
	})

	log.Printf("✅ Bluetooth devices retrieved: %d devices", len(devices))
}
