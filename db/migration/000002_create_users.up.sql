CREATE TABLE "users" (
  "id"                  SERIAL          PRIMARY KEY,
  "password"            bytea           NOT NULL,
  "first_name"          varchar(30)     NOT NULL,
  "last_name"           varchar(30)     NOT NULL,
  "role"                role            NOT NULL,
  "email"               varchar(320)    UNIQUE NOT NULL,
  "opt"                 varchar(128)    DEFAULT null,
  "created_at"          timestamptz     NOT NULL DEFAULT (now()),
  "updated_at"          timestamptz     NOT NULL DEFAULT (now()),
  "deleted_at"          timestamptz     DEFAULT null,
  "password_updated_at" timestamptz     NOT NULL DEFAULT (now())
);

CREATE INDEX ON "users" ("email");