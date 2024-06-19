package entities

import (
	"github.com/google/uuid"
	"time"
)

// User represents a user entity
type User struct {
	ID             uuid.UUID `json:"uuid" gorm:"primaryKey"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"hashed_password"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
