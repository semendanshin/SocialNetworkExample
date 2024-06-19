package entities

import (
	"github.com/google/uuid"
	"time"
)

// Comment is a comment on a post in gorm.
type Comment struct {
	ID        uuid.UUID  `json:"id" gorm:"primary_key"`
	PostID    uuid.UUID  `json:"postId"`
	ParentID  *uuid.UUID `json:"parentId"`
	AuthorID  uuid.UUID  `json:"authorId"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

// Post is a post in gorm.
type Post struct {
	ID            uuid.UUID `json:"id" gorm:"primary_key"`
	AuthorID      uuid.UUID `json:"authorId"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	AllowComments bool      `json:"allowComments"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// User is a user in gorm.
type User struct {
	ID   uuid.UUID `json:"id" gorm:"primary_key"`
	Name string    `json:"name"`
}
