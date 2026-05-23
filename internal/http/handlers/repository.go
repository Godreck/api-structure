package handlers

import (
	"errors"

	"api-structure/internal/models"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type departmentRepository struct {
	db *gorm.DB
}

func newDepartmentRepository(db *gorm.DB) *departmentRepository {
	return &departmentRepository{db: db}
}

func (r *departmentRepository) getDepartmentByID(id int) (models.Department, error) {
	var dept models.Department
	err := r.db.First(&dept, id).Error
	return dept, err
}

func (r *departmentRepository) getDepartmentParentRef(id int) (models.Department, error) {
	var dept models.Department
	err := r.db.Select("id", "parent_id").First(&dept, id).Error
	return dept, err
}

func (r *departmentRepository) createDepartment(dept *models.Department) error {
	return r.db.Create(dept).Error
}

func (r *departmentRepository) updateDepartment(dept *models.Department, updates map[string]any) error {
	return r.db.Model(dept).Updates(updates).Error
}

func (r *departmentRepository) createEmployee(emp *models.Employee) error {
	return r.db.Create(emp).Error
}

func (r *departmentRepository) listEmployeesByDepartmentID(departmentID int) ([]models.Employee, error) {
	var employees []models.Employee
	err := r.db.
		Where("department_id = ?", departmentID).
		Order("created_at asc, full_name asc").
		Find(&employees).Error
	return employees, err
}

func (r *departmentRepository) listChildrenByParentID(parentID int) ([]models.Department, error) {
	var children []models.Department
	err := r.db.
		Where("parent_id = ?", parentID).
		Order("created_at asc, name asc").
		Find(&children).Error
	return children, err
}

func (r *departmentRepository) deleteDepartmentCascade(dept *models.Department) error {
	return r.db.Select("Employees", "Children").Delete(dept).Error
}

func (r *departmentRepository) reassignChildDepartments(sourceID, targetID int) error {
	return r.db.Model(&models.Department{}).
		Where("parent_id = ?", sourceID).
		Update("parent_id", targetID).Error
}

func (r *departmentRepository) reassignEmployees(sourceID, targetID int) error {
	return r.db.Model(&models.Employee{}).
		Where("department_id = ?", sourceID).
		Update("department_id", targetID).Error
}

func (r *departmentRepository) deleteDepartmentByID(dept *models.Department) error {
	return r.db.Delete(dept).Error
}

func (r *departmentRepository) runInTx(fn func(txRepo *departmentRepository) error) (err error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
			panic(rec)
		}
	}()

	txRepo := &departmentRepository{db: tx}
	if err := fn(txRepo); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
