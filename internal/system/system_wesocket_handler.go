package system

import (
	systemBrightness "Quazaar/internal/system/brightness"
	systemVolume "Quazaar/internal/system/volume"
	"Quazaar/pkg/models"
	"fmt"
)

func HandleWebSocket(msg any) (models.ServerResponse, error) {
	sysMsg := msg.(map[string]interface{})
	msgOf, ok := sysMsg["msg_of"].(string)
	if !ok {
		return models.ServerResponse{}, fmt.Errorf("wtf , please  send msg_of in msg")
	}
	switch msgOf {
	case "volume":
		return systemVolume.HandleVolumeWS(msg)
	case "brightness":
		return systemBrightness.HandleBrightnessWS(msg)
	default:
		return models.ServerResponse{}, fmt.Errorf("invalid msg_of received %s", msgOf)
	}

}
