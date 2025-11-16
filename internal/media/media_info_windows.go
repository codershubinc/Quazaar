package media

import (
	"Quazaar/pkg/helpers"
	"encoding/json"
	"fmt"
	"strings"
)

// GetPlayerInfoViaWindows retrieves media info using Windows APIs
func GetPlayerInfoViaWindows() (WindowsMediaInfo, error) {
	// Try System Media Transport Controls first (Windows 10+)
	if winInfo, err := GetPlayerInfoViaSMTC(); err == nil && winInfo.Title != "" {
		return winInfo, nil
	}

	// Fallback to PowerShell method
	return GetPlayerInfoViaPowerShell()
}

// GetPlayerInfoViaSMTC retrieves media info using System Media Transport Controls
func GetPlayerInfoViaSMTC() (WindowsMediaInfo, error) {
	info := WindowsMediaInfo{
		Metadata: make(map[string]string),
	}

	// Use PowerShell to query System Media Transport Controls
	psScript := `
		Add-Type -AssemblyName System.Runtime.WindowsRuntime
		$null = [Windows.Media.SystemMediaTransportControls, Windows.Media, ContentType=WindowsRuntime]

		try {
			$smtc = [Windows.Media.SystemMediaTransportControls]::GetForCurrentView()
			$timeline = $smtc.GetTimelineProperties()
			$display = $smtc.GetDisplayUpdater()

			$result = @{
				Title = $display.MusicProperties.Title
				Artist = $display.MusicProperties.Artist
				Album = $display.MusicProperties.AlbumTitle
				Status = if ($smtc.PlaybackStatus -eq [Windows.Media.MediaPlaybackStatus]::Playing) { "Playing" } elseif ($smtc.PlaybackStatus -eq [Windows.Media.MediaPlaybackStatus]::Paused) { "Paused" } else { "Stopped" }
				Position = $timeline.Position.TotalMilliseconds
				EndTime = $timeline.EndTime.TotalMilliseconds
				CanPlay = $true
				CanPause = $true
				CanGoPrevious = $true
				CanGoNext = $true
				CanSeek = $true
			}

			$result | ConvertTo-Json -Compress
		} catch {
			"{}"
		}
	`

	output, err := helpers.SpawnProcess("powershell", []string{"-Command", psScript})
	if err != nil {
		return info, err
	}

	// Parse JSON output
	if strings.TrimSpace(string(output)) == "{}" {
		return info, fmt.Errorf("no media playing via SMTC")
	}

	err = json.Unmarshal(output, &info)
	if err != nil {
		return info, err
	}

	info.Player = "System Media Transport Controls"
	info.PlayerName = "SMTC"

	return info, nil
}

// GetPlayerInfoViaPowerShell retrieves media info using PowerShell and Windows APIs
func GetPlayerInfoViaPowerShell() (WindowsMediaInfo, error) {
	info := WindowsMediaInfo{
		Metadata: make(map[string]string),
	}

	// PowerShell script to get media info from various sources
	psScript := `
		# Try to get media info from Windows Media Player first
		try {
			$wmp = New-Object -ComObject WMPlayer.OCX.7
			if ($wmp.currentMedia) {
				$result = @{
					Title = $wmp.currentMedia.name
					Artist = $wmp.currentMedia.getItemInfo("Artist")
					Album = $wmp.currentMedia.getItemInfo("Album")
					Status = if ($wmp.playState -eq 3) { "Playing" } elseif ($wmp.playState -eq 2) { "Paused" } else { "Stopped" }
					Position = $wmp.controls.currentPosition * 1000
					Length = $wmp.currentMedia.duration * 1000
					Player = "Windows Media Player"
					PlayerName = "WMP"
					CanPlay = $true
					CanPause = $true
					CanGoPrevious = $true
					CanGoNext = $true
					CanSeek = $true
				}
				$result | ConvertTo-Json -Compress
				exit
			}
		} catch {}

		# Try Spotify (if running)
		try {
			$spotify = Get-Process spotify -ErrorAction SilentlyContinue
			if ($spotify) {
				# This is a simplified approach - in reality you'd need Spotify Web API
				$result = @{
					Title = "Spotify Track (API integration needed)"
					Artist = "Spotify Artist"
					Status = "Unknown"
					Player = "Spotify"
					PlayerName = "Spotify"
					CanPlay = $true
					CanPause = $true
					CanGoPrevious = $true
					CanGoNext = $true
					CanSeek = $true
				}
				$result | ConvertTo-Json -Compress
				exit
			}
		} catch {}

		# Fallback: Try to detect any media applications
		$mediaApps = @("wmplayer", "spotify", "vlc", "foobar2000", "itunes", "musicbee")
		$runningMediaApps = Get-Process | Where-Object { $mediaApps -contains $_.ProcessName }

		if ($runningMediaApps) {
			$app = $runningMediaApps[0]
			$result = @{
				Title = "Media playing in " + $app.ProcessName
				Status = "Unknown"
				Player = $app.ProcessName
				PlayerName = $app.ProcessName
				CanPlay = $true
				CanPause = $true
				CanGoPrevious = $true
				CanGoNext = $true
				CanSeek = $true
			}
			$result | ConvertTo-Json -Compress
		} else {
			"{}"
		}
	`

	output, err := helpers.SpawnProcess("powershell", []string{"-Command", psScript})
	if err != nil {
		return info, err
	}

	// Parse JSON output
	if strings.TrimSpace(string(output)) == "{}" {
		return info, fmt.Errorf("no media applications detected")
	}

	err = json.Unmarshal(output, &info)
	if err != nil {
		return info, err
	}

	return info, nil
}

// GetWindowsActivePlayers lists active media players on Windows
func GetWindowsActivePlayers() ([]string, error) {
	psScript := `
		$mediaApps = @("wmplayer", "spotify", "vlc", "foobar2000", "itunes", "musicbee", "aimp", "winamp")
		$running = Get-Process | Where-Object { $mediaApps -contains $_.ProcessName } | Select-Object -ExpandProperty ProcessName
		if ($running) {
			$running | ConvertTo-Json
		} else {
			"[]"
		}
	`

	output, err := helpers.SpawnProcess("powershell", []string{"-Command", psScript})
	if err != nil {
		return nil, err
	}

	var players []string
	err = json.Unmarshal(output, &players)
	if err != nil {
		return nil, err
	}

	return players, nil
}

// ControlWindowsMedia sends control commands to Windows media players
func ControlWindowsMedia(command string) error {
	var psScript string

	switch command {
	case "play":
		psScript = `
			try {
				$wmp = New-Object -ComObject WMPlayer.OCX.7
				$wmp.controls.play()
			} catch {
				# Try SMTC
				Add-Type -AssemblyName System.Runtime.WindowsRuntime
				$smtc = [Windows.Media.SystemMediaTransportControls]::GetForCurrentView()
				$smtc.RequestPlayAsync()
			}
		`
	case "pause":
		psScript = `
			try {
				$wmp = New-Object -ComObject WMPlayer.OCX.7
				$wmp.controls.pause()
			} catch {
				Add-Type -AssemblyName System.Runtime.WindowsRuntime
				$smtc = [Windows.Media.SystemMediaTransportControls]::GetForCurrentView()
				$smtc.RequestPauseAsync()
			}
		`
	case "play-pause":
		psScript = `
			try {
				$wmp = New-Object -ComObject WMPlayer.OCX.7
				if ($wmp.playState -eq 3) {
					$wmp.controls.pause()
				} else {
					$wmp.controls.play()
				}
			} catch {
				Add-Type -AssemblyName System.Runtime.WindowsRuntime
				$smtc = [Windows.Media.SystemMediaTransportControls]::GetForCurrentView()
				$smtc.RequestPlayPauseAsync()
			}
		`
	case "next":
		psScript = `
			try {
				$wmp = New-Object -ComObject WMPlayer.OCX.7
				$wmp.controls.next()
			} catch {
				Add-Type -AssemblyName System.Runtime.WindowsRuntime
				$smtc = [Windows.Media.SystemMediaTransportControls]::GetForCurrentView()
				$smtc.RequestNextAsync()
			}
		`
	case "previous":
		psScript = `
			try {
				$wmp = New-Object -ComObject WMPlayer.OCX.7
				$wmp.controls.previous()
			} catch {
				Add-Type -AssemblyName System.Runtime.WindowsRuntime
				$smtc = [Windows.Media.SystemMediaTransportControls]::GetForCurrentView()
				$smtc.RequestPreviousAsync()
			}
		`
	default:
		return fmt.Errorf("unsupported command: %s", command)
	}

	_, err := helpers.SpawnProcess("powershell", []string{"-Command", psScript})
	return err
}

// ConvertWindowsMediaToMediaInfo converts WindowsMediaInfo to MediaInfo format
func ConvertWindowsMediaToMediaInfo(winInfo WindowsMediaInfo) MediaInfo {
	positionStr := formatDuration(winInfo.Position)
	lengthStr := formatDuration(winInfo.Length)

	return MediaInfo{
		Title:    winInfo.Title,
		Artist:   winInfo.Artist,
		Album:    winInfo.Album,
		Artwork:  winInfo.AlbumArt,
		Position: positionStr,
		Length:   lengthStr,
		Status:   strings.ToLower(winInfo.Status),
		Player:   winInfo.PlayerName,
	}
}

// GetPlayerInfoViaWindowsWithFallback tries Windows APIs first, falls back to other methods
func GetPlayerInfoViaWindowsWithFallback() (MediaInfo, error) {
	// Try Windows-specific methods first
	winInfo, err := GetPlayerInfoViaWindows()
	if err == nil && winInfo.Title != "" {
		return ConvertWindowsMediaToMediaInfo(winInfo), nil
	}

	// This would be a fallback for cross-platform compatibility
	// For now, return error since we're focusing on Windows
	return MediaInfo{}, fmt.Errorf("no media information available on Windows")
}

// Windows-specific helper functions
func getWindowsMediaSessions() ([]string, error) {
	psScript := `
		try {
			$sessions = Get-WindowsMediaSession
			$sessions | Where-Object { $_.PlaybackStatus -ne $null } | Select-Object -ExpandProperty AppDisplayName
		} catch {
			@()
		}
	`

	output, err := helpers.SpawnProcess("powershell", []string{"-Command", psScript})
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var sessions []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			sessions = append(sessions, line)
		}
	}

	return sessions, nil
}

// GetWindowsMediaInfoByPlayer gets media info for a specific player
func GetWindowsMediaInfoByPlayer(playerName string) (WindowsMediaInfo, error) {
	info := WindowsMediaInfo{
		Metadata: make(map[string]string),
	}

	// Normalize player name for matching
	playerName = strings.ToLower(playerName)

	psScript := fmt.Sprintf(`
		try {
			# Try SMTC first
			Add-Type -AssemblyName System.Runtime.WindowsRuntime
			$smtc = [Windows.Media.SystemMediaTransportControls]::GetForCurrentView()
			$timeline = $smtc.GetTimelineProperties()
			$display = $smtc.GetDisplayUpdater()

			$result = @{
				Title = $display.MusicProperties.Title
				Artist = $display.MusicProperties.Artist
				Album = $display.MusicProperties.AlbumTitle
				Status = if ($smtc.PlaybackStatus -eq [Windows.Media.MediaPlaybackStatus]::Playing) { "Playing" } elseif ($smtc.PlaybackStatus -eq [Windows.Media.MediaPlaybackStatus]::Paused) { "Paused" } else { "Stopped" }
				Position = $timeline.Position.TotalMilliseconds
				EndTime = $timeline.EndTime.TotalMilliseconds
				Player = "%s"
				PlayerName = "%s"
				CanPlay = $true
				CanPause = $true
				CanGoPrevious = $true
				CanGoNext = $true
				CanSeek = $true
			}

			$result | ConvertTo-Json -Compress
		} catch {
			# Fallback to process-specific detection
			$processes = Get-Process | Where-Object { $_.ProcessName -like "*%s*" }
			if ($processes) {
				$result = @{
					Title = "Media playing in %s"
					Status = "Unknown"
					Player = "%s"
					PlayerName = "%s"
					CanPlay = $true
					CanPause = $true
					CanGoPrevious = $true
					CanGoNext = $true
					CanSeek = $true
				}
				$result | ConvertTo-Json -Compress
			} else {
				"{}"
			}
		}
	`, playerName, playerName, playerName, playerName, playerName, playerName)

	output, err := helpers.SpawnProcess("powershell", []string{"-Command", psScript})
	if err != nil {
		return info, err
	}

	if strings.TrimSpace(string(output)) == "{}" {
		return info, fmt.Errorf("player %s not found or not playing media", playerName)
	}

	err = json.Unmarshal(output, &info)
	if err != nil {
		return info, err
	}

	return info, nil
}
