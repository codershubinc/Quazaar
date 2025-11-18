//go:build !windows
// +build !windows

package media

import "fmt"

// Windows-specific functions - stubs for non-Windows platforms

func GetPlayerInfoViaWindows() (WindowsMediaInfo, error) {
	return WindowsMediaInfo{}, fmt.Errorf("windows media API not available on this platform")
}

func GetPlayerInfoViaSMTC() (WindowsMediaInfo, error) {
	return WindowsMediaInfo{}, fmt.Errorf("windows SMTC API not available on this platform")
}

func GetPlayerInfoViaPowerShell() (WindowsMediaInfo, error) {
	return WindowsMediaInfo{}, fmt.Errorf("powershell not available on this platform")
}

func GetWindowsActivePlayers() ([]string, error) {
	return nil, fmt.Errorf("windows media API not available on this platform")
}

func ControlWindowsMedia(command string) error {
	return fmt.Errorf("windows media control not available on this platform")
}

func ConvertWindowsMediaToMediaInfo(winInfo WindowsMediaInfo) MediaInfo {
	return MediaInfo{}
}

func GetPlayerInfoViaWindowsWithFallback() (MediaInfo, error) {
	return MediaInfo{}, fmt.Errorf("windows media API not available on this platform")
}

func GetWindowsMediaInfoByPlayer(playerName string) (WindowsMediaInfo, error) {
	return WindowsMediaInfo{}, fmt.Errorf("windows media API not available on this platform")
}
