CREATE TABLE "jobs" (
  "id"              integer         PRIMARY KEY,
  "owner"           integer         NOT NULL,
  "status"          job_status      NOT NULL,
  "src_api_id"      smallint        NOT NULL,
  "src_query"       varchar(2048)   NOT NULL,
  "llm_api_id"      smallint        NOT NULL,
  "llm_query"       varchar(2048)   NOT NULL,
  "created_at"      timestamp       NOT NULL DEFAULT (now()),
  "updated_at"      timestamp       NOT NULL DEFAULT (now()),
  "deleted_at"      timestamp       DEFAULT null,
  "completed_at"    timestamp       DEFAULT null
);

CREATE INDEX ON "jobs" ("owner", "status");

ALTER TABLE "jobs" ADD FOREIGN KEY ("owner") REFERENCES "users" ("id");

ALTER TABLE "jobs" ADD FOREIGN KEY ("src_api_id") REFERENCES "apis" ("id");

ALTER TABLE "jobs" ADD FOREIGN KEY ("llm_api_id") REFERENCES "apis" ("id");