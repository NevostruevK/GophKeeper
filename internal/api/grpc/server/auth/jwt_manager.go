package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

var (
	// ErrUnexpectedTokenSigningMethod error: unexpected token signing method.
	ErrUnexpectedTokenSigningMethod = errors.New("unexpected token signing method")
	// ErrInvalidTokenClaims error: invalid token claims.
	ErrInvalidTokenClaims = errors.New("invalid token claims")
)

// JWTManager JWT manager for gRPC AuthServer.
type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

// UserClaims claims for JWT manager.
type UserClaims struct {
	jwt.StandardClaims
}

// NewJWTManager returns JWTManager.
func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{secretKey, tokenDuration}
}

// Generate generates JWT token.
func (manager *JWTManager) Generate(id uuid.UUID) (string, error) {
	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(manager.tokenDuration).Unix(),
			Id:        id.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(manager.secretKey))
}

// Verify verifies JWT token and returns user ID.
func (manager *JWTManager) Verify(accessToken string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, ErrUnexpectedTokenSigningMethod
			}

			return []byte(manager.secretKey), nil
		},
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return uuid.Nil, ErrInvalidTokenClaims
	}

	return uuid.Parse(claims.Id)
}
