// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.0
// source: news.sql

package model

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createNews = `-- name: CreateNews :one
INSERT INTO news (
    md5_hash, guid, author, title, link, description, language,
    content, category, source, related_guid, publish_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
)
RETURNING id
`

type CreateNewsParams struct {
	Md5Hash     string             `json:"md5_hash"`
	Guid        string             `json:"guid"`
	Author      []string           `json:"author"`
	Title       string             `json:"title"`
	Link        string             `json:"link"`
	Description string             `json:"description"`
	Language    pgtype.Text        `json:"language"`
	Content     []string           `json:"content"`
	Category    string             `json:"category"`
	Source      string             `json:"source"`
	RelatedGuid []string           `json:"related_guid"`
	PublishAt   pgtype.Timestamptz `json:"publish_at"`
}

func (q *Queries) CreateNews(ctx context.Context, arg *CreateNewsParams) (int64, error) {
	row := q.db.QueryRow(ctx, createNews,
		arg.Md5Hash,
		arg.Guid,
		arg.Author,
		arg.Title,
		arg.Link,
		arg.Description,
		arg.Language,
		arg.Content,
		arg.Category,
		arg.Source,
		arg.RelatedGuid,
		arg.PublishAt,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const deleteNews = `-- name: DeleteNews :execrows
DELETE FROM news
 WHERE id = $1
`

func (q *Queries) DeleteNews(ctx context.Context, id int64) (int64, error) {
	result, err := q.db.Exec(ctx, deleteNews, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const deleteNewsPublishBefore = `-- name: DeleteNewsPublishBefore :execrows
DELETE FROM news
 WHERE publish_at < $1
`

func (q *Queries) DeleteNewsPublishBefore(ctx context.Context, beforeTime pgtype.Timestamptz) (int64, error) {
	result, err := q.db.Exec(ctx, deleteNewsPublishBefore, beforeTime)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const getContentById = `-- name: GetContentById :many
SELECT id, content
  FROM news
 WHERE id = ANY($1::int[]) 
 ORDER BY id
`

type GetContentByIdRow struct {
	ID      int64    `json:"id"`
	Content []string `json:"content"`
}

func (q *Queries) GetContentById(ctx context.Context, ids []int32) ([]*GetContentByIdRow, error) {
	rows, err := q.db.Query(ctx, getContentById, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetContentByIdRow
	for rows.Next() {
		var i GetContentByIdRow
		if err := rows.Scan(&i.ID, &i.Content); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getNewsByJob = `-- name: GetNewsByJob :many
SELECT id, title, description, source, related_guid, publish_at
  FROM news
 WHERE news.id = ANY(
    SELECT newsjobs.news_id
      FROM jobs
      LEFT JOIN newsjobs
        ON jobs.id = newsjobs.jobs_id
 )
 ORDER BY publish_at
`

type GetNewsByJobRow struct {
	ID          int64              `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Source      string             `json:"source"`
	RelatedGuid []string           `json:"related_guid"`
	PublishAt   pgtype.Timestamptz `json:"publish_at"`
}

func (q *Queries) GetNewsByJob(ctx context.Context) ([]*GetNewsByJobRow, error) {
	rows, err := q.db.Query(ctx, getNewsByJob)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetNewsByJobRow
	for rows.Next() {
		var i GetNewsByJobRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.Source,
			&i.RelatedGuid,
			&i.PublishAt,
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

const getNewsByKeywords = `-- name: GetNewsByKeywords :many
SELECT id, title, description, source, related_guid, publish_at
  FROM news
 WHERE id = ANY(
    SELECT news_id
      FROM keywords
    WHERE keyword = ANY($1::string[])
 )
 ORDER BY publish_at
`

type GetNewsByKeywordsRow struct {
	ID          int64              `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Source      string             `json:"source"`
	RelatedGuid []string           `json:"related_guid"`
	PublishAt   pgtype.Timestamptz `json:"publish_at"`
}

func (q *Queries) GetNewsByKeywords(ctx context.Context, keywords []string) ([]*GetNewsByKeywordsRow, error) {
	rows, err := q.db.Query(ctx, getNewsByKeywords, keywords)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetNewsByKeywordsRow
	for rows.Next() {
		var i GetNewsByKeywordsRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.Source,
			&i.RelatedGuid,
			&i.PublishAt,
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

const getNewsByMD5Hash = `-- name: GetNewsByMD5Hash :one
SELECT id, title, description, source, related_guid, publish_at
  FROM news
 WHERE md5_hash = $1
`

type GetNewsByMD5HashRow struct {
	ID          int64              `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Source      string             `json:"source"`
	RelatedGuid []string           `json:"related_guid"`
	PublishAt   pgtype.Timestamptz `json:"publish_at"`
}

func (q *Queries) GetNewsByMD5Hash(ctx context.Context, md5Hash string) (*GetNewsByMD5HashRow, error) {
	row := q.db.QueryRow(ctx, getNewsByMD5Hash, md5Hash)
	var i GetNewsByMD5HashRow
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.Source,
		&i.RelatedGuid,
		&i.PublishAt,
	)
	return &i, err
}

const getNewsPublishBetween = `-- name: GetNewsPublishBetween :many
SELECT id, title, description, source, related_guid, publish_at
  FROM news
 WHERE publish_at BETWEEN timestamp $1 AND $2
 ORDER BY publish_at
`

type GetNewsPublishBetweenParams struct {
	FromTime pgtype.Timestamptz `json:"from_time"`
	ToTime   pgtype.Timestamptz `json:"to_time"`
}

type GetNewsPublishBetweenRow struct {
	ID          int64              `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Source      string             `json:"source"`
	RelatedGuid []string           `json:"related_guid"`
	PublishAt   pgtype.Timestamptz `json:"publish_at"`
}

func (q *Queries) GetNewsPublishBetween(ctx context.Context, arg *GetNewsPublishBetweenParams) ([]*GetNewsPublishBetweenRow, error) {
	rows, err := q.db.Query(ctx, getNewsPublishBetween, arg.FromTime, arg.ToTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetNewsPublishBetweenRow
	for rows.Next() {
		var i GetNewsPublishBetweenRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.Source,
			&i.RelatedGuid,
			&i.PublishAt,
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

const listRecentNNews = `-- name: ListRecentNNews :many
SELECT id, title, description, source, related_guid, publish_at
  FROM news
 ORDER BY publish_at
 LIMIT $1
`

type ListRecentNNewsRow struct {
	ID          int64              `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Source      string             `json:"source"`
	RelatedGuid []string           `json:"related_guid"`
	PublishAt   pgtype.Timestamptz `json:"publish_at"`
}

func (q *Queries) ListRecentNNews(ctx context.Context, n int32) ([]*ListRecentNNewsRow, error) {
	rows, err := q.db.Query(ctx, listRecentNNews, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*ListRecentNNewsRow
	for rows.Next() {
		var i ListRecentNNewsRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.Source,
			&i.RelatedGuid,
			&i.PublishAt,
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