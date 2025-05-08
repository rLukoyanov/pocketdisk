package api

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

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type FilesHandler struct {
	cfg *config.Config
	db  *sql.DB
}

func NewFilesHandler(cfg *config.Config, db *sql.DB) *FilesHandler {
	return &FilesHandler{cfg: cfg, db: db}
}

func (h *FilesHandler) getUser(c echo.Context) (models.UserTokenInfo, error) {
	user, ok := c.Get("user").(models.UserTokenInfo)
	if !ok {
		logrus.Warn("User context missing or invalid")
		return models.UserTokenInfo{}, echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}
	return user, nil
}

func (h *FilesHandler) GetFiles(c echo.Context) error {
	user, err := h.getUser(c)
	if err != nil {
		return err
	}

	const query = `SELECT id, name, path, size FROM files WHERE user_id = ?`
	rows, err := h.db.Query(query, user.ID)
	if err != nil {
		logrus.Errorf("DB query failed: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve files")
	}
	defer rows.Close()

	var files []models.FileInfo
	for rows.Next() {
		var f models.FileInfo
		if err := rows.Scan(&f.ID, &f.Name, &f.Path, &f.Size); err != nil {
			logrus.Errorf("Row scan failed: %v", err)
			continue
		}
		files = append(files, f)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"files": files,
		"count": len(files),
	})
}

func (h *FilesHandler) Upload(c echo.Context) error {
	user, err := h.getUser(c)
	if err != nil {
		return err
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		logrus.Warnf("Form file read failed: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid file upload")
	}

	src, err := fileHeader.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to open uploaded file")
	}
	defer src.Close()

	tx, err := h.db.BeginTx(c.Request().Context(), nil)
	if err != nil {
		logrus.Errorf("BeginTx failed: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Transaction start failed")
	}
	defer tx.Rollback()

	var storageUsed, storageLimit int64
	err = tx.QueryRowContext(c.Request().Context(),
		"SELECT storage_used, storage_limit FROM users WHERE id = ?", user.ID).
		Scan(&storageUsed, &storageLimit)

	if err != nil {
		logrus.Errorf("Storage query failed: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve storage info")
	}

	if storageUsed+fileHeader.Size > storageLimit {
		return echo.NewHTTPError(http.StatusForbidden, "Insufficient storage space")
	}

	fileExt := filepath.Ext(fileHeader.Filename)
	fileName := pkg.HashFilename(fileHeader.Filename) + fileExt
	dstPath := filepath.Join("./uploads", fileName)

	dst, err := os.Create(dstPath)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create file")
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save file")
	}

	if _, err := tx.ExecContext(c.Request().Context(),
		`INSERT INTO files (user_id, name, path, size) VALUES (?, ?, ?, ?)`,
		user.ID, fileName, dstPath, fileHeader.Size); err != nil {
		logrus.Errorf("Insert file DB failed: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save file metadata")
	}

	if _, err := tx.ExecContext(c.Request().Context(),
		"UPDATE users SET storage_used = storage_used + ? WHERE id = ?", fileHeader.Size, user.ID); err != nil {
		logrus.Errorf("Storage update failed: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update storage")
	}

	if err := tx.Commit(); err != nil {
		logrus.Errorf("Transaction commit failed: %v", err)
		os.Remove(dstPath)
		return echo.NewHTTPError(http.StatusInternalServerError, "Upload commit failed")
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "File uploaded successfully",
		"name":    fileName,
		"size":    fileHeader.Size,
		"userId":  user.ID,
		"path":    dstPath,
		"storage": echo.Map{
			"used":  storageUsed + fileHeader.Size,
			"limit": storageLimit,
		},
	})
}

func (h *FilesHandler) DeleteFile(c echo.Context) error {
	user, err := h.getUser(c)
	if err != nil {
		return err
	}

	fileID := c.Param("id")
	if fileID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "File ID is required")
	}

	tx, err := h.db.BeginTx(c.Request().Context(), nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Transaction start failed")
	}
	defer tx.Rollback()

	var path string
	var size int64

	err = tx.QueryRowContext(c.Request().Context(),
		`SELECT path, size FROM files WHERE id = ? AND user_id = ?`, fileID, user.ID).
		Scan(&path, &size)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "File not found")
		}
		logrus.Errorf("Query file info failed: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch file data")
	}

	if _, err := tx.ExecContext(c.Request().Context(),
		"DELETE FROM files WHERE id = ?", fileID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete file record")
	}

	if _, err := tx.ExecContext(c.Request().Context(),
		"UPDATE users SET storage_used = storage_used - ? WHERE id = ?", size, user.ID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update storage")
	}

	if err := os.Remove(path); err != nil {
		logrus.Warnf("Failed to delete file from disk: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Deletion commit failed")
	}

	return c.NoContent(http.StatusNoContent)
}
