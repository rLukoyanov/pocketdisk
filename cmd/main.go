package main

import (
	"os"
	"os/signal"
	"pocketdisk/internal/config"
	"pocketdisk/internal/db"
	"pocketdisk/internal/handlers"
	"pocketdisk/internal/pkg"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

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

	e := echo.New()

	pkg.AddNewRender(e)

	renderHandlers := handlers.RenderHandlers{Cfg: cfg}
	// apiHandlers := handlers.ApiHandlers{Cfg: cfg}
	e.Use(middleware.Logger())

	e.GET("/", renderHandlers.DashboardPage)
	e.GET("/login", renderHandlers.LoginPage)
	e.Start(":8080")

	logrus.Info("Сервер запущен")

	signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	<-done
}

func initFolders() {
	if err := os.Mkdir("uploads", 0600); os.IsNotExist(err) {
		logrus.Error(err)
	}
}
