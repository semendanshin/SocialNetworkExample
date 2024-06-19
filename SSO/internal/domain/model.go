package domain

import "github.com/google/uuid"

// Model is an interface for models.
type Model interface {
	GetID() uuid.UUID
	SetID(id uuid.UUID)
}
