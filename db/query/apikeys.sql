-- name: GetAPIKey :one
SELECT id, owner, api_id, key FROM apikeys
 WHERE owner = $1 
   AND api_id = $2;

-- name: CreateAPIKey :exec
INSERT INTO apikeys (
    owner, api_id, key
) VALUES (
    $1, $2, $3
);

-- name: UpdateAPIKey :exec
UPDATE apikeys
   SET key = $1,
       api_id = $2,
       updated_at = CURRENT_TIMESTAMP
 WHERE owner = $3
   AND api_id = $4;

-- name: DeleteAPIKey :exec
UPDATE apikeys
   SET deleted_at = CURRENT_TIMESTAMP
 WHERE owner = $1
   AND api_id = $2;

-- name: HardAPIKey :exec
DELETE FROM apikeys
 WHERE owner = $1
   AND api_id = $2;