package routes

import (
	"database/sql"
	"pocketdisk/internal/config"
	"pocketdisk/internal/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitRoutes(e *echo.Echo, sqlite *sql.DB, cfg *config.Config) {

	renderHandlers := handlers.RenderHandlers{Cfg: cfg}

	e.Use(middleware.Logger())

	e.GET("/", renderHandlers.DashboardPage)
	e.GET("/login", renderHandlers.LoginPage)
}
