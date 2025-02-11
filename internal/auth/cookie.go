package auth

import "net/http"

func DeleteCookie(key string, w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     key,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})
}
