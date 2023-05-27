-- name: GetAPI :one
SELECT id, name, type FROM apis
 WHERE id = $1;

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
 WHERE id = $3;

-- name: DeleteAPI :exec
UPDATE apis
   SET deleted_at = CURRENT_TIMESTAMP
 WHERE id = $1;

-- name: HardDeleteAPI :exec
DELETE FROM apis
 WHERE id = $1;
