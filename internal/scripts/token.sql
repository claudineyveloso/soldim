-- name: CreateToken :exec
INSERT INTO tokens ( ID, access_token, expires_in, token_type, scope, refresh_token)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetToken :many
SELECT *
FROM tokens;

-- name: UpdateToken :exec
UPDATE tokens SET access_token = $2, 
  expires_in = $3, 
  token_type = $4, 
  scope = $5, 
  refresh_token = $6 
WHERE tokens.id = $1;

-- name: DeleteToken :exec
DELETE FROM 
tokens WHERE id = $1;
