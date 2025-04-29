package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./cloud.db")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT,
		password TEXT,
		is_admin BOOLEAN DEFAULT FALSE,
		CHECK (email LIKE '%_@_%._%')
	);
	
	CREATE TABLE IF NOT EXISTS files (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		name TEXT,
		path TEXT,
		size INTEGER,
		uploaded TEXT,
		downloads INTEGER DEFAULT 0,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);
	`
	_, err = db.Exec(query)
	if err != nil {
		log.Printf("%q: %s\n", err, query)
		return nil, err
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123"), bcrypt.DefaultCost)
	db.Exec("INSERT OR IGNORE INTO users (email, password, is_admin) VALUES (?, ?, ?)",
		"admin@test.com", hashedPassword, true)
	return db, nil
}
