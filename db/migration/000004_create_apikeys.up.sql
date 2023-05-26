CREATE TABLE "apikeys" (
  "id"          SERIAL      PRIMARY KEY,
  "owner"       integer     NOT NULL,
  "api_id"      smallint    NOT NULL,
  "key"         text        NOT NULL,
  "created_at"  timestamp   NOT NULL DEFAULT (now()),
  "updated_at"  timestamp   NOT NULL DEFAULT (now()),
  "deleted_at"  timestamp   DEFAULT null
);

ALTER TABLE "apikeys" ADD FOREIGN KEY ("owner") REFERENCES "users" ("id");

ALTER TABLE "apikeys" ADD FOREIGN KEY ("api_id") REFERENCES "apis" ("id");
