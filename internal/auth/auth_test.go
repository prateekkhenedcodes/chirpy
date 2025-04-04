package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "Correct password",
			password: password1,
			hash:     hash1,
			wantErr:  false,
		},
		{
			name:     "Incorrect password",
			password: "wrongpassword",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Password doesn't match different hash",
			password: password1,
			hash:     hash2,
			wantErr:  true,
		},
		{
			name:     "empty password",
			password: "",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "invalid hash",
			password: password1,
			hash:     "invalidhash",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.hash, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMakeValidateJWT(t *testing.T) {
	userId := uuid.New()
	tokenSecret := "your-secret-code"
	expIn := time.Second * 5

	tokenString, err := MakeJWT(userId, tokenSecret, expIn)
	if err != nil {
		t.Fatalf("could not create token: %v", err)
	}

	if tokenString == "" {
		t.Fatal("created token was empty")
	}

	userIdRet, err := ValidateJWT(tokenString, tokenSecret)
	if err != nil {
		t.Fatalf("Error validating the token: %v", err)
	}

	if userId != userIdRet {
		t.Errorf("Expected userId %v, but got %v", userId, userIdRet)
	}

}

func TestStripAuth(t *testing.T) {
	authHeader1 := "Bearer TOKEN_STRING"
	authHeader2 := "Bearer 456789133125"

	tests := []struct {
		name       string
		authHeader string
		wantErr    bool
	}{
		{
			name:       "correct bearer token ",
			authHeader: authHeader1,
			wantErr:    false,
		},
		{
			name:       "incorrect bearer token ",
			authHeader: authHeader1,
			wantErr:    true,
		},
		{
			name:       "empty bearer token ",
			authHeader: "",
			wantErr:    true,
		},
		{
			name:       "correct bearer token ",
			authHeader: authHeader2,
			wantErr:    false,
		},
		{
			name:       "incorrect bearer token",
			authHeader: authHeader2,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retString, err := stripAuth(tt.authHeader)
			if !strings.Contains(tt.authHeader, retString) {
				if (err != nil) != tt.wantErr {
					t.Errorf("stripAuth(%q) = %v; wantErr %v", tt.authHeader, err != nil, tt.wantErr)
				}
			}
		})
	}

}
