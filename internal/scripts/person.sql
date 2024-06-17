-- name: CreatePerson :exec
INSERT INTO people ( ID, first_name, last_name, email, phone, cell_phone, personable_id, personable_type, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);

-- name: UpdatePerson :exec
UPDATE people SET first_name = $2, last_name = $3, email = $4, phone = $5, cell_phone = $6, updated_at = $7 WHERE personable_id = $1;

-- name: GetPerson :one
SELECT *
FROM people
WHERE people.personable_id = $1;

-- name: DeletePerson :exec
DELETE FROM people
WHERE people.personable_id = $1;
