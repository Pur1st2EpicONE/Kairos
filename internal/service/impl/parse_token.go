package impl

import (
	"Kairos/internal/errs"
	"strconv"

	"github.com/golang-jwt/jwt"
)

// ParseToken validates a JWT token string and extracts the user ID.
// It uses the service's configured signing key. Returns the user ID
// or an error (ErrInvalidToken if the token is malformed or invalid,
// or ErrInvalidUserID if the subject claim cannot be parsed as int64).
func (a *AuthService) ParseToken(tokenString string) (int64, error) {

	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, a.keyFunc)
	if err != nil {
		return 0, errs.ErrInvalidToken
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok || !token.Valid {
		return 0, errs.ErrInvalidToken
	}

	userID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return 0, errs.ErrInvalidUserID
	}

	return userID, nil

}

// keyFunc returns the signing key for JWT verification.
// It is used as the key function for jwt.ParseWithClaims.
func (a *AuthService) keyFunc(token *jwt.Token) (any, error) {
	return []byte(a.config.TokenSignedString), nil
}
