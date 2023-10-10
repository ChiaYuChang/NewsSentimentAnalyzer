-- name: GetNewsPublishBetween :many
SELECT id, title, description, source, related_guid, publish_at
  FROM news
 WHERE publish_at BETWEEN timestamp @from_time AND @to_time
 ORDER BY publish_at;

-- name: GetNewsByMD5Hash :one
SELECT id, title, description, source, related_guid, publish_at
  FROM news
 WHERE md5_hash = $1;

-- name: GetNewsByKeywords :many
SELECT id, title, description, source, related_guid, publish_at
  FROM news
 WHERE id = ANY(
    SELECT news_id
      FROM keywords
    WHERE keyword = ANY(@keywords::string[])
 )
 ORDER BY publish_at;

-- name: GetNewsByJob :many
SELECT id, title, description, source, related_guid, publish_at
  FROM news
 WHERE news.id = ANY(
    SELECT newsjobs.news_id
      FROM jobs
      LEFT JOIN newsjobs
        ON jobs.id = newsjobs.jobs_id
 )
 ORDER BY publish_at;

-- name: GetContentById :many
SELECT id, content
  FROM news
 WHERE id = ANY(@ids::int[]) 
 ORDER BY id;
 
-- name: ListRecentNNews :many
SELECT id, title, description, source, related_guid, publish_at
  FROM news
 ORDER BY publish_at
 LIMIT @n;

-- name: CreateNews :one
INSERT INTO news (
    md5_hash, guid, author, title, link, description, language,
    content, category, source, related_guid, publish_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
)
RETURNING id;

-- name: DeleteNewsPublishBefore :execrows
DELETE FROM news
 WHERE publish_at < @before_time;

-- name: DeleteNews :execrows
DELETE FROM news
 WHERE id = $1;