// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: apis.sql

package models

import (
	"context"
)

const createAPI = `-- name: CreateAPI :exec
INSERT INTO apis (
    name, type
) VALUES (
    $1, $2
)
`

type CreateAPIParams struct {
	Name string  `json:"name"`
	Type ApiType `json:"type"`
}

func (q *Queries) CreateAPI(ctx context.Context, arg *CreateAPIParams) error {
	_, err := q.exec(ctx, q.createAPIStmt, createAPI, arg.Name, arg.Type)
	return err
}

const deleteAPI = `-- name: DeleteAPI :exec
UPDATE apis
   SET deleted_at = CURRENT_TIMESTAMP
 WHERE id = $1
`

func (q *Queries) DeleteAPI(ctx context.Context, id int16) error {
	_, err := q.exec(ctx, q.deleteAPIStmt, deleteAPI, id)
	return err
}

const getAPI = `-- name: GetAPI :one
SELECT id, name, type FROM apis
 WHERE id = $1
`

type GetAPIRow struct {
	ID   int16   `json:"id"`
	Name string  `json:"name"`
	Type ApiType `json:"type"`
}

func (q *Queries) GetAPI(ctx context.Context, id int16) (*GetAPIRow, error) {
	row := q.queryRow(ctx, q.getAPIStmt, getAPI, id)
	var i GetAPIRow
	err := row.Scan(&i.ID, &i.Name, &i.Type)
	return &i, err
}

const hardDeleteAPI = `-- name: HardDeleteAPI :exec
DELETE FROM apis
 WHERE id = $1
`

func (q *Queries) HardDeleteAPI(ctx context.Context, id int16) error {
	_, err := q.exec(ctx, q.hardDeleteAPIStmt, hardDeleteAPI, id)
	return err
}

const updateAPI = `-- name: UpdateAPI :exec
UPDATE apis
   SET name = $1,
       type = $2,
       updated_at = CURRENT_TIMESTAMP
 WHERE id = $3
`

type UpdateAPIParams struct {
	Name string  `json:"name"`
	Type ApiType `json:"type"`
	ID   int16   `json:"id"`
}

func (q *Queries) UpdateAPI(ctx context.Context, arg *UpdateAPIParams) error {
	_, err := q.exec(ctx, q.updateAPIStmt, updateAPI, arg.Name, arg.Type, arg.ID)
	return err
}
