
-- name: ListPersons :many
SELECT * FROM "person";

-- name: GetPerson :one
SELECT * FROM "person"
WHERE "id" = $1 LIMIT 1;

-- name: GetPersonByEmail :one
SELECT * FROM "person"
WHERE "email" = $1 LIMIT 1;

-- name: GetPersonByPersonName :one
SELECT * FROM "person"
WHERE "personname" = $1 LIMIT 1;

-- name: GetPersonByPersonNameAndPassword :one
SELECT * FROM "person"
WHERE "personname" = $1 AND "password_hash" = $2 LIMIT 1;

-- name: CreatePerson :one
INSERT INTO "person" (
  "personname",
  "email",
  "password_hash"
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: DeletePersonByID :exec
DELETE FROM "person"
WHERE "id" = $1;

-- name: DeletePersonByEmail :exec
DELETE FROM "person"
WHERE "email" = $1;

-- name: UpdatePerson :one
UPDATE "person"
SET "personname" = $2, "email" = $3, "password_hash" = $4
WHERE "id" = $1
RETURNING *;

-- name: UpdatePersonPasswordHashById :one
UPDATE "person"
SET "password_hash" = $2
WHERE "id" = $1
RETURNING *;

