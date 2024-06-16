-- name: CreateTransfer :exec
INSERT INTO transfers (from_account_id, to_account_id, amount)
VALUES (?, ?, ?);

-- name: GetLastTransfer :one
SELECT * FROM transfers WHERE id = LAST_INSERT_ID();

-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id = ? LIMIT 1;

-- name: ListTransfers :many
SELECT * FROM transfers
WHERE from_account_id = ? OR to_account_id = ?
ORDER BY id LIMIT ? OFFSET ?;
