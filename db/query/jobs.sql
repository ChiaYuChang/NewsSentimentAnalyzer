-- name: GetJobsByOwner :many
SELECT j.id, j.owner, j.status, asrc.name AS news_src, allm.name AS analyzer, j.created_at, j.updated_at
  FROM jobs AS j
 INNER JOIN apis AS asrc ON j.src_api_id = asrc.id
 INNER JOIN apis AS allm ON j.llm_api_id = allm.id 
 WHERE j.owner = $1
   AND j.id > @next::int
   AND j.deleted_at IS NULL
 ORDER BY 
       j.id DESC
 LIMIT @n::int;

-- name: GetJobsByJobId :one
SELECT j.id, j.owner, j.status, asrc.name AS news_src, j.src_query, 
       allm.name AS analyzer, j.llm_query, j.created_at, j.updated_at
  FROM jobs AS j
 INNER JOIN apis AS asrc ON j.src_api_id = asrc.id
 INNER JOIN apis AS allm ON j.llm_api_id = allm.id 
 WHERE j.owner = $1
   AND j.id = $2
   AND j.deleted_at IS NULL;

-- name: CreateJob :one
INSERT INTO jobs (
  owner, status, src_api_id, src_query, llm_api_id, llm_query
) VALUES (
    $1, $2, $3, $4, $5, $6
) 
RETURNING id;

-- name: UpdateJobStatus :execrows
UPDATE jobs
   SET status = $1,
       updated_at = CURRENT_TIMESTAMP
 WHERE id = $2
   AND owner = $3
   AND deleted_at IS NULL;

-- name: DeleteJob :execrows
UPDATE jobs
   SET deleted_at = CURRENT_TIMESTAMP
 WHERE id = $1
   AND owner = $2;

-- name: CleanUpJobs :execrows
DELETE FROM jobs
 WHERE deleted_at IS NOT NULL;