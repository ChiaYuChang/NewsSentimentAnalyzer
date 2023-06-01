-- name: ListAPIKey :many 
WITH k AS (
  SELECT id, owner, api_id, key
    FROM apikeys
   WHERE owner = $1
     AND deleted_at IS NULL
), a AS (
  SELECT id, name, type
    FROM apis
   WHERE deleted_at IS NULL
) 
SELECT k.id, k.owner, k.api_id, a.name, a.type, k.key 
  FROM k
  LEFT JOIN a
    ON k.api_id = a.id;

-- name: GetAPIKey :many
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