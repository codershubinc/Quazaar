package poller

import (
	"Quazaar/internal/logger"
	"Quazaar/internal/player"
	"Quazaar/internal/websocket"
	"Quazaar/pkg/models"
	"context"
	"time"
)

func Handle(ctx context.Context) {
	// fmt.Println("Started poller Handler ....")

	Poller(ctx, 1*time.Second, func() {
		msg, err := player.GetCurrentPlayerFunctions().GetCurrentPlayerMetadata()

		if err != nil {
			logger.Error("⚠️ Failed to get player info", "error", err)
			return
		}

		websocket.WriteChannelMessage(
			models.ServerResponse{
				Status:  "success",
				Message: "media_info",
				Data:    msg,
			},
		)
	})
}
