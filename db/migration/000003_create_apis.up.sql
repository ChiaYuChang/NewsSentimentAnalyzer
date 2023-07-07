CREATE TABLE "apis" (
  "id"           SMALLSERIAL  PRIMARY KEY,
  "name"         varchar(20)  NOT NULL,
  "type"         api_type     NOT NULL,
  "image"        varchar(128) NOT NULL DEFAULT 'logo_Default.svg',
  "icon"         varchar(128) NOT NULL DEFAULT 'favicon_Default.svg',
  "document_url" varchar(128) NOT NULL DEFAULT '#',
  "created_at"   timestamptz  NOT NULL DEFAULT (now()),
  "updated_at"   timestamptz  NOT NULL DEFAULT (now()),
  "deleted_at"   timestamptz  DEFAULT null
);