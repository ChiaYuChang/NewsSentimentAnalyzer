// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: keywords.sql

package model

import (
	"context"
)

const createKeyword = `-- name: CreateKeyword :one
INSERT INTO keywords (
    news_id, keyword
) VALUES (
    $1, $2
)
RETURNING id
`

type CreateKeywordParams struct {
	NewsID  int64  `json:"news_id"`
	Keyword string `json:"keyword"`
}

func (q *Queries) CreateKeyword(ctx context.Context, arg *CreateKeywordParams) (int64, error) {
	row := q.db.QueryRow(ctx, createKeyword, arg.NewsID, arg.Keyword)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const deleteKeyword = `-- name: DeleteKeyword :execrows
DELETE FROM keywords
 WHERE keyword = $1
`

func (q *Queries) DeleteKeyword(ctx context.Context, keyword string) (int64, error) {
	result, err := q.db.Exec(ctx, deleteKeyword, keyword)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const getKeywordsByNewsId = `-- name: GetKeywordsByNewsId :many
SELECT keyword
  FROM keywords
 WHERE news_id = ANY($1::int[]) 
   AND deleted_at IS NULL
`

func (q *Queries) GetKeywordsByNewsId(ctx context.Context, newsID []int32) ([]string, error) {
	rows, err := q.db.Query(ctx, getKeywordsByNewsId, newsID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var keyword string
		if err := rows.Scan(&keyword); err != nil {
			return nil, err
		}
		items = append(items, keyword)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
