-- 创建关联预约的违约记录
-- name: CreateViolationWithCheck :one
INSERT INTO violation (user_id, reservation_id, reason)
SELECT $1, $2, $3
FROM reservation 
WHERE id = $2
AND user_id = $1
RETURNING *;

-- 动态查询，可能参数reservation_id, user_id
-- name: ListViolation :many
SELECT v.* FROM violation v
WHERE
    (sqlc.narg(reservation_id)::UUID IS NULL OR v.reservation_id = sqlc.narg(reservation_id)) AND
    (sqlc.narg(user_id)::INT IS NULL OR v.user_id = sqlc.narg(user_id))
ORDER BY v.created_at DESC
LIMIT $1 OFFSET $2;

-- 更新违约记录
-- name: UpdateViolation :one
UPDATE violation SET
    reason = $2
WHERE id = $1
RETURNING *;

-- 获取违约记录
-- name: GetViolation :one
SELECT * FROM violation 
WHERE id = $1;
