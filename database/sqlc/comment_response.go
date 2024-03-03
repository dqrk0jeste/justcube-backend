package database

import (
	"time"

	"github.com/google/uuid"
)

type CommentResponse struct {
	ID              uuid.UUID    `json:"id"`
	Content         string       `json:"content"`
	User            UserResponse `json:"user"`
	NumberOfReplies int64        `json:"number_of_replies"`
	CreatedAt       time.Time    `json:"created_at"`
}

func (comment GetCommentByIdRow) MakeResponse() CommentResponse {
	user := UserResponse{
		ID:        comment.UserID,
		Username:  comment.Username,
		CreatedAt: comment.CreatedAt_2,
	}

	return CommentResponse{
		ID:              comment.ID,
		Content:         comment.Content,
		NumberOfReplies: comment.NumberOfReplies,
		CreatedAt:       comment.CreatedAt,
		User:            user,
	}
}

func (comment GetCommentsByPostRow) MakeResponse() CommentResponse {
	user := UserResponse{
		ID:        comment.UserID,
		Username:  comment.Username,
		CreatedAt: comment.CreatedAt_2,
	}

	return CommentResponse{
		ID:              comment.ID,
		Content:         comment.Content,
		NumberOfReplies: comment.NumberOfReplies,
		CreatedAt:       comment.CreatedAt,
		User:            user,
	}
}

type ReplyResponse struct {
	ID        uuid.UUID    `json:"id"`
	Content   string       `json:"content"`
	User      UserResponse `json:"user"`
	CreatedAt time.Time    `json:"created_at"`
}

func (reply GetReplyByIdRow) MakeResponse() ReplyResponse {
	user := UserResponse{
		ID:        reply.UserID,
		Username:  reply.Username,
		CreatedAt: reply.CreatedAt_2,
	}

	return ReplyResponse{
		ID:        reply.ID,
		Content:   reply.Content,
		CreatedAt: reply.CreatedAt,
		User:      user,
	}
}

func (reply GetRepliesByCommentRow) MakeResponse() ReplyResponse {
	user := UserResponse{
		ID:        reply.UserID,
		Username:  reply.Username,
		CreatedAt: reply.CreatedAt_2,
	}

	return ReplyResponse{
		ID:        reply.ID,
		Content:   reply.Content,
		CreatedAt: reply.CreatedAt,
		User:      user,
	}
}
