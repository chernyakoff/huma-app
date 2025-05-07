-- name: CreateUser :one
INSERT INTO users
(
    id,
    email,
    password
) VALUES (
    ?, ?, ?
) RETURNING id,email, role ;

-- name: GetUserById :one
SELECT id, email, role FROM users WHERE id = ? LIMIT 1;

-- name: GetUserByEmail :one
SELECT id, email, password, role FROM users WHERE email = ? LIMIT 1;

-- name: GetUsers :many
SELECT id, email, role, created_at FROM users;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?;

-- name: VerifyUser :exec
UPDATE users SET verified = 1 WHERE id = ?;