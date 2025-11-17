
package media

import (
	"Quazaar/pkg/helpers"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// DBusMediaInfo represents media information retrieved via D-Bus/MPRIS
type DBusMediaInfo struct {
	Title         string
	Artist        string
	Album         string
	AlbumArt      string
	Position      int64  // in microseconds
	Length        int64  // in microseconds
	Status        string // Playing, Paused, Stopped
	Player        string
	PlayerName    string
	CanPlay       bool
	CanPause      bool
	CanGoPrevious bool
	CanGoNext     bool
	CanSeek       bool
	Metadata      map[string]string
}

// GetPlayerInfoViaDBus retrieves media info using D-Bus MPRIS interface
func GetPlayerInfoViaDBus() (DBusMediaInfo, error) {
	info := DBusMediaInfo{
		Metadata: make(map[string]string),
	}

	// Get list of available MPRIS players
	players, err := GetMPRISPlayers()
	if err != nil || len(players) == 0 {
		return info, fmt.Errorf("no MPRIS players found via D-Bus")
	}

	// Use the first available player
	player := players[0]
	info.Player = player
	info.PlayerName = extractPlayerName(player)

	// Get properties
	props, err := GetMPRISProperties(player)
	if err != nil {
		return info, fmt.Errorf("failed to get MPRIS properties: %v", err)
	}

	// Parse properties
	parseDBusProperties(&info, props)

	// Get metadata
	metadata, err := GetMPRISMetadata(player)
	if err == nil {
		parseMetadata(&info, metadata)
	}

	return info, nil
}

// GetPlayerInfoViaDBusForPlayer retrieves media info for a specific D-Bus MPRIS player
func GetPlayerInfoViaDBusForPlayer(playerService string) (DBusMediaInfo, error) {
	info := DBusMediaInfo{
		Metadata: make(map[string]string),
	}

	// Validate player service name
	if !strings.HasPrefix(playerService, "org.mpris.MediaPlayer2.") {
		return info, fmt.Errorf("invalid MPRIS player service name: %s", playerService)
	}

	info.Player = playerService
	info.PlayerName = extractPlayerName(playerService)

	// Get properties
	props, err := GetMPRISProperties(playerService)
	if err != nil {
		return info, fmt.Errorf("failed to get MPRIS properties for %s: %v", playerService, err)
	}

	// Parse properties
	parseDBusProperties(&info, props)

	// Get metadata
	metadata, err := GetMPRISMetadata(playerService)
	if err == nil {
		parseMetadata(&info, metadata)
	}

	return info, nil
}

// GetMPRISPlayers retrieves list of available MPRIS players from D-Bus
func GetMPRISPlayers() ([]string, error) {
	output, err := helpers.SpawnProcess("dbus-send", []string{
		"--session",
		"--print-reply",
		"--dest=org.freedesktop.DBus",
		"/org/freedesktop/DBus",
		"org.freedesktop.DBus.ListNames",
	})

	if err != nil {
		return nil, err
	}

	var players []string
	lines := strings.Split(string(output), "\n")

	// Look for MPRIS player service names
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, `"org.mpris.MediaPlayer2`) {
			// Extract the service name
			serviceName := extractServiceName(line)
			if serviceName != "" {
				players = append(players, serviceName)
			}
		}
	}

	return players, nil
}

// GetMPRISProperties retrieves properties from a specific MPRIS player
func GetMPRISProperties(playerService string) (string, error) {
	output, err := helpers.SpawnProcess("dbus-send", []string{
		"--session",
		"--print-reply",
		fmt.Sprintf("--dest=%s", playerService),
		"/org/mpris/MediaPlayer2",
		"org.freedesktop.DBus.Properties.GetAll",
		"string:org.mpris.MediaPlayer2.Player",
	})

	return string(output), err
}

// GetMPRISMetadata retrieves metadata from a specific MPRIS player
func GetMPRISMetadata(playerService string) (string, error) {
	output, err := helpers.SpawnProcess("dbus-send", []string{
		"--session",
		"--print-reply",
		fmt.Sprintf("--dest=%s", playerService),
		"/org/mpris/MediaPlayer2",
		"org.freedesktop.DBus.Properties.Get",
		"string:org.mpris.MediaPlayer2.Player",
		"string:Metadata",
	})

	return string(output), err
}

// parseDBusProperties parses D-Bus output and extracts properties
func parseDBusProperties(info *DBusMediaInfo, output string) {
	lines := strings.Split(output, "\n")

	for i, line := range lines {
		line = strings.TrimSpace(line)

		// Parse PlaybackStatus
		if strings.Contains(line, "PlaybackStatus") {
			if i+1 < len(lines) {
				status := extractStringValue(lines[i+1])
				info.Status = status
			}
		}

		// Parse Position
		if strings.Contains(line, "Position") && !strings.Contains(line, "PlaybackStatus") {
			if i+1 < len(lines) {
				pos := extractInt64Value(lines[i+1])
				info.Position = pos
			}
		}

		// Parse CanPlay
		if strings.Contains(line, "CanPlay") && !strings.Contains(line, "PlaybackStatus") {
			if i+1 < len(lines) {
				info.CanPlay = extractBoolValue(lines[i+1])
			}
		}

		// Parse CanPause
		if strings.Contains(line, "CanPause") {
			if i+1 < len(lines) {
				info.CanPause = extractBoolValue(lines[i+1])
			}
		}

		// Parse CanGoPrevious
		if strings.Contains(line, "CanGoPrevious") {
			if i+1 < len(lines) {
				info.CanGoPrevious = extractBoolValue(lines[i+1])
			}
		}

		// Parse CanGoNext
		if strings.Contains(line, "CanGoNext") {
			if i+1 < len(lines) {
				info.CanGoNext = extractBoolValue(lines[i+1])
			}
		}

		// Parse CanSeek
		if strings.Contains(line, "CanSeek") {
			if i+1 < len(lines) {
				info.CanSeek = extractBoolValue(lines[i+1])
			}
		}
	}
}

// parseMetadata parses metadata from D-Bus output
func parseMetadata(info *DBusMediaInfo, output string) {
	lines := strings.Split(output, "\n")

	for i, line := range lines {
		line = strings.TrimSpace(line)

		// Parse Title
		if strings.Contains(line, "xesam:title") {
			if i+1 < len(lines) {
				info.Title = extractStringValue(lines[i+1])
			}
		}

		// Parse Artist
		if strings.Contains(line, "xesam:artist") {
			if i+1 < len(lines) {
				info.Artist = extractStringArrayValue(lines[i+1])
			}
		}

		// Parse Album
		if strings.Contains(line, "xesam:album") {
			if i+1 < len(lines) {
				info.Album = extractStringValue(lines[i+1])
			}
		}

		// Parse Art URL
		if strings.Contains(line, "mpris:artUrl") {
			if i+1 < len(lines) {
				artURL := extractStringValue(lines[i+1])
				if artURL != "" {
					// Process artwork URL if needed
					HandleArtworkRequest(artURL)
					info.AlbumArt = artURL
				}
			}
		}

		// Parse Length
		if strings.Contains(line, "mpris:length") {
			if i+1 < len(lines) {
				info.Length = extractInt64Value(lines[i+1])
			}
		}
	}
}

// Helper extraction functions
func extractServiceName(line string) string {
	re := regexp.MustCompile(`"(org\.mpris\.MediaPlayer2\.[^"]+)"`)
	if matches := re.FindStringSubmatch(line); len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func extractStringValue(line string) string {
	re := regexp.MustCompile(`string\s+"([^"]*)"`)
	if matches := re.FindStringSubmatch(line); len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func extractStringArrayValue(line string) string {
	re := regexp.MustCompile(`string\s+"([^"]*)"`)
	var values []string
	matches := re.FindAllStringSubmatch(line, -1)
	for _, match := range matches {
		if len(match) > 1 {
			values = append(values, match[1])
		}
	}
	return strings.Join(values, ", ")
}

func extractInt64Value(line string) int64 {
	re := regexp.MustCompile(`int64\s+(\d+)`)
	if matches := re.FindStringSubmatch(line); len(matches) > 1 {
		var val int64
		fmt.Sscanf(matches[1], "%d", &val)
		return val
	}
	return 0
}

func extractBoolValue(line string) bool {
	line = strings.ToLower(line)
	return strings.Contains(line, "true")
}

func extractPlayerName(serviceString string) string {
	parts := strings.Split(serviceString, ".")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return "Unknown"
}

// ConvertDBusMediaToMediaInfo converts DBusMediaInfo to MediaInfo format
func ConvertDBusMediaToMediaInfo(dbusInfo DBusMediaInfo) MediaInfo {
	positionStr := formatDuration(dbusInfo.Position / 1000)
	lengthStr := formatDuration(dbusInfo.Length / 1000)

	return MediaInfo{
		Title:    dbusInfo.Title,
		Artist:   dbusInfo.Artist,
		Album:    dbusInfo.Album,
		Artwork:  dbusInfo.AlbumArt,
		Position: positionStr,
		Length:   lengthStr,
		Status:   strings.ToLower(dbusInfo.Status),
		Player:   dbusInfo.PlayerName,
	}
}

// formatDuration converts milliseconds to MM:SS format
func formatDuration(ms int64) string {
	if ms <= 0 {
		return "0:00"
	}
	seconds := ms / 1000
	minutes := seconds / 60
	seconds = seconds % 60
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

// GetPlayerInfoViaDBusWithFallback tries D-Bus first, falls back to playerctl
func GetPlayerInfoViaDBusWithFallback() (MediaInfo, error) {
	// Try D-Bus method first
	dbusInfo, err := GetPlayerInfoViaDBus()
	if err == nil && dbusInfo.Title != "" {
		return ConvertDBusMediaToMediaInfo(dbusInfo), nil
	}

	// Fall back to playerctl
	return GetPlayerInfo()
}

// MonitorPlayerChangesViaDBus monitors for player changes via D-Bus signals
func MonitorPlayerChangesViaDBus(callback func(DBusMediaInfo)) error {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	var lastTitle string

	for range ticker.C {
		info, err := GetPlayerInfoViaDBus()
		if err != nil {
			continue
		}

		// Only call callback if track changed
		if info.Title != lastTitle {
			callback(info)
			lastTitle = info.Title
		}
	}

	return nil
}
