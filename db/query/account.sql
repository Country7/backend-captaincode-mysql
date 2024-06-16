-- name: CreateAccount :exec
INSERT INTO accounts (owner, balance, currency)
VALUES (?, ?, ?);

-- name: GetLastAccount :one
SELECT * FROM accounts WHERE id = LAST_INSERT_ID();

-- name: GetAccount :one
SELECT *
FROM accounts
WHERE id = ? LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT *
FROM accounts
WHERE id = ? LIMIT 1
FOR UPDATE;

-- name: ListAccounts :many
SELECT *
FROM accounts
WHERE owner = ?
ORDER BY id LIMIT ?
OFFSET ?;

-- name: UpdateAccount :exec
UPDATE accounts
set balance = ?
WHERE id = ?;

-- name: AddAccountBalance :exec
UPDATE accounts
set balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id);

-- name: DeleteAccount :exec
DELETE
FROM accounts
WHERE id = ?;