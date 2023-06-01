// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: apis.sql

package model

import (
	"context"
)

const cleanUpAPIs = `-- name: CleanUpAPIs :execrows
DELETE FROM apis
 WHERE deleted_at IS NOT NULL
`

func (q *Queries) CleanUpAPIs(ctx context.Context) (int64, error) {
	result, err := q.db.Exec(ctx, cleanUpAPIs)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const createAPI = `-- name: CreateAPI :one
INSERT INTO apis (
    name, type
) VALUES (
    $1, $2
)
RETURNING id
`

type CreateAPIParams struct {
	Name string  `json:"name"`
	Type ApiType `json:"type"`
}

func (q *Queries) CreateAPI(ctx context.Context, arg *CreateAPIParams) (int16, error) {
	row := q.db.QueryRow(ctx, createAPI, arg.Name, arg.Type)
	var id int16
	err := row.Scan(&id)
	return id, err
}

const deleteAPI = `-- name: DeleteAPI :execrows
UPDATE apis
   SET deleted_at = CURRENT_TIMESTAMP
 WHERE id = $1
`

func (q *Queries) DeleteAPI(ctx context.Context, id int16) (int64, error) {
	result, err := q.db.Exec(ctx, deleteAPI, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const listAPI = `-- name: ListAPI :many
SELECT id, name, type 
  FROM apis
 WHERE deleted_at IS NULL
 ORDER BY 
       type ASC,
       name ASC
 LIMIT $1::int
`

type ListAPIRow struct {
	ID   int16   `json:"id"`
	Name string  `json:"name"`
	Type ApiType `json:"type"`
}

func (q *Queries) ListAPI(ctx context.Context, n int32) ([]*ListAPIRow, error) {
	rows, err := q.db.Query(ctx, listAPI, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*ListAPIRow
	for rows.Next() {
		var i ListAPIRow
		if err := rows.Scan(&i.ID, &i.Name, &i.Type); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateAPI = `-- name: UpdateAPI :execrows
UPDATE apis
   SET name = $1,
       type = $2,
       updated_at = CURRENT_TIMESTAMP
 WHERE id = $3
   AND deleted_at IS NULL
`

type UpdateAPIParams struct {
	Name string  `json:"name"`
	Type ApiType `json:"type"`
	ID   int16   `json:"id"`
}

func (q *Queries) UpdateAPI(ctx context.Context, arg *UpdateAPIParams) (int64, error) {
	result, err := q.db.Exec(ctx, updateAPI, arg.Name, arg.Type, arg.ID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
