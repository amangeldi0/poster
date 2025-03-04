package response

import (
	"github.com/go-playground/validator/v10"
	"net/http"
	"sort"
)

type StatusType string

const StatusOK = StatusType("ok")
const StatusError = StatusType("error")

const (
	unauthorizedMsg = "You are not authenticated to perform the requested action."
	forbiddenMsg    = "Your request is in a bad format."
	serverMsg       = "Server error. Please try again later."
	badRequestMsg   = "Your request is in a bad format."
	notFoundMsg     = "Not found"
)

type invalidField struct {
	Field string `helpers:"field"`
	Error string `helpers:"error"`
}

type OKResp struct {
	Status     StatusType  `json:"status"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	StatusCode int         `json:"status_code"`
}

type ErrorResp struct {
	Status     StatusType  `json:"status,default=error"`
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Details    interface{} `json:"details,omitempty"`
}

func OkWMsg(msg string) OKResp {
	return OKResp{
		StatusCode: http.StatusOK,
		Message:    msg,
		Status:     StatusOK,
	}
}

func OkWDataAMsg(data interface{}, msg string) OKResp {
	return OKResp{
		StatusCode: http.StatusOK,
		Message:    msg,
		Data:       data,
		Status:     StatusOK,
	}
}

func OkWData(data interface{}) OKResp {
	return OKResp{
		StatusCode: http.StatusOK,
		Data:       data,
		Status:     StatusOK,
	}
}

func NotFound(msg string) ErrorResp {
	if msg == "" {
		msg = notFoundMsg
	}

	return ErrorResp{
		StatusCode: http.StatusNotFound,
		Message:    msg,
		Status:     StatusError,
	}
}

func BadRequest(msg string) ErrorResp {
	if msg == "" {
		msg = badRequestMsg
	}

	return ErrorResp{
		StatusCode: http.StatusBadRequest,
		Message:    msg,
		Status:     StatusError,
	}
}

func InternalServerError(msg string) ErrorResp {
	if msg == "" {
		msg = serverMsg
	}

	return ErrorResp{
		StatusCode: http.StatusInternalServerError,
		Message:    msg,
		Status:     StatusError,
	}
}

func Unauthorized(msg string) ErrorResp {
	if msg == "" {
		msg = unauthorizedMsg
	}
	return ErrorResp{
		Status:     StatusError,
		StatusCode: http.StatusUnauthorized,
		Message:    msg,
	}
}

func Forbidden(msg string) ErrorResp {
	if msg == "" {
		msg = forbiddenMsg
	}

	return ErrorResp{
		StatusCode: http.StatusForbidden,
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
