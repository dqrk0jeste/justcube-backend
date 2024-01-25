// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: posts.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createPost = `-- name: CreatePost :one
INSERT INTO posts(id, text_content, image_count, user_id)
VALUES ($1, $2, $3, $4)
RETURNING id, text_content, image_count, user_id, created_at
`

type CreatePostParams struct {
	ID          uuid.UUID `json:"id"`
	TextContent string    `json:"text_content"`
	ImageCount  int32     `json:"image_count"`
	UserID      uuid.UUID `json:"user_id"`
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, createPost,
		arg.ID,
		arg.TextContent,
		arg.ImageCount,
		arg.UserID,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.TextContent,
		&i.ImageCount,
		&i.UserID,
		&i.CreatedAt,
	)
	return i, err
}

const deletePost = `-- name: DeletePost :exec
DELETE FROM posts WHERE id = $1
`

func (q *Queries) DeletePost(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deletePost, id)
	return err
}

const getPostById = `-- name: GetPostById :one
SELECT posts.id, text_content, image_count, user_id, posts.created_at, users.id, username, password_hash, users.created_at
FROM posts
INNER JOIN users ON posts.user_id = users.id
WHERE posts.id = $1
LIMIT 1
`

type GetPostByIdRow struct {
	ID           uuid.UUID `json:"id"`
	TextContent  string    `json:"text_content"`
	ImageCount   int32     `json:"image_count"`
	UserID       uuid.UUID `json:"user_id"`
	CreatedAt    time.Time `json:"created_at"`
	ID_2         uuid.UUID `json:"id_2"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt_2  time.Time `json:"created_at_2"`
}

func (q *Queries) GetPostById(ctx context.Context, id uuid.UUID) (GetPostByIdRow, error) {
	row := q.db.QueryRowContext(ctx, getPostById, id)
	var i GetPostByIdRow
	err := row.Scan(
		&i.ID,
		&i.TextContent,
		&i.ImageCount,
		&i.UserID,
		&i.CreatedAt,
		&i.ID_2,
		&i.Username,
		&i.PasswordHash,
		&i.CreatedAt_2,
	)
	return i, err
}

const getPostsByUser = `-- name: GetPostsByUser :many
SELECT posts.id, text_content, image_count, user_id, posts.created_at, users.id, username, password_hash, users.created_at
FROM posts
INNER JOIN users ON posts.user_id = users.id
WHERE user_id = $1
ORDER BY posts.created_at DESC
LIMIT $2 OFFSET $3
`

type GetPostsByUserParams struct {
	UserID uuid.UUID `json:"user_id"`
	Limit  int32     `json:"limit"`
	Offset int32     `json:"offset"`
}

type GetPostsByUserRow struct {
	ID           uuid.UUID `json:"id"`
	TextContent  string    `json:"text_content"`
	ImageCount   int32     `json:"image_count"`
	UserID       uuid.UUID `json:"user_id"`
	CreatedAt    time.Time `json:"created_at"`
	ID_2         uuid.UUID `json:"id_2"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt_2  time.Time `json:"created_at_2"`
}

func (q *Queries) GetPostsByUser(ctx context.Context, arg GetPostsByUserParams) ([]GetPostsByUserRow, error) {
	rows, err := q.db.QueryContext(ctx, getPostsByUser, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPostsByUserRow{}
	for rows.Next() {
		var i GetPostsByUserRow
		if err := rows.Scan(
			&i.ID,
			&i.TextContent,
			&i.ImageCount,
			&i.UserID,
			&i.CreatedAt,
			&i.ID_2,
			&i.Username,
			&i.PasswordHash,
			&i.CreatedAt_2,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
