package systemInfo

import (
	"encoding/json"
	"net/http"
)

func GetSystemInfoApi(w http.ResponseWriter, r *http.Request) {
	info := getUser()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(info)

}
