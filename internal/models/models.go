package models

import (
	"time"
)

type DefaultResponse struct {
	Type    string `json:"type"`    // Error | Data | Message
	Message string `json:"message"` // Message
}

type Department struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	Name       string    `json:"name" gorm:"size:200;not null"`
	Parent_id  *int      `json:"parent_id" gorm:"column:parent_id"`
	Created_at time.Time `json:"created_at" gorm:"autoCreateTime; column:created_at"`
}

func (Department) TableName() string {
	return "department"
}

type Employee struct {
	ID            int        `json:"id" gorm:"primaryKey"`
	Department_id int        `json:"department_id" gorm:"column:department_id"`
	Full_name     string     `json:"full_name" gorm:"not null;size:200;column:full_name"`
	Position      string     `json:"position" gorm:"not null;size:200"`
	Hired_at      *time.Time `json:"hired_at" gorm:"column:hired_at"`
	Created_at    time.Time  `json:"created_at" gorm:"autoCreateTime; column:created_at"`
}

func (Employee) TableName() string {
	return "employee"
}

type DepartmentRequest struct {
	Name      string `json:"name"`
	Parent_id *int   `json:"parent_id"`
}
