-- 创建带时间冲突检测的预约
-- name: CreateReservation :one
INSERT INTO reservation (user_id, seat_id, start_time, end_time)
SELECT $1, $2, $3, $4
WHERE NOT EXISTS (
    SELECT 1 FROM reservation 
    WHERE seat_id = $2
    AND start_time < $4 
    AND end_time > $3
    AND status NOT IN ('canceled', 'completed', 'violated')    
)
RETURNING *;

-- 分页查询用户当前有效预约
-- name: ListUserActiveReservation :many
SELECT r.*, s.number as seat_number, rm.name as room_name
FROM reservation r
JOIN seat s ON r.seat_id = s.id
JOIN room rm ON s.room_id = rm.id
WHERE r.user_id = $1
AND r.status IN ('reserved', 'checked_in')
ORDER BY start_time DESC
LIMIT $2 OFFSET $3;

-- 更新预约状态（含自动签到时间）
-- name: UpdateReservationStatus :one
UPDATE reservation SET
    status = COALESCE($2, status),
    checkin_time = CASE 
        WHEN $2 = 'checked_in' THEN CURRENT_TIMESTAMP
        ELSE checkin_time 
    END
WHERE id = $1
RETURNING *;

