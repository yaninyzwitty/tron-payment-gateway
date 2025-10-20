-- name: CreateAccount :exec
INSERT INTO accounts (client_id, name) VALUES ($1, $2);

-- name: GetAccountsByClientID :many
SELECT id, client_id, name, created_at
FROM accounts
WHERE client_id = $1;

-- name: GetAccountByIDAndClientID :one
SELECT id, client_id, name, address_index, created_at
FROM accounts
WHERE id = $1 AND client_id = $2;