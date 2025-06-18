-- 创建自习室，不考虑二维码/签到码
-- name: CreateRoom :one
INSERT INTO room (name, department, code, qr_code, open_time, close_time)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- 获取自习室信息
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

-- 更新自习室信息
-- 签到码和二维码自动生成，无需更新
-- name: UpdateRoom :one
UPDATE room SET
    name = COALESCE(sqlc.narg(name), name),
    department = COALESCE(sqlc.narg(department), department),
    open_time = COALESCE(sqlc.narg(open_time), open_time),
    close_time = COALESCE(sqlc.narg(close_time), close_time),
    qr_code = COALESCE(sqlc.narg(qr_code), qr_code),
    code = COALESCE(sqlc.narg(code), code),
    is_active = COALESCE(sqlc.narg(is_active), is_active)
WHERE id = $1
RETURNING *;

-- 删除自习室
-- name: DeleteRoom :exec
DELETE FROM room
WHERE id = $1;

-- 获取自习室数量
-- 可选参数：is_active
-- name: CountRoom :one
SELECT COUNT(*) FROM room
WHERE
  (sqlc.narg(is_active)::BOOLEAN IS NULL OR is_active = sqlc.narg(is_active));


-- 更新所有自习室签到码
-- name: UpdateAllRoomCode :exec
UPDATE room SET code = unnest(@codes::text[]);