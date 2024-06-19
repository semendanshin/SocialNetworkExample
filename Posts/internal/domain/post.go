package domain

import (
	"github.com/google/uuid"
	"time"
)

// Post is a post in the domain.
type Post struct {
	ID            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	AuthorID      uuid.UUID `json:"author"`
	AllowComments bool      `json:"allow_comments"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// GetID returns the ID of the post.
func (p *Post) GetID() uuid.UUID {
	return p.ID
}

// SetID sets the ID of the post.
func (p *Post) SetID(id uuid.UUID) {
	p.ID = id
}

// DisableComments disables comments on the post.
func (p *Post) DisableComments() {
	p.AllowComments = false
}

// EnableComments enables comments on the post.
func (p *Post) EnableComments() {
	p.AllowComments = true
}
