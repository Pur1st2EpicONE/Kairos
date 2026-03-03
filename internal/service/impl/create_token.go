package impl

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

func (a *AuthService) CreateToken(userID int64) (string, error) {
	claims := jwt.StandardClaims{
		Subject:   strconv.FormatInt(userID, 10),
		ExpiresAt: time.Now().Add(30 * time.Minute).Unix(),
		IssuedAt:  time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("abobus"))
}
