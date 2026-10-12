-- +goose Up
CREATE UNIQUE INDEX monitors_user_id_url_idx ON monitors (user_id, url);

-- +goose Down
DROP INDEX monitors_user_id_url_idx;
