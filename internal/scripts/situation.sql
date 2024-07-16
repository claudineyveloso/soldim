-- name: CreateSituation :exec
INSERT INTO situations (id, descricao, created_at, updated_at)
VALUES ($1, $2, $3, $4);

-- name: GetSituations :many
SELECT id,
        descricao,
        created_at,
        updated_at
FROM situations;

-- name: GetSituation :one
SELECT id,
        descricao,
        created_at,
        updated_at
FROM situations
WHERE situations.id = $1;

-- name: GetSituationByDescroption :one
SELECT id,
        descricao,
        created_at,
        updated_at
FROM situations
WHERE situations.descricao = $1;


