package models

import (
	"encoding/json"
	"net/http"
)

func (r DefaultResponse) Response(w http.ResponseWriter, header int) {
	SetHeader(w, header)
	json.NewEncoder(w).Encode(DefaultResponse{Type: r.Type, Message: r.Message})
}

func SetHeader(w http.ResponseWriter, header int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(header)
}
