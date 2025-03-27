-- name: CreateUser :one
INSERT INTO "user" (username, password, role, department, email)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM "user"
WHERE username = $1 LIMIT 1;

-- name: GetUser :one
SELECT * FROM "user"
WHERE id = $1 LIMIT 1;

-- name: ListUser :many
SELECT * FROM "user"
ORDER BY id LIMIT $1 OFFSET $2;

-- name: UpdateUser :one
UPDATE "user" SET
    username = COALESCE(sqlc.narg(username), username),
    password = COALESCE(sqlc.narg(password), password),
    role = COALESCE(sqlc.narg(role), role),
    department = COALESCE(sqlc.narg(department), department),
    email = COALESCE(sqlc.narg(email), email)
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM "user"
WHERE id = $1;