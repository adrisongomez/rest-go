package utils

import (
	"errors"
	"time"

	"github.com/adrisongomez/project-go/models"
	"github.com/golang-jwt/jwt"
)

func GenerateToken(user *models.User, secret string) (string, error) {
	claims := models.AppClaims{
		UserId: user.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(2 * time.Hour * 24).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateToken(tokenString string, secret string) (*models.AppClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("Error parsing the claims")
}
