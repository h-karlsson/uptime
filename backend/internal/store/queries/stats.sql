-- name: MonitorBelongsToUser :one
SELECT EXISTS (
  SELECT 1 FROM monitors WHERE id = sqlc.arg(id) AND user_id = sqlc.arg(user_id)
);

-- Latency statistics only count successful checks, so timeouts don't skew them. sqlc can't
-- type aggregates as nullable, so latencies default to 0 and callers rely on ok_checks.
-- name: GetCheckSummary :one
SELECT count(*)::int AS checks,
       (count(*) FILTER (WHERE ok))::int AS ok_checks,
       coalesce(avg(latency_ms) FILTER (WHERE ok), 0)::float8 AS avg_latency_ms,
       coalesce(percentile_cont(0.95) WITHIN GROUP (ORDER BY latency_ms) FILTER (WHERE ok), 0)::float8 AS p95_latency_ms
FROM checks
WHERE monitor_id = sqlc.arg(monitor_id) AND checked_at >= sqlc.arg(since);

-- name: GetLastCheckedAt :one
SELECT checked_at
FROM checks
WHERE monitor_id = $1
ORDER BY checked_at DESC
LIMIT 1;

-- name: ListLatencyBuckets :many
SELECT date_bin(make_interval(secs => sqlc.arg(bucket_seconds)::int), checked_at, '2000-01-01 00:00:00+00')::timestamptz AS bucket_start,
       count(*)::int AS checks,
       (count(*) FILTER (WHERE ok))::int AS ok_checks,
       coalesce(avg(latency_ms) FILTER (WHERE ok), 0)::float8 AS avg_latency_ms
FROM checks
WHERE monitor_id = sqlc.arg(monitor_id) AND checked_at >= sqlc.arg(since)
GROUP BY bucket_start
ORDER BY bucket_start;

-- name: ListDailyUptime :many
SELECT (checked_at AT TIME ZONE 'UTC')::date AS day,
       count(*)::int AS checks,
       (count(*) FILTER (WHERE ok))::int AS ok_checks
FROM checks
WHERE monitor_id = sqlc.arg(monitor_id) AND checked_at >= sqlc.arg(since)
GROUP BY day
ORDER BY day;

-- name: ListIncidents :many
SELECT id, started_at, resolved_at, cause
FROM incidents
WHERE monitor_id = sqlc.arg(monitor_id)
ORDER BY started_at DESC
LIMIT 20;
