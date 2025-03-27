-- name: CreateRoom :one
INSERT INTO room (name, department, open_time, close_time, qr_code, is_active)
VALUES ($1, $2, $3, $4, $5, COALESCE($6, TRUE))
RETURNING *;

-- name: GetRoom :one
SELECT * FROM room WHERE id = $1 LIMIT 1;

-- department,is_active作为可能查询条件
-- name: ListRoom :many
SELECT * FROM room 
WHERE
  (sqlc.narg(department)::VARCHAR(50) IS NULL OR department = sqlc.narg(department)) AND
  (sqlc.narg(is_active)::BOOLEAN IS NULL OR is_active = sqlc.narg(is_active))
ORDER BY id DESC
LIMIT $1 OFFSET $2;


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
