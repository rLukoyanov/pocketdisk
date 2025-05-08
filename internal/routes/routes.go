package routes

import (
	"database/sql"
	"pocketdisk/internal/config"
	"pocketdisk/internal/customMiddleware"
	"pocketdisk/internal/handlers"
	"pocketdisk/internal/handlers/api"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitRoutes(e *echo.Echo, sqlite *sql.DB, cfg *config.Config) {

	renderHandlers := handlers.RenderHandlers{Cfg: cfg}
	apiHandlers := handlers.ApiHandlers{Cfg: cfg, DB: sqlite}

	authMiddleware := customMiddleware.AuthMiddleware{Cfg: cfg}

	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	private := e.Group("/")
	private.Use(authMiddleware.AuthMiddlewareRedirect)
	private.GET("", renderHandlers.DashboardPage)

	e.GET("/login", renderHandlers.LoginPage)
	e.Static("/static", "static")

	apiGroup := e.Group("/api")
	apiGroup.POST("/login", apiHandlers.Login)

	fileHandler := api.NewFilesHandler(cfg, sqlite)

	apiGroup.POST("/upload", fileHandler.Upload, authMiddleware.AuthMiddleware)
	apiGroup.GET("/files", fileHandler.GetFiles, authMiddleware.AuthMiddleware)
	apiGroup.DELETE("/files/:id", fileHandler.DeleteFile, authMiddleware.AuthMiddleware)

	apiGroup.GET("/me", apiHandlers.GetUser, authMiddleware.AuthMiddleware)
}
