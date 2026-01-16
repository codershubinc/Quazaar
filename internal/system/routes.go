package system

import (
	"Quazaar/internal/middleware"
	systemBattery "Quazaar/internal/system/battery"
	"Quazaar/internal/system/bluetooth"
	brightness "Quazaar/internal/system/brightness"
	systemInfo "Quazaar/internal/system/info"
	sound "Quazaar/internal/system/sound"
	"Quazaar/internal/system/usage"
	volume "Quazaar/internal/system/volume"
	"Quazaar/internal/wakatime"
	"net/http"
)

func SetupRoutes() {
	// API v0.1 - System Info
	http.HandleFunc("/api/v0.1/system/info", middleware.AuthenticationMiddleware(HandleGetSystemInfo))
	http.HandleFunc("/api/v0.1/system/wifi", middleware.AuthenticationMiddleware(HandleGetWiFiInfo))
	http.HandleFunc("/api/v0.1/system/bluetooth", middleware.AuthenticationMiddleware(bluetooth.HandleGetBluetoothDevices))
	http.HandleFunc("/api/v0.1/system/volume", middleware.AuthenticationMiddleware(volume.HandleGetVolume))
	http.HandleFunc("/api/v0.1/system/brightness", middleware.AuthenticationMiddleware(brightness.HandleGetBrightness))
	http.HandleFunc("/api/v0.1/system/sound/devices", middleware.AuthenticationMiddleware(sound.HandleListDevices))
	http.HandleFunc("/api/v0.1/system/sound/device", middleware.AuthenticationMiddleware(sound.HandleSetDevice))

	// API v0.1 - System Temp routes
	http.HandleFunc("/api/v0.1/system/info/t", systemInfo.GetSystemInfoApi)

	// API v0.1 - Usage
	http.HandleFunc("/api/v0.1/system/usage", middleware.AuthenticationMiddleware(usage.HandleGetSystemUsage))
	http.HandleFunc("/api/v0.1/system/storage", middleware.AuthenticationMiddleware(usage.HandleGetStorageUsage))

	// API v0.1 - WakaTime
	http.HandleFunc("/api/v0.1/system/wakatime", middleware.AuthenticationMiddleware(wakatime.HandleGetWakaTimeStats))

	// API v0.1 - System  TODO: Move to systemBattery package
	http.HandleFunc("/api/v0.1/system/battery", systemBattery.GetBatteryInfoApi)
}
