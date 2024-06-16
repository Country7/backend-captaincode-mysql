-- name: CreateEntry :exec
INSERT INTO entries (account_id, amount)
VALUES (?, ?);

-- name: GetLastEntry :one
SELECT * FROM entries WHERE id = LAST_INSERT_ID();

-- name: GetEntry :one
SELECT * FROM entries
WHERE id = ? LIMIT 1;

-- name: ListEntries :many
SELECT * FROM entries
WHERE account_id = ? 
ORDER BY id LIMIT ? OFFSET ?;