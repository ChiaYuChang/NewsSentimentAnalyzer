CREATE TABLE "news" (
  "id"           BIGSERIAL     PRIMARY KEY,
  "md5_hash"     char(128)     UNIQUE NOT NULL,
  "guid"         varchar       NOT NULL,
  "author"       text[],
  "title"        text          NOT NULL,
  "link"         text          NOT NULL,
  "description"  text          NOT NULL,
  "language"     varchar,
  "content"      text[]        NOT NULL,
  "category"     varchar       NOT NULL,
  "source"       text          NOT NULL,
  "related_guid" varchar[],
  "publish_at"   timestamptz   NOT NULL,
  "created_at"   timestamptz   NOT NULL DEFAULT (now())
);

CREATE INDEX ON "news" ("md5_hash");

CREATE INDEX ON "news" ("guid");

CREATE INDEX ON "news" ("publish_at");
