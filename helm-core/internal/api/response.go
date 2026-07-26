package api

import (
	"encoding/json"
	"net/http"
)

// apiError represents the standard error response wrapper.
type apiError struct {
	Error errorDetail `json:"error"`
}

type errorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// respondJSON writes a JSON response with specified status code and payload.
func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// respondError writes a standardized JSON error response.
func respondError(w http.ResponseWriter, status int, code string, message string) {
	respondJSON(w, status, apiError{
		Error: errorDetail{
			Code:    code,
			Message: message,
		},
	})
}
