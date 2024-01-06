-- name: CreatePost :one
INSERT INTO posts(id, text_content, image_count, user_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetPostById :one
SELECT * FROM posts WHERE id = $1 LIMIT 1;

-- name: GetPostsByUser :many
SELECT * FROM posts
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

