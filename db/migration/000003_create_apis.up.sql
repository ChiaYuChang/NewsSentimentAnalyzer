CREATE TABLE "apis" (
  "id"          SMALLSERIAL PRIMARY KEY,
  "name"        varchar(20) NOT NULL,
  "type"        api_type    NOT NULL,
  "created_at"  timestamptz NOT NULL DEFAULT (now()),
  "updated_at"  timestamptz NOT NULL DEFAULT (now()),
  "deleted_at"  timestamptz DEFAULT null
);