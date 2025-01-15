-- name: CreateUser :one
INSERT INTO users
(
    email,
    password,
    is_admin
) VALUES (
    ?, ?, 0
) RETURNING id,email ;

-- name: GetUserById :one
SELECT id, email FROM users WHERE id = ? LIMIT 1;

-- name: GetUserByEmail :one
SELECT id, email, password, is_admin FROM users WHERE email = ? LIMIT 1;

-- name: GetUsers :many
SELECT id, email, created_at FROM users;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?;