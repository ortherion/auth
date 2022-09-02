package models

import "errors"

var (
	ErrUserExist       = errors.New("user already exist")
	ErrUserNotFound    = errors.New("can't find user login")
	ErrInvalidPassword = errors.New("invalid password")

	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token expired")
	ErrTokenClaims  = errors.New("token claims are not of type *tokenClaims")
)
