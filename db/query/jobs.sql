-- name: GetJobsByOwner :many
SELECT id, owner, status, src_api_id, src_query, llm_api_id, llm_query, created_at, updated_at
  FROM jobs
 WHERE owner = $1
   AND deleted_at IS NULL
 ORDER BY 
       updated_at DESC,
       status     DESC
 LIMIT @n::int;

-- name: CreateJob :exec
INSERT INTO jobs (
  owner, status, src_api_id, src_query, llm_api_id, llm_query
) VALUES (
    $1, $2, $3, $4, $5, $6
);

-- name: UpdateJobStatus :exec
UPDATE jobs
   SET status = $1,
       updated_at = CURRENT_TIMESTAMP
 WHERE id = $2
   AND owner = $3
   AND deleted_at IS NULL;

-- name: DeleteJob :exec
UPDATE jobs
   SET deleted_at = CURRENT_TIMESTAMP
 WHERE id = $1
   AND owner = $2;

-- name: CleanUpJobs :exec
DELETE FROM jobs
 WHERE deleted_at IS NOT NULL;