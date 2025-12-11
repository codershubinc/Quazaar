package player

import (
	"Quazaar/pkg/models"
)

var CurrentPlayerFunc models.PlayerFunctions

func GetCurrentPlayerFunctions() models.PlayerFunctions {
	if CurrentPlayerFunc.PlayPause == nil {
		CurrentPlayerFunc = initializeDefaultPlayer()
	}
	return CurrentPlayerFunc
}
