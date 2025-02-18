package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gopkg.in/gomail.v2"
	"log/slog"
	"net/http"
	"poster/internal/auth"
	"poster/internal/database"
	"poster/internal/lib/http/json"
	"poster/internal/lib/http/response"
	"poster/internal/lib/logger/sl"
	"poster/internal/lib/sql/sqlhelpers"
	"time"
)

type userRegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

func verificationCodeTemplate(verificationCode string, email string) *gomail.Message {
	m := gomail.NewMessage()
	m.SetHeader("To", email)
	m.SetHeader("Subject", "🔐 Подтверждение регистрации")

	htmlBody := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>Подтверждение регистрации</title>
			<style>
				body { font-family: Arial, sans-serif; background-color: #f4f4f4; padding: 20px; text-align: center; }
				.container { background: white; padding: 20px; border-radius: 8px; box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.1); display: inline-block; }
				h2 { color: #333; }
				p { font-size: 16px; color: #555; }
				.code { font-size: 24px; font-weight: bold; color: #007bff; background: #e7f3ff; padding: 10px 20px; border-radius: 5px; display: inline-block; }
			</style>
		</head>
		<body>
			<div class="container">
				<h2>🔐 Подтверждение регистрации</h2>
				<p>Спасибо за регистрацию! Ваш код подтверждения:</p>
				<p class="code">%s</p>
				<p>Введите этот код в приложении, чтобы завершить регистрацию.</p>
				<p>Если вы не регистрировались, просто проигнорируйте это письмо.</p>
				<p>С уважением,<br>Ваша команда</p>
			</div>
		</body>
		</html>
	`, verificationCode)

	m.SetBody("text/html", htmlBody)

	return m
}

func (h *Handler) isUserCanRegister(ctx context.Context, email string) (response.ErrorResp, error) {
	u, err := h.query.GetUserByEmail(ctx, email)

	if errors.Is(err, sql.ErrNoRows) {
		return response.ErrorResp{}, nil
	}
	if err != nil {
		return response.ErrorResp{}, err
	}

	if u.IsVerified.Bool {
		return response.ErrorResp{
			Status:     response.StatusError,
			StatusCode: http.StatusConflict,
			Message:    "user is already registered",
		}, errors.New("user already exists")
	}

	err = h.query.DeleteUserByEmail(ctx, u.Email)
	if err != nil {
		return sqlhelpers.GetDBError(err, "isUserCanRegister"), err
	}

	return response.ErrorResp{}, nil
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	const op = "user.Register"

	var req userRegisterRequest

	h.logger.Debug("Incoming registration request", slog.String("op", op))

	if details, err := json.DecodeJSONBody(w, r, &req); err != nil {
		h.logger.Warn("Invalid JSON body", slog.String("op", op), sl.Err(err))
		json.WriteJSON(w, details.StatusCode, details)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			h.logger.Warn("Validation failed", slog.String("op", op), sl.Err(err))
			json.WriteJSON(w, http.StatusBadRequest, response.InvalidInput(validationErrors))
			return
		}

		h.logger.Warn("Invalid input data", slog.String("op", op), sl.Err(err))
		response.BadRequest("invalid input data")
		return
	}

	if errD, err := h.isUserCanRegister(r.Context(), req.Email); err != nil {
		h.logger.Warn("User already exists", slog.String("op", op), sl.Err(err))
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	password, err := auth.HashPassword(req.Password)
	if err != nil {
		h.logger.Error("Failed to hash password", slog.String("op", op), sl.Err(err))
		json.WriteJSON(w, http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	code, err := auth.GenerateCode()
	if err != nil {
		h.logger.Error("Failed to generate verification code", slog.String("op", op), sl.Err(err))
		json.WriteJSON(w, http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	emailMessage := verificationCodeTemplate(code, req.Email)
	if err := h.mailer.Send(emailMessage); err != nil {
		h.logger.Warn("Failed to send verification email", slog.String("op", op), sl.Err(err))
		json.WriteJSON(w, http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	h.logger.Info("Verification email sent", slog.String("op", op), slog.String("email", req.Email))

	_, err = h.query.CreateUser(r.Context(), database.CreateUserParams{
		ID:           uuid.New(),
		Username:     req.Username,
		Email:        req.Email,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		PasswordHash: password,
		VerifyCode:   sql.NullString{String: code, Valid: true},
	})

	if err != nil {
		errD := sqlhelpers.GetDBError(err, label)
		h.logger.Error("Failed to create user", slog.String("op", op), sl.Err(err))
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	h.logger.Info("User registered successfully", slog.String("op", op), slog.String("email", req.Email))

	json.WriteJSON(w, http.StatusOK, response.OkWMsg("User is registered, please verify your email"))
}
