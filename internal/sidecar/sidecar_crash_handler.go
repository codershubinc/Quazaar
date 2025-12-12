package sidecar

import (
	"encoding/json"
	"log"
)

// HandleCrash logs the crash report as JSON so it can be picked up by monitoring tools
func HandleCrash(report map[string]string) {
	data, _ := json.Marshal(report)
	log.Printf("SIDECAR_CRASH: %s", string(data))
}
