package auth

import (
	"database/sql"
	"github.com/google/uuid"
	"log/slog"
	"money-manager/internal/auth"
	"money-manager/internal/database"
	"money-manager/internal/lib/http/json"
	"money-manager/internal/lib/http/response"
	"money-manager/internal/lib/logger/sl"
	"money-manager/internal/lib/sql/sqlhelpers"
	"net/http"
	"strings"
	"time"
)

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token,omitempty"`
}

func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	const op = "auth.RefreshToken"

	var refreshToken string
	var req refreshTokenRequest

	cookie, err := r.Cookie("refresh_token")
	if err == nil {
		refreshToken = cookie.Value
	} else {
		if details, err := json.DecodeJSONBody(w, r, &req); err != nil {
			h.logger.Warn("Invalid JSON body", slog.String("op", op), sl.Err(err))
			json.WriteJSON(w, details.StatusCode, details)
			return
		}
		refreshToken = req.RefreshToken
	}

	if refreshToken == "" {
		h.logger.Warn("Missing refresh token", slog.String("op", op), sl.Err(err))
		json.WriteJSON(w, http.StatusUnauthorized, response.Unauthorized("Session expired, please login again"))
		return
	}

	claims, err := auth.VerifyToken(refreshToken)
	if err != nil {
		if strings.Contains(err.Error(), "token expired") {
			h.logger.Warn("Refresh token expired", slog.String("op", op))
			json.WriteJSON(w, http.StatusUnauthorized, response.Unauthorized("Refresh token expired, please login again"))
			return
		}

		h.logger.Warn("Invalid refresh token", slog.String("op", op), sl.Err(err))
		json.WriteJSON(w, http.StatusUnauthorized, response.Unauthorized("Invalid refresh token"))
		return
	}

	uId, err := uuid.Parse(claims.UserID)
	if err != nil {
		h.logger.Warn("Invalid user ID", slog.String("op", op), sl.Err(err))
		json.WriteJSON(w, http.StatusUnauthorized, response.Unauthorized("Invalid user ID"))
		return
	}

	u, err := h.query.GetUserByUUID(r.Context(), uId)
	if err != nil {
		h.logger.Warn("Invalid user query", slog.String("op", op), sl.Err(err))
		errD := sqlhelpers.GetDBError(err, label)
		json.WriteJSON(w, http.StatusUnauthorized, errD)
		return
	}

	accessToken, err := auth.GenerateAccessToken(u.ID.String())
	if err != nil {
		h.logger.Error("Failed to generate access token", slog.String("op", op), sl.Err(err))
		json.WriteJSON(w, http.StatusInternalServerError, response.InternalServerError("Failed to generate token"))
		return
	}

	newRefreshToken, err := auth.GenerateRefreshToken(u.ID.String())
	if err != nil {

		h.logger.Error("Failed to generate refresh token", slog.String("op", op), sl.Err(err))
		json.WriteJSON(w, http.StatusInternalServerError, response.InternalServerError("Failed to generate token"))
		return
	}

	err = h.query.UpdateUserRefreshToken(r.Context(), database.UpdateUserRefreshTokenParams{
		RefreshToken: sql.NullString{String: newRefreshToken, Valid: true},
		ID:           u.ID,
	})

	if err != nil {
		h.logger.Warn("Invalid update user refresh token", slog.String("op", op), sl.Err(err))
		errD := sqlhelpers.GetDBError(err, label)
		json.WriteJSON(w, http.StatusUnauthorized, errD)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(15 * time.Minute),
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})

	json.WriteJSON(w, http.StatusOK, map[string]string{
		"access_token":  accessToken,
		"refresh_token": newRefreshToken,
		"message":       "Tokens refreshed",
	})
}
