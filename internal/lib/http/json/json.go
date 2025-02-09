package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"money-manager/internal/lib/http/response"
	"net/http"
	"strings"
)

func WriteJSON(w http.ResponseWriter, statusCode int, resp interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(resp)
}

func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) (response.ErrorResp, error) {
	ct := r.Header.Get("Content-Type")
	if ct != "" {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
		if mediaType != "application/helpers" {

			return response.ErrorResp{
				Message:    "content-Type header is not application/helpers",
				StatusCode: http.StatusUnsupportedMediaType,
				Status:     response.StatusError,
			}, errors.New("content-Type header is not application/helpers")
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return response.BadRequest(msg), err

		case errors.Is(err, io.ErrUnexpectedEOF):
			return response.BadRequest("Request body contains badly-formed JSON"), err

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return response.BadRequest(msg), err

		case strings.HasPrefix(err.Error(), "helpers: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "helpers: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return response.BadRequest(msg), err

		case errors.Is(err, io.EOF):
			return response.BadRequest("Request body must not be empty"), err

		case errors.As(err, &maxBytesError):
			msg := fmt.Sprintf("Request body must not be larger than %d bytes", maxBytesError.Limit)
			return response.BadRequest(msg), err

		default:
			return response.InternalServerError(err.Error()), err
		}
	}

	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return response.BadRequest("Request body must only contain a single JSON object"), err

	}

	return response.ErrorResp{}, nil
}
