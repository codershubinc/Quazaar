package player

import (
	"Quazaar/pkg/models"
	"fmt"
	"strconv"
)

func HandleWebSocket(msg any) (models.UniServerResponse, error) {

	currentPlayerFunc := GetCurrentPlayerFunctions()
	playerMsg := msg.(map[string]interface{})
	action, ok := playerMsg["action"].(string)
	if !ok {
		return models.UniServerResponse{}, fmt.Errorf("cant get the action")
	}

	switch action {
	case "toggle_play_pause":
		success, err := currentPlayerFunc.PlayPause()
		if err != nil {
			return models.UniServerResponse{}, err
		}
		return models.UniServerResponse{
			Distributor: "responce",
			Type:        "player",
			MsgOf:       "toggle_play_pause",
			Action:      "toggle_play_pause",
			Data:        nil,
			Success:     strconv.FormatBool(success),
		}, nil
	case "next":
		success, err := currentPlayerFunc.Next()
		if err != nil {
			return models.UniServerResponse{}, err
		}
		return models.UniServerResponse{
			Distributor: "responce",
			Type:        "player",
			MsgOf:       "next",
			Action:      "next",
			Data:        nil,
			Success:     strconv.FormatBool(success),
		}, nil
	case "prev":
		success, err := currentPlayerFunc.Prev()
		if err != nil {
			return models.UniServerResponse{}, err
		}
		return models.UniServerResponse{
			Distributor: "responce",
			Type:        "player",
			MsgOf:       "prev",
			Action:      "prev",
			Data:        nil,
			Success:     strconv.FormatBool(success),
		}, nil
	case "seek_prev":
		success, err := currentPlayerFunc.SeekBackward()
		if err != nil {
			return models.UniServerResponse{}, err
		}
		return models.UniServerResponse{
			Distributor: "responce",
			Type:        "player",
			MsgOf:       "seek_prev",
			Action:      "seek_prev",
			Data:        nil,
			Success:     strconv.FormatBool(success),
		}, nil
	case "seek_for":
		success, err := currentPlayerFunc.SeekBackward()
		if err != nil {
			return models.UniServerResponse{}, err
		}
		return models.UniServerResponse{
			Distributor: "responce",
			Type:        "player",
			MsgOf:       "seek_for",
			Action:      "seek_for",
			Data:        nil,
			Success:     strconv.FormatBool(success),
		}, nil

	default:
		return models.UniServerResponse{}, fmt.Errorf("something went wrong , check docs for action ")
	}

}
