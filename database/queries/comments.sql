-- name: PostComment :one
INSERT INTO comments(id, content, user_id, post_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeleteComment :exec
DELETE FROM comments WHERE id = $1;

-- name: GetCommentById :one
SELECT comments.*, users.* , COUNT(replies.id) as number_of_replies
FROM comments
INNER JOIN users ON posts.user_id = users.id
LEFT JOIN replies ON replies.comment_id = comments.id
WHERE comments.id = $1
LIMIT 1;

-- name: GetCommentsByPost :many
SELECT *
FROM 
  ( SELECT comments.*, COUNT(replies.id) as number_of_replies
  FROM comments
  LEFT JOIN replies ON replies.comment_id = comments.id
  WHERE post_id = $1
  GROUP BY comments.id ) as c
INNER JOIN users ON c.user_id = users.id
ORDER BY number_of_replies DESC, c.created_at DESC
LIMIT $2 OFFSET $3;

-- name: PostReply :one
INSERT INTO replies(id, content, user_id, comment_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeleteReply :exec
DELETE FROM replies WHERE id = $1;

-- name: GetReplyById :one
SELECT *
FROM replies
INNER JOIN users ON replies.user_id = users.id
WHERE replies.id = $1
LIMIT 1;

-- name: GetRepliesByComment :many
SELECT *
FROM replies
INNER JOIN users ON replies.user_id = users.id
WHERE comment_id = $1
ORDER BY replies.created_at ASC
LIMIT $2 OFFSET $3;

