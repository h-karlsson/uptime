-- name: CreateSession :exec
INSERT INTO sessions (token_hash, user_id, expires_at)
VALUES ($1, $2, $3);

-- name: GetSession :one
SELECT u.id, u.email, u.created_at, s.expires_at
FROM sessions s
JOIN users u ON u.id = s.user_id
WHERE s.token_hash = $1 AND s.expires_at > now();

-- name: ExtendSession :exec
UPDATE sessions
SET expires_at = $2
WHERE token_hash = $1;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE token_hash = $1;

-- name: DeleteUserSessions :exec
DELETE FROM sessions
WHERE user_id = $1;

-- name: DeleteExpiredSessions :execrows
DELETE FROM sessions
WHERE expires_at <= now();
