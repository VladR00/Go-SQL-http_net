package handlers

import (
	"Go-SQL-http_net/internal/models"
	"Go-SQL-http_net/internal/storage/postgre"
	"encoding/json"
	"fmt"
	"net/http"
)

type StorageHandler struct {
	Storage *postgre.Storage
}

func NewStorageHandler(storage *postgre.Storage) *StorageHandler {
	return &StorageHandler{Storage: storage}
}

func (s *StorageHandler) HandlerCreateDepartment(w http.ResponseWriter, r *http.Request) (string, error) {
	if r.Method != http.MethodPost {
		models.DefaultResponse{Type: "Error", Message: "Only POST method allowed"}.Response(w, http.StatusMethodNotAllowed)
		return "Method", fmt.Errorf("Attempting to use the %s method when %s is permitted", r.Method, http.MethodPost)
	}
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

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newDepartment)
	return fmt.Sprintf("Department created successfully. id:%d, name:%s, parrent_id:%d, created_at:%v", newDepartment.ID, newDepartment.Name, newDepartment.Parent_id, newDepartment.Created_at), nil
}
