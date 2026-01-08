package bluetooth

import (
	"Quazaar/pkg/helpers"
	"fmt"
	"regexp"
	"strings"

	"github.com/godbus/dbus/v5"
)

type BluetoothDevice struct {
	Name         string `json:"name"`
	MACAddress   string `json:"mac"`
	Battery      int    `json:"battery"`      // Average battery, -1 if not available
	BatteryLeft  int    `json:"batteryLeft"`  // Left earbud battery, -1 if not available
	BatteryRight int    `json:"batteryRight"` // Right earbud battery, -1 if not available
	BatteryCase  int    `json:"batteryCase"`  // Case battery, -1 if not available
	Icon         string `json:"icon"`
	Connected    bool   `json:"connected"`
}

// GetBluetoothDevices returns a list of connected Bluetooth devices with battery info using DBus
func GetBluetoothDevices() ([]BluetoothDevice, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to system bus: %v", err)
	}

	// Get all objects managed by BlueZ
	var objects map[dbus.ObjectPath]map[string]map[string]dbus.Variant
	err = conn.Object("org.bluez", "/").Call("org.freedesktop.DBus.ObjectManager.GetManagedObjects", 0).Store(&objects)
	if err != nil {
		return nil, fmt.Errorf("failed to get managed objects: %v", err)
	}

	var devices []BluetoothDevice

	for _, interfaces := range objects {
		// Check if it's a device
		deviceProps, ok := interfaces["org.bluez.Device1"]
		if !ok {
			continue
		}

		// Check if connected
		if connected, ok := deviceProps["Connected"].Value().(bool); !ok || !connected {
			continue
		}

		device := BluetoothDevice{
			Battery:      -1,
			BatteryLeft:  -1,
			BatteryRight: -1,
			BatteryCase:  -1,
			Connected:    true,
		}

		// Get Name
		if name, ok := deviceProps["Name"].Value().(string); ok {
			device.Name = name
		} else if alias, ok := deviceProps["Alias"].Value().(string); ok {
			device.Name = alias
		} else {
			device.Name = "Unknown Device"
		}

		// Get Address
		if addr, ok := deviceProps["Address"].Value().(string); ok {
			device.MACAddress = addr
		}

		// Get Icon
		if icon, ok := deviceProps["Icon"].Value().(string); ok {
			device.Icon = icon
		} else {
			device.Icon = "bluetooth"
		}

		// Check for Battery1 interface on the same object
		if batteryProps, ok := interfaces["org.bluez.Battery1"]; ok {
			if percentage, ok := batteryProps["Percentage"].Value().(byte); ok {
				device.Battery = int(percentage)
			}
		}

		// Try to get individual battery info using GalaxyBudsClient or earbuds CLI
		if strings.Contains(strings.ToLower(device.Name), "galaxy buds") ||
			strings.Contains(strings.ToLower(device.Name), "buds") {
			tryGalaxyBudsTools(&device, device.MACAddress)
		}

		devices = append(devices, device)
	}

	return devices, nil
}

// tryGalaxyBudsTools attempts to get individual battery info using specialized Galaxy Buds tools
func tryGalaxyBudsTools(device *BluetoothDevice, mac string) {
	// Try GalaxyBudsClient CLI if available (https://github.com/ThePBone/GalaxyBudsClient)
	// Install: yay -S galaxybudsclient-bin
	output, err := helpers.SpawnProcess("galaxybudsclient", []string{"--address", mac, "--get-battery"})
	if err == nil {
		parseGalaxyBudsClientOutput(device, string(output))
		return
	}
}

// parseGalaxyBudsClientOutput parses output from GalaxyBudsClient
func parseGalaxyBudsClientOutput(device *BluetoothDevice, output string) {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		lower := strings.ToLower(line)
		re := regexp.MustCompile(`(\d+)%?`)

		if strings.Contains(lower, "left") {
			if matches := re.FindStringSubmatch(line); len(matches) > 1 {
				var percent int
				fmt.Sscanf(matches[1], "%d", &percent)
				device.BatteryLeft = percent
			}
		} else if strings.Contains(lower, "right") {
			if matches := re.FindStringSubmatch(line); len(matches) > 1 {
				var percent int
				fmt.Sscanf(matches[1], "%d", &percent)
				device.BatteryRight = percent
			}
		} else if strings.Contains(lower, "case") {
			if matches := re.FindStringSubmatch(line); len(matches) > 1 {
				var percent int
				fmt.Sscanf(matches[1], "%d", &percent)
				device.BatteryCase = percent
			}
		}
	}
}

// parseEarbudsOutput parses JSON output from earbuds tools (legacy - kept for compatibility)
func parseEarbudsOutput(device *BluetoothDevice, output string) {
	// Simple parsing - look for battery values in output
	// earbuds tool outputs format varies, so we'll parse common patterns
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		lower := strings.ToLower(line)
		if strings.Contains(lower, "left") && strings.Contains(lower, "battery") {
			re := regexp.MustCompile(`(\d+)%?`)
			if matches := re.FindStringSubmatch(line); len(matches) > 1 {
				var percent int
				fmt.Sscanf(matches[1], "%d", &percent)
				device.BatteryLeft = percent
			}
		} else if strings.Contains(lower, "right") && strings.Contains(lower, "battery") {
			re := regexp.MustCompile(`(\d+)%?`)
			if matches := re.FindStringSubmatch(line); len(matches) > 1 {
				var percent int
				fmt.Sscanf(matches[1], "%d", &percent)
				device.BatteryRight = percent
			}
		} else if strings.Contains(lower, "case") && strings.Contains(lower, "battery") {
			re := regexp.MustCompile(`(\d+)%?`)
			if matches := re.FindStringSubmatch(line); len(matches) > 1 {
				var percent int
				fmt.Sscanf(matches[1], "%d", &percent)
				device.BatteryCase = percent
			}
		}
	}
}

// tryDirectGalaxyBudsRead is deprecated - use GalaxyBudsClient instead
func tryDirectGalaxyBudsRead(device *BluetoothDevice, mac string) {
	// Placeholder - users should install GalaxyBudsClient for full functionality
	// Install on Arch: yay -S galaxybudsclient-bin
	_ = device
	_ = mac
}
