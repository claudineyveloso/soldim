-- name: CreateAddress :exec
INSERT INTO addresses ( ID, public_place, complement, neighborhood, city, state, zip_code, addressable_id, addressable_type, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);

-- name: UpdateAddress :exec
UPDATE addresses SET public_place = $2, complement = $3, neighborhood = $4, city = $5, state = $6, zip_code = $7, updated_at = $8 WHERE addressable_id = $1;

-- name: GetAddress :one
SELECT *
FROM addresses
WHERE addresses.addressable_id = $1;

-- name: DeleteAddress :exec
DELETE FROM addresses
WHERE addresses.addressable_id = $1;
