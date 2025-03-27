-- 自习室实时使用统计（含插座统计）
-- name: GetRoomUtilization :many
SELECT 
    r.id,
    r.name,
    COUNT(s.id) FILTER (WHERE s.is_available) as available_seats,
    COUNT(s.id) FILTER (WHERE NOT s.is_available) as occupied_seats,
    COUNT(s.id) FILTER (WHERE s.has_socket) as socket_seats
FROM room r
LEFT JOIN seat s ON r.id = s.room_id
GROUP BY r.id, r.name;

-- 用户行为分析（含预约和违约统计）
-- name: GetUserBehaviorStats :one
SELECT 
    u.id,
    COUNT(r.id) FILTER (WHERE r.status = 'completed') as completed_count,
    COUNT(v.id) as violation_count,
    MAX(r.end_time) as last_reservation_time
FROM "user" u
LEFT JOIN reservation r ON u.id = r.user_id
LEFT JOIN violation v ON u.id = v.user_id
WHERE u.id = $1
GROUP BY u.id;

-- 参数化超时时间（分钟）
-- name: ExpireReservations :exec
WITH expired_reservations AS (
    SELECT id, user_id
    FROM reservation 
    WHERE status = 'reserved'
    AND start_time < (NOW() - ($1 || ' minutes')::INTERVAL)
    FOR UPDATE SKIP LOCKED
    LIMIT 1000
)
INSERT INTO violation (user_id, reservation_id, reason)
SELECT user_id, id, '超时未签到' 
FROM expired_reservations
RETURNING reservation_id;

UPDATE reservation r
SET status = 'violated'
FROM expired_reservations e
WHERE r.id = e.id;