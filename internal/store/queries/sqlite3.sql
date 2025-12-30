-- name: CreateUser :one
INSERT INTO users (
    id, name, email, password_hash,
    is_paused, check_in_interval, trigger_interval_num, buffer_period, verifier_quorum,
    last_check_in, current_status
) VALUES (
    ?, ?, ?, ?,
    ?, ?, ?, ?, ?,
    ?, ?
) RETURNING *;

-- name: CreateContactMethod :one
INSERT INTO contact_methods (
    id, user_id, beneficiary_id, channel, target, metadata, created_at
) VALUES (
    ?, ?, ?, ?, ?, ?, ?
) RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = ? LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = ? LIMIT 1;

-- name: ListContactMethodsByUserID :many
SELECT * FROM contact_methods
WHERE user_id = ?;

-- name: UpdateUserCheckIn :exec
UPDATE users
SET last_check_in = ?, current_status = 'ALIVE'
WHERE id = ?;