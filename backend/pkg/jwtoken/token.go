package jwtoken

import (
	"time"
)

type Purpose int

const (
	PurposeAccess = Purpose(iota)
	PurposeRefresh
)

type JWTokenData struct {
	Purpose   Purpose
	Role      int
	ID        int64
	Number    int64
	ExpiresAt time.Time
	Secret    string
}

type tokenManager struct {
	signingKey string
}

//go:generate  mockgen -source token.go -destination ../../mocks/pkg/jwtoken/jwtoken.go

type JWTokenManager interface {
	GenerateToken(data *JWTokenData) (string, error)
	ParseToken(jwtoken string) (*JWTokenData, error)
}

func NewTokenManager(signingKey string) JWTokenManager {
	return &tokenManager{
		signingKey: signingKey,
	}
}
