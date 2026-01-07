-- name: Signup :one
INSERT INTO Users (name, password_hash)
VALUES ($1, $2)
RETURNING id, name;

-- name: Login :one
SELECT id, name, password_hash, COALESCE(totp, '') AS totp
FROM Users
WHERE name = $1;

-- name: EnableTotp :one
UPDATE Users
SET totp = $1
WHERE id = $2
RETURNING id, name;

-- name: DisableTotp :one
UPDATE Users
SET totp = NULL
WHERE id = $1
RETURNING id, name;

-- name: GetTotpSecret :one
SELECT id, name, COALESCE(totp, '') AS totp
FROM Users
WHERE id = $1;
