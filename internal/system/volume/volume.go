package systemVolume

import (
	"Quazaar/pkg/helpers"
	"fmt"
	"strconv"
	"strings"
)

func CurrentSystemVolume() (int, error) {
	/*
		~ it is for , getting the system volume %
		~ Using D-Bus,
		~ D-Bus  is the default
		~ return -1 || int , err
	*/

	out, err := helpers.SpawnProcess("pactl", []string{"get-sink-volume", "@DEFAULT_SINK@"})
	if err != nil {
		fmt.Println("Error occurred for getting system volume", err)
		return -1, err
	}
	str := string(out)
	idx := strings.Index(str, "%")
	if idx == -1 {
		return -1, fmt.Errorf("could not parse the volume")
	}
	sub := str[:idx]

	slashIdx := strings.LastIndex(sub, "/")
	if slashIdx == -1 {
		return -1, fmt.Errorf("could not parse the volume")
	}
	vol, err := strconv.Atoi(strings.TrimSpace(sub[slashIdx+1:]))
	if err != nil {
		return -1, fmt.Errorf("failed to convert the volume string")
	}
	return vol, nil
}

// SetVolume sets the volume to a specific percentage (0-100)
func SetVolume(percent int) (bool, int, error) {
	if percent < 0 {
		percent = 0
	}
	// if percent > 100 {
	// 	percent = 100
	// }
	//removed cap of 100%

	_, err := helpers.SpawnProcess("pactl", []string{
		"set-sink-volume",
		"@DEFAULT_SINK@",
		fmt.Sprintf("%d%%", percent),
	})
	if err != nil {
		return false, -1, err
	}
	return true, percent, nil
}

func IncreaseSystemVolume() (bool, int, error) {
	curr, err := CurrentSystemVolume()
	if err != nil {
		return false, -1, err
	}
	return SetVolume(curr + 5)
}

// DecreaseSystemVolume decreases volume by 5% (floored at 0%)
func DecreaseSystemVolume() (bool, int, error) {
	curr, err := CurrentSystemVolume()
	if err != nil {
		return false, -1, err
	}
	return SetVolume(curr - 5)
}

func ToggleMute() (success bool, isMuted bool, err error) {
	_, err = helpers.SpawnProcess("pactl", []string{"set-sink-mute", "@DEFAULT_SINK@", "toggle"})
	if err != nil {
		return false, false, err
	}

	// Check new mute status
	out, err := helpers.SpawnProcess("pactl", []string{"get-sink-mute", "@DEFAULT_SINK@"})
	if err != nil {
		return true, false, nil
	}

	isMuted = strings.Contains(string(out), "yes")
	return true, isMuted, nil
}
