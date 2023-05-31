-- name: GetKeywordsByNewsId :many
SELECT keyword
  FROM keywords
 WHERE news_id = ANY(@news_id::int[]) 
   AND deleted_at IS NULL;

-- name: CreateKeyword :exec
INSERT INTO keywords (
    news_id, keyword
) VALUES (
    $1, $2
);

-- name: DeleteKeyword :exec
UPDATE keywords
   SET deleted_at = CURRENT_TIMESTAMP
 WHERE keyword = $1;

-- name: DeleteKeywordByNewsId :exec
UPDATE keywords
   SET deleted_at = CURRENT_TIMESTAMP
 WHERE news_id = $1;

-- name: CleanUpKeywords :exec
DELETE FROM keywords
 WHERE deleted_at IS NOT NULL;