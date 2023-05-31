CREATE TYPE "role" AS ENUM (
  'user',
  'admin'
);

CREATE TYPE "job_status" AS ENUM (
  'created',
  'running',
  'done',
  'failure',
  'canceled'
);

CREATE TYPE "api_type" AS ENUM (
  'language_model',
  'source'
);

CREATE TYPE "event_type" AS ENUM (
  'sign-in',
  'sign-out',
  'authorization',
  'api-key',
  'query'
);