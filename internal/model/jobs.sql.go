// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.0
// source: jobs.sql

package model

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const cleanUpJobs = `-- name: CleanUpJobs :execrows
DELETE FROM jobs
 WHERE deleted_at IS NOT NULL
`

func (q *Queries) CleanUpJobs(ctx context.Context) (int64, error) {
	result, err := q.db.Exec(ctx, cleanUpJobs)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const countJob = `-- name: CountJob :many
SELECT status, COUNT(*) AS n_job
  FROM jobs
 WHERE owner = $1
 GROUP BY status
 ORDER BY 
        status ASC
`

type CountJobRow struct {
	Status JobStatus `json:"status"`
	NJob   int64     `json:"n_job"`
}

func (q *Queries) CountJob(ctx context.Context, owner uuid.UUID) ([]*CountJobRow, error) {
	rows, err := q.db.Query(ctx, countJob, owner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*CountJobRow
	for rows.Next() {
		var i CountJobRow
		if err := rows.Scan(&i.Status, &i.NJob); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const createJob = `-- name: CreateJob :one
INSERT INTO jobs (
  owner, status, src_api_id, src_query, llm_api_id, llm_query
) VALUES (
    $1, $2, $3, $4, $5, $6
) 
RETURNING id
`

type CreateJobParams struct {
	Owner    uuid.UUID `json:"owner"`
	Status   JobStatus `json:"status"`
	SrcApiID int16     `json:"src_api_id"`
	SrcQuery string    `json:"src_query"`
	LlmApiID int16     `json:"llm_api_id"`
	LlmQuery []byte    `json:"llm_query"`
}

func (q *Queries) CreateJob(ctx context.Context, arg *CreateJobParams) (int32, error) {
	row := q.db.QueryRow(ctx, createJob,
		arg.Owner,
		arg.Status,
		arg.SrcApiID,
		arg.SrcQuery,
		arg.LlmApiID,
		arg.LlmQuery,
	)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const deleteJob = `-- name: DeleteJob :execrows
UPDATE jobs
   SET deleted_at = CURRENT_TIMESTAMP
 WHERE id = $1
   AND owner = $2
`

type DeleteJobParams struct {
	ID    int32     `json:"id"`
	Owner uuid.UUID `json:"owner"`
}

func (q *Queries) DeleteJob(ctx context.Context, arg *DeleteJobParams) (int64, error) {
	result, err := q.db.Exec(ctx, deleteJob, arg.ID, arg.Owner)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const getJobByOwnerFilterByJIdAndStatus = `-- name: GetJobByOwnerFilterByJIdAndStatus :many
SELECT j.id, j.owner, j.status, asrc.name AS news_src, allm.name AS analyzer, j.created_at, j.updated_at
  FROM jobs AS j
 INNER JOIN apis AS asrc ON j.src_api_id = asrc.id
 INNER JOIN apis AS allm ON j.llm_api_id = allm.id 
 WHERE j.owner = $1
   AND j.id BETWEEN $2::int AND $3::int
   AND j.status = $4
   AND j.deleted_at IS NULL
 ORDER BY 
       j.id DESC
 LIMIT $5::int
`

type GetJobByOwnerFilterByJIdAndStatusParams struct {
	Owner   uuid.UUID `json:"owner"`
	FJid    int32     `json:"f_jid"`
	TJid    int32     `json:"t_jid"`
	JStatus JobStatus `json:"j_status"`
	N       int32     `json:"n"`
}

type GetJobByOwnerFilterByJIdAndStatusRow struct {
	ID        int32              `json:"id"`
	Owner     uuid.UUID          `json:"owner"`
	Status    JobStatus          `json:"status"`
	NewsSrc   string             `json:"news_src"`
	Analyzer  string             `json:"analyzer"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}

func (q *Queries) GetJobByOwnerFilterByJIdAndStatus(ctx context.Context, arg *GetJobByOwnerFilterByJIdAndStatusParams) ([]*GetJobByOwnerFilterByJIdAndStatusRow, error) {
	rows, err := q.db.Query(ctx, getJobByOwnerFilterByJIdAndStatus,
		arg.Owner,
		arg.FJid,
		arg.TJid,
		arg.JStatus,
		arg.N,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetJobByOwnerFilterByJIdAndStatusRow
	for rows.Next() {
		var i GetJobByOwnerFilterByJIdAndStatusRow
		if err := rows.Scan(
			&i.ID,
			&i.Owner,
			&i.Status,
			&i.NewsSrc,
			&i.Analyzer,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const getJobByOwnerFilterByJIdRange = `-- name: GetJobByOwnerFilterByJIdRange :many
SELECT j.id, j.owner, j.status, asrc.name AS news_src, allm.name AS analyzer, j.created_at, j.updated_at
  FROM jobs AS j
 INNER JOIN apis AS asrc ON j.src_api_id = asrc.id
 INNER JOIN apis AS allm ON j.llm_api_id = allm.id 
 WHERE j.owner = $1
   AND j.id BETWEEN $2::int AND $3::int
   AND j.deleted_at IS NULL
 ORDER BY 
       j.id DESC
 LIMIT $4::int
`

type GetJobByOwnerFilterByJIdRangeParams struct {
	Owner uuid.UUID `json:"owner"`
	FJid  int32     `json:"f_jid"`
	TJid  int32     `json:"t_jid"`
	N     int32     `json:"n"`
}

type GetJobByOwnerFilterByJIdRangeRow struct {
	ID        int32              `json:"id"`
	Owner     uuid.UUID          `json:"owner"`
	Status    JobStatus          `json:"status"`
	NewsSrc   string             `json:"news_src"`
	Analyzer  string             `json:"analyzer"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}

func (q *Queries) GetJobByOwnerFilterByJIdRange(ctx context.Context, arg *GetJobByOwnerFilterByJIdRangeParams) ([]*GetJobByOwnerFilterByJIdRangeRow, error) {
	rows, err := q.db.Query(ctx, getJobByOwnerFilterByJIdRange,
		arg.Owner,
		arg.FJid,
		arg.TJid,
		arg.N,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetJobByOwnerFilterByJIdRangeRow
	for rows.Next() {
		var i GetJobByOwnerFilterByJIdRangeRow
		if err := rows.Scan(
			&i.ID,
			&i.Owner,
			&i.Status,
			&i.NewsSrc,
			&i.Analyzer,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const getJobByOwnerFilterByJIds = `-- name: GetJobByOwnerFilterByJIds :many
SELECT j.id, j.owner, j.status, asrc.name AS news_src, allm.name AS analyzer, j.created_at, j.updated_at
  FROM jobs AS j
 INNER JOIN apis AS asrc ON j.src_api_id = asrc.id
 INNER JOIN apis AS allm ON j.llm_api_id = allm.id 
 WHERE j.owner = $1
   AND j.id = ANY($2::int[])
   AND j.deleted_at IS NULL
 ORDER BY 
       j.id DESC
 LIMIT $3::int
`

type GetJobByOwnerFilterByJIdsParams struct {
	Owner uuid.UUID `json:"owner"`
	Ids   []int32   `json:"ids"`
	N     int32     `json:"n"`
}

type GetJobByOwnerFilterByJIdsRow struct {
	ID        int32              `json:"id"`
	Owner     uuid.UUID          `json:"owner"`
	Status    JobStatus          `json:"status"`
	NewsSrc   string             `json:"news_src"`
	Analyzer  string             `json:"analyzer"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}

func (q *Queries) GetJobByOwnerFilterByJIds(ctx context.Context, arg *GetJobByOwnerFilterByJIdsParams) ([]*GetJobByOwnerFilterByJIdsRow, error) {
	rows, err := q.db.Query(ctx, getJobByOwnerFilterByJIds, arg.Owner, arg.Ids, arg.N)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetJobByOwnerFilterByJIdsRow
	for rows.Next() {
		var i GetJobByOwnerFilterByJIdsRow
		if err := rows.Scan(
			&i.ID,
			&i.Owner,
			&i.Status,
			&i.NewsSrc,
			&i.Analyzer,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const getJobsByJobId = `-- name: GetJobsByJobId :one
SELECT j.id, j.owner, j.status, asrc.name AS news_src, j.src_query, 
       allm.name AS analyzer, j.llm_query, j.created_at, j.updated_at
  FROM jobs AS j
 INNER JOIN apis AS asrc ON j.src_api_id = asrc.id
 INNER JOIN apis AS allm ON j.llm_api_id = allm.id 
 WHERE j.owner = $1
   AND j.id = $2
   AND j.deleted_at IS NULL
`

type GetJobsByJobIdParams struct {
	Owner uuid.UUID `json:"owner"`
	ID    int32     `json:"id"`
}

type GetJobsByJobIdRow struct {
	ID        int32              `json:"id"`
	Owner     uuid.UUID          `json:"owner"`
	Status    JobStatus          `json:"status"`
	NewsSrc   string             `json:"news_src"`
	SrcQuery  string             `json:"src_query"`
	Analyzer  string             `json:"analyzer"`
	LlmQuery  []byte             `json:"llm_query"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}

func (q *Queries) GetJobsByJobId(ctx context.Context, arg *GetJobsByJobIdParams) (*GetJobsByJobIdRow, error) {
	row := q.db.QueryRow(ctx, getJobsByJobId, arg.Owner, arg.ID)
	var i GetJobsByJobIdRow
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Status,
		&i.NewsSrc,
		&i.SrcQuery,
		&i.Analyzer,
		&i.LlmQuery,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const getJobsByOwner = `-- name: GetJobsByOwner :many
SELECT j.id, j.owner, j.status, asrc.name AS news_src, allm.name AS analyzer, j.created_at, j.updated_at
  FROM jobs AS j
 INNER JOIN apis AS asrc ON j.src_api_id = asrc.id
 INNER JOIN apis AS allm ON j.llm_api_id = allm.id 
 WHERE j.owner = $1
   AND j.id < $2::int
   AND j.deleted_at IS NULL
 ORDER BY 
       j.id DESC
 LIMIT $3::int
`

type GetJobsByOwnerParams struct {
	Owner uuid.UUID `json:"owner"`
	Next  int32     `json:"next"`
	N     int32     `json:"n"`
}

type GetJobsByOwnerRow struct {
	ID        int32              `json:"id"`
	Owner     uuid.UUID          `json:"owner"`
	Status    JobStatus          `json:"status"`
	NewsSrc   string             `json:"news_src"`
	Analyzer  string             `json:"analyzer"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}

func (q *Queries) GetJobsByOwner(ctx context.Context, arg *GetJobsByOwnerParams) ([]*GetJobsByOwnerRow, error) {
	rows, err := q.db.Query(ctx, getJobsByOwner, arg.Owner, arg.Next, arg.N)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetJobsByOwnerRow
	for rows.Next() {
		var i GetJobsByOwnerRow
		if err := rows.Scan(
			&i.ID,
			&i.Owner,
			&i.Status,
			&i.NewsSrc,
			&i.Analyzer,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const getJobsByOwnerFilterByStatus = `-- name: GetJobsByOwnerFilterByStatus :many
SELECT j.id, j.owner, j.status, asrc.name AS news_src, allm.name AS analyzer, j.created_at, j.updated_at
  FROM jobs AS j
 INNER JOIN apis AS asrc ON j.src_api_id = asrc.id
 INNER JOIN apis AS allm ON j.llm_api_id = allm.id 
 WHERE j.owner = $1
   AND j.id < $2::int
   AND j.status = $3
   AND j.deleted_at IS NULL
 ORDER BY 
       j.id DESC
 LIMIT $4::int
`

type GetJobsByOwnerFilterByStatusParams struct {
	Owner   uuid.UUID `json:"owner"`
	Next    int32     `json:"next"`
	JStatus JobStatus `json:"j_status"`
	N       int32     `json:"n"`
}

type GetJobsByOwnerFilterByStatusRow struct {
	ID        int32              `json:"id"`
	Owner     uuid.UUID          `json:"owner"`
	Status    JobStatus          `json:"status"`
	NewsSrc   string             `json:"news_src"`
	Analyzer  string             `json:"analyzer"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}

func (q *Queries) GetJobsByOwnerFilterByStatus(ctx context.Context, arg *GetJobsByOwnerFilterByStatusParams) ([]*GetJobsByOwnerFilterByStatusRow, error) {
	rows, err := q.db.Query(ctx, getJobsByOwnerFilterByStatus,
		arg.Owner,
		arg.Next,
		arg.JStatus,
		arg.N,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetJobsByOwnerFilterByStatusRow
	for rows.Next() {
		var i GetJobsByOwnerFilterByStatusRow
		if err := rows.Scan(
			&i.ID,
			&i.Owner,
			&i.Status,
			&i.NewsSrc,
			&i.Analyzer,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const getLastJobId = `-- name: GetLastJobId :many
SELECT DISTINCT ON (status)
       id, status
  FROM jobs
 WHERE owner = $1
 ORDER BY 
        status ASC,
        id DESC
`

type GetLastJobIdRow struct {
	ID     int32     `json:"id"`
	Status JobStatus `json:"status"`
}

func (q *Queries) GetLastJobId(ctx context.Context, owner uuid.UUID) ([]*GetLastJobIdRow, error) {
	rows, err := q.db.Query(ctx, getLastJobId, owner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetLastJobIdRow
	for rows.Next() {
		var i GetLastJobIdRow
		if err := rows.Scan(&i.ID, &i.Status); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateJobStatus = `-- name: UpdateJobStatus :execrows
UPDATE jobs
   SET status = $1,
       updated_at = CURRENT_TIMESTAMP
 WHERE id = $2
   AND owner = $3
   AND deleted_at IS NULL
`

type UpdateJobStatusParams struct {
	Status JobStatus `json:"status"`
	ID     int32     `json:"id"`
	Owner  uuid.UUID `json:"owner"`
}

func (q *Queries) UpdateJobStatus(ctx context.Context, arg *UpdateJobStatusParams) (int64, error) {
	result, err := q.db.Exec(ctx, updateJobStatus, arg.Status, arg.ID, arg.Owner)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
