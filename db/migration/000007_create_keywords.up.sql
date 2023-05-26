CREATE TABLE "keywords" (
  "id"      BIGSERIAL   PRIMARY KEY,
  "news_id" bigint      NOT NULL,
  "keyword" varchar(50) NOT NULL
);

CREATE INDEX ON "keywords" ("keyword");

ALTER TABLE "keywords" ADD FOREIGN KEY ("news_id") REFERENCES "news" ("id");

