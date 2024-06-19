package jwtservice

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// Errors
var (
	ErrInvalidToken = errors.New("invalid token")
	ErrInvalidPair  = errors.New("invalid pair")
)

// Service is a JWT token generator and parser
type Service struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

// NewGenerator creates a new Service
func NewGenerator(secret string, accessTTL time.Duration, refreshTTL time.Duration) *Service {
	return &Service{
		secret:     []byte(secret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

// NewPair generates a pair of access and refresh tokens
func (g *Service) NewPair(sub string) (accessToken string, refreshToken string, err error) {
	iat := time.Now().UnixNano()

	accessClaims := jwt.MapClaims{
		"sub": sub,
		"exp": time.Now().Add(g.accessTTL).Unix(),
		"iat": iat,
	}

	refreshClaims := jwt.MapClaims{
		"sub":  sub,
		"exp":  time.Now().Add(g.refreshTTL).Unix(),
		"iat":  iat,
		"type": "refresh",
	}

	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(g.secret)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(g.secret)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// Parse parses a token string
func (g *Service) Parse(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return g.secret, nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	return token, nil
}

// EnsurePair ensures that the access and refresh tokens are a pair
func (g *Service) EnsurePair(accessToken *jwt.Token, refreshToken *jwt.Token) error {
	accessTokenClaims := accessToken.Claims.(jwt.MapClaims)
	refreshTokenClaims := refreshToken.Claims.(jwt.MapClaims)

	if accessTokenClaims["sub"] != refreshTokenClaims["sub"] || accessTokenClaims["iat"] != refreshTokenClaims["iat"] {
		return ErrInvalidPair
	}

	if refreshTokenClaims["type"] != "refresh" {
		return ErrInvalidPair
	}

	return nil
}
