-- 创建座位表
-- name: CreateSeat :one
INSERT INTO seat (room_id, number, has_socket, is_available)
VALUES ($1, $2, COALESCE($3, FALSE), COALESCE($4, TRUE))
RETURNING *;

-- 获取所有座位信息
-- name: GetSeat :one
SELECT * FROM seat
WHERE id = $1 LIMIT 1;

-- 动态查询座位，可能参数room_id, has_socket, is_available
-- name: ListRoomSeat :many
SELECT s.* FROM seat s
WHERE 
    (sqlc.narg(room_id)::INT IS NULL OR s.room_id = sqlc.narg(room_id)) AND
    (sqlc.narg(has_socket)::BOOLEAN IS NULL OR s.has_socket = sqlc.narg(has_socket)) AND
    (sqlc.narg(is_available)::BOOLEAN IS NULL OR s.is_available = sqlc.narg(is_available))
ORDER BY s.room_id, s.number
LIMIT $1 OFFSET $2;

-- 更新座位信息
-- name: UpdateSeat :one
UPDATE seat SET
    number = COALESCE(sqlc.narg(number), number),
    has_socket = COALESCE(sqlc.narg(has_socket), has_socket),
    is_available = COALESCE(sqlc.narg(is_available), is_available)
WHERE id = $1
RETURNING *;

-- 删除座位
-- name: DeleteSeat :exec
DELETE FROM seat
WHERE id = $1;

-- 批量更新座位
-- name: UpdateSeats :many
UPDATE seat SET
    has_socket = COALESCE(sqlc.narg(has_socket), has_socket),
    is_available = COALESCE(sqlc.narg(is_available), is_available)
WHERE id IN (SELECT * FROM unnest($1::int[]))
RETURNING *;

-- 获取座位详细信息（含自习室信息）
-- name: GetSeatWithRoom :one
SELECT s.*, r.name as room_name, r.open_time, r.close_time 
FROM seat s
JOIN room r ON s.room_id = r.id
WHERE s.id = $1;


