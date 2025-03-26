-- name: CreateClient :exec
INSERT INTO clients (client_name, client_token, role)
VALUES ($1, $2, $3);

-- name: GetClientByName :one
SELECT id, client_name, client_token, role, created_at
FROM clients
WHERE client_name = $1;

-- name: GetClientByToken :one
SELECT id, client_name, client_token, role, created_at
FROM clients
WHERE client_token = $1;

-- name: UpdateClientRoles :exec
UPDATE clients SET role = $2
WHERE client_token = $1;

-- name: AddClientRole :exec
UPDATE clients
SET role = array_append(role, sqlc.arg(roleToAdd)::text)
WHERE client_token = $1 AND NOT (sqlc.arg(roleToAdd) = ANY(role));


-- name: RemoveClientRole :exec
UPDATE clients
SET role = array_remove(role, sqlc.arg(roleToRemove)::text)
WHERE client_token = $1;


-- name: UpdateClientToken :exec
UPDATE clients SET client_token = $2
WHERE client_name = $1;


-- name: GetClientsByRole :many
SELECT id, client_name, client_token, role, created_at
FROM clients
WHERE $1 = ANY(role);


-- name: HasRole :one
SELECT EXISTS (
  SELECT 1 FROM clients
  WHERE client_token = $1
    AND $2 = ANY(role)
) AS has_role;


-- name: GetClientRoles :one
SELECT role
FROM clients
WHERE client_token = $1;


-- name: DeleteClient :exec
DELETE FROM clients WHERE client_token = $1;

-- name: ListClients :many
SELECT id, client_name, client_token, role, created_at
FROM clients
ORDER BY created_at DESC;


-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (client_token, token, expires_at)
VALUES ($1, $2, $3);

-- name: GetRefreshToken :one
SELECT id, client_token, token, expires_at, created_at
FROM refresh_tokens
WHERE token = $1;

-- name: IsRefreshTokenExpired :one
SELECT expires_at < now() AS is_expired
FROM refresh_tokens
WHERE token = $1;


-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens WHERE token = $1;

-- name: DeleteAllRefreshTokensForClient :exec
DELETE FROM refresh_tokens WHERE client_token = $1;

-- name: DeleteExpiredRefreshTokens :exec
DELETE FROM refresh_tokens WHERE expires_at < now();

-- name: CreateRevokedToken :exec
INSERT INTO revoked_tokens (client_token)
VALUES ($1);

-- name: IsRevokedToken :one
SELECT EXISTS (
    SELECT 1 FROM revoked_tokens WHERE client_token = $1
) AS revoked;

-- name: DeleteRevokedToken :exec
DELETE FROM revoked_tokens WHERE client_token = $1;

-- name: CreateRevokedRefreshToken :exec
INSERT INTO revoked_refresh_tokens (token)
VALUES ($1);

-- name: IsRevokedRefreshToken :one
SELECT EXISTS (
    SELECT 1 FROM revoked_refresh_tokens WHERE token = $1
) AS revoked;

-- name: DeleteRevokedRefreshToken :exec
DELETE FROM revoked_refresh_tokens WHERE token = $1;

-- name: InsertRevokedRefreshToken :exec
INSERT INTO revoked_refresh_tokens (token)
VALUES ($1);

-- name: DeleteRefreshTokenByToken :exec
DELETE FROM refresh_tokens
WHERE token = $1;

-- name: InsertNewRefreshToken :exec
INSERT INTO refresh_tokens (client_token, token, expires_at)
VALUES ($1, $2, $3);

-- name: InsertTokenLog :exec
INSERT INTO token_logs (
    client_token,
    token_type,
    action,
    token,
    ip_address,
    user_agent,
    metadata
) VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: GetLogsByClient :many
SELECT *
FROM token_logs
WHERE client_token = $1
ORDER BY created_at DESC;

-- name: GetLogsByTokenType :many
SELECT *
FROM token_logs
WHERE token_type = $1
ORDER BY created_at DESC;

-- name: GetLogsByAction :many
SELECT *
FROM token_logs
WHERE action = $1
ORDER BY created_at DESC;


