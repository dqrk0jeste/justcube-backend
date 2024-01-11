-- name: FollowUser :one
INSERT INTO follows(user_id, followed_user_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetFollowers :many
SELECT following_users.* FROM follows
INNER JOIN users as followed_users on follows.followed_user_id = followed_users.id
INNER JOIN users as following_users on follows.user_id = following_users.id
WHERE follows.followed_user_id = $1
LIMIT $2 OFFSET $3;

-- name: GetFollowing :many
SELECT followed_users.* FROM follows
INNER JOIN users as followed_users on follows.followed_user_id = followed_users.id
INNER JOIN users as following_users on follows.user_id = following_users.id
WHERE follows.user_id = $1
LIMIT $2 OFFSET $3;

-- name: UnfollowUser :exec
DELETE FROM follows WHERE user_id = $1 AND followed_user_id = $2;

