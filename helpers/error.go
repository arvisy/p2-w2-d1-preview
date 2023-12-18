package helpers

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Status     string `json:"status"`
	Title      string `json:"title"`
	Detail     string `json:"detail"`
	StatusCode int    `json:"-"`
}

func SendErrorResponse(w http.ResponseWriter, errResponse ErrorResponse, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(errResponse)
	if err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
	}
}
