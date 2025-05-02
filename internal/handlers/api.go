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
		logrus.Info("Can't get user")
		return errors.New("can't get user")
	}

	file, err := c.FormFile("file")
	if err != nil {
		logrus.Info(err)
		return err
	}

	logrus.Infof("Upload file: %v, size: %v", file.Filename, file.Size)

	tx, err := h.DB.BeginTx(c.Request().Context(), nil)
	if err != nil {
		logrus.Errorf("Failed to begin transaction: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to start transaction")
	}
	defer tx.Rollback()

	var storageUsed, storageLimit int64
	err = tx.QueryRowContext(c.Request().Context(),
		"SELECT storage_used, storage_limit FROM users WHERE id = ?",
		user.ID).Scan(&storageUsed, &storageLimit)

	if err != nil {
		logrus.Errorf("Failed to get user storage info: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "ошибка получения информации об хранилище")
	}

	if storageUsed+file.Size > storageLimit {
		return echo.NewHTTPError(http.StatusForbidden,
			"У вас не достаточно места")
	}

	src, err := file.Open()
	if err != nil {
		logrus.Info(err)
		return err
	}
	defer src.Close()

	fileExt := filepath.Ext(file.Filename)
	fileName := pkg.HashFilename(file.Filename) + fileExt
	dstPath := filepath.Join("./uploads", fileName)

	dst, err := os.Create(dstPath)
	if err != nil {
		logrus.Info(err)
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	_, err = tx.ExecContext(c.Request().Context(), `
        INSERT INTO files (
            user_id, 
            name, 
            path,
            size  
        ) VALUES (?, ?, ?, ?)`,
		user.ID,
		fileName,
		dstPath,
		file.Size,
	)

	if err != nil {
		logrus.Errorf("Failed to save file info to DB: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to save file info")
	}

	_, err = tx.ExecContext(c.Request().Context(),
		"UPDATE users SET storage_used = storage_used + ? WHERE id = ?",
		file.Size, user.ID)

	if err != nil {
		logrus.Errorf("Failed to update user storage: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update storage")
	}

	if err := tx.Commit(); err != nil {
		logrus.Errorf("Failed to commit transaction: %v", err)
		os.Remove(dstPath)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to complete upload")
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "file uploaded successfully",
		"name":    fileName,
		"size":    file.Size,
		"userId":  user.ID,
		"path":    dstPath,
		"storage": echo.Map{
			"used":  storageUsed + file.Size,
			"limit": storageLimit,
		},
	})
}

func (h *ApiHandlers) GetFiles(c echo.Context) error {
	user, ok := c.Get("user").(models.UserTokenInfo)
	if !ok {
		logrus.Info("Can't get user")
		return errors.New("can't get user")
	}

	query := `SELECT id, name, path, size FROM files WHERE user_id = ?`

	rows, err := h.DB.Query(query, user.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to find files")
	}

	defer rows.Close()

	var files []models.FileInfo
	for rows.Next() {
		var f models.FileInfo

		if err := rows.Scan(&f.ID, &f.Name, &f.Path, &f.Size); err != nil {
			logrus.Error(err)
			continue
		}

		files = append(files, f)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"files": files,
		"count": len(files),
	})
}
