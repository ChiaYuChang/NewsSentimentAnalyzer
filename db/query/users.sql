-- name: GetUserAuth :one
SELECT id, email, password FROM users
 WHERE email = $1;

-- name: CreateUser :exec
INSERT INTO users (
    password, first_name, last_name, role, email
) VALUES (
    $1, $2, $3, $4, $5
);

-- name: UpdatePassword :exec
UPDATE users
   SET password = $1,
       password_updated_at = CURRENT_TIMESTAMP
 WHERE id = $2;

-- name: DeleteUser :exec
UPDATE users
   SET deleted_at = CURRENT_TIMESTAMP
 WHERE id = $1;

-- name: HardDeleteUser :exec
DELETE FROM users
 WHERE id = $1;