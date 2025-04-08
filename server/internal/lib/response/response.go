package response

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse structure for a standard JSON error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// RespondError sends a JSON response with the error and the specified status code.
func RespondError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

// RespondJSON sends data as a JSON response with the specified status code.
func RespondJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}
