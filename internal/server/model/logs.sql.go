// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: logs.sql

package model

import (
	"context"
)

const createLog = `-- name: CreateLog :exec
INSERT INTO logs (
    user_id, type, message
) VALUES (
    $1, $2, $3
)
`

type CreateLogParams struct {
	UserID  int32     `json:"user_id"`
	Type    EventType `json:"type"`
	Message string    `json:"message"`
}

func (q *Queries) CreateLog(ctx context.Context, arg *CreateLogParams) error {
	_, err := q.db.Exec(ctx, createLog, arg.UserID, arg.Type, arg.Message)
	return err
}

const getLogByUserId = `-- name: GetLogByUserId :many
SELECT id, user_id, type, message, created_at
  FROM logs
 WHERE user_id = $1
 ORDER BY
       id DESC,
       created_at DESC,
       type       DESC
 LIMIT $2::int
`

type GetLogByUserIdParams struct {
	UserID int32 `json:"user_id"`
	N      int32 `json:"n"`
}

func (q *Queries) GetLogByUserId(ctx context.Context, arg *GetLogByUserIdParams) ([]*Log, error) {
	rows, err := q.db.Query(ctx, getLogByUserId, arg.UserID, arg.N)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*Log
	for rows.Next() {
		var i Log
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Type,
			&i.Message,
			&i.CreatedAt,
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

const getLogByUserIdNext = `-- name: GetLogByUserIdNext :many
SELECT id, user_id, type, message, created_at
  FROM logs
 WHERE user_id = $1
   AND id > $2
 ORDER BY
       id DESC,
       created_at DESC,
       type       DESC
 LIMIT $3::int
`

type GetLogByUserIdNextParams struct {
	UserID int32 `json:"user_id"`
	ID     int64 `json:"id"`
	N      int32 `json:"n"`
}

func (q *Queries) GetLogByUserIdNext(ctx context.Context, arg *GetLogByUserIdNextParams) ([]*Log, error) {
	rows, err := q.db.Query(ctx, getLogByUserIdNext, arg.UserID, arg.ID, arg.N)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*Log
	for rows.Next() {
		var i Log
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Type,
			&i.Message,
			&i.CreatedAt,
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
