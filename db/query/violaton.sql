-- 创建关联预约的违约记录
-- name: CreateViolationWithCheck :one
INSERT INTO violation (user_id, reservation_id, reason)
SELECT $1, $2, $3
FROM reservation 
WHERE id = $2
AND user_id = $1
RETURNING *;

-- 获取用户违约详情（含预约信息）
-- name: GetUserViolations :many
SELECT v.*, r.start_time, r.end_time, s.number as seat_number
FROM violation v
JOIN reservation r ON v.reservation_id = r.id
JOIN seat s ON r.seat_id = s.id
WHERE v.user_id = $1;