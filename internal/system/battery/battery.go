package systemBattery

import "github.com/godbus/dbus/v5"

type BatteryInfo struct {
	Percentage  float64 `json:"percentage"`
	State       string  `json:"state"`
	TimeToEmpty int64   `json:"time_to_empty"` // Seconds
	TimeToFull  int64   `json:"time_to_full"`  // Seconds
	EnergyRate  float64 `json:"energy_rate"`
}

func GetSystemBatteryInfo() (*BatteryInfo, error) {
	conn, err := dbus.ConnectSystemBus()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	obj := conn.Object("org.freedesktop.UPower", "/org/freedesktop/UPower/devices/DisplayDevice")

	var props map[string]dbus.Variant
	err = obj.Call("org.freedesktop.DBus.Properties.GetAll", 0, "org.freedesktop.UPower.Device").Store(&props)
	if err != nil {
		return nil, err
	}

	status := &BatteryInfo{}

	if v, ok := props["Percentage"]; ok {
		status.Percentage = v.Value().(float64)
	}
	if v, ok := props["State"]; ok {
		if val, ok := v.Value().(uint32); ok {
			status.State = mapState(val)
		}
	}
	if v, ok := props["TimeToEmpty"]; ok {
		status.TimeToEmpty = v.Value().(int64)
	}
	if v, ok := props["TimeToFull"]; ok {
		status.TimeToFull = v.Value().(int64)
	}
	if v, ok := props["EnergyRate"]; ok {
		status.EnergyRate = v.Value().(float64)
	}

	return status, nil
}

func mapState(state uint32) string {
	switch state {
	case 1:
		return "charging"
	case 2:
		return "discharging"
	case 3:
		return "empty"
	case 4:
		return "fully charged"
	case 5:
		return "pending charge"
	case 6:
		return "pending discharge"
	default:
		return "unknown"
	}
}
