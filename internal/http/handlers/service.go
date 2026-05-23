package handlers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"api-structure/internal/dto"
	"api-structure/internal/models"

	"gorm.io/gorm"
)

type appError struct {
	status  int
	message string
	cause   error
}

func (e *appError) Error() string {
	if e == nil {
		return ""
	}
	return e.message
}

type departmentService struct {
	repo *departmentRepository
}

func newDepartmentService(repo *departmentRepository) *departmentService {
	return &departmentService{repo: repo}
}

func (s *departmentService) createDepartment(req dto.DepartmentCreateRequest) (dto.DepartmentResponse, *appError) {
	name := strings.TrimSpace(req.Name)
	if err := validateDepartmentName(name); err != nil {
		return dto.DepartmentResponse{}, &appError{status: http.StatusBadRequest, message: err.Error(), cause: err}
	}

	if req.ParentID != nil {
		if _, err := s.repo.getDepartmentByID(*req.ParentID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return dto.DepartmentResponse{}, &appError{status: http.StatusNotFound, message: "parent department not found", cause: err}
			}
			return dto.DepartmentResponse{}, &appError{status: http.StatusInternalServerError, message: "database error", cause: err}
		}
	}

	dept := models.Department{Name: name, ParentID: req.ParentID}
	if err := s.repo.createDepartment(&dept); err != nil {
		if isUniqueViolation(err) {
			return dto.DepartmentResponse{}, &appError{status: http.StatusConflict, message: "department name must be unique within the same parent", cause: err}
		}
		return dto.DepartmentResponse{}, &appError{status: http.StatusInternalServerError, message: "failed to create department", cause: err}
	}

	return departmentResponseFromModel(dept), nil
}

func (s *departmentService) createEmployeeInDepartment(deptID int, req dto.EmployeeCreateRequest) (dto.EmployeeResponse, *appError) {
	if _, err := s.repo.getDepartmentByID(deptID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.EmployeeResponse{}, &appError{status: http.StatusNotFound, message: "department not found", cause: err}
		}
		return dto.EmployeeResponse{}, &appError{status: http.StatusInternalServerError, message: "database error", cause: err}
	}

	fullName := strings.TrimSpace(req.FullName)
	position := strings.TrimSpace(req.Position)
	if err := validateEmployeeFields(fullName, position); err != nil {
		return dto.EmployeeResponse{}, &appError{status: http.StatusBadRequest, message: err.Error(), cause: err}
	}

	hiredAt, err := parseDatePtr(req.HiredAt)
	if err != nil {
		return dto.EmployeeResponse{}, &appError{status: http.StatusBadRequest, message: "hired_at must be YYYY-MM-DD", cause: err}
	}

	emp := models.Employee{
		DepartmentID: deptID,
		FullName:     fullName,
		Position:     position,
		HiredAt:      hiredAt,
	}
	if err := s.repo.createEmployee(&emp); err != nil {
		return dto.EmployeeResponse{}, &appError{status: http.StatusInternalServerError, message: "failed to create employee", cause: err}
	}

	return employeeResponseFromModel(emp), nil
}

func (s *departmentService) getDepartment(deptID, depth int, includeEmployees bool) (dto.DepartmentResponse, *appError) {
	if depth < 1 {
		depth = 1
	}
	if depth > 5 {
		depth = 5
	}

	dept, err := s.repo.getDepartmentByID(deptID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.DepartmentResponse{}, &appError{status: http.StatusNotFound, message: "department not found", cause: err}
		}
		return dto.DepartmentResponse{}, &appError{status: http.StatusInternalServerError, message: "database error", cause: err}
	}

	resp, treeErr := s.buildDepartmentTree(dept, depth, includeEmployees)
	if treeErr != nil {
		return dto.DepartmentResponse{}, &appError{status: http.StatusInternalServerError, message: "failed to build tree", cause: treeErr}
	}

	return resp, nil
}

func (s *departmentService) updateDepartment(deptID int, req dto.DepartmentUpdateRequest) (dto.DepartmentResponse, *appError) {
	dept, err := s.repo.getDepartmentByID(deptID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.DepartmentResponse{}, &appError{status: http.StatusNotFound, message: "department not found", cause: err}
		}
		return dto.DepartmentResponse{}, &appError{status: http.StatusInternalServerError, message: "database error", cause: err}
	}

	updates := map[string]any{}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if err := validateDepartmentName(name); err != nil {
			return dto.DepartmentResponse{}, &appError{status: http.StatusBadRequest, message: err.Error(), cause: err}
		}
		updates["name"] = name
	}

	if req.ParentID != nil {
		if *req.ParentID == deptID {
			return dto.DepartmentResponse{}, &appError{status: http.StatusBadRequest, message: "department cannot be its own parent"}
		}

		if err := s.ensureNoCycle(deptID, *req.ParentID); err != nil {
			return dto.DepartmentResponse{}, &appError{status: http.StatusConflict, message: err.Error(), cause: err}
		}

		if *req.ParentID == 0 {
			updates["parent_id"] = nil
		} else {
			if _, err := s.repo.getDepartmentByID(*req.ParentID); err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return dto.DepartmentResponse{}, &appError{status: http.StatusNotFound, message: "parent department not found", cause: err}
				}
				return dto.DepartmentResponse{}, &appError{status: http.StatusInternalServerError, message: "database error", cause: err}
			}
			updates["parent_id"] = req.ParentID
		}
	}

	if len(updates) > 0 {
		if err := s.repo.updateDepartment(&dept, updates); err != nil {
			if isUniqueViolation(err) {
				return dto.DepartmentResponse{}, &appError{status: http.StatusConflict, message: "department name must be unique within the same parent", cause: err}
			}
			return dto.DepartmentResponse{}, &appError{status: http.StatusInternalServerError, message: "failed to update department", cause: err}
		}
	}

	dept, err = s.repo.getDepartmentByID(deptID)
	if err != nil {
		return dto.DepartmentResponse{}, &appError{status: http.StatusInternalServerError, message: "database error", cause: err}
	}

	return departmentResponseFromModel(dept), nil
}

func (s *departmentService) deleteDepartment(deptID int, mode string, reassignToDepartmentID *int) *appError {
	err := s.repo.runInTx(func(txRepo *departmentRepository) error {
		dept, err := txRepo.getDepartmentByID(deptID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &appError{status: http.StatusNotFound, message: "department not found", cause: err}
			}
			return &appError{status: http.StatusInternalServerError, message: "database error", cause: err}
		}

		switch mode {
		case "cascade":
			if err := txRepo.deleteDepartmentCascade(&dept); err != nil {
				return &appError{status: http.StatusInternalServerError, message: "failed to delete department", cause: err}
			}
			return nil
		case "reassign":
			if reassignToDepartmentID == nil || *reassignToDepartmentID <= 0 {
				return &appError{status: http.StatusBadRequest, message: "reassign_to_department_id is required for reassign mode"}
			}
			targetID := *reassignToDepartmentID
			if targetID == deptID {
				return &appError{status: http.StatusBadRequest, message: "cannot reassign to same department"}
			}

			if err := s.ensureNoCycleWithRepo(txRepo, deptID, targetID); err != nil {
				return &appError{status: http.StatusConflict, message: err.Error(), cause: err}
			}

			if _, err := txRepo.getDepartmentByID(targetID); err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return &appError{status: http.StatusNotFound, message: "reassign target not found", cause: err}
				}
				return &appError{status: http.StatusInternalServerError, message: "database error", cause: err}
			}

			if err := txRepo.reassignChildDepartments(deptID, targetID); err != nil {
				return &appError{status: http.StatusInternalServerError, message: "failed to reassign child departments", cause: err}
			}
			if err := txRepo.reassignEmployees(deptID, targetID); err != nil {
				return &appError{status: http.StatusInternalServerError, message: "failed to reassign employees", cause: err}
			}
			if err := txRepo.deleteDepartmentByID(&dept); err != nil {
				return &appError{status: http.StatusInternalServerError, message: "failed to delete department", cause: err}
			}
			return nil
		default:
			return &appError{status: http.StatusBadRequest, message: "mode must be cascade or reassign"}
		}
	})

	if err == nil {
		return nil
	}

	if domainErr, ok := err.(*appError); ok {
		return domainErr
	}
	return &appError{status: http.StatusInternalServerError, message: "database error", cause: err}
}

func (s *departmentService) buildDepartmentTree(dept models.Department, depth int, includeEmployees bool) (dto.DepartmentResponse, error) {
	resp := dto.DepartmentResponse{
		ID:        dept.ID,
		Name:      dept.Name,
		ParentID:  dept.ParentID,
		CreatedAt: dept.CreatedAt.Format(time.RFC3339),
	}

	if includeEmployees {
		employees, err := s.repo.listEmployeesByDepartmentID(dept.ID)
		if err != nil {
			return dto.DepartmentResponse{}, err
		}
		resp.Employees = make([]dto.EmployeeResponse, 0, len(employees))
		for _, e := range employees {
			resp.Employees = append(resp.Employees, employeeResponseFromModel(e))
		}
	}

	if depth <= 1 {
		return resp, nil
	}

	children, err := s.repo.listChildrenByParentID(dept.ID)
	if err != nil {
		return dto.DepartmentResponse{}, err
	}

	resp.Children = make([]dto.DepartmentTreeNode, 0, len(children))
	for _, child := range children {
		childResp, err := s.buildDepartmentTree(child, depth-1, includeEmployees)
		if err != nil {
			return dto.DepartmentResponse{}, err
		}
		resp.Children = append(resp.Children, dto.DepartmentTreeNode{
			ID:        childResp.ID,
			Name:      childResp.Name,
			ParentID:  childResp.ParentID,
			Employees: childResp.Employees,
			Children:  childResp.Children,
		})
	}

	return resp, nil
}

func (s *departmentService) ensureNoCycle(deptID, newParentID int) error {
	return s.ensureNoCycleWithRepo(s.repo, deptID, newParentID)
}

func (s *departmentService) ensureNoCycleWithRepo(repo *departmentRepository, deptID, newParentID int) error {
	if newParentID == 0 {
		return nil
	}

	current := newParentID
	for current != 0 {
		if current == deptID {
			return errors.New("cycle in department tree is not allowed")
		}

		parent, err := repo.getDepartmentParentRef(current)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("parent department not found")
			}
			return err
		}

		if parent.ParentID == nil {
			break
		}
		current = *parent.ParentID
	}

	return nil
}

func validateDepartmentName(name string) error {
	if len(name) < 1 || len(name) > 200 {
		return errors.New("department name must be 1..200 characters")
	}
	return nil
}

func validateEmployeeFields(fullName, position string) error {
	if len(fullName) < 1 || len(fullName) > 200 {
		return errors.New("full_name must be 1..200 characters")
	}
	if len(position) < 1 || len(position) > 200 {
		return errors.New("position must be 1..200 characters")
	}
	return nil
}

func parseDatePtr(v *string) (*time.Time, error) {
	if v == nil || strings.TrimSpace(*v) == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", strings.TrimSpace(*v))
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func departmentResponseFromModel(dept models.Department) dto.DepartmentResponse {
	return dto.DepartmentResponse{
		ID:        dept.ID,
		Name:      dept.Name,
		ParentID:  dept.ParentID,
		CreatedAt: dept.CreatedAt.Format(time.RFC3339),
	}
}

func employeeResponseFromModel(emp models.Employee) dto.EmployeeResponse {
	var hiredAt *string
	if emp.HiredAt != nil {
		s := emp.HiredAt.Format("2006-01-02")
		hiredAt = &s
	}
	return dto.EmployeeResponse{
		ID:           emp.ID,
		DepartmentID: emp.DepartmentID,
		FullName:     emp.FullName,
		Position:     emp.Position,
		HiredAt:      hiredAt,
		CreatedAt:    emp.CreatedAt.Format(time.RFC3339),
	}
}
