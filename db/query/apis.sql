-- name: ListAPI :many
SELECT id, name, type 
  FROM apis
 WHERE deleted_at IS NULL
 ORDER BY 
       type ASC,
       name ASC
 LIMIT @n::int;

-- name: CreateAPI :exec
INSERT INTO apis (
    name, type
) VALUES (
    $1, $2
);

-- name: UpdateAPI :exec
UPDATE apis
   SET name = $1,
       type = $2,
       updated_at = CURRENT_TIMESTAMP
 WHERE id = $3
   AND deleted_at IS NULL;

-- name: DeleteAPI :exec
UPDATE apis
   SET deleted_at = CURRENT_TIMESTAMP
 WHERE id = $1;

-- name: CleanUpAPIs :exec
DELETE FROM apis
 WHERE deleted_at IS NOT NULL;
