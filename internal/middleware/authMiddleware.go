package middleware

import (
	"context"
	"net/http"
	"pocketdisk/internal/models"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusUnauthorized)
			return
		}

		// проверить токен получить пользака

		_ = token
		ctx := context.WithValue(r.Context(), "user",
			models.User{
				Name: "chebureki",
			})
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
