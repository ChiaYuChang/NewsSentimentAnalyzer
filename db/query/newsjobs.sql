-- name: CreateNewsJob :one
INSERT INTO newsjobs (
    job_id, news_id
) VALUES (
    $1, $2
)
RETURNING id;
