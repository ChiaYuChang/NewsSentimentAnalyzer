-- name: CreateEmbedding :one
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
    ) RETURNING id;

-- name: GetEmbeddingByNewsIdsAndModel :many
SELECT
    id,
    news_id,
    model,
    embedding,
    sentiment
FROM embeddings
WHERE
    news_id = ANY(@news_ids:: int [])
    AND model = $1
    AND deleted_at IS NULL;

-- name: GetEmbeddingByJobId :many
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
    AND e.deleted_at IS NULL;