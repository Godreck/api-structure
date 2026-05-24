// api-structure/internal/http/handlers
package handlers

import (
	"errors"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

func TestIsUniqueViolation(t *testing.T) {
	t.Run("postgres unique violation", func(t *testing.T) {
		err := &pgconn.PgError{Code: "23505"}
		assert.True(t, isUniqueViolation(err))
	})

	t.Run("wrapped postgres unique violation", func(t *testing.T) {
		err := fmt.Errorf("wrapped: %w", &pgconn.PgError{Code: "23505"})
		assert.True(t, isUniqueViolation(err))
	})

	t.Run("other postgres code", func(t *testing.T) {
		err := &pgconn.PgError{Code: "23503"}
		assert.False(t, isUniqueViolation(err))
	})

	t.Run("generic error", func(t *testing.T) {
		err := errors.New("some error")
		assert.False(t, isUniqueViolation(err))
	})

	t.Run("nil", func(t *testing.T) {
		assert.False(t, isUniqueViolation(nil))
	})
}
