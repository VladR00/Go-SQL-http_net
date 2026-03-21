package models

import (
	"time"
)

// DefaultResponse represents a generic API response
type DefaultResponse struct {
	Type    string      `json:"type"`    // Error | Data | Message
	Message string      `json:"message"` // Message
	Data    interface{} `json:"data,omitempty"`
}

// Department represents a department in the organization
type Department struct {
	ID        int          `json:"id" gorm:"primaryKey"`
	Name      string       `json:"name" gorm:"size:200;not null;uniqueIndex:idx_dept_name_parent,composite:parent_id"`
	ParentID  *int         `json:"parent_id" gorm:"column:parent_id;uniqueIndex:idx_dept_name_parent,composite:name"`
	CreatedAt time.Time    `json:"created_at" gorm:"autoCreateTime"`
	Employees []Employee   `json:"employees,omitempty" gorm:"foreignKey:DepartmentID;constraint:OnDelete:CASCADE"`
	Children  []Department `json:"children,omitempty" gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE"`
}

func (Department) TableName() string {
	return "department"
}

// Employee represents an employee in a department
type Employee struct {
	ID           int        `json:"id" gorm:"primaryKey"`
	DepartmentID int        `json:"department_id" gorm:"column:department_id;not null"`
	FullName     string     `json:"full_name" gorm:"size:200;not null"`
	Position     string     `json:"position" gorm:"size:200;not null"`
	HiredAt      *time.Time `json:"hired_at" gorm:"column:hired_at"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
}

func (Employee) TableName() string {
	return "employee"
}

// DepartmentRequest represents a request to create or update a department
type DepartmentRequest struct {
	Name     string `json:"name"`
	ParentID *int   `json:"parent_id"`
}

// DepartmentUpdateRequest represents a request to update a department
type DepartmentUpdateRequest struct {
	Name     *string `json:"name"`
	ParentID *int    `json:"parent_id"`
}

// EmployeeRequest represents a request to create an employee
type EmployeeRequest struct {
	FullName string     `json:"full_name"`
	Position string     `json:"position"`
	HiredAt  *time.Time `json:"hired_at"`
}

// DepartmentDeleteRequest represents parameters for deleting a department
type DepartmentDeleteRequest struct {
	Mode                   string `json:"mode" schema:"mode"` // cascade or reassign
	ReassignToDepartmentID *int   `json:"reassign_to_department_id" schema:"reassign_to_department_id"`
}

// DepartmentResponse represents a department with optional children and employees
type DepartmentResponse struct {
	ID        int                  `json:"id"`
	Name      string               `json:"name"`
	ParentID  *int                 `json:"parent_id"`
	CreatedAt time.Time            `json:"created_at"`
	Employees []Employee           `json:"employees,omitempty"`
	Children  []DepartmentResponse `json:"children,omitempty"`
}
