package utils

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"os"
	"time"
)

var ErrInvalidToken = errors.New("invalid token")

type tokenClaims struct {
	jwt.StandardClaims
	UserLogin string
}

func GenerateToken(UserLogin string, expiredTime time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiredTime).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserLogin: UserLogin,
	})
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ParseToken(tokenString string) (string, error) {
	claims := new(tokenClaims)
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", ErrInvalidToken
	}
	return claims.UserLogin, nil
}
