package response

import (
	"github.com/go-playground/validator/v10"
	"net/http"
	"sort"
)

type StatusType string

const StatusOK = StatusType("OK")
const StatusError = StatusType("Error")

type invalidField struct {
	Field string `helpers:"field"`
	Error string `helpers:"error"`
}

type OKResp struct {
	Status     StatusType  `helpers:"status"`
	Message    string      `helpers:"message,omitempty"`
	Data       interface{} `helpers:"data,omitempty"`
	StatusCode int         `helpers:"status_code"`
}

type ErrorResp struct {
	Status     StatusType  `helpers:"status,default=error"`
	StatusCode int         `helpers:"status_code"`
	Message    string      `helpers:"message"`
	Details    interface{} `helpers:"details,omitempty"`
}

func Ok(data interface{}, msg string) OKResp {
	return OKResp{
		StatusCode: http.StatusOK,
		Message:    msg,
		Data:       data,
		Status:     StatusOK,
	}
}

func NotFound(msg string) ErrorResp {
	if msg == "" {
		msg = "Not found"
	}

	return ErrorResp{
		StatusCode: http.StatusNotFound,
		Message:    msg,
		Status:     StatusOK,
	}
}

func BadRequest(msg string) ErrorResp {
	if msg == "" {
		msg = "Your request is in a bad format."
	}

	return ErrorResp{
		StatusCode: http.StatusBadRequest,
		Message:    msg,
		Status:     StatusError,
	}
}

func InternalServerError(msg string) ErrorResp {
	if msg == "" {
		msg = "Server error. Please try again later."
	}

	return ErrorResp{
		StatusCode: http.StatusInternalServerError,
		Message:    msg,
		Status:     StatusError,
	}
}

func Unauthorized(msg string) ErrorResp {
	if msg == "" {
		msg = "You are not authenticated to perform the requested action."
	}
	return ErrorResp{
		Status:     StatusError,
		StatusCode: http.StatusUnauthorized,
		Message:    msg,
	}
}

func Forbidden(msg string) ErrorResp {
	if msg == "" {
		msg = "Your request is in a bad format."
	}

	return ErrorResp{
		StatusCode: http.StatusBadRequest,
		Message:    msg,
		Status:     StatusError,
	}
}

func InvalidInput(errs validator.ValidationErrors) ErrorResp {
	var details []invalidField

	fields := make([]string, len(errs))
	for i, e := range errs {
		fields[i] = e.Field()
	}
	sort.Strings(fields)

	for _, e := range errs {
		details = append(details, invalidField{
			Field: e.Field(),
			Error: e.Error(),
		})
	}

	return ErrorResp{
		Status:     StatusError,
		StatusCode: http.StatusBadRequest,
		Message:    "There is some problem with the data you submitted.",
		Details:    details,
	}
}
