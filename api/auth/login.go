package auth

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

type userLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	const op = "auth.Login"
	var req userLoginRequest

	h.logger.Debug("Incoming login request", slog.String("op", op))

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

	u, err := h.query.GetUserByEmail(r.Context(), req.Email)

	if err != nil {
		errD := sqlhelpers.GetDBError(err, label)
		h.logger.Error("Failed to find user", slog.String("op", op), sl.Err(err))
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	if !u.IsVerified.Bool {
		h.logger.Warn("User is not verified", slog.String("op", op), slog.String("email", u.Email))
		json.WriteJSON(w, http.StatusForbidden, response.ErrorResp{
			Status:     "error",
			StatusCode: http.StatusForbidden,
			Message:    "User is not verified. Please check your email.",
		})
		return
	}

	if err = auth.CheckPasswordHash(req.Password, u.PasswordHash); err != nil {
		h.logger.Warn("Invalid password attempt", slog.String("op", op), slog.String("email", u.Email))
		json.WriteJSON(w, http.StatusUnauthorized, response.Unauthorized("Invalid email or password"))
		return
	}

	accessToken, err := auth.GenerateAccessToken(u.ID.String())
	if err != nil {
		h.logger.Error("Failed to generate access token", slog.String("op", op), sl.Err(err))
		json.WriteJSON(w, http.StatusInternalServerError, response.InternalServerError("Failed to generate token"))
		return
	}

	refreshToken, err := auth.GenerateRefreshToken(u.ID.String())
	if err != nil {
		h.logger.Error("Failed to generate refresh token", slog.String("op", op), sl.Err(err))
		json.WriteJSON(w, http.StatusInternalServerError, response.InternalServerError("Failed to generate token"))
		return
	}

	err = h.query.UpdateUserRefreshToken(r.Context(), database.UpdateUserRefreshTokenParams{
		ID:           u.ID,
		RefreshToken: sql.NullString{String: refreshToken, Valid: true},
	})

	if err != nil {
		h.logger.Error("Failed to update refresh token", slog.String("op", op), sl.Err(err))
		errD := sqlhelpers.GetDBError(err, label)
		json.WriteJSON(w, errD.StatusCode, errD)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "access_token",
		Value: accessToken,
		Path:  "/",
		// TODO Поменять потом на https
		HttpOnly: true,  // 🔐 JS не сможет прочитать
		Secure:   false, // 🔒 Только HTTPS
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(15 * time.Minute),
	})

	http.SetCookie(w, &http.Cookie{
		Name:  "refresh_token",
		Value: refreshToken,
		Path:  "/",
		// TODO Поменять потом на https
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})

	json.WriteJSON(w, http.StatusOK, map[string]any{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"message":       "Logged in successfully",
		"status":        response.StatusOK,
	})

}
