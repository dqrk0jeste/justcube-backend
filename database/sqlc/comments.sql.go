// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: comments.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const deleteComment = `-- name: DeleteComment :exec
DELETE FROM COMMENTS WHERE id = $1
`

func (q *Queries) DeleteComment(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteComment, id)
	return err
}

const getCommentById = `-- name: GetCommentById :one
SELECT id, content, user_id, post_id, created_at FROM comments WHERE id = $1 LIMIT 1
`

func (q *Queries) GetCommentById(ctx context.Context, id uuid.UUID) (Comment, error) {
	row := q.db.QueryRowContext(ctx, getCommentById, id)
	var i Comment
	err := row.Scan(
		&i.ID,
		&i.Content,
		&i.UserID,
		&i.PostID,
		&i.CreatedAt,
	)
	return i, err
}

const sendComment = `-- name: SendComment :one
INSERT INTO comments(id, content, user_id, post_id)
VALUES ($1, $2, $3, $4)
RETURNING id, content, user_id, post_id, created_at
`

type SendCommentParams struct {
	ID      uuid.UUID `json:"id"`
	Content string    `json:"content"`
	UserID  uuid.UUID `json:"user_id"`
	PostID  uuid.UUID `json:"post_id"`
}

func (q *Queries) SendComment(ctx context.Context, arg SendCommentParams) (Comment, error) {
	row := q.db.QueryRowContext(ctx, sendComment,
		arg.ID,
		arg.Content,
		arg.UserID,
		arg.PostID,
	)
	var i Comment
	err := row.Scan(
		&i.ID,
		&i.Content,
		&i.UserID,
		&i.PostID,
		&i.CreatedAt,
	)
	return i, err
}
