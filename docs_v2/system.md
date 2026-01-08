# System API

## `/api/v0.1/system/wifi` [GET]

**Description:** Retrieves current WiFi connection status and network speed.

**Required Headers:**
- `Authorization`: `Bearer <token>`

**Response Data (JSON):**
```json
{
  "success": true,
  "wifi": {
    "ssid": "MyHomeNetwork",
    "signalStrength": 85,
    "linkSpeed": 866,
    "frequency": "5 GHz",
    "security": "WPA2",
    "ipAddress": "192.168.1.105",
    "connected": true,
    "downloadSpeed": 12.5, // Mbps
    "uploadSpeed": 5.2,    // Mbps
    "interface": "wlan0",
    "unitOfSpeed": "Mbps"
  }
}
```

## `/api/v0.1/system/bluetooth` [GET]

**Description:** Lists connected Bluetooth devices and their battery levels (if available).

**Required Headers:**
- `Authorization`: `Bearer <token>`

**Response Data (JSON):**
```json
{
  "success": true,
  "count": 1,
  "devices": [
    {
      "name": "Galaxy Buds Pro",
      "mac": "AA:BB:CC:11:22:33",
      "battery": 85,
      "batteryLeft": 80,
      "batteryRight": 90,
      "batteryCase": 50,
      "icon": "audio-headset",
      "connected": true
    }
  ]
}
```

## `/api/v0.1/system/volume` [GET]

**Description:** Gets the current system audio volume percentage.

**Required Headers:**
- `Authorization`: `Bearer <token>`

**Response Data (JSON):**
```json
{
  "success": true,
  "volume": 65
}
```

## `/api/v0.1/system/brightness` [GET]

**Description:** Gets the current screen brightness percentage.

**Required Headers:**
- `Authorization`: `Bearer <token>`

**Response Data (JSON):**
```json
{
  "success": true,
  "brightness": 70
}
```

## `/api/v0.1/system/sound/devices` [GET]

**Description:** Lists available audio output devices (sinks).

**Required Headers:**
- `Authorization`: `Bearer <token>`

**Response Data (JSON):**
```json
{
  "success": true,
  "devices": [
    {
      "id": "1",
      "name": "alsa_output.pci-0000_00_1f.3.analog-stereo",
      "description": "Built-in Audio Analog Stereo",
      "active": true
    },
    {
      "id": "2",
      "name": "bluez_sink.AA_BB_CC_DD_EE_FF.a2dp_sink",
      "description": "Galaxy Buds Pro",
      "active": false
    }
  ]
}
```

## `/api/v0.1/system/sound/device` [POST]

**Description:** Sets the default audio output device.

**Required Headers:**
- `Authorization`: `Bearer <token>`

**Required Body (JSON):**
```json
{
  "deviceName": "alsa_output.pci-0000_00_1f.3.analog-stereo"
}
```

**Response Data (JSON):**
```json
{
  "success": true,
  "message": "Audio device updated successfully"
}
```

## `/api/v0.1/system/wakatime` [GET]

**Description:** Retrieves today's coding stats from WakaTime (requires WAKATIME_API_KEY env var).

**Required Headers:**
- `Authorization`: `Bearer <token>`

**Response Data (JSON):**
```json
{
  "success": true,
  "stats": {
    "cummulative_total": {
      "text": "3 hrs 45 mins",
      "seconds": 13500
    },
    "data": [ "..." ]
  }
}
```

## `/api/v0.1/system/battery` [GET]

**Description:** Gets system battery status (No Auth Required).

**Response Data (JSON):**
```json
{
  "percentage": 85.5,
  "state": "discharging",
  "time_to_empty": 10800, // seconds
  "time_to_full": 0,
  "energy_rate": 12.5
}
```
