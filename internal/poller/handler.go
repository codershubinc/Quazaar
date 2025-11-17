
package poller

import (
	"Quazaar/pkg/models"
	"Quazaar/internal/media"
	"Quazaar/internal/websocket"
	"fmt"
	"time"
)

func Handle() {
	// fmt.Println("Started poller Handler ....")

	Poller(1*time.Second, make(chan struct{}), func() {
		msg, err := media.GetPlayerInfo()

		if err != nil {
			fmt.Printf("⚠️ Failed to get player info: %v\n", err)
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

func QuiteChan() chan struct{} {
	quit := make(chan struct{})
	// close(quit)
	return quit
}
