-- name: CreateClient :exec
INSERT INTO clients (name, api_key) VALUES ($1, $2);

-- name: GetClientByAPIKey :one
SELECT id, name, api_key, is_active, created_at
FROM clients
WHERE api_key = $1 AND is_active = TRUE
LIMIT 1;

-- name: GetClientByID :one
SELECT id, name, api_key, is_active, created_at
FROM clients
WHERE id = $1
LIMIT 1;
