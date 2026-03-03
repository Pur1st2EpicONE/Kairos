package impl

import (
	"errors"
	"strconv"

	"github.com/golang-jwt/jwt"
)

func (a *AuthService) ParseToken(tokenString string) (int64, error) {

	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (any, error) {
		return []byte("abobus"), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok || !token.Valid {
		return 0, errors.New("invalid token")
	}

	userID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return 0, err
	}

	return userID, nil
}
