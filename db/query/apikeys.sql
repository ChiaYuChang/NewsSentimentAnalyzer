-- name: ListAPIKey :many
WITH k AS (
  SELECT id, owner, api_id, key
    FROM apikeys
   WHERE owner = $1
     AND deleted_at IS NULL
)
SELECT k.id AS api_key_id, k.owner, k.key, a.id AS api_id, a.type, a.name
  FROM apis AS a
  LEFT JOIN k
    ON a.id = k.api_id
 WHERE a.deleted_at IS NULL;

-- name: GetAPIKey :one
SELECT id, owner, api_id, key 
  FROM apikeys
 WHERE owner = $1 
   AND api_id = $2
   AND deleted_at IS NULL;

-- name: CreateAPIKey :one
INSERT INTO apikeys (
    owner, api_id, key
) VALUES (
    $1, $2, $3
)
RETURNING id;

-- name: UpdateAPIKey :execrows
UPDATE apikeys
   SET key = $1,
       api_id = @old_api_id,
       updated_at = CURRENT_TIMESTAMP
 WHERE owner = $2
   AND api_id = @new_api_id
   AND deleted_at IS NULL;

-- name: DeleteAPIKey :execrows
UPDATE apikeys
   SET deleted_at = CURRENT_TIMESTAMP
 WHERE owner = $1
   AND api_id = $2;

-- name: CleanUpAPIKey :execrows
DELETE FROM apikeys
 WHERE deleted_at IS NOT NULL;