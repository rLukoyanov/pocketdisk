package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"pocketdisk/internal/config"
	"pocketdisk/internal/models"
	"pocketdisk/internal/pkg"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type ApiHandlers struct {
	Cfg *config.Config
	DB  *sql.DB
}

const CookieName = "token"

func (h *ApiHandlers) Login(c echo.Context) error {
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	resetCookie := &http.Cookie{
		Name:     CookieName,
		HttpOnly: true,
		MaxAge:   -1,
	}

	req := &LoginRequest{}
	if err := c.Bind(req); err != nil {
		logrus.Error(err)
		c.SetCookie(resetCookie)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Не все поля заполнены",
		})
	}

	var dbPassword string
	var userID string
	err := h.DB.QueryRow("SELECT id, password FROM users WHERE email = ?", req.Email).Scan(&userID, &dbPassword)
	if err != nil {
		logrus.Error(err)
		if err == sql.ErrNoRows {
			c.SetCookie(resetCookie)
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "Пользователь не найден",
			})
		}
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(req.Password))
	if err != nil {
		logrus.Error(err)
		c.SetCookie(resetCookie)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Не правильный логин или пароль",
		})
	}

	token, err := pkg.GenerateJWT(h.Cfg, userID, "chernorabochiy")
	if err != nil {
		logrus.Error(err)
		c.SetCookie(resetCookie)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Не правильный логин или пароль",
		})
	}

	c.SetCookie(&http.Cookie{
		Name:     CookieName,
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
	})

	return c.JSON(http.StatusOK, "authorized")
}

func (h *ApiHandlers) GetUser(c echo.Context) error {
	user, ok := c.Get("user").(models.UserTokenInfo)
	if !ok {
		logrus.Info("Can't get user")
		return errors.New("can't get user")
	}

	query := `SELECT email, storage_used, storage_limit FROM users WHERE id = ?`

	var dbUser models.User
	err := h.DB.QueryRowContext(c.Request().Context(),
		query,
		user.ID).Scan(&dbUser.Email, &dbUser.StorageUsed, &dbUser.StorageLimit)
	if err != nil {
		logrus.Errorf("Failed to get user info: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get user info")
	}

	return c.JSON(http.StatusOK, echo.Map{
		"user": dbUser,
	})
}
