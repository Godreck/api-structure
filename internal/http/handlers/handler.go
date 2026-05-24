// api-structure/internal/http/handlers
package handlers

import "gorm.io/gorm"

type DepartmentHandler struct {
	service *departmentService
}

func NewDepartmentHandler(db *gorm.DB) *DepartmentHandler {
	repo := newDepartmentRepository(db)
	svc := newDepartmentService(repo)
	return &DepartmentHandler{service: svc}
}
