-- name: CreateUser :one
INSERT INTO users(id, username, password_hash, email)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1 LIMIT 1;

-- name: GetUsersByUsername :many
SELECT * FROM users
WHERE username LIKE @input
ORDER BY username ASC
LIMIT $1 OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET username = $1, password_hash = $2
WHERE id = $3
RETURNING *;

-- name: UpdateUsersUsername :one
UPDATE users
SET username = $1
WHERE id = $2
RETURNING *;

-- name: UpdateUsersPassword :one
UPDATE users
SET password_hash = $1
WHERE id = $2
RETURNING *;

