-- 插入用户数据到"user"表
-- name: CreateUser :one
INSERT INTO "user" (username, password, role, department, email)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- 根据用户名获取用户信息
-- name: GetUserByUsername :one
SELECT * FROM "user"
WHERE username = $1 LIMIT 1;

-- 获取用户信息（通过ID）
-- name: GetUser :one
SELECT * FROM "user"
WHERE id = $1 LIMIT 1;

-- 动态查询，可能参数role, department
-- name: ListUser :many
SELECT * FROM "user"
WHERE
    (sqlc.narg(role)::VARCHAR(10) IS NULL OR role = sqlc.narg(role)) AND
    (sqlc.narg(department)::VARCHAR(50) IS NULL OR department = sqlc.narg(department))
ORDER BY id LIMIT $1 OFFSET $2;

-- 更新用户信息
-- name: UpdateUser :one
UPDATE "user" SET
    username = COALESCE(sqlc.narg(username), username),
    password = COALESCE(sqlc.narg(password), password),
    role = COALESCE(sqlc.narg(role), role),
    department = COALESCE(sqlc.narg(department), department),
    email = COALESCE(sqlc.narg(email), email)
WHERE id = $1
RETURNING *;

-- 删除用户
-- name: DeleteUser :exec
DELETE FROM "user"
WHERE id = $1;