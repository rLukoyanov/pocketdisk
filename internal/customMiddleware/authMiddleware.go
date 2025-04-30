package custommiddleware

import (
	"net/http"
	"pocketdisk/internal/config"
	"pocketdisk/internal/handlers"
	"pocketdisk/internal/models"
	"pocketdisk/internal/pkg"

	"github.com/labstack/echo/v4"
)

type Middleware struct {
	Cfg *config.Config
}

func (m *Middleware) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie(handlers.CookieName)
		if err != nil {
			c.Redirect(http.StatusUnauthorized, "/login")
		}

		token := cookie.Value

		claims, err := pkg.GetJWTClaims(m.Cfg, token)
		if err != nil {
			c.Redirect(http.StatusUnauthorized, "/login")
		}

		user := models.UserTokenInfo{
			ID:   claims["userID"].(string),
			Role: claims["role"].(string),
		}

		c.Set("user", user)
		return next(c)
	}
}
