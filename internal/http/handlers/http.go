// api-structure/internal/http/handlers
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"api-structure/internal/dto"
)

// CreateDepartment godoc
// @Summary Create department
// @Tags departments
// @Accept json
// @Produce json
// @Param body body dto.DepartmentCreateRequest true "Department payload"
// @Success 201 {object} dto.DepartmentResponse
// @Failure 400 {string} string
// @Failure 409 {string} string
// @Router /departments/ [post]
func (h *DepartmentHandler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	var req dto.DepartmentCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	resp, appErr := h.service.createDepartment(req)
	if appErr != nil {
		writeAppError(w, appErr)
		return
	}
	writeJSON(w, http.StatusCreated, resp)
}

// CreateEmployeeInDepartment godoc
// @Summary Create employee in department
// @Tags departments
// @Accept json
// @Produce json
// @Param id path int true "Department ID"
// @Param body body dto.EmployeeCreateRequest true "Employee payload"
// @Success 201 {object} dto.EmployeeResponse
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Router /departments/{id}/employees/ [post]
func (h *DepartmentHandler) CreateEmployeeInDepartment(w http.ResponseWriter, r *http.Request) {
	deptID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid department id", http.StatusBadRequest)
		return
	}

	var req dto.EmployeeCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	resp, appErr := h.service.createEmployeeInDepartment(deptID, req)
	if appErr != nil {
		writeAppError(w, appErr)
		return
	}
	writeJSON(w, http.StatusCreated, resp)
}

// GetDepartment godoc
// @Summary Get department tree
// @Tags departments
// @Produce json
// @Param id path int true "Department ID"
// @Param depth query int false "Tree depth, max 5" default(1) minimum(1) maximum(5)
// @Param include_employees query bool false "Include employees" default(true)
// @Success 200 {object} dto.DepartmentResponse
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Router /departments/{id} [get]
func (h *DepartmentHandler) GetDepartment(w http.ResponseWriter, r *http.Request) {
	deptID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid department id", http.StatusBadRequest)
		return
	}

	depth := 1
	if v := r.URL.Query().Get("depth"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			depth = parsed
		}
	}

	includeEmployees := true
	if v := r.URL.Query().Get("include_employees"); v != "" {
		if parsed, err := strconv.ParseBool(v); err == nil {
			includeEmployees = parsed
		}
	}

	resp, appErr := h.service.getDepartment(deptID, depth, includeEmployees)
	if appErr != nil {
		writeAppError(w, appErr)
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// UpdateDepartment godoc
// @Summary Update department
// @Tags departments
// @Accept json
// @Produce json
// @Param id path int true "Department ID"
// @Param body body dto.DepartmentUpdateRequest true "Update payload"
// @Success 200 {object} dto.DepartmentResponse
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 409 {string} string
// @Router /departments/{id} [patch]
func (h *DepartmentHandler) UpdateDepartment(w http.ResponseWriter, r *http.Request) {
	deptID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid department id", http.StatusBadRequest)
		return
	}

	var req dto.DepartmentUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	resp, appErr := h.service.updateDepartment(deptID, req)
	if appErr != nil {
		writeAppError(w, appErr)
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// DeleteDepartment godoc
// @Summary Delete department
// @Tags departments
// @Produce json
// @Param id path int true "Department ID"
// @Param mode query string false "Delete mode: cascade or reassign" Enums(cascade,reassign) default(cascade)
// @Param reassign_to_department_id query int false "Target department for reassign mode"
// @Success 204
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Router /departments/{id} [delete]
func (h *DepartmentHandler) DeleteDepartment(w http.ResponseWriter, r *http.Request) {
	deptID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid department id", http.StatusBadRequest)
		return
	}

	mode := r.URL.Query().Get("mode")
	if mode == "" {
		mode = "cascade"
	}

	var reassignToDepartmentID *int
	if v := r.URL.Query().Get("reassign_to_department_id"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			http.Error(w, "reassign_to_department_id must be integer", http.StatusBadRequest)
			return
		}
		reassignToDepartmentID = &parsed
	}

	if appErr := h.service.deleteDepartment(deptID, mode, reassignToDepartmentID); appErr != nil {
		writeAppError(w, appErr)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeAppError(w http.ResponseWriter, appErr *appError) {
	if appErr == nil {
		return
	}
	http.Error(w, appErr.message, appErr.status)
}
