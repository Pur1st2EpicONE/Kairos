package impl

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

// CreateToken generates a JWT token for the given user ID.
// It uses the HS256 signing method and sets standard claims:
// subject (user ID as string), issued at, and expiration based on config.TokenTTL.
// Returns the signed token string or an error.
func (a *AuthService) CreateToken(userID int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   strconv.FormatInt(userID, 10),
		ExpiresAt: time.Now().Add(a.config.TokenTTL).Unix(),
		IssuedAt:  time.Now().Unix()})
	return token.SignedString([]byte(a.config.TokenSignedString))
}
