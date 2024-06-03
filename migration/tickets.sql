-- name: GetTicket :one
SELECT
  *
FROM
  tickets
  JOIN statuses ON tickets.status_id = statuses.id
  JOIN boards ON statuses.board_id = boards.id
WHERE
  tickets.id = ?
  AND tickets.status_id = coalesce(sqlc.narg('status_id'), tickets.status_id)
  AND statuses.board_id = coalesce(sqlc.narg('board_id'), tickets.board_id)
  AND boards.user_id = coalesce(sqlc.narg('user_id'), tickets.user_id);

-- name: GetTicketsWithMinimumSortOrder :many
SELECT
  *
FROM
  tickets
WHERE
  status_id = ?
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

-- name: CreateTicket :exec
INSERT INTO
  tickets (
    status_id,
    title,
    description,
    contact,
    sort_order,
    created_at
  )
VALUES
  (?, ?, ?, ?, ?, NOW());

-- name: UpdateTicket :exec
UPDATE
  tickets
SET
  status_id = ?,
  title = ?,
  description = ?,
  contact = ?,
  sort_order = ?,
  updated_at = NOW()
WHERE
  id = ?;

-- name: UpdateTicketSortOrder :exec
UPDATE
  tickets
SET
  sort_order = ?
WHERE
  id = ?;

-- name: CountTicketByStatusID :one
SELECT
  COUNT(*)
FROM
  tickets
WHERE
  status_id = ?;

-- name: GetLastInsertTicketByStatusID :one
SELECT
  *
FROM
  tickets
WHERE
  tickets.status_id = ?
  AND id = (
    SELECT
      LAST_INSERT_ID()
    FROM
      tickets AS t
    LIMIT
      1
  );