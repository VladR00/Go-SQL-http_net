package models

import (
	"encoding/json"
	"net/http"
	"time"
)

func (r DefaultResponse) Response(w http.ResponseWriter, header int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(header)
	json.NewEncoder(w).Encode(DefaultResponse{Type: r.Type, Message: r.Message})
}

type DefaultResponse struct {
	Type    string `json:"type"`    // Error | Data | Message
	Message string `json:"message"` // Message
}

type Department struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	Name       string    `json:"name" gorm:"size:200;not null"`
	Parent_id  *int      `json:"parent_id" gorm:"column:parent_id"`
	Created_at time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (Department) TableName() string {
	return "department"
}

type DepartmentRequest struct {
	Name      string `json:"name"`
	Parent_id *int   `json:"parent_id"`
}
