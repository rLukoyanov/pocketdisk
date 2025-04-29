package main

import (
	"embed"
	"net/http"
	"os"
	"os/signal"
	"pocketdisk/internal/config"
	"pocketdisk/internal/db"
	"pocketdisk/internal/handlers"
	"pocketdisk/internal/middleware"
	"syscall"

	"github.com/sirupsen/logrus"
)

//go:embed templates/*
var templateFS embed.FS

func main() {
	done := make(chan os.Signal, 1)

	cfg, err := config.NewConfig()
	if err != nil {
		logrus.Error(err)
		panic(err)
	}

	initFolders()

	sqlite, err := db.InitDB()
	if err != nil {
		logrus.Error(err)
		panic(err)
	}
	_ = sqlite

	go serve(cfg)
	logrus.Info("Сервер запущен")

	signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	<-done
}

func serve(cfg *config.Config) {
	renderHandlers := handlers.RenderHandlers{TemplateFS: templateFS, Cfg: cfg}
	apiHandlers := handlers.ApiHandlers{Cfg: cfg}

	http.HandleFunc("/api/login", apiHandlers.Login)

	http.HandleFunc("/", middleware.AuthMiddleware(renderHandlers.DashboardPage))
	http.HandleFunc("/login", renderHandlers.LoginPage)
	http.ListenAndServe(":8080", nil)
}

func initFolders() {
	if err := os.Mkdir("uploads", 0600); os.IsNotExist(err) {
		logrus.Error(err)
	}

	if err := os.Mkdir("templates", 0700); os.IsNotExist(err) {
		logrus.Error(err)
	}
}
