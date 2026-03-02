package impl

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

func (a *AuthService) CreateToken(userID int) (string, error) {
	claims := jwt.StandardClaims{
		Subject:   strconv.Itoa(userID),
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		IssuedAt:  time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("abobus"))
}
