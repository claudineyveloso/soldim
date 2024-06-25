-- name: CreateToken :exec
INSERT INTO tokens ( ID, access_token, expires_in, token_type, scope, refresh_token, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: UpdateToken :exec
UPDATE tokens SET access_token = $2, expires_in = $3, token_type = $4, scope = $5, refresh_token = $6, updated_at = $7 WHERE tokens.id = $1;

