-- name: CreateUser :one
INSERT INTO users (
    name,
    email,
    password_hash
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: CreateSession :one
INSERT INTO sessions (
    user_id,
    token_hash,
    expires_at
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- GetValidSessionWithUser returns an unexpired session and the user who owns it.
-- The database clock decides expiry, so the app and database never disagree on "now".
-- Only safe user columns are selected; password_hash never leaves this query.
-- name: GetValidSessionWithUser :one
SELECT
    sqlc.embed(sessions),
    users.name,
    users.email
FROM sessions
JOIN users ON users.id = sessions.user_id
WHERE sessions.token_hash = sqlc.arg(token_hash)
  AND sessions.expires_at > now();

-- DeleteSession removes a session only when it belongs to the given user.
-- name: DeleteSession :exec
DELETE FROM sessions
WHERE token_hash = sqlc.arg(token_hash)
  AND user_id = sqlc.arg(user_id);
