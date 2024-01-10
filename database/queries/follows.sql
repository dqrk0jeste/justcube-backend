-- name: FollowUser :one
INSERT INTO follows(user_id, followed_user_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetFollows :many
SELECT *, count(followed_user_id) FROM follows
INNER JOIN users as followed_users on follows.followed_user_id = users.id
INNER JOIN users as following_users on follows.user_id = users.id
WHERE follows.followed_user_id = $1;

-- name: GetFollowing :many
SELECT *, count(user_id) FROM follows
INNER JOIN users as followed_users on follows.followed_user_id = users.id
INNER JOIN users as following_users on follows.user_id = users.id
WHERE follows.user_id = $1;

-- name: UnfollowUser :exec
DELETE FROM follows WHERE user_id = $1 AND followed_user_id = $2;

