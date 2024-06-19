package domain

import "github.com/google/uuid"

// Model is a model in the domain.
type Model interface {
	GetID() uuid.UUID
	SetID(id uuid.UUID)
}
