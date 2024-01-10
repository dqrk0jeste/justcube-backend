// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package database

import (
	"time"

	"github.com/google/uuid"
)

type Follow struct {
	UserID         uuid.UUID `json:"user_id"`
	FollowedUserID uuid.UUID `json:"followed_user_id"`
	CreatedAt      time.Time `json:"created_at"`
}

type Message struct {
	ID         uuid.UUID `json:"id"`
	Content    string    `json:"content"`
	FromUserID uuid.UUID `json:"from_user_id"`
	ToUserID   uuid.UUID `json:"to_user_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type Post struct {
	ID          uuid.UUID `json:"id"`
	TextContent string    `json:"text_content"`
	ImageCount  int32     `json:"image_count"`
	UserID      uuid.UUID `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
}
