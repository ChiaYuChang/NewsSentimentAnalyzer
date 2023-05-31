CREATE TABLE "newsjobs" (
  "id"      BIGSERIAL PRIMARY KEY,
  "job_id"  bigint,
  "news_id" bigint
);

CREATE INDEX ON "newsjobs" ("job_id", "news_id");

ALTER TABLE "newsjobs" ADD FOREIGN KEY ("job_id") REFERENCES "jobs" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "newsjobs" ADD FOREIGN KEY ("news_id") REFERENCES "news" ("id") ON DELETE CASCADE ON UPDATE CASCADE;
