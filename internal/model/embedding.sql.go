// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: embedding.sql

package model

import (
	"context"

	pgv "github.com/pgvector/pgvector-go"
)

const createEmbedding = `-- name: CreateEmbedding :one
INSERT INTO
    embeddings (
        news_id,
        model,
        embedding,
        sentiment,
        created_at,
        updated_at
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ) RETURNING id
`

type CreateEmbeddingParams struct {
	NewsID    int64      `json:"news_id"`
	Model     string     `json:"model"`
	Embedding pgv.Vector `json:"embedding"`
	Sentiment Sentiment  `json:"sentiment"`
}

func (q *Queries) CreateEmbedding(ctx context.Context, arg *CreateEmbeddingParams) (int64, error) {
	row := q.db.QueryRow(ctx, createEmbedding,
		arg.NewsID,
		arg.Model,
		arg.Embedding,
		arg.Sentiment,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const getEmbeddingByJobId = `-- name: GetEmbeddingByJobId :many
SELECT
    e.id,
    nj.job_id,
    e.news_id,
    e.model,
    e.embedding,
    e.sentiment
FROM embeddings AS e
    INNER JOIN newsjobs AS nj ON e.news_id = nj.news_id
WHERE
    nj.job_id = $1
    AND e.deleted_at IS NULL
`

type GetEmbeddingByJobIdRow struct {
	ID        int64      `json:"id"`
	JobID     int64      `json:"job_id"`
	NewsID    int64      `json:"news_id"`
	Model     string     `json:"model"`
	Embedding pgv.Vector `json:"embedding"`
	Sentiment Sentiment  `json:"sentiment"`
}

func (q *Queries) GetEmbeddingByJobId(ctx context.Context, jobID int64) ([]*GetEmbeddingByJobIdRow, error) {
	rows, err := q.db.Query(ctx, getEmbeddingByJobId, jobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetEmbeddingByJobIdRow
	for rows.Next() {
		var i GetEmbeddingByJobIdRow
		if err := rows.Scan(
			&i.ID,
			&i.JobID,
			&i.NewsID,
			&i.Model,
			&i.Embedding,
			&i.Sentiment,
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

const getEmbeddingByNewsIdsAndModel = `-- name: GetEmbeddingByNewsIdsAndModel :many
SELECT
    id,
    news_id,
    model,
    embedding,
    sentiment
FROM embeddings
WHERE
    news_id = ANY($2:: int [])
    AND model = $1
    AND deleted_at IS NULL
`

type GetEmbeddingByNewsIdsAndModelParams struct {
	Model   string  `json:"model"`
	NewsIds []int32 `json:"news_ids"`
}

type GetEmbeddingByNewsIdsAndModelRow struct {
	ID        int64      `json:"id"`
	NewsID    int64      `json:"news_id"`
	Model     string     `json:"model"`
	Embedding pgv.Vector `json:"embedding"`
	Sentiment Sentiment  `json:"sentiment"`
}

func (q *Queries) GetEmbeddingByNewsIdsAndModel(ctx context.Context, arg *GetEmbeddingByNewsIdsAndModelParams) ([]*GetEmbeddingByNewsIdsAndModelRow, error) {
	rows, err := q.db.Query(ctx, getEmbeddingByNewsIdsAndModel, arg.Model, arg.NewsIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetEmbeddingByNewsIdsAndModelRow
	for rows.Next() {
		var i GetEmbeddingByNewsIdsAndModelRow
		if err := rows.Scan(
			&i.ID,
			&i.NewsID,
			&i.Model,
			&i.Embedding,
			&i.Sentiment,
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
