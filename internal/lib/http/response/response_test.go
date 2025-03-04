package response

import (
	"errors"
	"github.com/go-playground/validator/v10"
	assertP "github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestForbidden(t *testing.T) {
	assert := assertP.New(t)

	t.Run("Default message", func(t *testing.T) {
		resp := Forbidden("")
		expected := ErrorResp{
			StatusCode: http.StatusForbidden,
			Message:    forbiddenMsg,
			Status:     StatusError,
		}

		assert.Equal(expected, resp, "Should return default forbidden response")
	})

	t.Run("Custom message", func(t *testing.T) {
		customMsg := "Access denied"
		resp := Forbidden(customMsg)

		assert.Equal(http.StatusForbidden, resp.StatusCode, "Status should be 403")
		assert.Equal(customMsg, resp.Message, "Should return the provided message")
		assert.Equal(StatusError, resp.Status, "Status should be error")
	})

	t.Run("Wrong expected status", func(t *testing.T) {
		resp := Forbidden("")
		assert.NotEqual(http.StatusBadRequest, resp.StatusCode, "Forbidden should not return 400")
	})
}

func TestUnauthorized(t *testing.T) {
	assert := assertP.New(t)

	t.Run("Default message", func(t *testing.T) {
		resp := Unauthorized("")
		expected := ErrorResp{
			StatusCode: http.StatusUnauthorized,
			Message:    unauthorizedMsg,
			Status:     StatusError,
		}

		assert.Equal(expected, resp, "Should return default unauthorized response")
	})

	t.Run("Custom message", func(t *testing.T) {
		customMsg := "you are unauthorized"
		resp := Unauthorized(customMsg)

		assert.Equal(http.StatusUnauthorized, resp.StatusCode, "Status should be 401")
		assert.Equal(customMsg, resp.Message, "Should return the provided message")
		assert.Equal(StatusError, resp.Status, "Status should be error")
	})

	t.Run("Wrong expected status", func(t *testing.T) {
		resp := Unauthorized("")
		assert.NotEqual(http.StatusForbidden, resp.StatusCode, "Unauthorized should not return 403")
	})
}

func TestOkWMsg(t *testing.T) {
	assert := assertP.New(t)

	t.Run("Returns OK response with message", func(t *testing.T) {
		msg := "Operation successful"
		resp := OkWMsg(msg)

		expected := OKResp{
			StatusCode: http.StatusOK,
			Message:    msg,
			Status:     StatusOK,
		}

		assert.Equal(expected, resp, "Should return OK response with message")
	})
}

func TestOkWDataAMsg(t *testing.T) {
	assert := assertP.New(t)

	t.Run("Returns OK response with data and message", func(t *testing.T) {
		data := map[string]string{"key": "value"}
		msg := "Operation successful"
		resp := OkWDataAMsg(data, msg)

		expected := OKResp{
			StatusCode: http.StatusOK,
			Message:    msg,
			Data:       data,
			Status:     StatusOK,
		}

		assert.Equal(expected, resp, "Should return OK response with data and message")
	})
}

func TestOkWData(t *testing.T) {
	assert := assertP.New(t)

	t.Run("Returns OK response with data only", func(t *testing.T) {
		data := []int{1, 2, 3}
		resp := OkWData(data)

		expected := OKResp{
			StatusCode: http.StatusOK,
			Data:       data,
			Status:     StatusOK,
		}

		assert.Equal(expected, resp, "Should return OK response with data only")
	})
}

func TestNotFound(t *testing.T) {
	assert := assertP.New(t)

	t.Run("Returns Not Found response with default message", func(t *testing.T) {
		resp := NotFound("")
		assert.Equal(http.StatusNotFound, resp.StatusCode)
		assert.Equal(notFoundMsg, resp.Message)
	})

	t.Run("Returns Not Found response with custom message", func(t *testing.T) {
		customMsg := "Post not found"
		resp := NotFound(customMsg)

		assert.Equal(http.StatusNotFound, resp.StatusCode)
		assert.Equal(customMsg, resp.Message)
	})
}

func TestBadRequest(t *testing.T) {
	assert := assertP.New(t)

	t.Run("Returns Bad Request response with default message", func(t *testing.T) {
		resp := BadRequest("")
		assert.Equal(http.StatusBadRequest, resp.StatusCode)
		assert.Equal(badRequestMsg, resp.Message)
	})

	t.Run("Returns Bad Request response with custom message", func(t *testing.T) {
		customMsg := "Invalid input data"
		resp := BadRequest(customMsg)

		assert.Equal(http.StatusBadRequest, resp.StatusCode)
		assert.Equal(customMsg, resp.Message)
	})
}

func TestInternalServerError(t *testing.T) {
	assert := assertP.New(t)

	t.Run("Returns Internal Server Error response with default message", func(t *testing.T) {
		resp := InternalServerError("")
		assert.Equal(http.StatusInternalServerError, resp.StatusCode)
		assert.Equal(serverMsg, resp.Message)
	})

	t.Run("Returns Internal Server Error response with custom message", func(t *testing.T) {
		customMsg := "Unexpected failure"
		resp := InternalServerError(customMsg)

		assert.Equal(http.StatusInternalServerError, resp.StatusCode)
		assert.Equal(customMsg, resp.Message)
	})
}

func TestInvalidInput(t *testing.T) {
	assert := assertP.New(t)

	t.Run("Returns Invalid Input response with validation errors", func(t *testing.T) {
		validate := validator.New()
		type TestStruct struct {
			Name string `validate:"required"`
		}

		testData := TestStruct{}
		err := validate.Struct(testData)
		var validationErrors validator.ValidationErrors
		errors.As(err, &validationErrors)

		resp := InvalidInput(validationErrors)

		assert.Equal(http.StatusBadRequest, resp.StatusCode)
		assert.Equal("There is some problem with the data you submitted.", resp.Message)
		assert.NotNil(resp.Details, "Details should not be nil")
		assert.Greater(len(resp.Details.([]invalidField)), 0, "There should be at least one validation error")
	})
}
