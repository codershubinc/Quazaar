package spotify

import "fmt"

func HandleWebSocketMessage(message string, req any) {
	switch message {
	case "user":
		fmt.Println("User got")
	case "devices":
		fmt.Println("Devices got")
	case "playback":
		fmt.Println("Playback got")
	case "isSaved":
		fmt.Println("IsSaved got")
	case "saveTrack":
		fmt.Println("SaveTrack got")
	case "removeTrack":
		fmt.Println("RemoveTrack got")
	case "savedTracks":
		fmt.Println("SavedTracks got")
	default:
		fmt.Println("Unknown message")
	}
}
