package handlers

import (
	"net/http"
	"pocketdisk/internal/config"
	"pocketdisk/internal/models"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type RenderHandlers struct {
	Cfg *config.Config
}

func (h *RenderHandlers) DashboardPage(c echo.Context) error {
	data, ok := c.Get("user").(models.UserTokenInfo)
	if !ok {
		c.Redirect(http.StatusUnauthorized, "/login")
	}

	logrus.Info(data)

	return c.Render(http.StatusOK, "index.html", data)
}

func (h *RenderHandlers) LoginPage(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", nil)
}
