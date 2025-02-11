package authmiddleware

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"log/slog"
	"money-manager/internal/auth"
	"money-manager/internal/lib/http/json"
	"money-manager/internal/lib/http/response"
	"money-manager/internal/lib/logger/sl"
	"net/http"
	"strings"
)

type key string

const UserIDKey key = "user_id"

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenString string

		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			cookie, err := r.Cookie("access_token")
			if err != nil {
				json.WriteJSON(w, http.StatusUnauthorized, response.Unauthorized("Missing token"))
				return
			}
			tokenString = cookie.Value
		}

		claims, err := auth.VerifyToken(tokenString)
		if err != nil {
			if errors.Is(err, auth.ErrJwtExpired) {
				json.WriteJSON(w, http.StatusUnauthorized, response.Unauthorized("Token expired, please refresh"))
				return
			}
			json.WriteJSON(w, http.StatusUnauthorized, response.Unauthorized("Invalid token"))
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Identify(r context.Context, w http.ResponseWriter, h *slog.Logger, op string) (uuid.UUID, response.ErrorResp, error) {
	mId, ok := r.Value(UserIDKey).(string)
	if !ok || mId == "" {
		h.Warn(op, "UserID is missing or not a string")
		auth.DeleteCookie("access_token", w)
		auth.DeleteCookie("refresh_token", w)
		return uuid.Nil, response.Unauthorized("invalid account jwt"), errors.New("missing or invalid user ID")
	}

	userId, err := uuid.Parse(mId)
	if err != nil {
		h.Warn(op, "failed to parse id as a valid uuid from context", sl.Err(err))
		auth.DeleteCookie("access_token", w)
		auth.DeleteCookie("refresh_token", w)
		return uuid.Nil, response.Unauthorized("invalid account jwt"), err
	}

	return userId, response.ErrorResp{}, nil
}
