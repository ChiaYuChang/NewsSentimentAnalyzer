CREATE TABLE "jobs" (
  "id"              integer         PRIMARY KEY,
  "owner"           integer         NOT NULL,
  "status"          job_status      NOT NULL,
  "src_api_id"      smallint        NOT NULL,
  "src_query"       text            NOT NULL,
  "llm_api_id"      smallint        NOT NULL,
  "llm_query"       text            NOT NULL,
  "created_at"      timestamptz     NOT NULL DEFAULT (now()),
  "updated_at"      timestamptz     NOT NULL DEFAULT (now()),
  "deleted_at"      timestamptz     DEFAULT null
);

CREATE INDEX ON "jobs" ("owner", "status");

ALTER TABLE "jobs" ADD FOREIGN KEY ("owner") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "jobs" ADD FOREIGN KEY ("src_api_id") REFERENCES "apis" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "jobs" ADD FOREIGN KEY ("llm_api_id") REFERENCES "apis" ("id") ON DELETE CASCADE ON UPDATE CASCADE;
