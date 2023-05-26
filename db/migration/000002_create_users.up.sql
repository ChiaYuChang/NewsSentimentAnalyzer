CREATE TABLE "users" (
  "id"                  SERIAL          PRIMARY KEY,
  "password"            bytea           NOT NULL,
  "first_name"          varchar(30)     NOT NULL,
  "last_name"           varchar(30)     NOT NULL,
  "role"                role            NOT NULL,
  "email"               varchar(320)    NOT NULL,
  "created_at"          timestamp       NOT NULL DEFAULT (now()),
  "updated_at"          timestamp       NOT NULL DEFAULT (now()),
  "deleted_at"          timestamp       DEFAULT null,
  "password_updated_at" timestamp       NOT NULL DEFAULT (now())
);