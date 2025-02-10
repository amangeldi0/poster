-- name: GetUsers :many
SELECT * FROM users;

-- name: CreateUser :one
INSERT INTO users (
    id, username, email, password_hash, created_at, updated_at, verify_code
) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByUUID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserVerifyStatus :one
SELECT is_verified FROM users WHERE email = $1;

-- name: UpdateUserRefreshToken :exec
UPDATE users SET refresh_token = $1 WHERE id = $2;

-- name: DeleteUserByEmail :exec
DELETE FROM users WHERE email = $1;
