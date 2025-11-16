package system

import "Quazaar/pkg/helpers"

func LaunchApp(appName string) (string, error) {
	output, err := helpers.SpawnProcess(
		`gtk-launch`,
		[]string{appName},
	)
	if err != nil {
		return "", err
	}

	return string(output), nil
}
