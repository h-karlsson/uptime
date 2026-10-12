-- +goose Up
CREATE TABLE users (
  id            bigserial PRIMARY KEY,
  email         text UNIQUE NOT NULL,
  password_hash text NOT NULL,
  created_at    timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE monitors (
  id                   bigserial PRIMARY KEY,
  user_id              bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  name                 text NOT NULL,
  url                  text NOT NULL,
  method               text NOT NULL DEFAULT 'GET',
  interval_seconds     int  NOT NULL DEFAULT 60 CHECK (interval_seconds >= 30),
  timeout_ms           int  NOT NULL DEFAULT 10000 CHECK (timeout_ms > 0),
  expected_status      int  NOT NULL DEFAULT 200,
  enabled              boolean NOT NULL DEFAULT true,
  current_status       text NOT NULL DEFAULT 'unknown' CHECK (current_status IN ('up', 'down', 'unknown')),
  consecutive_failures int  NOT NULL DEFAULT 0,
  next_check_at        timestamptz NOT NULL DEFAULT now(),
  created_at           timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX monitors_user_id_idx ON monitors (user_id);
CREATE INDEX monitors_due_idx ON monitors (next_check_at) WHERE enabled;

-- Partitioned by month so retention is just dropping old partitions.
-- Monthly partitions are created by the checker (milestone 3).
CREATE TABLE checks (
  monitor_id  bigint NOT NULL REFERENCES monitors(id) ON DELETE CASCADE,
  checked_at  timestamptz NOT NULL,
  ok          boolean NOT NULL,
  status_code int,
  latency_ms  int,
  error       text,
  PRIMARY KEY (monitor_id, checked_at)
) PARTITION BY RANGE (checked_at);

CREATE TABLE incidents (
  id          bigserial PRIMARY KEY,
  monitor_id  bigint NOT NULL REFERENCES monitors(id) ON DELETE CASCADE,
  started_at  timestamptz NOT NULL,
  resolved_at timestamptz,
  cause       text
);
CREATE INDEX incidents_monitor_id_started_at_idx ON incidents (monitor_id, started_at DESC);
-- At most one open incident per monitor.
CREATE UNIQUE INDEX incidents_one_open_per_monitor_idx ON incidents (monitor_id) WHERE resolved_at IS NULL;

CREATE TABLE notification_channels (
  id         bigserial PRIMARY KEY,
  user_id    bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  type       text NOT NULL CHECK (type IN ('webhook', 'discord', 'email')),
  config     jsonb NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX notification_channels_user_id_idx ON notification_channels (user_id);

CREATE TABLE status_pages (
  id         bigserial PRIMARY KEY,
  user_id    bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  slug       text UNIQUE NOT NULL,
  title      text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE status_page_monitors (
  status_page_id bigint NOT NULL REFERENCES status_pages(id) ON DELETE CASCADE,
  monitor_id     bigint NOT NULL REFERENCES monitors(id) ON DELETE CASCADE,
  PRIMARY KEY (status_page_id, monitor_id)
);

-- +goose Down
DROP TABLE status_page_monitors;
DROP TABLE status_pages;
DROP TABLE notification_channels;
DROP TABLE incidents;
DROP TABLE checks;
DROP TABLE monitors;
DROP TABLE users;
