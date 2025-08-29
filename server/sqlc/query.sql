-- name: GetCache :one
SELECT *
FROM tbl_cache
WHERE key = $1
LIMIT 1;
-- name: ListSymbols :many
SELECT *
FROM tbl_seven_fifty
ORDER BY id ASC;
-- name: CreateCache :exec
INSERT INTO tbl_cache (key, value, created)
VALUES ($1, $2, $3);
-- name: UpsertCache :exec
INSERT INTO tbl_cache (key, value, created)
VALUES ($1, $2, $3) ON CONFLICT (key) DO
UPDATE
SET value = EXCLUDED.value,
  created = EXCLUDED.created;
-- name: UpdateCache :exec
UPDATE tbl_cache
SET value = $2
WHERE key = $1;
-- name: DeleteCache :exec
DELETE FROM tbl_cache
WHERE key = $1;