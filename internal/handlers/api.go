package handlers

import (
	"database/sql"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
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
	err := h.DB.QueryRow("SELECT id, password FROM users WHERE email = ?", req.Email).Scan(&userID, &dbPassword) // заменить на squirell
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

func (h *ApiHandlers) Upload(c echo.Context) error {
	user, ok := c.Get("user").(models.UserTokenInfo)
	if !ok {
		logrus.Info("Cant get user")
		return errors.New("cant get user")
	}

	file, err := c.FormFile("file")
	if err != nil {
		logrus.Info(err)
		return err
	}

	logrus.Infof("Upload file: %v, size: %v", file.Filename, file.Size)

	src, err := file.Open()
	if err != nil {
		logrus.Info(err)
		return err
	}
	defer src.Close()

	dstPath := filepath.Join("./uploads", file.Filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		logrus.Info(err)
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message":  "file created",
		"filename": file.Filename,
		"size":     file.Size,
		"forUser":  user.ID,
	})
}
