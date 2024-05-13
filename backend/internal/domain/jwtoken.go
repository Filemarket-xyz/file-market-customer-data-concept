package domain

import (
	"time"
)

const (
	AccessTokenCookieName  = "access_token"
	RefreshTokenCookieName = "refresh_token"
)

type Token struct {
	Token     string
	ExpiresAt time.Time
}

type PairOfTokens struct {
	RefreshToken *Token
	AccessToken  *Token
}

type DropTokensData struct {
	ID     int64
	Number int
}

type DropAllTokensData struct {
	ID int64
}
