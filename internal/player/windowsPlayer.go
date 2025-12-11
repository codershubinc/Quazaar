//go:build windows

package player

import (
	"Quazaar/pkg/models"
	"fmt"
)

func playPauseWindows() (bool, error) {
	// Implement Windows-specific play/pause functionality here
	return true, nil
}
func nextWindows() (bool, error) {
	// Implement Windows-specific next functionality here
	return true, nil
}
func prevWindows() (bool, error) {
	// Implement Windows-specific previous functionality here
	return true, nil
}
func seekBackwardWindows() (bool, error) {
	// Implement Windows-specific seek backward functionality here
	return true, nil
}
func seekForwardWindows() (bool, error) {
	// Implement Windows-specific seek forward functionality here
	return true, nil
}
func setVolumeWindows(volume int) (bool, error) {
	// Implement Windows-specific set volume functionality here
	return true, nil
}

func getCurrentPlayerMetadataWindows() (models.MediaInfo, error) {
	return models.MediaInfo{}, fmt.Errorf("not implemented on windows")
}

func getAllPlayersWindows() ([]string, error) {
	return nil, fmt.Errorf("not implemented on windows")
}

var WindowsPlayerHandler = models.PlayerFunctions{
	PlayPause:                playPauseWindows,
	Next:                     nextWindows,
	Prev:                     prevWindows,
	SeekBackward:             seekBackwardWindows,
	SeekForward:              seekForwardWindows,
	SetVolume:                setVolumeWindows,
	GetCurrentPlayerMetadata: getCurrentPlayerMetadataWindows,
	GetAllPlayers:            getAllPlayersWindows,
}

func initializeDefaultPlayer() models.PlayerFunctions {
	return WindowsPlayerHandler
}
