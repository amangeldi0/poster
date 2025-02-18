package auth

import (
	"database/sql"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"poster/internal/database"
	"poster/internal/lib/http/json"
	"poster/internal/lib/http/response"
	"poster/internal/lib/logger/sl"
)

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	const op = "auth.Logout"

	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		h.logger.Warn("Unauthorized logout attempt", slog.String("op", op))
		json.WriteJSON(w, http.StatusUnauthorized, response.Unauthorized("Invalid token"))
		return
	}

	uID, err := uuid.Parse(userID)
	if err != nil {
		h.logger.Warn("Invalid user ID format", slog.String("op", op), sl.Err(err))
		json.WriteJSON(w, http.StatusInternalServerError, response.InternalServerError("Failed to logout"))
		return
	}

	err = h.query.UpdateUserRefreshToken(r.Context(), database.UpdateUserRefreshTokenParams{
		RefreshToken: sql.NullString{Valid: false},
		ID:           uID,
	})

	if err != nil {
		h.logger.Error("Failed to logout user", slog.String("op", op), sl.Err(err))
		json.WriteJSON(w, http.StatusInternalServerError, response.InternalServerError("Failed to logout"))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})

	json.WriteJSON(w, http.StatusOK, response.OkWMsg("Successfully logged out"))
}
