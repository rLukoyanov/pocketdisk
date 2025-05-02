package main

import (
	"os"
	"os/signal"
	"pocketdisk/internal/config"
	"pocketdisk/internal/db"
	"pocketdisk/internal/logger"
	"pocketdisk/internal/pkg"
	"pocketdisk/internal/routes"
	"syscall"

	"github.com/labstack/echo/v4"
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
	logger.InitLogger()

	sqlite, err := db.InitDB(cfg)

	if err != nil {
		logrus.Error(err)
		panic(err)
	}

	e := echo.New()
	pkg.AddNewRender(e)
	routes.InitRoutes(e, sqlite, cfg)

	e.Start(":8080")
	logrus.Info("Сервер запущен")

	signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	<-done
}

func initFolders() {
	if err := os.Mkdir("uploads", 0777); os.IsNotExist(err) {
		logrus.Error(err)
	}
}
