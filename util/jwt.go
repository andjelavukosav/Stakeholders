package util

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JwtKey = []byte("tajna_lozinka") // prebaci u .env u produkciji

type Claims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"http://schemas.microsoft.com/ws/2008/06/identity/claims/role"`
	jwt.RegisteredClaims
}

func GenerateToken(id string, username, role string) (string, error) {
	expiration := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		ID:       id,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtKey)
}
