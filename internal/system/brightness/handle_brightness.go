package systemBrightness

import (
	"Quazaar/pkg/models"
	"fmt"
)

type BrightnessResponse struct {
	CurrentBrightness int    `json:"current_brightness"` // 0-100
	Action            string `json:"action"`             // "inc", "dec", "set", "get"
	Success           bool   `json:"success"`
	Error             string `json:"error,omitempty"`
}

func HandleBrightnessWS(msg any) (models.ServerResponse, error) {
	brMsg := msg.(map[string]interface{})
	action, ok := brMsg["action"].(string)

	if !ok {
		return models.ServerResponse{Status: "error"}, fmt.Errorf("no action found in brightness msg")
	}

	switch action {
	case "inc":
		err := IncreaseBrightness()
		if err != nil {
			return models.ServerResponse{Status: "error", Data: BrightnessResponse{Error: err.Error()}}, err
		}
		curr, _ := GetCurrent()
		return models.ServerResponse{Status: "success", Message: "system", Data: BrightnessResponse{
			CurrentBrightness: curr,
			Action:            "inc",
			Success:           true,
		}}, nil

	case "dec":
		err := DecreaseBrightness()
		if err != nil {
			return models.ServerResponse{Status: "error", Data: BrightnessResponse{Error: err.Error()}}, err
		}
		curr, _ := GetCurrent()
		return models.ServerResponse{Status: "success", Message: "system", Data: BrightnessResponse{
			CurrentBrightness: curr,
			Action:            "dec",
			Success:           true,
		}}, nil

	case "set":
		val, ok := brMsg["set_to"].(float64)
		if !ok {
			return models.ServerResponse{Status: "error"}, fmt.Errorf("invalid set_to value")
		}
		err := SetBrightness(int(val))
		if err != nil {
			return models.ServerResponse{Status: "error", Data: BrightnessResponse{Error: err.Error()}}, err
		}
		return models.ServerResponse{Status: "success", Message: "system", Data: BrightnessResponse{
			CurrentBrightness: int(val),
			Action:            "set",
			Success:           true,
		}}, nil

	case "get":
		curr, err := GetCurrent()
		if err != nil {
			return models.ServerResponse{Status: "error", Data: BrightnessResponse{Error: err.Error()}}, err
		}
		return models.ServerResponse{Status: "success", Message: "system", Data: BrightnessResponse{
			CurrentBrightness: curr,
			Action:            "get",
			Success:           true,
		}}, nil

	default:
		return models.ServerResponse{Status: "error"}, fmt.Errorf("unknown brightness action: %s", action)
	}
}
