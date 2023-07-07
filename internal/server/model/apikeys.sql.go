// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: apikeys.sql

package model

import (
	"context"
)

const cleanUpAPIKey = `-- name: CleanUpAPIKey :execrows
DELETE FROM apikeys
 WHERE deleted_at IS NOT NULL
`

func (q *Queries) CleanUpAPIKey(ctx context.Context) (int64, error) {
	result, err := q.db.Exec(ctx, cleanUpAPIKey)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const createAPIKey = `-- name: CreateAPIKey :one
INSERT INTO apikeys (
    owner, api_id, key
) VALUES (
    $1, $2, $3
)
RETURNING id
`

type CreateAPIKeyParams struct {
	Owner int32  `json:"owner"`
	ApiID int16  `json:"api_id"`
	Key   string `json:"key"`
}

func (q *Queries) CreateAPIKey(ctx context.Context, arg *CreateAPIKeyParams) (int32, error) {
	row := q.db.QueryRow(ctx, createAPIKey, arg.Owner, arg.ApiID, arg.Key)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const deleteAPIKey = `-- name: DeleteAPIKey :execrows
UPDATE apikeys
   SET deleted_at = CURRENT_TIMESTAMP
 WHERE owner = $1
   AND api_id = $2
`

type DeleteAPIKeyParams struct {
	Owner int32 `json:"owner"`
	ApiID int16 `json:"api_id"`
}

func (q *Queries) DeleteAPIKey(ctx context.Context, arg *DeleteAPIKeyParams) (int64, error) {
	result, err := q.db.Exec(ctx, deleteAPIKey, arg.Owner, arg.ApiID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const getAPIKey = `-- name: GetAPIKey :one
SELECT id, owner, api_id, key 
  FROM apikeys
 WHERE owner = $1 
   AND api_id = $2
   AND deleted_at IS NULL
`

type GetAPIKeyParams struct {
	Owner int32 `json:"owner"`
	ApiID int16 `json:"api_id"`
}

type GetAPIKeyRow struct {
	ID    int32  `json:"id"`
	Owner int32  `json:"owner"`
	ApiID int16  `json:"api_id"`
	Key   string `json:"key"`
}

func (q *Queries) GetAPIKey(ctx context.Context, arg *GetAPIKeyParams) (*GetAPIKeyRow, error) {
	row := q.db.QueryRow(ctx, getAPIKey, arg.Owner, arg.ApiID)
	var i GetAPIKeyRow
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.ApiID,
		&i.Key,
	)
	return &i, err
}

const listAPIKey = `-- name: ListAPIKey :many
WITH k AS (
  SELECT id, owner, api_id, key
    FROM apikeys
   WHERE owner = $1
     AND deleted_at IS NULL
)
SELECT k.id AS api_key_id, k.owner, k.key, 
       a.id AS api_id, a.type, a.name, a.image, a.icon
  FROM apis AS a
 INNER JOIN k
    ON a.id = k.api_id
 WHERE a.deleted_at IS NULL
`

type ListAPIKeyRow struct {
	ApiKeyID int32   `json:"api_key_id"`
	Owner    int32   `json:"owner"`
	Key      string  `json:"key"`
	ApiID    int16   `json:"api_id"`
	Type     ApiType `json:"type"`
	Name     string  `json:"name"`
	Image    string  `json:"image"`
	Icon     string  `json:"icon"`
}

func (q *Queries) ListAPIKey(ctx context.Context, owner int32) ([]*ListAPIKeyRow, error) {
	rows, err := q.db.Query(ctx, listAPIKey, owner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*ListAPIKeyRow
	for rows.Next() {
		var i ListAPIKeyRow
		if err := rows.Scan(
			&i.ApiKeyID,
			&i.Owner,
			&i.Key,
			&i.ApiID,
			&i.Type,
			&i.Name,
			&i.Image,
			&i.Icon,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateAPIKey = `-- name: UpdateAPIKey :execrows
UPDATE apikeys
   SET key = $1,
       api_id = $3,
       updated_at = CURRENT_TIMESTAMP
 WHERE owner = $2
   AND api_id = $4
   AND deleted_at IS NULL
`

type UpdateAPIKeyParams struct {
	Key      string `json:"key"`
	Owner    int32  `json:"owner"`
	OldApiID int16  `json:"old_api_id"`
	NewApiID int16  `json:"new_api_id"`
}

func (q *Queries) UpdateAPIKey(ctx context.Context, arg *UpdateAPIKeyParams) (int64, error) {
	result, err := q.db.Exec(ctx, updateAPIKey,
		arg.Key,
		arg.Owner,
		arg.OldApiID,
		arg.NewApiID,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
