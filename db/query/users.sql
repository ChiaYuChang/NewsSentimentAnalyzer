-- name: GetUserAuth :one
SELECT id, email, password, role FROM users
 WHERE email = $1
   AND deleted_at IS NULl;

-- name: CreateUser :one
INSERT INTO users (
    password, first_name, last_name, role, email
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING id;

-- name: UpdatePassword :execrows
UPDATE users
   SET password = $1,
       password_updated_at = CURRENT_TIMESTAMP
 WHERE id = $2
   AND deleted_at IS NULL;

-- name: DeleteUser :execrows
UPDATE users
   SET deleted_at = CURRENT_TIMESTAMP
 WHERE id = $1;

-- name: HardDeleteUser :execrows
DELETE FROM users
 WHERE id = $1;

-- name: CleanUpUsers :execrows
DELETE FROM users
 WHERE deleted_at IS NOT NULL;