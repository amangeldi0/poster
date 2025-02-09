package sql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"money-manager/internal/lib/http/response"
	"net/http"
	"regexp"
)

func extractDuplicateField(detail string) string {
	re := regexp.MustCompile(`Key \((.*?)\)=`)
	matches := re.FindStringSubmatch(detail)

	if len(matches) > 1 {
		return matches[1]
	}

	return "unknown field"
}

func GetDBError(err error, label string) response.ErrorResp {

	if errors.Is(err, sql.ErrNoRows) {
		return response.NotFound(fmt.Sprintf("%s not found", label))
	}

	var pgErr *pq.Error

	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":

			duplicateField := extractDuplicateField(pgErr.Detail)

			return response.ErrorResp{
				Status:     response.StatusError,
				Message:    fmt.Sprintf("%s with this %s already exists", label, duplicateField),
				StatusCode: http.StatusConflict,
			}

		case "08006":
			return response.InternalServerError("database connection error")
		}
	}

	return response.InternalServerError(err.Error())
}
