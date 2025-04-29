package handlers

import (
	"database/sql"
	"net/http"
	"pocketdisk/internal/config"

	"golang.org/x/crypto/bcrypt"
)

type ApiHandlers struct {
	Cfg *config.Config
	db  *sql.DB
}

func (h *ApiHandlers) Login(w http.ResponseWriter, r *http.Request) {

	var dbPassword string
	var userID int

	username := r.FormValue("username") // заменить на данные с post
	password := r.FormValue("password")

	err := h.db.QueryRow("SELECT id, password FROM users WHERE username = ?", username).Scan(&userID, &dbPassword)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "session",
		Value: username,
		Path:  "/",
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
