-- Claims due monitors and schedules their next check in one statement. SKIP LOCKED lets
-- several checker instances run at once without ever claiming the same monitor.
-- name: ClaimDueMonitors :many
WITH due AS (
  SELECT id
  FROM monitors
  WHERE enabled AND next_check_at <= now()
  ORDER BY next_check_at
  LIMIT sqlc.arg(batch_size)::int
  FOR UPDATE SKIP LOCKED
)
UPDATE monitors m
SET next_check_at = now() + make_interval(secs => m.interval_seconds)
FROM due
WHERE m.id = due.id
RETURNING m.id, m.url, m.method, m.timeout_ms, m.expected_status;

-- Joining the row to itself exposes its values from before the update, so the caller can
-- tell whether the status changed.
-- name: UpdateMonitorHealth :one
UPDATE monitors m
SET consecutive_failures = CASE WHEN sqlc.arg(ok)::bool THEN 0 ELSE m.consecutive_failures + 1 END,
    current_status = CASE
      WHEN sqlc.arg(ok)::bool THEN 'up'
      WHEN m.consecutive_failures + 1 >= sqlc.arg(failure_threshold)::int THEN 'down'
      ELSE m.current_status
    END
FROM monitors previous
WHERE m.id = sqlc.arg(id) AND previous.id = m.id
RETURNING previous.current_status AS previous_status, m.current_status;

-- name: InsertCheck :exec
INSERT INTO checks (monitor_id, checked_at, ok, status_code, latency_ms, error)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: OpenIncident :exec
INSERT INTO incidents (monitor_id, started_at, cause)
VALUES ($1, $2, $3)
ON CONFLICT (monitor_id) WHERE resolved_at IS NULL DO NOTHING;

-- name: ResolveIncident :exec
UPDATE incidents
SET resolved_at = $2
WHERE monitor_id = $1 AND resolved_at IS NULL;
