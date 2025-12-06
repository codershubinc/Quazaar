package systemBrightness

import (
	"Quazaar/pkg/helpers"
	"fmt"
	"strconv"
	"strings"
)

func GetCurrent() (int, error) {
	maxB, err := helpers.SpawnProcess("brightnessctl", []string{"max"})
	if err != nil {
		fmt.Println("Error while getting max brightness")
		return -1, err
	}
	curB, err := helpers.SpawnProcess("brightnessctl", []string{"get"})
	if err != nil {
		fmt.Println("Error while getting current brightness")
		return -1, err
	}
	maxBInt, err := strconv.Atoi(strings.TrimSpace(string(maxB)))
	if err != nil {
		fmt.Println("Error while converting max brightness")
		return -1, err
	}
	curBInt, err := strconv.Atoi(strings.TrimSpace(string(curB)))
	if err != nil {
		fmt.Println("Error while converting current brightness")
		return -1, err
	}
	return (curBInt * 100) / maxBInt, nil
}

func SetBrightness(percent int) error {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}

	_, err := helpers.SpawnProcess("brightnessctl", []string{"set", fmt.Sprintf("%d%%", percent)})
	return err
}

func IncreaseBrightness() error {
	_, err := helpers.SpawnProcess("brightnessctl", []string{"set", "+5%"})
	return err
}

func DecreaseBrightness() error {
	_, err := helpers.SpawnProcess("brightnessctl", []string{"set", "5%-"})
	return err
}
