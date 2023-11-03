
-- name: ListTitles :many
SELECT * FROM "title";

-- name: GetTitle :one
SELECT * FROM "title"
WHERE "id" = $1 LIMIT 1;

-- name: GetTitleByName :one
SELECT * FROM "title"
WHERE "name" = $1 LIMIT 1;

-- name: CreateTitle :one
INSERT INTO "title" (
    "name",
    "description"
) VALUES (
    $1 , $2
) RETURNING *;

-- name: DeleteTitleByID :exec
DELETE FROM "title"
WHERE "id" = $1;

-- name: DeleteTitleByName :exec
DELETE FROM "title"
WHERE "name" = $1;

-- name: UpdateTitle :one
UPDATE "title"
SET "name" = $2, "description" = $3
WHERE "id" = $1
RETURNING *;

