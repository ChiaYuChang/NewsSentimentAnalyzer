CREATE EXTENSION IF NOT EXISTS "vector";

CREATE TYPE "sentiment" AS ENUM (
    'positive',
    'neutral',
    'negative'
);

CREATE TABLE
    embeddings (
        id bigserial PRIMARY KEY,
        model varchar(32) NOT NULL,
        news_id bigint NOT NULL,
        embedding vector(1536),
        sentiment sentiment NOT NULL,
        created_at timestamptz NOT NULL DEFAULT (now()),
        updated_at timestamptz NOT NULL DEFAULT (now()),
        deleted_at timestamptz DEFAULT null
    );

ALTER TABLE embeddings
ADD
    FOREIGN KEY (news_id) REFERENCES news (id) ON DELETE CASCADE ON UPDATE CASCADE;

CREATE INDEX ON embeddings USING hnsw (embedding vector_ip_ops);

CREATE INDEX ON embeddings (model, sentiment);