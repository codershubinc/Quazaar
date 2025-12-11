package player

import "errors"

type Backend interface {
	Play() error
	Pause() error
	Toggle() error
	Next() error
	Previous() error
	GetState() (MediaInfo, error)
}

// activeBackend holds the OS-specific implementation
var activeBackend Backend

// Public wrappers for the rest of the app to call
func Play() error {
	if activeBackend == nil {
		return errors.New("no player backend initialized")
	}
	return activeBackend.Play()
}

func Pause() error {
	if activeBackend == nil {
		return errors.New("no player backend initialized")
	}
	return activeBackend.Pause()
}

func Toggle() error {
	if activeBackend == nil {
		return errors.New("no player backend initialized")
	}
	return activeBackend.Toggle()
}

func Next() error {
	if activeBackend == nil {
		return errors.New("no player backend initialized")
	}
	return activeBackend.Next()
}

func Previous() error {
	if activeBackend == nil {
		return errors.New("no player backend initialized")
	}
	return activeBackend.Previous()
}

func GetState() (MediaInfo, error) {
	if activeBackend == nil {
		return MediaInfo{}, errors.New("no player backend initialized")
	}
	return activeBackend.GetState()
}
