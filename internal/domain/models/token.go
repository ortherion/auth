package models

import "time"

// TokenDetails swagger: model TokenDetails
type TokenDetails struct {
	TokenPair
	AtExpires time.Time `json:"-"`
	RtExpires time.Time `json:"-"`
}

// TokenPair swagger: model TokenPair
type TokenPair struct {
	// AccessToken at
	AccessToken string `json:"accessToken"`
	// RefreshToken rt
	RefreshToken string `json:"refreshToken"`
}
