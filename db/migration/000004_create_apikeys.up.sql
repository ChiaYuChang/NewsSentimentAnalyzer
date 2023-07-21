CREATE TABLE "apikeys" (
  "id"          SERIAL      PRIMARY KEY,
  "owner"       uuid        NOT NULL,
  "api_id"      smallint    NOT NULL,
  "key"         text        NOT NULL,
  "created_at"  timestamptz NOT NULL DEFAULT (now()),
  "updated_at"  timestamptz NOT NULL DEFAULT (now()),
  "deleted_at"  timestamptz DEFAULT null
);

CREATE INDEX ON "apikeys" ("owner", "api_id");

ALTER TABLE "apikeys" ADD FOREIGN KEY ("owner") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "apikeys" ADD FOREIGN KEY ("api_id") REFERENCES "apis" ("id") ON DELETE CASCADE ON UPDATE CASCADE;
