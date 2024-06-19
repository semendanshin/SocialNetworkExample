package domain

import (
	"github.com/google/uuid"
	"time"
)

// User is a user in the domain.
type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetID returns the ID of the user.
func (u *User) GetID() uuid.UUID {
	return u.ID
}

// SetID sets the ID of the user.
func (u *User) SetID(id uuid.UUID) {
	u.ID = id
}
