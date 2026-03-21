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

func (s *Storage) CreateEmployee(params models.DepartmentRequest) (*models.Department, error) {

}

func (s *Storage) GetDepartment(params models.DepartmentRequest) (*models.Department, error) {
	return "", nil
}

func (s *Storage) RelocateDepartment(params models.DepartmentRequest) (*models.Department, error) {
	return "", nil
}

func (s *Storage) DeleteDepartment(params models.DepartmentRequest) (*models.Department, error) {
	return "", nil
}
