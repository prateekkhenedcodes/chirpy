package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return "", fmt.Errorf("unexpected signing method, %s", t.Header["alg"])
		}
		return []byte(tokenSecret), nil

	})

	if err != nil {
		return uuid.Nil, fmt.Errorf("error parsing the token %s", err)
	}

	if !token.Valid {
		return uuid.Nil, fmt.Errorf("token is not valid ")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid claims")
	}

	userIDStr, ok := claims["sub"].(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("user ID claim not found or invalid")
	}

	userId, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID: %s", err)
	}

	return userId, nil

}
