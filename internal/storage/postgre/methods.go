package postgre

import (
	"Go-SQL-http_net/internal/models"
)

func (s *Storage) CreateDepartment(params models.DepartmentRequest) (*models.Department, error) {
	department := &models.Department{
		Name:      params.Name,
		Parent_id: params.Parent_id,
	}
	result := s.Db.Create(department)
	if result.Error != nil {
		return &models.Department{}, result.Error
	}
	return department, nil
}
