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
SET totp = @totp::text
WHERE id = $1
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

-- name: InsertRecoveryCodes :exec
INSERT INTO user_recovery_codes (user_id, code_hash)
SELECT $1, unnest($2::text[]);

-- name: ClearRecoveryCodes :exec
DELETE FROM User_recovery_codes
WHERE user_id = $1;

-- name: GetUnusedRecoveryCodes :many
SELECT id, code_hash
FROM User_recovery_codes
WHERE used_at IS NULL AND user_id = $1;

-- name: SetRecoveryCodeAsUsed :exec
UPDATE user_recovery_codes
SET used_at = NOW()
WHERE id = $1;
