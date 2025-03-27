-- name: CreateRoom :one
INSERT INTO room (name, department, open_time, close_time, qr_code, is_active)
VALUES ($1, $2, $3, $4, $5, COALESCE($6, TRUE))
RETURNING *;

-- name: GetRoom :one
SELECT * FROM room WHERE id = $1 LIMIT 1;

-- name: ListRoom :many
SELECT * FROM room ORDER BY id LIMIT $1 OFFSET $2;

-- name: ListActiveRoom :many
SELECT * FROM room WHERE is_active = TRUE ORDER BY id LIMIT $1 OFFSET $2;

-- name: UpdateRoom :one
UPDATE room SET
    name = COALESCE(sqlc.narg(name), name),
    department = COALESCE(sqlc.narg(department), department),
    open_time = COALESCE(sqlc.narg(open_time), open_time),
    close_time = COALESCE(sqlc.narg(close_time), close_time),
    qr_code = COALESCE(sqlc.narg(qr_code), qr_code),
    is_active = COALESCE(sqlc.narg(is_active), is_active)
WHERE id = $1
RETURNING *;
