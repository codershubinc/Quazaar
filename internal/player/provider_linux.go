//go:build linux

package player

import (
	"Quazaar/internal/media"
	"Quazaar/pkg/models"
	"fmt"
	"strings"
	"sync"

	"github.com/godbus/dbus/v5"
)

var (
	dbusConn         *dbus.Conn
	dbusOnce         sync.Once
	lastActivePlayer string
)

func GetDBusConnection() (*dbus.Conn, error) {
	var err error
	dbusOnce.Do(func() {
		dbusConn, err = dbus.SessionBus()
	})
	return dbusConn, err
}

func findActivePlayer(conn *dbus.Conn) (string, error) {
	var names []string
	err := conn.BusObject().Call("org.freedesktop.DBus.ListNames", 0).Store(&names)
	if err != nil {
		return "", err
	}

	var players []string
	for _, name := range names {
		if strings.HasPrefix(name, "org.mpris.MediaPlayer2.") {
			players = append(players, name)
		}
	}

	if len(players) == 0 {
		return "", fmt.Errorf("no players found")
	}

	var pausedPlayer string
	var spotifyPlayer string

	// Priority: Playing > Paused > Spotify > Stopped
	for _, player := range players {
		if strings.Contains(strings.ToLower(player), "spotify") {
			spotifyPlayer = player
		}

		obj := conn.Object(player, "/org/mpris/MediaPlayer2")
		statusVar, err := obj.GetProperty("org.mpris.MediaPlayer2.Player.PlaybackStatus")
		if err != nil {
			continue
		}
		status := statusVar.Value().(string)
		if status == "Playing" {
			lastActivePlayer = player
			return player, nil
		}

		if status == "Paused" && pausedPlayer == "" {
			pausedPlayer = player
		}
	}

	if pausedPlayer != "" {
		lastActivePlayer = pausedPlayer
		return pausedPlayer, nil
	}

	// Try last active player if it still exists
	if lastActivePlayer != "" {
		for _, p := range players {
			if p == lastActivePlayer {
				return lastActivePlayer, nil
			}
		}
	}

	if spotifyPlayer != "" {
		return spotifyPlayer, nil
	}

	// If none playing, return the first one
	return players[0], nil
}

func getMetadata() (models.MediaInfo, error) {
	conn, err := GetDBusConnection()
	if err != nil {
		return models.MediaInfo{}, err
	}

	player, err := findActivePlayer(conn)
	if err != nil {
		return models.MediaInfo{}, err
	}

	obj := conn.Object(player, "/org/mpris/MediaPlayer2")

	variant, err := obj.GetProperty("org.mpris.MediaPlayer2.Player.Metadata")
	if err != nil {
		return models.MediaInfo{}, err
	}

	statusVar, _ := obj.GetProperty("org.mpris.MediaPlayer2.Player.PlaybackStatus")
	positionVar, _ := obj.GetProperty("org.mpris.MediaPlayer2.Player.Position")

	metadata := variant.Value().(map[string]dbus.Variant)

	getString := func(key string) string {
		if val, ok := metadata[key]; ok {
			if s, ok := val.Value().(string); ok {
				return s
			}
			if s, ok := val.Value().([]string); ok {
				return strings.Join(s, ", ")
			}
		}
		return ""
	}

	getInt64 := func(key string) string {
		if val, ok := metadata[key]; ok {
			return fmt.Sprintf("%v", val.Value())
		}
		return "0"
	}

	playerName := strings.TrimPrefix(player, "org.mpris.MediaPlayer2.")

	artwork, _ := media.HandleArtworkRequest(getString("mpris:artUrl"))

	return models.MediaInfo{
		Title:    getString("xesam:title"),
		Artist:   getString("xesam:artist"),
		Album:    getString("xesam:album"),
		Artwork:  artwork,
		Length:   getInt64("mpris:length"),
		Status:   statusVar.Value().(string),
		Position: fmt.Sprintf("%v", positionVar.Value()),
		Player:   playerName,
	}, nil
}

func next() (success bool, err error) {
	dbusConn, err := GetDBusConnection()
	if err != nil {
		fmt.Println("Error connecting to D-Bus:", err)
		return false, err
	}

	player, err := findActivePlayer(dbusConn)
	if err != nil {
		fmt.Println("No active player found:", err)
		return false, err
	}
	obj := dbusConn.Object(player, "/org/mpris/MediaPlayer2")
	err = obj.Call("org.mpris.MediaPlayer2.Player.Next", 0).Store()
	return err == nil, err
}

func previous() (success bool, err error) {
	dbusConn, err := GetDBusConnection()
	if err != nil {
		fmt.Println("Error connecting to D-Bus:", err)
		return false, err
	}

	player, err := findActivePlayer(dbusConn)
	if err != nil {
		fmt.Println("No active player found:", err)
		return false, err
	}
	obj := dbusConn.Object(player, "/org/mpris/MediaPlayer2")
	err = obj.Call("org.mpris.MediaPlayer2.Player.Previous", 0).Store()
	return err == nil, err
}

func playPause() (success bool, err error) {
	dbusConn, err := GetDBusConnection()
	if err != nil {
		fmt.Println("Error connecting to D-Bus:", err)
		return false, err
	}

	player, err := findActivePlayer(dbusConn)
	if err != nil {
		fmt.Println("No active player found:", err)
		return false, err
	}
	obj := dbusConn.Object(player, "/org/mpris/MediaPlayer2")
	err = obj.Call("org.mpris.MediaPlayer2.Player.PlayPause", 0).Store()
	return err == nil, err
}

func seekForward() (success bool, err error) {
	dbusConn, err := GetDBusConnection()
	if err != nil {
		fmt.Println("Error connecting to D-Bus:", err)
		return false, err
	}

	player, err := findActivePlayer(dbusConn)
	if err != nil {
		fmt.Println("No active player found:", err)
		return false, err
	}

	obj := dbusConn.Object(player, "/org/mpris/MediaPlayer2")
	err = obj.Call("org.mpris.MediaPlayer2.Player.Seek", 0, int64(5000000)).Store()
	return err == nil, err
}

func seekBackward() (success bool, err error) {
	dbusConn, err := GetDBusConnection()
	if err != nil {
		fmt.Println("Error connecting to D-Bus:", err)
		return false, err
	}

	player, err := findActivePlayer(dbusConn)
	if err != nil {
		fmt.Println("No active player found:", err)
		return false, err
	}
	obj := dbusConn.Object(player, "/org/mpris/MediaPlayer2")
	err = obj.Call("org.mpris.MediaPlayer2.Player.Seek", 0, int64(-5000000)).Store()
	return err == nil, err
}

func seekToPosition(position int64) (success bool, err error) {
	dbusConn, err := GetDBusConnection()
	if err != nil {
		fmt.Println("Error connecting to D-Bus:", err)
		return false, err
	}
	player, err := findActivePlayer(dbusConn)
	if err != nil {
		fmt.Println("No active player found:", err)
		return false, err
	}
	obj := dbusConn.Object(player, "/org/mpris/MediaPlayer2")
	err = obj.Call("org.mpris.MediaPlayer2.Player.SetPosition", 0, dbus.ObjectPath("/org/mpris/MediaPlayer2/TrackList/NoTrack"), position).Store()
	return err == nil, err
}

func setVolume(volume int) (success bool, err error) {
	dbusConn, err := GetDBusConnection()
	if err != nil {
		return false, err
	}

	player, err := findActivePlayer(dbusConn)
	if err != nil {
		return false, err
	}
	obj := dbusConn.Object(player, "/org/mpris/MediaPlayer2")

	// FIX: Volume is a Property, so we use SetProperty, not Call
	volFloat := float64(volume) / 100.0
	err = obj.SetProperty("org.mpris.MediaPlayer2.Player.Volume", dbus.MakeVariant(volFloat))

	return err == nil, err
}

var LinuxDBusPlayer = models.PlayerFunctions{
	PlayPause:                playPause,
	SeekForward:              seekForward,
	SeekBackward:             seekBackward,
	Next:                     next,
	Prev:                     previous,
	SeekTo:                   seekToPosition,
	SetVolume:                setVolume,
	GetCurrentPlayerMetadata: getMetadata,
	GetAllPlayers: func() ([]string, error) {
		conn, err := GetDBusConnection()
		if err != nil {
			return nil, err
		}
		var names []string
		err = conn.BusObject().Call("org.freedesktop.DBus.ListNames", 0).Store(&names)
		if err != nil {
			return nil, err
		}
		var players []string
		for _, name := range names {
			if strings.HasPrefix(name, "org.mpris.MediaPlayer2.") {
				players = append(players, strings.TrimPrefix(name, "org.mpris.MediaPlayer2."))
			}
		}
		return players, nil
	},
}

func initializeDefaultPlayer() models.PlayerFunctions {
	return LinuxDBusPlayer
}
