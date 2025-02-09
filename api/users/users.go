package users

import (
	"database/sql"
	"errors"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"money-manager/internal/auth"
	"money-manager/internal/database"
	"money-manager/internal/lib/http/json"
	"money-manager/internal/lib/http/response"
	"money-manager/internal/lib/logger/sl"
	"money-manager/internal/lib/sql/sqlhelpers"
	"net/http"
	"time"
)

const label = "user"

type Handler struct {
	log      *slog.Logger
	query    *database.Queries
	validate *validator.Validate
}

type userRegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

func NewUserHandler(log *slog.Logger, db *database.Queries) *Handler {
	return &Handler{
		log:      log,
		query:    db,
		validate: validator.New(),
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	const op = "user.Register"

	var req userRegisterRequest

	if details, err := json.DecodeJSONBody(w, r, &req); err != nil {
		h.log.Warn("invalid JSON body", slog.String("op", op), sl.Err(err))
		json.WriteJSON(w, details.StatusCode, details)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			h.log.Warn("validation failed", slog.String("op", op), sl.Err(err))
			json.WriteJSON(w, http.StatusBadRequest, response.InvalidInput(validationErrors))
			return
		}

		h.log.Warn("invalid input data", slog.String("op", op), sl.Err(err))
		response.BadRequest("invalid input data")
		return
	}

	password, err := auth.HashPassword(req.Password)

	if err != nil {
		h.log.Warn("invalid input data", slog.String("op", op), sl.Err(err))
		json.WriteJSON(w, http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	_, err = h.query.CreateUser(r.Context(), database.CreateUserParams{
		Username:     req.Username,
		Email:        req.Email,
		CreatedAt:    sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt:    sql.NullTime{Time: time.Now(), Valid: true},
		PasswordHash: password,
	})

	if err != nil {
		errD := sqlhelpers.GetDBError(err, label)
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

}
