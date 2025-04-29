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
	apiHandlers := handlers.ApiHandlers{Cfg: cfg, DB: sqlite}

	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	private := e.Group("/")
	// private.Use(echojwt.JWT([]byte(cfg.SECRET)))
	private.GET("", renderHandlers.DashboardPage)

	e.GET("/login", renderHandlers.LoginPage)
	e.Static("/static", "static")

	api := e.Group("/api")
	api.POST("/login", apiHandlers.Login)
}
