package systemSound

import (
	"Quazaar/pkg/helpers"
	"fmt"
	"strings"
)

type AudioDevice struct {
	ID          string `json:"id"`          // Use string for ID to handle both index and name if needed, but usually index for pcts
	Name        string `json:"name"`        // internal name (alsa_output...)
	Description string `json:"description"` // Friendly name
	Active      bool   `json:"active"`
}

// ListAudioDevices returns available audio output devices (sinks)
func ListAudioDevices() ([]AudioDevice, error) {
	// 1. Get default sink
	defaultSinkBytes, err := helpers.SpawnProcess("pactl", []string{"get-default-sink"})
	if err != nil {
		fmt.Printf("Error getting default sink: %v\n", err)
		// Continue even if getting default fails, though it's weird
	}
	defaultSink := strings.TrimSpace(string(defaultSinkBytes))

	// 2. List sinks
	// using 'pactl list sinks' provides descriptions
	outBytes, err := helpers.SpawnProcess("pactl", []string{"list", "sinks"})
	if err != nil {
		return nil, err
	}
	output := string(outBytes)

	var devices []AudioDevice

	// Manual parsing of the pactl output
	// Blocks start with "Sink #N"
	// Properties are indented

	lines := strings.Split(output, "\n")
	var currentDevice *AudioDevice

	// A more robust parsing loop considering structure
	for _, rawLine := range lines {
		trimmed := strings.TrimSpace(rawLine)

		if strings.HasPrefix(rawLine, "Sink #") {
			// Save previous device if exists
			if currentDevice != nil {
				// Filter out HDMI/DisplayPort devices
				if !isHiddenDevice(currentDevice) {
					// Check if it was active
					if currentDevice.Name == defaultSink {
						currentDevice.Active = true
					}
					devices = append(devices, *currentDevice)
				}
			}

			// Start new device
			ids := strings.TrimPrefix(rawLine, "Sink #")
			currentDevice = &AudioDevice{
				ID: ids,
			}
			continue
		}

		if currentDevice == nil {
			continue
		}

		// Property parsing
		if strings.HasPrefix(trimmed, "Name: ") {
			currentDevice.Name = strings.TrimPrefix(trimmed, "Name: ")
		} else if strings.HasPrefix(trimmed, "Description: ") {
			currentDevice.Description = strings.TrimPrefix(trimmed, "Description: ")
		}
	}

	// Append the last one
	if currentDevice != nil {
		if !isHiddenDevice(currentDevice) {
			if currentDevice.Name == defaultSink {
				currentDevice.Active = true
			}
			devices = append(devices, *currentDevice)
		}
	}

	return devices, nil
}

// isHiddenDevice checks if the device should be hidden (HDMI, DisplayPort)
func isHiddenDevice(d *AudioDevice) bool {
	lowerName := strings.ToLower(d.Name)
	lowerDesc := strings.ToLower(d.Description)
	return strings.Contains(lowerName, "hdmi") ||
		strings.Contains(lowerDesc, "hdmi") ||
		strings.Contains(lowerName, "displayport") ||
		strings.Contains(lowerDesc, "displayport")
}

// SetAudioDevice sets the default output device
func SetAudioDevice(deviceName string) error {
	_, err := helpers.SpawnProcess("pactl", []string{"set-default-sink", deviceName})
	return err
}
