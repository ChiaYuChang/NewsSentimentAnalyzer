-- name: ListAPI :many
SELECT id, name, type 
  FROM apis
 WHERE deleted_at IS NULL
 ORDER BY 
       type ASC,
       name ASC
 LIMIT @n::int;

-- name: ListAPIByType :many
SELECT id, name, type
  FROM apis
 WHERE type = @APIType
   AND deleted_at IS NULL
 ORDER BY name ASC;

-- name: GetAPI :one
SELECT *
  FROM apis
 WHERE id = $1
   AND deleted_at IS NULL;

-- name: CreateAPI :one
INSERT INTO apis (
    name, type
) VALUES (
    $1, $2
)
RETURNING id;

-- name: UpdateAPI :execrows
UPDATE apis
   SET name = $1,
       type = $2,
       updated_at = CURRENT_TIMESTAMP
 WHERE id = $3
   AND deleted_at IS NULL;

-- name: DeleteAPI :execrows
UPDATE apis
   SET deleted_at = CURRENT_TIMESTAMP
 WHERE id = $1;

-- name: CleanUpAPIs :execrows
DELETE FROM apis
 WHERE deleted_at IS NOT NULL;
