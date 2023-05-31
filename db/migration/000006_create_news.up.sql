CREATE TABLE "news" (
  "id"          BIGSERIAL     PRIMARY KEY,
  "md5_hash"    char(128)     UNIQUE NOT NULL,
  "title"       text          NOT NULL,
  "url"         varchar(256)  NOT NULL,
  "description" text          NOT NULL,
  "content"     text          NOT NULL,
  "source"      varchar(256),
  "publish_at"  timestamptz   NOT NULL,
  "created_at"  timestamptz   NOT NULL DEFAULT (now()),
  "updated_at"  timestamptz   NOT NULL DEFAULT (now()),
  "deleted_at"  timestamptz   DEFAULT null
);

CREATE INDEX ON "news" ("md5_hash");

CREATE INDEX ON "news" ("publish_at");
