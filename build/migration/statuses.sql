-- name: GetStatusesByBoardId :many
SELECT
  *
FROM
  statuses
WHERE
  board_id = ?;

-- name: CreateStatus :exec
INSERT INTO
  statuses (board_id, title, created_at)
VALUES
  (?, ?, NOW());

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