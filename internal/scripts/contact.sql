-- name: CreateContact :exec
INSERT INTO contacts ( ID, nome, codigo, situacao, numeroDocumento, telefone, celular, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: GetContact :one
SELECT *
FROM contacts
WHERE contacts.id = $1;

-- name: GetContacts :many
SELECT *
FROM contacts;

-- name: GetUserByName:one
SELECT *
FROM contacts
WHERE contacts.nome = $1;

