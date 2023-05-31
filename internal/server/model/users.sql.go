// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: users.sql

package model

import (
	"context"
)

const cleanUpUsers = `-- name: CleanUpUsers :exec
DELETE FROM users
 WHERE deleted_at IS NOT NULL
`

func (q *Queries) CleanUpUsers(ctx context.Context) error {
	_, err := q.db.Exec(ctx, cleanUpUsers)
	return err
}

const createUser = `-- name: CreateUser :exec
INSERT INTO users (
    password, first_name, last_name, role, email
) VALUES (
    $1, $2, $3, $4, $5
)
`

type CreateUserParams struct {
	Password  []byte `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      Role   `json:"role"`
	Email     string `json:"email"`
}

func (q *Queries) CreateUser(ctx context.Context, arg *CreateUserParams) error {
	_, err := q.db.Exec(ctx, createUser,
		arg.Password,
		arg.FirstName,
		arg.LastName,
		arg.Role,
		arg.Email,
	)
	return err
}

const deleteUser = `-- name: DeleteUser :exec
UPDATE users
   SET deleted_at = CURRENT_TIMESTAMP
 WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteUser, id)
	return err
}

const getUserAuth = `-- name: GetUserAuth :one
SELECT id, email, password FROM users
 WHERE email = $1
   AND deleted_at IS NULl
`

type GetUserAuthRow struct {
	ID       int32  `json:"id"`
	Email    string `json:"email"`
	Password []byte `json:"password"`
}

func (q *Queries) GetUserAuth(ctx context.Context, email string) (*GetUserAuthRow, error) {
	row := q.db.QueryRow(ctx, getUserAuth, email)
	var i GetUserAuthRow
	err := row.Scan(&i.ID, &i.Email, &i.Password)
	return &i, err
}

const updatePassword = `-- name: UpdatePassword :exec
UPDATE users
   SET password = $1,
       password_updated_at = CURRENT_TIMESTAMP
 WHERE id = $2
   AND deleted_at IS NULL
`

type UpdatePasswordParams struct {
	Password []byte `json:"password"`
	ID       int32  `json:"id"`
}

func (q *Queries) UpdatePassword(ctx context.Context, arg *UpdatePasswordParams) error {
	_, err := q.db.Exec(ctx, updatePassword, arg.Password, arg.ID)
	return err
}