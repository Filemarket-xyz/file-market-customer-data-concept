package jwtoken

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

func parseTokenIntClaim(claims jwt.MapClaims, key string) (int, error) {
	if parsedValue, ok := claims[key].(float64); !ok {
		return 0, fmt.Errorf("parseTokenIntClaim: error: invalid token claim: %s", key)
	} else {
		return int(parsedValue), nil
	}
}

func parseTokenInt64Claim(claims jwt.MapClaims, key string) (int64, error) {
	if parsedValue, ok := claims[key].(float64); !ok {
		return 0, fmt.Errorf("parseTokenInt64Claim: error: invalid token claim: %s", key)
	} else {
		return int64(parsedValue), nil
	}
}

func parseTokenStringClaim(claims jwt.MapClaims, key string) (string, error) {
	if stringValue, ok := claims[key].(string); !ok {
		return "", fmt.Errorf("parseTokenStringClaim: error: invalid token claim: %s", key)
	} else {
		return stringValue, nil
	}
}

func (t *tokenManager) ParseToken(token string) (*JWTokenData, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (i interface{}, e error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("ParseToken/Parse: error: unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(t.signingKey), nil
	})
	if err != nil && parsedToken == nil {
		return nil, fmt.Errorf("ParseToken/Parse: parse token failed: %w", err)
	}
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("ParseToken/parsedToken.Claims: error: token wrong claims")
	}
	purpose, err := parseTokenIntClaim(claims, "purpose")
	if err != nil {
		return nil, fmt.Errorf("ParseToken/parseTokenIntClaim/purpose: %w", err)
	}
	if purpose != int(PurposeAccess) && purpose != int(PurposeRefresh) {
		return nil, fmt.Errorf("ParseToken: error: invalid purpose: %d", purpose)
	}
	id, err := parseTokenInt64Claim(claims, "id")
	if err != nil {
		return nil, fmt.Errorf("ParseToken/parseTokenIntClaim/id: %w", err)
	}
	number, err := parseTokenInt64Claim(claims, "number")
	if err != nil {
		return nil, fmt.Errorf("ParseToken/parseTokenIntClaim/number: %w", err)
	}
	role, err := parseTokenIntClaim(claims, "role")
	if err != nil {
		return nil, fmt.Errorf("ParseToken/parseTokenIntClaim/role: %w", err)
	}
	secret, err := parseTokenStringClaim(claims, "secret")
	if err != nil {
		return nil, fmt.Errorf("ParseToken/parseTokenIntClaim/secret: %w", err)
	}
	expiresAt, err := parseTokenInt64Claim(claims, "exp")
	if err != nil {
		return nil, fmt.Errorf("ParseToken/parseTokenIntClaim/exp: %w", err)
	}
	return &JWTokenData{
		Purpose:   Purpose(purpose),
		Role:      int(role),
		ID:        id,
		Number:    number,
		ExpiresAt: time.Unix(expiresAt, 0),
		Secret:    secret,
	}, nil
}
