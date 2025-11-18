package helpers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func SendRequest(url string, requestType string, headers map[string]string, data any) (*http.Response, error) {
	if requestType == "" {
		requestType = http.MethodGet
	}

	client := &http.Client{}

	// Initialize headers if nil
	if headers == nil {
		headers = make(map[string]string)
	}

	var body io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonData)

		// Set Content-Type if not already set
		if _, hasContentType := headers["Content-Type"]; !hasContentType {
			headers["Content-Type"] = "application/json"
		}
	}

	req, err := http.NewRequest(requestType, url, body)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return client.Do(req)
}

func GetParsedRequestData(resp *http.Response) (any, error) {
	var data any
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func SendJsonDataToClient(w http.ResponseWriter, req *http.Response, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(req.StatusCode)
	defer w.(http.Flusher).Flush()
	return json.NewEncoder(w).Encode(data)
}

// HandleHTTPError logs the error and sends an HTTP error response
func HandleHTTPError(w http.ResponseWriter, err error, statusCode int, message string) bool {
	if err != nil {
		LogMessage(ERROR, "%s: %v", message, err)
		http.Error(w, message, statusCode)
		return true
	}
	return false
}
