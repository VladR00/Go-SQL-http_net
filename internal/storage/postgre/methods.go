package postgre

import (
	"Go-SQL-http_net/internal/models"
	"errors"
	"strings"

	"gorm.io/gorm"
)

// CreateDepartment creates a new department
func (s *Storage) CreateDepartment(params models.DepartmentRequest) (*models.Department, error) {
	// Validate name
	if err := validateDepartmentName(params.Name); err != nil {
		return nil, err
	}

	// If parent_id is provided, verify it exists
	if params.ParentID != nil {
		if exists, err := s.departmentExists(*params.ParentID); err != nil {
			return nil, err
		} else if !exists {
			return nil, errors.New("parent department not found")
		}
	}

	department := &models.Department{
		Name:     params.Name,
		ParentID: params.ParentID,
	}

	result := s.Db.Create(department)
	if result.Error != nil {
		// Check if it's a unique constraint violation
		if strings.Contains(result.Error.Error(), "duplicate key") || strings.Contains(result.Error.Error(), "violates unique") {
			return nil, errors.New("department with this name already exists under the same parent")
		}
		return nil, result.Error
	}

	return department, nil
}

// CreateEmployee creates a new employee in a department
func (s *Storage) CreateEmployee(departmentID int, params models.EmployeeRequest) (*models.Employee, error) {
	// Verify department exists
	if exists, err := s.departmentExists(departmentID); err != nil {
		return nil, err
	} else if !exists {
		return nil, errors.New("department not found")
	}

	// Validate employee data
	if err := validateEmployeeName(params.FullName); err != nil {
		return nil, err
	}

	if err := validatePosition(params.Position); err != nil {
		return nil, err
	}

	employee := &models.Employee{
		DepartmentID: departmentID,
		FullName:     strings.TrimSpace(params.FullName),
		Position:     strings.TrimSpace(params.Position),
		HiredAt:      params.HiredAt,
	}

	result := s.Db.Create(employee)
	if result.Error != nil {
		return nil, result.Error
	}

	return employee, nil
}

// GetDepartment retrieves a department with optional children and employees
func (s *Storage) GetDepartment(id int, depth int, includeEmployees bool) (*models.DepartmentResponse, error) {
	if depth < 0 || depth > 5 {
		depth = 1
	}

	dept := &models.Department{}
	result := s.Db.First(dept, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("department not found")
		}
		return nil, result.Error
	}

	response := &models.DepartmentResponse{
		ID:        dept.ID,
		Name:      dept.Name,
		ParentID:  dept.ParentID,
		CreatedAt: dept.CreatedAt,
	}

	// Get employees if requested
	if includeEmployees {
		employees := []models.Employee{}
		s.Db.Where("department_id = ?", id).Order("created_at ASC").Find(&employees)
		response.Employees = employees
	}

	// Get children if depth > 0
	if depth > 0 {
		response.Children = []models.DepartmentResponse{}
		children := []models.Department{}
		s.Db.Where("parent_id = ?", id).Order("created_at ASC").Find(&children)

		for _, child := range children {
			childResp, err := s.GetDepartment(child.ID, depth-1, includeEmployees)
			if err != nil {
				continue
			}
			response.Children = append(response.Children, *childResp)
		}
	}

	return response, nil
}

// UpdateDepartment updates a department's name and/or parent
func (s *Storage) UpdateDepartment(id int, params models.DepartmentUpdateRequest) (*models.Department, error) {
	dept := &models.Department{}
	result := s.Db.First(dept, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("department not found")
		}
		return nil, result.Error
	}

	// Check if trying to move to itself
	if params.ParentID != nil && *params.ParentID == id {
		return nil, errors.New("cannot move department to itself")
	}

	// If changing parent, check for cycles
	if params.ParentID != nil && (dept.ParentID == nil || *dept.ParentID != *params.ParentID) {
		if hasCycle, err := s.checkCycle(id, *params.ParentID); err != nil {
			return nil, err
		} else if hasCycle {
			return nil, errors.New("cycle detected in department hierarchy")
		}

		// Verify new parent exists
		if exists, err := s.departmentExists(*params.ParentID); err != nil {
			return nil, err
		} else if !exists {
			return nil, errors.New("parent department not found")
		}
	}

	// Update name if provided
	if params.Name != nil {
		if err := validateDepartmentName(*params.Name); err != nil {
			return nil, err
		}
		dept.Name = strings.TrimSpace(*params.Name)
	}

	// Update parent if provided
	if params.ParentID != nil {
		dept.ParentID = params.ParentID
	}

	result = s.Db.Save(dept)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") || strings.Contains(result.Error.Error(), "violates unique") {
			return nil, errors.New("department with this name already exists under the same parent")
		}
		return nil, result.Error
	}

	return dept, nil
}

// DeleteDepartment deletes a department with specified mode
func (s *Storage) DeleteDepartment(id int, mode string, reassignToID *int) error {
	// Verify department exists
	dept := &models.Department{}
	result := s.Db.First(dept, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("department not found")
		}
		return result.Error
	}

	switch mode {
	case "cascade":
		// Delete all departments and employees in cascade (GORM handles this with ON DELETE CASCADE)
		result := s.Db.Delete(dept)
		if result.Error != nil {
			return result.Error
		}
		return nil

	case "reassign":
		if reassignToID == nil {
			return errors.New("reassign_to_department_id is required when mode is reassign")
		}

		// Verify target department exists
		targetDept := &models.Department{}
		result := s.Db.First(targetDept, *reassignToID)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return errors.New("target department not found")
			}
			return result.Error
		}

		// Update all child departments to have new parent
		result = s.Db.Model(&models.Department{}).
			Where("parent_id = ?", id).
			Update("parent_id", *reassignToID)
		if result.Error != nil {
			return result.Error
		}

		// Update all employees to new department
		result = s.Db.Model(&models.Employee{}).
			Where("department_id = ?", id).
			Update("department_id", *reassignToID)
		if result.Error != nil {
			return result.Error
		}

		// Delete the department
		result = s.Db.Delete(dept)
		if result.Error != nil {
			return result.Error
		}

		return nil

	default:
		return errors.New("invalid delete mode: must be 'cascade' or 'reassign'")
	}
}

// Helper functions

// departmentExists checks if a department exists
func (s *Storage) departmentExists(id int) (bool, error) {
	var count int64
	result := s.Db.Model(&models.Department{}).Where("id = ?", id).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

// checkCycle checks if moving a department would create a cycle
func (s *Storage) checkCycle(departmentID, newParentID int) (bool, error) {
	// Navigate up the chain from newParentID to see if we find departmentID
	currentID := newParentID

	for i := 0; i < 1000; i++ { // Safety limit to prevent infinite loops
		if currentID == departmentID {
			return true, nil
		}

		dept := &models.Department{}
		result := s.Db.Select("parent_id").First(dept, currentID)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return false, nil
			}
			return false, result.Error
		}

		if dept.ParentID == nil {
			return false, nil
		}

		currentID = *dept.ParentID
	}

	return false, errors.New("cycle detection timeout")
}

// Validation functions

func validateDepartmentName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("department name cannot be empty")
	}
	if len(name) > 200 {
		return errors.New("department name cannot exceed 200 characters")
	}
	return nil
}

func validateEmployeeName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("full name cannot be empty")
	}
	if len(name) > 200 {
		return errors.New("full name cannot exceed 200 characters")
	}
	return nil
}

func validatePosition(position string) error {
	position = strings.TrimSpace(position)
	if position == "" {
		return errors.New("position cannot be empty")
	}
	if len(position) > 200 {
		return errors.New("position cannot exceed 200 characters")
	}
	return nil
}

// GetAllDepartments retrieves all departments (for testing/debugging)
func (s *Storage) GetAllDepartments() ([]models.Department, error) {
	var departments []models.Department
	result := s.Db.Find(&departments)
	if result.Error != nil {
		return nil, result.Error
	}
	return departments, nil
}

// GetAllEmployees retrieves all employees (for testing/debugging)
func (s *Storage) GetAllEmployees() ([]models.Employee, error) {
	var employees []models.Employee
	result := s.Db.Find(&employees)
	if result.Error != nil {
		return nil, result.Error
	}
	return employees, nil
}
