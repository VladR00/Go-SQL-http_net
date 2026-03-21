package handlers

import (
	"Go-SQL-http_net/internal/models"
	"Go-SQL-http_net/internal/storage/postgre"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

type StorageHandler struct {
	Storage *postgre.Storage
}

func NewStorageHandler(storage *postgre.Storage) *StorageHandler {
	return &StorageHandler{Storage: storage}
}

func (s *StorageHandler) Handler(w http.ResponseWriter, r *http.Request) (string, error) {
	var err error
	var response string
	path := strings.TrimPrefix(r.URL.Path, "/departments/")
	path = strings.TrimSuffix(path, "/")

	switch r.Method {
	case http.MethodPost:
		if path == "" {
			response, err = s.createDepartment(w, r)
			response = fmt.Sprintf("POST, %s", response)
		}
		sl := strings.Split(path, "/")
		slog.Debug("POST", "Split slice", sl)
		if len(sl) == 2 {
			id, err := strconv.Atoi(sl[0])
			if err != nil || id < 0 || sl[1] != "employees" {
				models.DefaultResponse{Type: "Error", Message: "POST. URL Parse. Want /departments/{id}/employees/"}.Response(w, http.StatusBadRequest)
				return "POST. URL Parse", fmt.Errorf("Want /departments/{id}/employees/")
			}
			slog.Debug("POST", "id", id)
			models.DefaultResponse{Type: "Message", Message: "POST. Create employee successfull."}.Response(w, http.StatusOK)
			return "POST. Create employee", nil
		}

	case http.MethodGet: // /departments/{id}
		id, err := strconv.Atoi(path)
		if err != nil || id < 0 {
			models.DefaultResponse{Type: "Error", Message: "GET. URL Parse. Want /departments/{id}"}.Response(w, http.StatusBadRequest)
			return "GET. URL Parse", fmt.Errorf("Want /departments/{id}")
		}
		slog.Debug("GET", "id", id)
		models.DefaultResponse{Type: "Message", Message: "GET. Get department successfull."}.Response(w, http.StatusOK)
		return "GET. Get department", nil
	case http.MethodPatch: // /departments/{id}
		id, err := strconv.Atoi(path)
		if err != nil || id < 0 {
			models.DefaultResponse{Type: "Error", Message: "PATCH. URL Parse. Want /departments/{id}"}.Response(w, http.StatusBadRequest)
			return "PATCH. URL Parse", fmt.Errorf("Want /departments/{id}")
		}
		slog.Debug("GET", "id", id)
		models.DefaultResponse{Type: "Message", Message: "PATCH. Relocate department successfull."}.Response(w, http.StatusOK)
		return "PATCH. Relocate department", nil
	case http.MethodDelete:
		id, err := strconv.Atoi(path)
		if err != nil || id < 0 {
			models.DefaultResponse{Type: "Error", Message: "DELETE. URL Parse. Want /departments/{id}"}.Response(w, http.StatusBadRequest)
			return "DELETE. URL Parse", fmt.Errorf("Want /departments/{id}")
		}
		slog.Debug("DELETE", "id", id)
		models.DefaultResponse{Type: "Message", Message: "DELETE. DELETE department successfull."}.Response(w, http.StatusOK)
		return "DELETE. DELETE department", nil
	default:
		models.DefaultResponse{Type: "Error", Message: "Method not allowed"}.Response(w, http.StatusMethodNotAllowed)
		return "Method", fmt.Errorf("Attempting to use the %s method", r.Method)
	}
	if err != nil {
		return response, err
	}
	return response, nil
}

func (s *StorageHandler) createDepartment(w http.ResponseWriter, r *http.Request) (string, error) {
	var request models.DepartmentRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		models.DefaultResponse{Type: "Error", Message: "Invalid JSON. Want 'name':'string', 'parent_id':'int | null'"}.Response(w, http.StatusBadRequest)
		return "Invalid JSON", fmt.Errorf("Decode error: %w", err)
	}

	if request.Name == "" {
		models.DefaultResponse{Type: "Error", Message: "Decode. Want 'name':'string', 'parent_id':'int | null'"}.Response(w, http.StatusBadRequest)
		return "Invalid JSON", fmt.Errorf("empty param")
	}

	newDepartment, err := s.Storage.CreateDepartment(request)
	if err != nil {
		models.DefaultResponse{Type: "Error", Message: fmt.Sprintf("Create department: %v", err)}.Response(w, http.StatusForbidden)
		return "Create Department", err
	}

	models.SetHeader(w, http.StatusOK)
	json.NewEncoder(w).Encode(newDepartment)
	return fmt.Sprintf("Department created successfully. id:%d, name:%s, parrent_id:%d, created_at:%v", newDepartment.ID, newDepartment.Name, newDepartment.Parent_id, newDepartment.Created_at), nil
}

func (s *StorageHandler) createEmployee(w http.ResponseWriter, r *http.Request) (string, error) {
	return "", nil
}

func (s *StorageHandler) getDepartment(w http.ResponseWriter, r *http.Request) (string, error) {
	return "", nil
}

func (s *StorageHandler) relocateDepartment(w http.ResponseWriter, r *http.Request) (string, error) {
	return "", nil
}

func (s *StorageHandler) deleteDepartment(w http.ResponseWriter, r *http.Request) (string, error) {
	return "", nil
}
