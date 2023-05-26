CREATE TABLE "apis" (
  "id"          SMALLSERIAL PRIMARY KEY,
  "name"        varchar(20) NOT NULL,
  "type"        api_type    NOT NULL,
  "created_at"  timestamp   NOT NULL DEFAULT (now()),
  "updated_at"  timestamp   NOT NULL DEFAULT (now()),
  "deleted_at"  timestamp   DEFAULT null
);