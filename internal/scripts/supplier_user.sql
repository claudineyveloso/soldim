
-- name: CreateSuppliersUser :exec
INSERT INTO suppliers_users (supplier_id, user_id, created_at, updated_at)
VALUES ($1, $2, $3, $4);

-- name: GetSuppliersUsers :many
SELECT  supplier_id,
          user_id,
          created_at,
          updated_at
FROM suppliers_users;

-- name: UpdateSuppliersUsers :exec
UPDATE suppliers_users SET supplier_id = $2,
  updated_at = $3
WHERE suppliers_users.user_id = $1;


