-- name: GetTransactionById :one
SELECT * FROM transactions
WHERE id = $1 LIMIT 1;

-- name: GetTransactionsByToAccount :many
SELECT * FROM transactions
WHERE to_account = $1;

-- name: GetTransactionsByFromAccount :many
SELECT * FROM transactions
WHERE from_account = $1;

-- name: GetAllTransactionsByAccount :many
SELECT * FROM transactions
WHERE to_account = $1 OR from_account = $1;

-- name: GetAllTransactions :many
SELECT * FROM transactions;

-- name: GetDepositTotal :one
SELECT COALESCE(SUM(amount),0) as total
FROM transactions WHERE to_account = $1;

-- name: GetWithdrawalTotal :one
SELECT COALESCE(SUM(amount),0) as total
FROM transactions WHERE from_account = $1;

-- name: CreateTransaction :one
INSERT INTO transactions(
    amount, from_account, to_account, created_at
) VALUES (
             $1, $2, $3, $4
         )
    RETURNING *;

-- name: DeleteTransaction :exec
DELETE FROM transactions
WHERE id = $1;


