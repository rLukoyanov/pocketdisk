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
	data, ok := c.Request().Context().Value("user").(models.User)
	if !ok {
		c.Redirect(http.StatusUnauthorized, "/login")
	}

	logrus.Info(data)

	return c.Render(http.StatusOK, "index.html", data)
}

type LoginPageData struct {
	Title             string
	Error             string
	Email             string
	Action            string
	ForgotPasswordURL string
	RegisterURL       string
}

func (h *RenderHandlers) LoginPage(c echo.Context) error {
	data := LoginPageData{
		Title:             "Вход",
		Action:            "/api/login",
		ForgotPasswordURL: "/forgot-password",
		RegisterURL:       "/register",
	}

	return c.Render(http.StatusOK, "login.html", data)
}
