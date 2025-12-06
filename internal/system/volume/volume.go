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

func SickSystemSetVolume(action string, volTo int) (success bool, currentVolume int, err error) {
	var act string
	sysVol, err := CurrentSystemVolume()
	if err != nil {
		return false, -1, err
	}
	if sysVol >= 100 && action == "inc" {
		return false, -1, fmt.Errorf("volume out of bound , cant increase above 100")
	}
	if sysVol <= 0 && action == "dec" {
		return false, -1, fmt.Errorf("volume out of bound , cant decrease of the below 0")
	}
	if action == "" {
		return false, -1, fmt.Errorf("need the action for doing work, f*ck u")
	}

	switch action {
	case "dec":
		sysVol -= volTo
		if sysVol < 0 {
			sysVol = 0
		}
		act = fmt.Sprintf("%d%%", sysVol)
	case "inc":
		sysVol += volTo
		if sysVol > 100 {
			sysVol = 100
		}
		act = fmt.Sprintf("%d%%", sysVol)
	case "sick":
		sysVol = volTo
		if sysVol > 100 {
			sysVol = 100
		} else if sysVol < 0 {
			sysVol = 0
		}
		act = fmt.Sprintf("%d%%", sysVol)
	default:
		return false, -1, fmt.Errorf("are u doing mEth 🎋 ; what action  you send just see")
	}

	_, err = helpers.SpawnProcess("pactl", []string{
		"set-sink-volume",
		"@DEFAULT_SINK@",
		act,
	})
	if err != nil {
		return false, -1, err
	}

	return true, sysVol, nil
}

func IncreaseSystemVolume() (success bool, currentVolume int, err error) {
	success, cVol, err := SickSystemSetVolume("inc", 5)
	if err != nil {
		return false, -1, err
	}
	return success, cVol, nil
}
func DecreaseSystemVolume() (success bool, currentVolume int, err error) {
	success, cVol, err := SickSystemSetVolume("dec", 5)
	if err != nil {
		return false, -1, err
	}
	return success, cVol, nil
}
