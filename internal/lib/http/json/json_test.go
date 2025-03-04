package json

import (
	"encoding/json"
	assert2 "github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"poster/internal/lib/http/response"
	"strings"
	"testing"
)

func TestWriteJSON(t *testing.T) {

	t.Run("Ok response", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		assert := assert2.New(t)
		expected := response.OkWMsg("Hello world!")
		WriteJSON(recorder, expected.StatusCode, expected)

		t.Run("Check status code", func(t *testing.T) {
			assert.Equal(expected.StatusCode, recorder.Result().StatusCode)
		})

		t.Run("Check content type", func(t *testing.T) {
			expectedContentType := "application/json"
			contentType := recorder.Result().Header.Get("Content-Type")
			assert.Equal(expectedContentType, contentType)
		})

		t.Run("Check response body", func(t *testing.T) {
			var actualData response.OKResp
			err := json.NewDecoder(recorder.Body).Decode(&actualData)
			assert.NoError(err, "Response should be valid JSON")
			assert.Equal(expected, actualData, "Response body should match expected JSON")
		})

	})

	t.Run("Error response", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		assert := assert2.New(t)
		expected := response.Forbidden("error")
		WriteJSON(recorder, expected.StatusCode, expected)

		t.Run("Check status code", func(t *testing.T) {
			assert.Equal(expected.StatusCode, recorder.Result().StatusCode)
		})

		t.Run("Check content type", func(t *testing.T) {
			expectedContentType := "application/json"
			contentType := recorder.Result().Header.Get("Content-Type")
			assert.Equal(expectedContentType, contentType)
		})

		t.Run("Check response body", func(t *testing.T) {
			var actualData response.ErrorResp
			err := json.NewDecoder(recorder.Body).Decode(&actualData)
			assert.NoError(err, "Response should be valid JSON")
			assert.Equal(expected, actualData, "Response body should match expected JSON")
		})

	})

}

func TestDecodeJSONBody(t *testing.T) {
	assert := assert2.New(t)
	exampleData := `{"email": "email@example.com", "password": "1234"}`

	t.Run("Success decode", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodPost, "/api", strings.NewReader(exampleData))
		w := httptest.NewRecorder()

		var data map[string]string

		_, err := DecodeJSONBody(w, req, &data)

		assert.Equal("email@example.com", data["email"])
		assert.Equal("1234", data["password"])

		assert.NoError(err)
	})

	t.Run("Error decode wrong media", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api", strings.NewReader(exampleData))
		w := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/error")

		var data map[string]string
		_, err := DecodeJSONBody(w, req, &data)

		assert.Error(err)
	})
}
