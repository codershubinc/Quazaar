//go:build !linux && !windows

package player

import (
	"Quazaar/pkg/models"
	"fmt"
	"runtime"
)

func initializeDefaultPlayer() models.PlayerFunctions {
	return noopPlayer
}

var noopPlayer = models.PlayerFunctions{
	PlayPause:    func() (bool, error) { return false, fmt.Errorf("player unsupported on %s", runtime.GOOS) },
	Next:         func() (bool, error) { return false, fmt.Errorf("player unsupported on %s", runtime.GOOS) },
	Prev:         func() (bool, error) { return false, fmt.Errorf("player unsupported on %s", runtime.GOOS) },
	SeekForward:  func() (bool, error) { return false, fmt.Errorf("player unsupported on %s", runtime.GOOS) },
	SeekBackward: func() (bool, error) { return false, fmt.Errorf("player unsupported on %s", runtime.GOOS) },
	SeekTo:       func(position int64) (bool, error) { return false, fmt.Errorf("player unsupported on %s", runtime.GOOS) },
	SetVolume:    func(int) (bool, error) { return false, fmt.Errorf("player unsupported on %s", runtime.GOOS) },
	GetCurrentPlayerMetadata: func() (models.MediaInfo, error) {
		return models.MediaInfo{}, fmt.Errorf("player unsupported on %s", runtime.GOOS)
	},
	GetAllPlayers: func() ([]string, error) { return nil, fmt.Errorf("player unsupported on %s", runtime.GOOS) },
}
