-- name: GetBoardsByUserId :many
SELECT
  *
FROM
  boards
WHERE
  user_id = ?;

-- name: CreateBoard :exec
INSERT INTO
  boards (user_id, title, created_at)
VALUES
  (?, ?, NOW());

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