package routes

import (
	"database/sql"
	"pocketdisk/internal/config"
	custommiddleware "pocketdisk/internal/customMiddleware"
	"pocketdisk/internal/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitRoutes(e *echo.Echo, sqlite *sql.DB, cfg *config.Config) {

	renderHandlers := handlers.RenderHandlers{Cfg: cfg}
	apiHandlers := handlers.ApiHandlers{Cfg: cfg, DB: sqlite}

	authMiddleware := custommiddleware.AuthMiddleware{Cfg: cfg}

	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	private := e.Group("/")
	private.Use(authMiddleware.AuthMiddlewareRedirect)
	private.GET("", renderHandlers.DashboardPage)

	e.GET("/login", renderHandlers.LoginPage)
	e.Static("/static", "static")

	api := e.Group("/api")
	api.POST("/login", apiHandlers.Login)
	api.POST("/upload", apiHandlers.Upload, authMiddleware.AuthMiddleware)

	// logout
}
