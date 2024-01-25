package database

import (
	"time"

	"github.com/google/uuid"
)

type PostResponse struct {
	ID          uuid.UUID    `json:"id"`
	TextContent string       `json:"text_content"`
	ImageCount  int32        `json:"image_count"`
	User        UserResponse `json:"user"`
	CreatedAt   time.Time    `json:"created_at"`
}

func (post GetPostByIdRow) MakeResponse() PostResponse {
	user := UserResponse{
		ID:        post.UserID,
		Username:  post.Username,
		CreatedAt: post.CreatedAt_2,
	}

	return PostResponse{
		ID:          post.ID,
		TextContent: post.TextContent,
		ImageCount:  post.ImageCount,
		CreatedAt:   post.CreatedAt,
		User:        user,
	}
}

func (post GetPostsByUserRow) MakeResponse() PostResponse {
	user := UserResponse{
		ID:        post.UserID,
		Username:  post.Username,
		CreatedAt: post.CreatedAt_2,
	}

	return PostResponse{
		ID:          post.ID,
		TextContent: post.TextContent,
		ImageCount:  post.ImageCount,
		CreatedAt:   post.CreatedAt,
		User:        user,
	}
}
