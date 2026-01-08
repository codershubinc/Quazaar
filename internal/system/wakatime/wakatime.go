package wakatime

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const WakaTimeBaseURL = "https://wakatime.com/api/v1/users/current/summaries?range=Today"

func GetWakaTimeStats() (map[string]interface{}, error) {
	apiKey := os.Getenv("WAKATIME_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("WAKATIME_API_KEY environment variable is not set")
	}

	url := WakaTimeBaseURL + "&api_key=" + apiKey
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wakatime api returned status: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}
