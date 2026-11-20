-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: CreateFeed :one
INSERT INTO feeds(id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE name = $1 LIMIT 1;

-- name: Reset :exec
DELETE FROM users;

-- name: GetUsers :many
SELECT name FROM users;

-- name: ListFeeds :many
SELECT f.name, f.url, u.name FROM feeds f JOIN users u ON f.user_id = u.id;

-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows(id, created_at, updated_at, user_id, feed_id)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING *
)
SELECT inserted_feed_follow.*, feeds.name AS feed_name, users.name AS user_name
FROM inserted_feed_follow
INNER JOIN users
  ON inserted_feed_follow.user_id = users.id
INNER JOIN feeds
  ON inserted_feed_follow.user_id = feeds.user_id;

-- name: LookUpFeedByURL :one
SELECT id, name
FROM feeds
WHERE url = $1 LIMIT 1;

-- name: GetFeedFollowsForUser :many
WITH idNameTemp AS (
    SELECT feed_follows.user_id, feeds.name
    FROM feed_follows
    INNER JOIN feeds
      ON feeds.id = feed_follows.feed_id
    WHERE feed_follows.user_id = $1
)
SELECT u.name AS UserName, temp.name AS Blog
FROM users u
INNER JOIN idNameTemp temp
  ON temp.user_id = u.id;

-- name: DeleteFeedFollow :exec
DELETE
FROM feed_follows
WHERE (user_id, feed_id) IN (($1, $2));

-- name: MarkFeedFetched :exec
UPDATE feeds
SET updated_at = $2, last_fetched_at = $2
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT *
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;
