-- name: CreateUser :one
INSERT INTO users(id, username, password_hash, profile_image)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: GetUsersByUsername :many
SELECT * FROM users
WHERE username LIKE @input
ORDER BY username ASC
LIMIT $1 OFFSET $2;

