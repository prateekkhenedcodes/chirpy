package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")

	if !strings.HasPrefix(authHeader, "ApiKey ") {
		return "", fmt.Errorf("invalid authorization header format")
	}

	retString := strings.TrimSpace(strings.TrimPrefix(authHeader, "ApiKey"))

	if retString == "" {
		return "", fmt.Errorf("no apikey")
	}

	return retString, nil

}
