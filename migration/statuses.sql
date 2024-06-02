-- name: GetStatusesByBoardID :many
SELECT
  *
FROM
  statuses
WHERE
  board_id = ?;

-- name: CreateStatus :exec
INSERT INTO
  statuses (board_id, title, sort_order, created_at)
VALUES
  (?, ?, ?, NOW());

-- name: UpdateStatus :exec
UPDATE
  statuses
SET
  title = ?,
  updated_at = NOW()
WHERE
  id = ?;

-- name: DeleteStatus :exec
DELETE FROM
  statuses
WHERE
  id = ?;

-- name: CountStatusByBoardID :one
SELECT
  COUNT(*)
FROM
  statuses
WHERE
  board_id = ?;

-- name: GetLastInsertStatusViewByBoardID :one
SELECT
  *
FROM
  status_view
WHERE
  status_view.board_id = ?
  AND id = (
    SELECT
      LAST_INSERT_ID()
    FROM
      statuses
    LIMIT
      1
  );