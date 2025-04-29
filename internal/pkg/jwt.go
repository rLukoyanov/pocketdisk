package pkg

import (
	"errors"
	"pocketdisk/internal/config"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

func GenerateJWT(cfg *config.Config, userID, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"role":   role,
		"exp":    time.Now().Add(time.Hour * 24).Unix(), // 1 день
	})

	logrus.Info(cfg.SECRET)
	return token.SignedString([]byte(cfg.SECRET))
}

func GetJWTClaims(cfg *config.Config, token string) (jwt.Claims, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(cfg.SECRET), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok || !t.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
