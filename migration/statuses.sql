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

-- name: GetStatusesWithBoard :many
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

-- name: GetStatus :one
SELECT
  *
FROM
  statuses
WHERE
  id = coalesce(sqlc.narg('id'), id)
  AND board_id = coalesce(sqlc.narg('board_id'), board_id);

-- name: GetStatuses :many
SELECT
  *
FROM
  statuses
WHERE
  board_id = coalesce(sqlc.narg('board_id'), board_id)
ORDER BY
  board_id ASC,
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

-- name: GetLastInsertStatusID :one
SELECT
  LAST_INSERT_ID()
FROM
  statuses
LIMIT
  1;