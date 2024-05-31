-- name: GetBoardsByUserId :many
SELECT
  *
FROM
  boards
WHERE
  user_id = ?;

-- name: GetBoardById :one
SELECT
  *
FROM
  boards
WHERE
  id = ?;

-- name: CreateBoard :exec
INSERT INTO
  boards (user_id, title, sort_order, created_at)
VALUES
  (?, ?, ?, NOW());

-- name: UpdateBoard :exec
UPDATE
  boards
SET
  title = ?,
  updated_at = NOW()
WHERE
  id = ?;

-- name: DeleteBoard :exec
DELETE FROM
  boards
WHERE
  id = ?;

-- name: GetLastBoardByUserId :one
SELECT
  *
FROM
  boards
WHERE
  user_id = ?
ORDER BY
  created_at DESC
LIMIT
  1;

-- name: GetLastInsertBoardByUserId :one
SELECT
  *
FROM
  boards
WHERE
  user_id = ?
ORDER BY
  created_at DESC
LIMIT
  1;