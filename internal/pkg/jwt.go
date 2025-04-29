package pkg

import (
	"errors"
	"pocketdisk/internal/config"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateJWT(cfg *config.Config, user_id, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"user_id": user_id,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 1 день
	})

	return token.SignedString(cfg.SECRET)
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
