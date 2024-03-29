CREATE TABLE
    "jobs" (
        "id" BIGSERIAL PRIMARY KEY,
        "ulid" CHAR(26) UNIQUE NOT NULL,
        "owner" uuid NOT NULL,
        "status" job_status NOT NULL,
        "src_api_id" smallint NOT NULL,
        "src_query" text NOT NULL,
        "llm_api_id" smallint NOT NULL,
        "llm_query" json NOT NULL,
        "created_at" timestamptz NOT NULL DEFAULT (now()),
        "updated_at" timestamptz NOT NULL DEFAULT (now()),
        "deleted_at" timestamptz DEFAULT null
    );

CREATE INDEX ON "jobs" ("owner", "status");

CREATE INDEX ON "jobs" ("ulid");

ALTER TABLE "jobs"
ADD
    FOREIGN KEY ("owner") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "jobs"
ADD
    FOREIGN KEY ("src_api_id") REFERENCES "apis" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "jobs"
ADD
    FOREIGN KEY ("llm_api_id") REFERENCES "apis" ("id") ON DELETE CASCADE ON UPDATE CASCADE;