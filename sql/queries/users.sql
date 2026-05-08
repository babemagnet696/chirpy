-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, hashed_password, email)
VALUES(
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: GetUser :one
SELECT *
FROM users
WHERE email = $1;

-- name: UpdateUser :one
UPDATE users
SET updated_at = NOW(), hashed_password = $1, email = $2
WHERE id = $3
RETURNING *;

-- name: UpdateChirpyRed :exec
UPDATE users
SET is_chirpy_red = true
WHERE id = $1;