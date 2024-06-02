-- name: GetBoardsByUserID :many
SELECT
  *
FROM
  boards
WHERE
  user_id = ?;

-- name: GetBoardByID :one
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

-- name: CountBoardByUserID :one
SELECT
  COUNT(*)
FROM
  boards
WHERE
  user_id = ?;

-- name: ListBoardViewByUserID :many
SELECT
  *
FROM
  board_view
WHERE
  user_id = ?
ORDER BY
  sort_order ASC;

-- name: GetLastInsertBoardViewByUserID :one
SELECT
  *
FROM
  board_view
WHERE
  board_view.user_id = ?
  AND id = (
    SELECT
      LAST_INSERT_ID()
    FROM
      boards
    LIMIT
      1
  );