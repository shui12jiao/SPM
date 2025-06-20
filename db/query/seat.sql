-- 创建座位表
-- name: CreateSeat :one
INSERT INTO seat (room_id, number, has_socket, is_available)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- 创建座位表（批量插入）
-- name: CreateSeats :many
INSERT INTO seat (room_id, number, has_socket, is_available)
SELECT 
    @room_id, 
    data.number, 
    data.has_socket, 
    data.is_available
FROM (
    SELECT 
        unnest(@numbers::text[]) AS number,
        unnest(@has_sockets::boolean[]) AS has_socket,
        unnest(@is_availables::boolean[]) AS is_available
) AS data
RETURNING *;

-- 获取所有座位信息
-- name: GetSeat :one
SELECT * FROM seat
WHERE id = $1 LIMIT 1;

-- 动态查询座位，可能参数room_id, has_socket, is_available
-- name: ListSeat :many
SELECT s.* FROM seat s
WHERE 
    (sqlc.narg(room_id)::INT IS NULL OR s.room_id = sqlc.narg(room_id)) AND
    (sqlc.narg(has_socket)::BOOLEAN IS NULL OR s.has_socket = sqlc.narg(has_socket)) AND
    (sqlc.narg(is_available)::BOOLEAN IS NULL OR s.is_available = sqlc.narg(is_available))
ORDER BY s.room_id, s.id
LIMIT $1 OFFSET $2;

-- 更新座位信息
-- 主要用于用户预约座位时，更新座位的可用性
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

-- 批量更新座位, 不允许更改room_id
-- name: UpdateSeats :many
UPDATE seat
SET
    number = data.number,
    has_socket = data.has_socket,
    is_available = data.is_available
FROM (
    SELECT unnest(@ids::int[]) AS id,
           unnest(@numbers::text[]) AS number,
           unnest(@has_sockets::boolean[]) AS has_socket,
           unnest(@is_availables::boolean[]) AS is_available
) AS data
WHERE seat.id = data.id
RETURNING seat.*;


-- 获取座位详细信息（含自习室信息）
-- name: GetSeatWithRoom :one
SELECT s.*, r.name as room_name, r.open_time, r.close_time 
FROM seat s
JOIN room r ON s.room_id = r.id
WHERE s.id = $1;


