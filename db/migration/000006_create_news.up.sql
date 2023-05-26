CREATE TABLE "news" (
  "id"          BIGSERIAL       PRIMARY KEY,
  "job_id"      integer         NOT NULL,
  "md5_hash"    char(128)       NOT NULL,
  "title"       text            NOT NULL,
  "url"         varchar(100)    NOT NULL,
  "description" text            NOT NULL,
  "content"     text            NOT NULL,
  "source"      varchar(256)    NOT NULL,     
  "publish_at"  timestamp       NOT NULL,
  "created_at"  timestamp       DEFAULT (now()),
  "updated_at"  timestamp       DEFAULT (now()),
  "deleted_at"  timestamp       DEFAULT null
);

ALTER TABLE "news" ADD FOREIGN KEY ("job_id") REFERENCES "jobs" ("id");
