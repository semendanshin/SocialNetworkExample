package domain

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	UUID           uuid.UUID `json:"uuid"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	HashedPassword [32]byte  `json:"hashed_password"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (u *User) GetID() uuid.UUID {
	return u.UUID
}

func (u *User) SetID(id uuid.UUID) {
	u.UUID = id
}
