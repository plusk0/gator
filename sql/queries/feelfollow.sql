-- name: CreateFeedFollow :many
WITH inserted_feed_follow AS (
INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING * 
)
  SELECT
    inserted_feed_follow.*,
    feeds.name AS feedname,
    users.name AS username
  FROM inserted_feed_follow
  INNER JOIN feeds
    ON feeds.ID = $5
  INNER JOIN users
    ON users.ID = $4;

-- name: GetFeedFollowsForUser :many
SELECT feed_follows.*,
  users.name AS username,
  feeds.name AS feedname
FROM feed_follows
INNER JOIN users ON users.ID = feed_follows.user_id
INNER JOIN feeds ON feeds.ID = feed_follows.feed_id
WHERE users.name = $1;

-- name: UnFollow :one
DELETE FROM feed_follows 
WHERE user_id = $1 AND feed_id = $2
RETURNING *;
