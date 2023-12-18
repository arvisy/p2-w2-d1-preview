package helpers

import (
	"encoding/json"
	"net/http"
)

type SuccessResponse struct {
	Status     string `json:"status"`
	Title      string `json:"title"`
	Detail     string `json:"detail"`
	StatusCode int    `json:"-"`
}

func SendSuccessResponse(w http.ResponseWriter, successResponse SuccessResponse, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(successResponse)
	if err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
	}
}
