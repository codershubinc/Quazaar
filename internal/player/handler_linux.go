package player

import (
	"Quazaar/pkg/models"
	"fmt"
	"strings"
	"sync"

	"github.com/godbus/dbus/v5"
)

var (
	dbusConn *dbus.Conn
	dbusOnce sync.Once
)

func GetDBusConnection() (*dbus.Conn, error) {
	var err error
	dbusOnce.Do(func() {
		dbusConn, err = dbus.SessionBus()
	})
	return dbusConn, err
}

func getMetadata() (models.MediaInfo, error) {
	conn, err := GetDBusConnection()
	if err != nil {
		return models.MediaInfo{}, err
	}

	player := "org.mpris.MediaPlayer2.spotify"
	obj := conn.Object(player, "/org/mpris/MediaPlayer2")

	// ... rest of your code ...
	// ...existing code...
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

	return models.MediaInfo{
		Title:    getString("xesam:title"),
		Artist:   getString("xesam:artist"),
		Album:    getString("xesam:album"),
		Artwork:  getString("mpris:artUrl"),
		Length:   getInt64("mpris:length"),
		Status:   statusVar.Value().(string),
		Position: fmt.Sprintf("%v", positionVar.Value()),
		Player:   "spotify",
	}, nil
}

func next() (success bool, err error) {
	dbusConn, err := GetDBusConnection()
	if err != nil {
		fmt.Println("Error connecting to D-Bus:", err)
		return false, err
	}

	player := "org.mpris.MediaPlayer2.spotify"
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

	player := "org.mpris.MediaPlayer2.spotify"
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

	player := "org.mpris.MediaPlayer2.spotify"
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

	player := "org.mpris.MediaPlayer2.spotify"
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

	player := "org.mpris.MediaPlayer2.spotify"
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
	player := "org.mpris.MediaPlayer2.spotify"
	obj := dbusConn.Object(player, "/org/mpris/MediaPlayer2")
	err = obj.Call("org.mpris.MediaPlayer2.Player.SetPosition", 0, dbus.ObjectPath("/org/mpris/MediaPlayer2/TrackList/NoTrack"), position).Store()
	return err == nil, err
}

func setVolume(volume int) (success bool, err error) {
	dbusConn, err := GetDBusConnection()
	if err != nil {
		return false, err
	}

	player := "org.mpris.MediaPlayer2.spotify"
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
	// for now just focus on spotify buddy ..........
	GetAllPlayers: func() ([]string, error) {
		return []string{"spotify"}, nil
	},
}
