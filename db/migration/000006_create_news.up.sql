CREATE TABLE "news" (
  "id"          BIGSERIAL     PRIMARY KEY,
  "md5_hash"    char(128)     UNIQUE NOT NULL,
  "title"       text          NOT NULL,
  "url"         text          NOT NULL,
  "description" text          NOT NULL,
  "content"     text          NOT NULL,
  "source"      text,
  "publish_at"  timestamptz   NOT NULL,
  "created_at"  timestamptz   NOT NULL DEFAULT (now()),
  "updated_at"  timestamptz   NOT NULL DEFAULT (now())
);

CREATE INDEX ON "news" ("md5_hash");

CREATE INDEX ON "news" ("publish_at");
