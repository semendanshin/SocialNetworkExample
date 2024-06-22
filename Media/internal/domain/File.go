package domain

import "github.com/google/uuid"

// File is a file.
type File struct {
	ID       uuid.UUID `json:"id"`
	AuthorID uuid.UUID `json:"authorId"`
	Name     string    `json:"name"`
	Size     int64     `json:"size"`
	Content  []byte    `json:"content"`
}

// GetID returns the ID of the file.
func (f *File) GetID() uuid.UUID {
	return f.ID
}

// SetID sets the ID of the file.
func (f *File) SetID(id uuid.UUID) {
	f.ID = id
}
