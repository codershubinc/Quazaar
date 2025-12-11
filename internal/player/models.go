package player

// Re-export models for convenience
// Main player types are defined in pkg/models

import (
	"Quazaar/pkg/models"
)

// PlayerFunctions - re-exported from models
type PlayerFunctions = models.PlayerFunctions

// MediaInfo - re-exported from models
type MediaInfo = models.MediaInfo

// PlayerClientRequest - re-exported from models
type PlayerClientRequest = models.PlayerClientRequest

// Local types specific to player package can be added here