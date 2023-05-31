CREATE TABLE "logs" (
  "id"          BIGSERIAL    PRIMARY KEY,
  "user_id"     integer      NOT NULL,
  "type"        event_type   NOT NULL,
  "message"     varchar(256) NOT NULL,
  "created_at"  timestamptz  NOT NULL DEFAULT (now())
);

CREATE INDEX ON "logs" ("user_id", "type");

ALTER TABLE "logs" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE NO ACTION ON UPDATE CASCADE;

