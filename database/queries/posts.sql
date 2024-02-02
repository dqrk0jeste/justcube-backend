-- name: CreatePost :one
INSERT INTO posts(id, text_content, image_count, user_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetPostById :one
SELECT *
FROM posts
INNER JOIN users ON posts.user_id = users.id
WHERE posts.id = $1
LIMIT 1;

-- name: GetPostsByUser :many
SELECT *
FROM posts
INNER JOIN users ON posts.user_id = users.id
WHERE user_id = $1
ORDER BY posts.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetFeed :many
SELECT *
FROM posts
INNER JOIN users ON posts.user_id = users.id
WHERE user_id IN (
  SELECT followed_user.id
  FROM follows
  INNER JOIN users as followed_user ON followed_user.id = follows.followed_user_id
  WHERE follows.user_id = $1
)
ORDER BY posts.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetGuestFeed :many
SELECT *
FROM posts
INNER JOIN users ON posts.user_id = users.id
ORDER BY posts.created_at DESC
LIMIT $1 OFFSET $2;

-- name: DeletePost :exec
DELETE FROM posts WHERE id = $1;

