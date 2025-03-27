-- name: CreateSeat :one
INSERT INTO seat (room_id, number, has_socket, is_available)
VALUES ($1, $2, COALESCE($3, FALSE), COALESCE($4, TRUE))
RETURNING *;

--name: GetSeat :one
SELECT * FROM seat
WHERE id = $1 LIMIT 1;

--name: UpdateSeat :one
UPDATE seat SET
    number = COALESCE(sqlc.narg(number), number),
    has_socket = COALESCE(sqlc.narg(has_socket), has_socket),
    is_available = COALESCE(sqlc.narg(is_available), is_available)
WHERE id = $1
RETURNING *;

--name: DeleteSeat :exec
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

-- 获取自习室的座位列表
-- name: ListRoomSeat :many
SELECT s.* FROM seat s
WHERE s.room_id = $1
ORDER BY s.id;
