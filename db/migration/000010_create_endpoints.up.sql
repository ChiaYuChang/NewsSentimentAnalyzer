CREATE TABLE "endpoints" (
    "id"            SERIAL          PRIMARY KEY,
    "name"          varchar(32)     NOT NULL,
    "api_id"        smallint        NOT NULL,
    "template_name" varchar(32)     NOT NULL,
    "created_at"    timestamptz     NOT NULL DEFAULT (now()),
    "updated_at"    timestamptz     NOT NULL DEFAULT (now()),
    "deleted_at"    timestamptz     DEFAULT null
);

ALTER TABLE "endpoints" ADD FOREIGN KEY ("api_id") REFERENCES "apis" ("id") ON DELETE CASCADE ON UPDATE CASCADE;
