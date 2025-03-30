package auth

import (
	"crypto/rand"
	"encoding/hex"
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

func MakeRefreshToken() (string, error) {
	bs := make([]byte, 32)

	_, err := rand.Read(bs)
	if err != nil {
		return "", fmt.Errorf("random 32 bytes of data not generated: %v", err)
	}

	tokenString := hex.EncodeToString(bs)

	return tokenString, nil
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
