-- name: CreateDeposit :exec
INSERT INTO deposits (ID, descricao, situacao, padrao, desconsiderarSaldo, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: GetDeposits :many
SELECT *
FROM deposits;

