-- +goose Up
-- The cookie holds a random token; only its SHA-256 hash is stored, so a leaked
-- database can't be used to hijack sessions.
CREATE TABLE sessions (
  token_hash bytea PRIMARY KEY,
  user_id    bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at timestamptz NOT NULL DEFAULT now(),
  expires_at timestamptz NOT NULL
);
CREATE INDEX sessions_user_id_idx ON sessions (user_id);
CREATE INDEX sessions_expires_at_idx ON sessions (expires_at);

-- +goose Down
DROP TABLE sessions;
