
-- name: ListTitlePerson :many
SELECT * FROM "title_person";

-- name: CreateTitlePerson :one
INSERT INTO "title_person" ("person_id", "title_id") VALUES ($1, $2) RETURNING *;

-- name: DeleteTitlePersonByPersonIDAndTitleID :exec
DELETE FROM "title_person" WHERE "person_id" = $1 AND "title_id" = $2;

-- name: DeleteAllPersonTitles :exec
DELETE FROM "title_person"
WHERE "person_id" = $1;

-- name: GetPersonTitles :many
SELECT * FROM "title_person" WHERE "person_id" = $1;