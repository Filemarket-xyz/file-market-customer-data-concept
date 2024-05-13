package jwtoken

import (
	"fmt"

	"github.com/golang-jwt/jwt"
)

func (t *tokenManager) GenerateToken(data *JWTokenData) (string, error) {
	claims := jwt.MapClaims{
		"id":      data.ID,
		"role":    data.Role,
		"purpose": int64(data.Purpose),
		"secret":  data.Secret,
		"exp":     data.ExpiresAt.Unix(),
		"number":  data.Number,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	res, err := token.SignedString([]byte(t.signingKey))
	if err != nil {
		return "", fmt.Errorf("GenerateToken/SignedString: sign token failed: %w", err)
	}
	return res, nil
}
