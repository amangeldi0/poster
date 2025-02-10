package authmiddleware

import (
	"context"
	"errors"
	"fmt"
	"money-manager/internal/auth"
	"money-manager/internal/lib/http/json"
	"money-manager/internal/lib/http/response"
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

		fmt.Println("Token:", tokenString)

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
