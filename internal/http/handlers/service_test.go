package handlers

import (
	"testing"
	"time"

	"api-structure/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateDepartmentName(t *testing.T) {
	t.Run("valid name", func(t *testing.T) {
		err := validateDepartmentName("IT")
		require.NoError(t, err)
	})

	t.Run("empty name", func(t *testing.T) {
		err := validateDepartmentName("")
		require.Error(t, err)
		assert.Equal(t, "department name must be 1..200 characters", err.Error())
	})

	t.Run("too long", func(t *testing.T) {
		err := validateDepartmentName(string(make([]byte, 201)))
		require.Error(t, err)
		assert.Equal(t, "department name must be 1..200 characters", err.Error())
	})
}

func TestValidateEmployeeFields(t *testing.T) {
	t.Run("valid fields", func(t *testing.T) {
		err := validateEmployeeFields("Ivan Petrov", "Backend Developer")
		require.NoError(t, err)
	})

	t.Run("empty full_name", func(t *testing.T) {
		err := validateEmployeeFields("", "Backend Developer")
		require.Error(t, err)
		assert.Equal(t, "full_name must be 1..200 characters", err.Error())
	})

	t.Run("empty position", func(t *testing.T) {
		err := validateEmployeeFields("Ivan Petrov", "")
		require.Error(t, err)
		assert.Equal(t, "position must be 1..200 characters", err.Error())
	})
}

func TestParseDatePtr(t *testing.T) {
	t.Run("nil value", func(t *testing.T) {
		got, err := parseDatePtr(nil)
		require.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("empty string", func(t *testing.T) {
		value := "   "
		got, err := parseDatePtr(&value)
		require.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("valid date", func(t *testing.T) {
		value := "2026-05-23"
		got, err := parseDatePtr(&value)
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, 2026, got.Year())
		assert.Equal(t, time.May, got.Month())
		assert.Equal(t, 23, got.Day())
	})

	t.Run("invalid date", func(t *testing.T) {
		value := "23-05-2026"
		got, err := parseDatePtr(&value)
		require.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestDepartmentResponseFromModel(t *testing.T) {
	parentID := 10
	createdAt := time.Date(2026, time.May, 23, 12, 0, 0, 0, time.UTC)
	dept := models.Department{
		ID:        1,
		Name:      "Backend",
		ParentID:  &parentID,
		CreatedAt: createdAt,
	}

	got := departmentResponseFromModel(dept)

	assert.Equal(t, 1, got.ID)
	assert.Equal(t, "Backend", got.Name)
	require.NotNil(t, got.ParentID)
	assert.Equal(t, 10, *got.ParentID)
	assert.Equal(t, createdAt.Format(time.RFC3339), got.CreatedAt)
}

func TestEmployeeResponseFromModel(t *testing.T) {
	hiredAtValue := time.Date(2026, time.May, 1, 0, 0, 0, 0, time.UTC)
	createdAt := time.Date(2026, time.May, 23, 12, 0, 0, 0, time.UTC)
	emp := models.Employee{
		ID:           7,
		DepartmentID: 3,
		FullName:     "Ivan Petrov",
		Position:     "Backend Developer",
		HiredAt:      &hiredAtValue,
		CreatedAt:    createdAt,
	}

	got := employeeResponseFromModel(emp)

	assert.Equal(t, 7, got.ID)
	assert.Equal(t, 3, got.DepartmentID)
	assert.Equal(t, "Ivan Petrov", got.FullName)
	assert.Equal(t, "Backend Developer", got.Position)
	require.NotNil(t, got.HiredAt)
	assert.Equal(t, "2026-05-01", *got.HiredAt)
	assert.Equal(t, createdAt.Format(time.RFC3339), got.CreatedAt)
}

func TestAppErrorError(t *testing.T) {
	t.Run("nil receiver", func(t *testing.T) {
		var err *appError
		assert.Equal(t, "", err.Error())
	})

	t.Run("non nil receiver", func(t *testing.T) {
		err := &appError{status: 400, message: "bad request"}
		assert.Equal(t, "bad request", err.Error())
	})
}
