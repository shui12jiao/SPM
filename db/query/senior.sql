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
-- name: ExpireReservations :many
WITH expired_data AS (
    SELECT 
        r.id AS reservation_id,
        u.id AS user_id,
        u.email
    FROM reservation r
    JOIN "user" u ON r.user_id = u.id
    WHERE r.status = 'reserved'
    AND r.start_time < (NOW() - ($1 || ' minutes')::INTERVAL)
    FOR UPDATE SKIP LOCKED
    LIMIT 1000
),
insert_violation AS (
    INSERT INTO violation (user_id, reservation_id, reason)
    SELECT user_id, reservation_id, '超时未签到'
    FROM expired_data
    RETURNING reservation_id
)
UPDATE reservation r
SET status = 'violated'
FROM expired_data e
WHERE r.id = e.reservation_id
RETURNING e.reservation_id, e.user_id, e.email;