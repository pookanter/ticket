-- name: GetStatusWithBoard :one
SELECT
  sqlc.embed(statuses),
  sqlc.embed(boards)
FROM
  statuses
  JOIN boards ON statuses.board_id = boards.id
WHERE
  statuses.id = ?
  AND statuses.board_id = ?
  AND boards.user_id = ?;

-- name: GetStatusesByBoardID :many
SELECT
  *
FROM
  statuses
WHERE
  board_id = ?
ORDER BY
  sort_order ASC;

-- name: GetStatusesWithMinimumSortOrder :many
SELECT
  *
FROM
  statuses
WHERE
  board_id = ?
  AND sort_order >= ?
ORDER BY
  (
    CASE
      WHEN sqlc.arg('sort_order_direction') = 'asc' THEN sort_order
    END
  ) ASC,
  (
    CASE
      WHEN sqlc.arg('sort_order_direction') = 'desc' THEN sort_order
    END
  ) DESC;

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
  sort_order = ?,
  updated_at = NOW()
WHERE
  id = ?;

-- name: UpdateStatusSortOrder :exec
UPDATE
  statuses
SET
  sort_order = ?
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

-- name: GetStatusView :one
SELECT
  *
FROM
  status_view
WHERE
  id = ?
ORDER BY
  sort_order ASC;

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