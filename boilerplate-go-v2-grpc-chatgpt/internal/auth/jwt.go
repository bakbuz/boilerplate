package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

func ValidateToken(tokenStr, secret string) (string, error) {
	t, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenMalformed
		}
		return []byte(secret), nil
	})
	if err != nil || !t.Valid {
		return "", err
	}
	claims := t.Claims.(jwt.MapClaims)
	sub := claims["sub"].(string)
	return sub, nil
}
