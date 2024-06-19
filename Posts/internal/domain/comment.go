package domain

import (
	"github.com/google/uuid"
	"time"
)

// Comment is a comment on a post in the domain.
type Comment struct {
	ID        uuid.UUID  `json:"id"`
	PostID    uuid.UUID  `json:"post"`
	ParentID  *uuid.UUID `json:"parent"`
	Content   string     `json:"content"`
	AuthorID  uuid.UUID  `json:"author"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// GetID returns the ID of the comment.
func (c *Comment) GetID() uuid.UUID {
	return c.ID
}

// SetID sets the ID of the comment.
func (c *Comment) SetID(id uuid.UUID) {
	c.ID = id
}
