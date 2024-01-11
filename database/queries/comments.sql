-- name: SendComment :one
INSERT INTO comments(id, content, user_id, post_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeleteComment :exec
DELETE FROM COMMENTS WHERE id = $1;

-- name: GetCommentById :one
SELECT * FROM comments WHERE id = $1 LIMIT 1;