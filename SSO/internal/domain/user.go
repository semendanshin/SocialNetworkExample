package domain

import (
	"github.com/google/uuid"
	"time"
)

// User is a domain model for users.
type User struct {
	UUID           uuid.UUID `json:"uuid"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	HashedPassword [32]byte  `json:"hashed_password"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// GetID returns the ID of the user.
func (u *User) GetID() uuid.UUID {
	return u.UUID
}

// SetID sets the ID of the user.
func (u *User) SetID(id uuid.UUID) {
	u.UUID = id
}
