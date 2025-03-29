package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")

	if authHeader == "" {
		return "", fmt.Errorf("authorization header is missing ")
	}

	token, err := stripAuth(authHeader)
	if err != nil {
		return "", err
	}

	return token, nil

}

func stripAuth(s string) (string, error) {

	if !strings.HasPrefix(s, "Bearer ") {
		return "", fmt.Errorf("invalid authorization header format")
	}

	retString := strings.TrimSpace(strings.TrimPrefix(s, "Bearer"))

	if retString == "" {
		return "", fmt.Errorf("empty bearer token")
	}

	return retString, nil
}
