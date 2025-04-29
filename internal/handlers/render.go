package handlers

import (
	"net/http"
	"pocketdisk/internal/config"
	"pocketdisk/internal/models"

	"github.com/labstack/echo/v4"
)

type RenderHandlers struct {
	Cfg *config.Config
}

func (h *RenderHandlers) DashboardPage(c echo.Context) error {
	data, ok := c.Request().Context().Value("user").(models.User)
	if !ok {
		c.Redirect(http.StatusUnauthorized, "/login")
	}

	return c.Render(http.StatusOK, "index.html", data)
}

func (h *RenderHandlers) LoginPage(c echo.Context) error {
	data := map[string]any{"Error": "Хуй login"}

	return c.Render(http.StatusOK, "login.html", data)
}
