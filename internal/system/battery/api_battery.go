package systemBattery

import (
	"encoding/json"
	"net/http"
)

func GetBatteryInfoApi(w http.ResponseWriter, r *http.Request) {
	batteryInfo, err := GetSystemBatteryInfo()
	if err != nil {
		http.Error(w, "Failed to get battery info", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(batteryInfo)
	if err != nil {
		http.Error(w, "Failed to encode battery info", http.StatusInternalServerError)
		return
	}
}
