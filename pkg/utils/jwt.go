package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"instagram-clone.com/m/config"
	"instagram-clone.com/m/internal/models"
)

type Claims struct {
	Email string `json:"email"`
	ID    string `json:"id"`
	jwt.RegisteredClaims
}

func GenerateJWTToken(user *models.User, cfg *config.Config) (string, error) {
	claims := &Claims{
		Email: user.Email,
		ID:    user.ID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(cfg.JWT.Secret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
