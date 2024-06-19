package jwt

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

// Parser is a JWT token reader
type Parser interface {
	Parse(tokenString string) (*jwt.Token, error)
	ValidatePair(accessToken string, refreshToken string) (string, error)
}

// Generator is a JWT token writer
type Generator interface {
	GeneratePair(userGUID string) (*Pair, error)
}

// Errors
var (
	ErrInvalidToken = errors.New("invalid token")
	ErrInvalidPair  = errors.New("invalid pair")
)

// Pair is a pair of access and refresh tokens
type Pair struct {
	AccessToken  string
	RefreshToken string
}

// ManagerOptions is a set of options for the Manager
type ManagerOptions struct {
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

// Manager is a JWT token generator and parser
type Manager struct {
	Options ManagerOptions
}

// NewManager creates a new Manager
func NewManager(options ManagerOptions) *Manager {
	return &Manager{
		Options: options,
	}
}

// GeneratePair generates a pair of access and refresh tokens
func (m *Manager) GeneratePair(userGUID string) (*Pair, error) {
	iat := time.Now()

	accessClaims := jwt.MapClaims{
		"sub": userGUID,
		"exp": iat.Add(m.Options.AccessTTL).Unix(),
		"iat": iat.UnixNano(),
	}

	refreshClaims := jwt.MapClaims{
		"sub":  userGUID,
		"exp":  iat.Add(m.Options.RefreshTTL).Unix(),
		"iat":  iat.UnixNano(),
		"type": "refresh",
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodRS256, accessClaims).SignedString(m.Options.PrivateKey)
	if err != nil {
		return nil, err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshClaims).SignedString(m.Options.PrivateKey)
	if err != nil {
		return nil, err
	}

	return &Pair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// ValidatePair parses a token string
func (m *Manager) ValidatePair(accessToken string, refreshToken string) (string, error) {
	accessClaims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(accessToken, accessClaims, func(token *jwt.Token) (interface{}, error) {
		return m.Options.PublicKey, nil
	})
	if err != nil {
		return "", err
	}

	refreshClaims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(refreshToken, refreshClaims, func(token *jwt.Token) (interface{}, error) {
		return m.Options.PublicKey, nil
	})
	if err != nil {
		return "", err
	}

	if accessClaims["sub"] != refreshClaims["sub"] {
		return "", ErrInvalidPair
	}

	return accessClaims["sub"].(string), nil
}

// Parse parses a token string
func (m *Manager) Parse(tokenString string) (*jwt.Token, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return m.Options.PublicKey, nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	return token, nil
}

// ReadPublicKey reads a public key from a file
func ReadPublicKey(publicKeyPath string) (*rsa.PublicKey, error) {
	f, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key: %w", err)
	}
	key, err := jwt.ParseRSAPublicKeyFromPEM(f)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	return key, nil
}

// ReadPrivateKey reads a private key from a file
func ReadPrivateKey(privateKeyPath string) (*rsa.PrivateKey, error) {
	f, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(f)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return key, nil
}
