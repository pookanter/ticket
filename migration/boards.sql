-- name: GetBoardsByUserID :many
SELECT
  *
FROM
  boards
WHERE
  user_id = ?;

-- name: GetBoard :one
SELECT
  *
FROM
  boards
WHERE
  id = ?
  AND user_id = ?;
 
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

-- name: CountBoardByUserID :one
SELECT
  COUNT(*)
FROM
  boards
WHERE
  user_id = ?;

-- name: GetLastInsertBoardID :one
SELECT
  LAST_INSERT_ID()
FROM
  statuses
LIMIT
  1;