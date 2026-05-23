// api-structure/internal/dto
package dto

type DepartmentCreateRequest struct {
	Name     string `json:"name" example:"IT"`
	ParentID *int   `json:"parent_id,omitempty" example:"1"`
}

type DepartmentUpdateRequest struct {
	Name     *string `json:"name,omitempty" example:"Backend"`
	ParentID *int    `json:"parent_id,omitempty" example:"2"`
}

type EmployeeCreateRequest struct {
	FullName string  `json:"full_name" example:"Ivan Petrov"`
	Position string  `json:"position" example:"Backend Developer"`
	HiredAt  *string `json:"hired_at,omitempty" example:"2026-05-23"`
}

type EmployeeResponse struct {
	ID           int     `json:"id" example:"1"`
	DepartmentID int     `json:"department_id" example:"10"`
	FullName     string  `json:"full_name" example:"Ivan Petrov"`
	Position     string  `json:"position" example:"Backend Developer"`
	HiredAt      *string `json:"hired_at,omitempty" example:"2026-05-23"`
	CreatedAt    string  `json:"created_at" example:"2026-05-23T12:00:00Z"`
}

type DepartmentTreeNode struct {
	ID        int                  `json:"id" example:"1"`
	Name      string               `json:"name" example:"IT"`
	ParentID  *int                 `json:"parent_id,omitempty" example:"null"`
	Employees []EmployeeResponse   `json:"employees,omitempty"`
	Children  []DepartmentTreeNode `json:"children,omitempty"`
}

type DepartmentResponse struct {
	ID        int                  `json:"id" example:"1"`
	Name      string               `json:"name" example:"IT"`
	ParentID  *int                 `json:"parent_id,omitempty" example:"null"`
	CreatedAt string               `json:"created_at" example:"2026-05-23T12:00:00Z"`
	Employees []EmployeeResponse   `json:"employees,omitempty"`
	Children  []DepartmentTreeNode `json:"children,omitempty"`
}
