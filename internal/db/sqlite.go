package db

import (
	"database/sql"
	"log"
	"pocketdisk/internal/config"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func InitDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./cloud.db")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT UNIQUE,
		password TEXT,
		is_admin BOOLEAN DEFAULT FALSE,
		storage_limit BIGINT DEFAULT 1073741824,
		storage_used BIGINT DEFAULT 0,
		CHECK (email LIKE '%_@_%._%')
	);
	
	CREATE TABLE IF NOT EXISTS files (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		name TEXT,
		path TEXT,
		size INTEGER,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);
	`
	_, err = db.Exec(query)
	if err != nil {
		log.Printf("%q: %s\n", err, query)
		return nil, err
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(cfg.SECRET), bcrypt.DefaultCost)
	db.Exec("INSERT OR IGNORE INTO users (email, password, is_admin) VALUES (?, ?, ?)",
		"admin@test.com", hashedPassword, true)
	return db, nil
}
