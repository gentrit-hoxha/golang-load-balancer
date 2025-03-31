package jwt

import (
	"errors"
	"strings"

	"github.com/golang-jwt/jwt"
)

// Claims structure for JWT
type Claims struct {
	Role string `json:"role"`
	jwt.StandardClaims
}

// Validator for JWT tokens
type Validator struct {
	secretKey []byte
}

// Create new JWT validator
func NewValidator(secret string) *Validator {
	return &Validator{secretKey: []byte(secret)}
}

// Validate incoming JWT token
func (v *Validator) ValidateToken(authHeader string) (*Claims, error) {
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == "" {
		return nil, errors.New("missing token")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return v.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
