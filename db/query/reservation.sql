-- 创建预约
-- 预约时间段内不能有其他预约，且座位必须可用
-- name: CreateReservation :one
INSERT INTO reservation (user_id, seat_id, start_time, end_time)
SELECT $1, $2, $3, $4
WHERE NOT EXISTS (
    SELECT 1 FROM reservation 
    WHERE seat_id = $2
    AND start_time < $4 
    AND end_time > $3
    AND status IN ('reserved', 'completed')
)
AND EXISTS (
    SELECT 1 FROM seat 
    WHERE id = $2 AND is_available = TRUE
)
RETURNING *;

-- name: GetReservation :one
SELECT * FROM reservation WHERE id = $1;

-- 获取预约信息时附带座位所在房间的代码，用于签到
-- name: GetReservationWithRoomCode :one
SELECT
    r.*,
    rm.code AS room_code
FROM reservation r
JOIN seat s ON r.seat_id = s.id
JOIN room rm ON s.room_id = rm.id
WHERE r.id = $1;

-- 动态查询, 可能参数start_time, end_time, limit, offset, user_id, seat_id, status
-- name: ListReservation :many
SELECT * FROM reservation 
WHERE
  (sqlc.narg(start_time)::TIMESTAMP IS NULL OR start_time >= sqlc.narg(start_time)) AND
  (sqlc.narg(end_time)::TIMESTAMP IS NULL OR end_time <= sqlc.narg(end_time)) AND
  (sqlc.narg(user_id)::INT IS NULL OR user_id = sqlc.narg(user_id)) AND
  (sqlc.narg(seat_id)::INT IS NULL OR seat_id = sqlc.narg(seat_id)) AND
  (sqlc.narg(status)::VARCHAR(20) IS NULL OR status = sqlc.narg(status))
ORDER BY
  CASE WHEN sqlc.narg(sort_by) = 'start_time' THEN start_time END DESC,
  created_at DESC
LIMIT $1 OFFSET $2;


-- 更新预约状态（含自动签到时间）
-- name: UpdateReservationStatus :one
UPDATE reservation SET
    status = $2,
    checkin_time = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;


-- 取消预约（只能取消未开始的预约）
-- name: DeleteReservation :execresult
UPDATE reservation
SET status = 'canceled'
WHERE id = $1 AND start_time > CURRENT_TIMESTAMP AND status = 'reserved';
