-- name: GetNewsPublishBetween :many
SELECT *
  FROM news
 WHERE publish_at BETWEEN timestamp @from_time AND @to_time
   AND deleted_at IS NULL
 ORDER BY publish_at;

-- name: GetNewsByMD5Hash :one
SELECT *
  FROM news
 WHERE md5_hash = $1 
   AND deleted_at IS NULL;

-- name: GetNewsByKeywords :many
SELECT *
  FROM news
 WHERE id = ANY(
  SELECT news_id
    FROM keywords
   WHERE keyword = ANY(@keywords::string[])
 );

-- name: ListRecentNNews :many
SELECT *
  FROM news
 WHERE deleted_at IS NULL
 ORDER BY publish_at
 LIMIT @n;

-- name: CreateNews :exec
INSERT INTO news (
    md5_hash, title, url, description, content, source, publish_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
);

-- name: DeleteNews :exec
UPDATE news
   SET deleted_at = CURRENT_TIMESTAMP
 WHERE id = $1;

-- name: DeleteNewsPublishBefore :exec
UPDATE news
   SET deleted_at = CURRENT_TIMESTAMP
 WHERE publish_at < @before_time;

-- name: HardDeleteNews :exec
DELETE FROM news
 WHERE deleted_at IS NOT NULL;