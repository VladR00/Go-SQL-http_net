package handlers

import (
	"Go-SQL-http_net/internal/models"
	"Go-SQL-http_net/internal/storage/postgre"
	"encoding/json"
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

// Handler routes requests to appropriate handlers
func (sh *StorageHandler) Handler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/departments")
	path = strings.TrimSuffix(path, "/")

	switch r.Method {
	case http.MethodPost:
		sh.handlePost(w, r, path)
	case http.MethodGet:
		sh.handleGet(w, r, path)
	case http.MethodPatch:
		sh.handlePatch(w, r, path)
	case http.MethodDelete:
		sh.handleDelete(w, r, path)
	default:
		models.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// handlePost routes POST requests
func (sh *StorageHandler) handlePost(w http.ResponseWriter, r *http.Request, path string) {
	if path == "" {
		// POST /departments/ - create department
		sh.createDepartment(w, r)
		return
	}

	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")

	if len(parts) == 2 && parts[1] == "employees" {
		// POST /departments/{id}/employees/ - create employee
		deptID, err := strconv.Atoi(parts[0])
		if err != nil || deptID < 1 {
			models.ErrorResponse(w, http.StatusBadRequest, "Invalid department ID")
			return
		}
		sh.createEmployee(w, r, deptID)
		return
	}

	models.ErrorResponse(w, http.StatusBadRequest, "Invalid path")
}

// handleGet routes GET requests
func (sh *StorageHandler) handleGet(w http.ResponseWriter, r *http.Request, path string) {
	if path == "" {
		models.ErrorResponse(w, http.StatusBadRequest, "Department ID required")
		return
	}

	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
	deptID, err := strconv.Atoi(parts[0])
	if err != nil || deptID < 1 {
		models.ErrorResponse(w, http.StatusBadRequest, "Invalid department ID")
		return
	}

	sh.getDepartment(w, r, deptID)
}

// handlePatch routes PATCH requests
func (sh *StorageHandler) handlePatch(w http.ResponseWriter, r *http.Request, path string) {
	if path == "" {
		models.ErrorResponse(w, http.StatusBadRequest, "Department ID required")
		return
	}

	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
	deptID, err := strconv.Atoi(parts[0])
	if err != nil || deptID < 1 {
		models.ErrorResponse(w, http.StatusBadRequest, "Invalid department ID")
		return
	}

	sh.updateDepartment(w, r, deptID)
}

// handleDelete routes DELETE requests
func (sh *StorageHandler) handleDelete(w http.ResponseWriter, r *http.Request, path string) {
	if path == "" {
		models.ErrorResponse(w, http.StatusBadRequest, "Department ID required")
		return
	}

	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
	deptID, err := strconv.Atoi(parts[0])
	if err != nil || deptID < 1 {
		models.ErrorResponse(w, http.StatusBadRequest, "Invalid department ID")
		return
	}

	sh.deleteDepartment(w, r, deptID)
}

// createDepartment handles POST /departments/
func (sh *StorageHandler) createDepartment(w http.ResponseWriter, r *http.Request) {
	var req models.DepartmentRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		slog.Error("createDepartment: decode error", "error", err)
		return
	}

	// Trim spaces
	req.Name = strings.TrimSpace(req.Name)

	dept, err := sh.Storage.CreateDepartment(req)
	if err != nil {
		slog.Error("createDepartment: storage error", "error", err)

		if strings.Contains(err.Error(), "parent department not found") {
			models.ErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		if strings.Contains(err.Error(), "already exists") {
			models.ErrorResponse(w, http.StatusConflict, err.Error())
			return
		}

		models.ErrorResponse(w, http.StatusInternalServerError, "Failed to create department")
		return
	}

	models.JSONResponse(w, http.StatusCreated, dept)
	slog.Info("createDepartment: success", "id", dept.ID, "name", dept.Name)
}

// createEmployee handles POST /departments/{id}/employees/
func (sh *StorageHandler) createEmployee(w http.ResponseWriter, r *http.Request, deptID int) {
	var req models.EmployeeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		slog.Error("createEmployee: decode error", "error", err)
		return
	}

	emp, err := sh.Storage.CreateEmployee(deptID, req)
	if err != nil {
		slog.Error("createEmployee: storage error", "error", err)

		if strings.Contains(err.Error(), "department not found") {
			models.ErrorResponse(w, http.StatusNotFound, "Department not found")
			return
		}

		models.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	models.JSONResponse(w, http.StatusCreated, emp)
	slog.Info("createEmployee: success", "id", emp.ID, "name", emp.FullName, "department", deptID)
}

// getDepartment handles GET /departments/{id}
func (sh *StorageHandler) getDepartment(w http.ResponseWriter, r *http.Request, deptID int) {
	// Parse query parameters
	query := r.URL.Query()

	depth := 1
	if d := query.Get("depth"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil {
			if parsed >= 0 && parsed <= 5 {
				depth = parsed
			}
		}
	}

	includeEmployees := true
	if ie := query.Get("include_employees"); ie != "" {
		includeEmployees = ie != "false"
	}

	dept, err := sh.Storage.GetDepartment(deptID, depth, includeEmployees)
	if err != nil {
		slog.Error("getDepartment: storage error", "error", err)

		if strings.Contains(err.Error(), "not found") {
			models.ErrorResponse(w, http.StatusNotFound, "Department not found")
			return
		}

		models.ErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve department")
		return
	}

	models.JSONResponse(w, http.StatusOK, dept)
	slog.Info("getDepartment: success", "id", deptID)
}

// updateDepartment handles PATCH /departments/{id}
func (sh *StorageHandler) updateDepartment(w http.ResponseWriter, r *http.Request, deptID int) {
	var req models.DepartmentUpdateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		slog.Error("updateDepartment: decode error", "error", err)
		return
	}

	// Trim spaces if name is provided
	if req.Name != nil {
		trimmed := strings.TrimSpace(*req.Name)
		req.Name = &trimmed
	}

	dept, err := sh.Storage.UpdateDepartment(deptID, req)
	if err != nil {
		slog.Error("updateDepartment: storage error", "error", err)

		if strings.Contains(err.Error(), "not found") {
			models.ErrorResponse(w, http.StatusNotFound, "Department not found")
			return
		}
		if strings.Contains(err.Error(), "cycle") {
			models.ErrorResponse(w, http.StatusConflict, err.Error())
			return
		}
		if strings.Contains(err.Error(), "cannot move") {
			models.ErrorResponse(w, http.StatusConflict, err.Error())
			return
		}
		if strings.Contains(err.Error(), "already exists") {
			models.ErrorResponse(w, http.StatusConflict, err.Error())
			return
		}

		models.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	models.JSONResponse(w, http.StatusOK, dept)
	slog.Info("updateDepartment: success", "id", deptID)
}

// deleteDepartment handles DELETE /departments/{id}
func (sh *StorageHandler) deleteDepartment(w http.ResponseWriter, r *http.Request, deptID int) {
	// Parse query parameters
	query := r.URL.Query()
	mode := query.Get("mode")

	if mode == "" {
		models.ErrorResponse(w, http.StatusBadRequest, "mode parameter is required (cascade or reassign)")
		return
	}

	var reassignToID *int
	if mode == "reassign" {
		reassignStr := query.Get("reassign_to_department_id")
		if reassignStr == "" {
			models.ErrorResponse(w, http.StatusBadRequest, "reassign_to_department_id is required when mode is reassign")
			return
		}

		reassignParsed, err := strconv.Atoi(reassignStr)
		if err != nil || reassignParsed < 1 {
			models.ErrorResponse(w, http.StatusBadRequest, "Invalid reassign_to_department_id")
			return
		}
		reassignToID = &reassignParsed
	}

	err := sh.Storage.DeleteDepartment(deptID, mode, reassignToID)
	if err != nil {
		slog.Error("deleteDepartment: storage error", "error", err)

		if strings.Contains(err.Error(), "not found") {
			models.ErrorResponse(w, http.StatusNotFound, "Department not found")
			return
		}
		if strings.Contains(err.Error(), "invalid delete mode") {
			models.ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		models.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
	slog.Info("deleteDepartment: success", "id", deptID, "mode", mode)
}
