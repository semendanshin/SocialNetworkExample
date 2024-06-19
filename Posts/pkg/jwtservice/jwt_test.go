package jwtservice

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	secret = "secret"

	accessTTL  = 15 * time.Minute
	refreshTTL = 24 * time.Hour
)

func TestService_EnsurePair(t *testing.T) {
	service := NewGenerator(secret, accessTTL, refreshTTL)

	accessToken, refreshToken, err := service.NewPair("sub")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	t.Logf("access token: %s", accessToken)
	t.Logf("refresh token: %s", refreshToken)

	at, err := service.Parse(accessToken)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	rt, err := service.Parse(refreshToken)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	err = service.EnsurePair(at, rt)

	assert.NoError(t, err)
}

func TestService_EnsurePair_Invalid(t *testing.T) {
	service := NewGenerator(secret, accessTTL, refreshTTL)

	accessToken, _, err := service.NewPair("sub")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	_, refreshToken, err := service.NewPair("sub")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	at, err := service.Parse(accessToken)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	rt, err := service.Parse(refreshToken)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	err = service.EnsurePair(at, rt)

	assert.Equal(t, ErrInvalidPair, err)
}

func TestService_Parse(t *testing.T) {
	service := NewGenerator(secret, accessTTL, refreshTTL)

	accessToken, _, err := service.NewPair("sub")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	token, err := service.Parse(accessToken)

	assert.NoError(t, err)
	assert.NotNil(t, token)

	sub, err := token.Claims.GetSubject()

	assert.NoError(t, err)
	assert.Equal(t, "sub", sub)
}

func TestService_Parse_Invalid(t *testing.T) {
	service := NewGenerator(secret, accessTTL, refreshTTL)

	_, err := service.Parse("invalid")

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidToken, err)
}

func TestService_Parse_Expired(t *testing.T) {
	service := NewGenerator(secret, 1*time.Millisecond, refreshTTL)

	accessToken, _, err := service.NewPair("sub")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	time.Sleep(2 * time.Millisecond)

	_, err = service.Parse(accessToken)

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidToken, err)
}

func TestService_NewPair(t *testing.T) {
	service := NewGenerator(secret, accessTTL, refreshTTL)

	accessToken, refreshToken, err := service.NewPair("sub")

	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)
}
