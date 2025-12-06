package systemVolume

import (
	"Quazaar/pkg/models"
	"fmt"
)

type VolumeWSMessage struct {
	Action string                 `json:"action"`
	Volume int                    `json:"volume"`
	Step   int                    `json:"step"`
	Meta   map[string]interface{} `json:"meta,omitempty"`
}
type VolumeResponse struct {
	CurrentVolume  int    `json:"current_volume"` // 0-100
	PreviousVolume int    `json:"previous_volume"`
	Action         string `json:"action"` // "inc", "dec", "set", "get"
	Success        bool   `json:"success"`
	Error          string `json:"error,omitempty"`
	Timestamp      int64  `json:"timestamp"`
}

func HandleVolumeWS(msg any) (models.ServerResponse, error) {
	volMsg := msg.(map[string]interface{})
	action, ok := volMsg["action"].(string)

	if !ok {
		return models.ServerResponse{}, fmt.Errorf("no action found in msg ..l..l,,l,l,l,,")
	}
	switch action {
	case "inc":
		success, cVol, err := IncreaseSystemVolume()
		if err != nil {
			return models.ServerResponse{Status: "error", Data: VolumeResponse{Error: err.Error()}}, err
		}
		return models.ServerResponse{Status: "success", Message: "system", Data: VolumeResponse{
			CurrentVolume:  cVol,
			PreviousVolume: cVol - 5,
			Action:         "inc",
			Success:        success,
		}}, nil
	case "dec":
		success, cVol, err := DecreaseSystemVolume()
		if err != nil {
			return models.ServerResponse{Status: "error", Data: VolumeResponse{Error: err.Error()}}, err
		}
		return models.ServerResponse{Status: "success", Message: "system", Data: VolumeResponse{
			CurrentVolume:  cVol,
			PreviousVolume: cVol + 5,
			Action:         "dec",
			Success:        success,
		}}, nil
	case "set":
		sickVol, ok := volMsg["set_to_vol"].(float64)
		if !ok {
			fmt.Println("err to get  set_to_vol")
			return models.ServerResponse{Status: "error"}, fmt.Errorf("invalid set_to_val")
		}
		sickVolInt := int(sickVol)
		success, cVol, err := SickSystemSetVolume("set", sickVolInt)
		if err != nil {
			return models.ServerResponse{Status: "error", Data: VolumeResponse{Error: err.Error()}}, err
		}
		return models.ServerResponse{Status: "success", Message: "system", Data: VolumeResponse{
			CurrentVolume:  cVol,
			PreviousVolume: -1,
			Action:         "set",
			Success:        success,
		}}, nil
	case "get":
		cVol, err := CurrentSystemVolume()
		if err != nil {
			return models.ServerResponse{Status: "error", Data: VolumeResponse{Error: err.Error()}}, err
		}
		return models.ServerResponse{Status: "success", Message: "system", Data: VolumeResponse{
			CurrentVolume: cVol,
			Action:        "get",
			Success:       true,
		}}, nil

	default:
		return models.ServerResponse{Status: "error"}, fmt.Errorf("no action found or invalid action")
	}
}
