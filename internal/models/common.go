package models

import (
	"encoding/json"
	"net/http"
)

// Response sends a JSON response with the given status code
func (r DefaultResponse) Response(w http.ResponseWriter, header int) {
	SetHeader(w, header)
	json.NewEncoder(w).Encode(DefaultResponse{Type: r.Type, Message: r.Message, Data: r.Data})
}

// SetHeader sets the Content-Type header and status code
func SetHeader(w http.ResponseWriter, header int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(header)
}

// JSONResponse is a helper to send JSON data with a status code
func JSONResponse(w http.ResponseWriter, status int, data interface{}) {
	SetHeader(w, status)
	json.NewEncoder(w).Encode(data)
}

// ErrorResponse sends an error response
func ErrorResponse(w http.ResponseWriter, status int, message string) {
	SetHeader(w, status)
	response := DefaultResponse{
		Type:    "Error",
		Message: message,
	}
	json.NewEncoder(w).Encode(response)
}
