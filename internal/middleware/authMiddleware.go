package middleware

import (
	"context"
	"net/http"
	"pocketdisk/internal/models"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusUnauthorized)
			return
		}

		_ = cookie
		ctx := context.WithValue(r.Context(), "user",
			models.User{
				Name: "chebureki",
			})
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
