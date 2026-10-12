-- name: ListMonitors :many
SELECT *
FROM monitors
WHERE user_id = $1
ORDER BY created_at, id;

-- name: GetMonitor :one
SELECT *
FROM monitors
WHERE id = $1 AND user_id = $2;

-- name: CountMonitors :one
SELECT count(*)
FROM monitors
WHERE user_id = $1;

-- name: CreateMonitor :one
INSERT INTO monitors (user_id, name, url, method, interval_seconds, timeout_ms, expected_status, enabled)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- An edit is checked right away instead of waiting out the old interval.
-- name: UpdateMonitor :one
UPDATE monitors
SET name = $3,
    url = $4,
    method = $5,
    interval_seconds = $6,
    timeout_ms = $7,
    expected_status = $8,
    enabled = $9,
    next_check_at = now()
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: DeleteMonitor :execrows
DELETE FROM monitors
WHERE id = $1 AND user_id = $2;
