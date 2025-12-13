package config

// Version information constants
const (
	Version  = "v0.2.1-nightly"
	Codename = "Pulsar"
	AppName  = "Quazaar"
)

// VersionInfo holds the application version details
type VersionInfo struct {
	AppName  string `json:"appName"`
	Version  string `json:"version"`
	Codename string `json:"codename"`
}

// GetVersionInfo returns the version information struct
func GetVersionInfo() VersionInfo {
	return VersionInfo{
		AppName:  AppName,
		Version:  Version,
		Codename: Codename,
	}
}

// GetVersionString returns the formatted version string
func GetVersionString() string {
	return AppName + " " + Version + " (" + Codename + ")"
}
