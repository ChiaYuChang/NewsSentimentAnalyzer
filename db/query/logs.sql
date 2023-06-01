-- name: GetLogByUserId :many
SELECT *
  FROM logs
 WHERE user_id = $1
 ORDER BY
       id DESC,
       created_at DESC,
       type       DESC
 LIMIT @n::int;

-- name: GetLogByUserIdNext :many
SELECT *
  FROM logs
 WHERE user_id = $1
   AND id > $2
 ORDER BY
       id DESC,
       created_at DESC,
       type       DESC
 LIMIT @n::int;

-- name: CreateLog :one
INSERT INTO logs (
    user_id, type, message
) VALUES (
    $1, $2, $3
)
RETURNING id;
