package main

import (
	"embed"
	"log"
	"net/http"
	"os"
	"pocketdisk/internal/db"
	"pocketdisk/internal/handlers"
	"pocketdisk/internal/middleware"
)

//go:embed templates/*
var templateFS embed.FS

func main() {
	done := make(chan os.Signal, 1)

	initFolders()
	sqlite, err := db.InitDB()
	if err != nil {
		panic(err)
	}
	_ = sqlite
	go serve()
	log.Println("Сервер запущен")
	<-done
}

func serve() {

	h := handlers.Handlers{TemplateFS: templateFS}

	http.HandleFunc("/", middleware.AuthMiddleware(h.Dashboard))
	http.HandleFunc("/login", h.Login)
	http.ListenAndServe(":8080", nil)
}

func initFolders() {
	if err := os.Mkdir("uploads", 0600); os.IsNotExist(err) {
		log.Println(err)
	}

	if err := os.Mkdir("templates", 0700); os.IsNotExist(err) {
		log.Println(err)
	}
}
