package custommiddleware

import (
	"net/http"
	"pocketdisk/internal/config"
	"pocketdisk/internal/handlers"
	"pocketdisk/internal/models"
	"pocketdisk/internal/pkg"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type AuthMiddleware struct {
	Cfg *config.Config
}

func (m *AuthMiddleware) AuthMiddlewareRedirect(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie(handlers.CookieName)
		if err != nil {
			logrus.Error("Cookie not found: ", err)
			return c.Redirect(http.StatusFound, "/login")
		}

		token := cookie.Value

		claims, err := pkg.GetJWTClaims(m.Cfg, token)
		if err != nil {
			logrus.Error("JWT validation failed: ", err)
			return c.Redirect(http.StatusFound, "/login")
		}

		userID, ok := claims["userID"].(string)
		if !ok || userID == "" {
			logrus.Error("Invalid role in claims")
			return c.Redirect(http.StatusFound, "/login")
		}

		role, ok := claims["role"].(string)
		if !ok || role == "" {
			logrus.Error("Invalid role in claims")
			return c.Redirect(http.StatusFound, "/login")
		}

		user := models.UserTokenInfo{
			ID:   userID,
			Role: role,
		}

		c.Set("user", user)
		return next(c)
	}
}

func (m *AuthMiddleware) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie(handlers.CookieName)
		if err != nil {
			logrus.Error("Cookie not found: ", err)
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Unauthorized: missing or invalid cookie",
			})
		}

		token := cookie.Value
		claims, err := pkg.GetJWTClaims(m.Cfg, token)
		if err != nil {
			logrus.Error("JWT validation failed: ", err)
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Unauthorized: invalid token",
			})
		}

		userID, ok := claims["userID"].(string)
		if !ok || userID == "" {
			logrus.Error("Invalid userID in claims")
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Unauthorized: invalid user ID",
			})
		}

		role, ok := claims["role"].(string)
		if !ok || role == "" {
			logrus.Error("Invalid role in claims")
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Unauthorized: invalid role",
			})
		}

		user := models.UserTokenInfo{
			ID:   userID,
			Role: role,
		}
		c.Set("user", user)

		return next(c)
	}
}
