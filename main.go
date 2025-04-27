package main

import (
	"embed"
	"log"
	"net/http"
	"os"
	"pocketdisk/internal/handlers"
)

//go:embed templates/*
var templateFS embed.FS

func main() {

	initFolders()
	serve()
}

func serve() {

	h := handlers.Handlers{TemplateFS: templateFS}

	http.HandleFunc("/", h.Render)
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
