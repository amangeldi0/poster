package sqlhelpers

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
	assert2 "github.com/stretchr/testify/assert"
	"net/http"
	"poster/internal/lib/http/response"
	"testing"
)

func TestGetDBError(t *testing.T) {
	assert := assert2.New(t)

	t.Run("not found error", func(t *testing.T) {
		resp := GetDBError(sql.ErrNoRows, "user_test")
		assert.Equal("user_test not found", resp.Message, "Should return correct not found message")
		assert.Equal(http.StatusNotFound, resp.StatusCode, "Should return 404")
		assert.Equal(response.StatusError, resp.Status, "Should return error status")
	})

	t.Run("Duplicate key error", func(t *testing.T) {
		err := &pq.Error{
			Code:   "23505",
			Detail: `Key (email)=("test@example.com") already exists.`,
		}

		resp := GetDBError(err, "user_test")

		assert.Equal("user_test with this email already exists", resp.Message, "Should return duplicate error message")
		assert.Equal(http.StatusConflict, resp.StatusCode, "Should return 409 Conflict")
		assert.Equal(response.StatusError, resp.Status, "Should return error status")
	})

	t.Run("Database connection error", func(t *testing.T) {
		err := &pq.Error{
			Code: "08006",
		}

		resp := GetDBError(err, "user_test")

		assert.Equal("database connection error", resp.Message, "Should return connection error message")
		assert.Equal(http.StatusInternalServerError, resp.StatusCode, "Should return 500 Internal Server Error")
		assert.Equal(response.StatusError, resp.Status, "Should return error status")
	})

	t.Run("Unknown error", func(t *testing.T) {
		err := errors.New("some unexpected error")

		resp := GetDBError(err, "user_test")

		assert.Equal("some unexpected error", resp.Message, "Should return original error message")
		assert.Equal(http.StatusInternalServerError, resp.StatusCode, "Should return 500 Internal Server Error")
		assert.Equal(response.StatusError, resp.Status, "Should return error status")
	})
}

func TestExtractDuplicateFields(t *testing.T) {
	assert := assert2.New(t)

	t.Run("success extract duplicate fields", func(t *testing.T) {
		field := "Key (email)=(\"test@example.com\") already exists."
		details := extractDuplicateField(field)
		assert.Equal("email", details)
	})

	t.Run("extract unknown fields", func(t *testing.T) {
		field := "unknown field"
		details := extractDuplicateField(field)
		assert.Equal("unknown field", details)
	})
}
