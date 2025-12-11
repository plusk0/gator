-- name: AddFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
  )
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetFeed :one
SELECT * FROM feeds 
WHERE url = $1 LIMIT 1;

-- name: MarkFetchedFeeds :one
UPDATE feeds
SET last_fetched = NOW()
WHERE ID = $1
RETURNING *;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
ORDER BY last_fetched ASC NULLS FIRST
LIMIT 1;
