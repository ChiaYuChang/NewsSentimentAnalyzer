-- name: GetKeywordsByNewsId :many
SELECT keyword
  FROM keywords
 WHERE news_id = ANY(@news_id::int[]) 
   AND deleted_at IS NULL;

-- name: CreateKeyword :one
INSERT INTO keywords (
    news_id, keyword
) VALUES (
    $1, $2
)
RETURNING id;

-- name: DeleteKeyword :execrows
DELETE FROM keywords
 WHERE keyword = $1;
