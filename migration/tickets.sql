-- name: GetTicketByID :one
SELECT
  *
FROM
  tickets
WHERE
  id = ?;

-- name: GetTicketWithBoard :one
SELECT
  sqlc.embed(tickets),
  sqlc.embed(boards)
FROM
  tickets
  JOIN statuses ON tickets.status_id = statuses.id
  JOIN boards ON statuses.board_id = boards.id
WHERE
  tickets.id = ?
  AND statuses.board_id = ?
  AND boards.user_id = ?;

-- name: GetTicketsWithBoard :many
SELECT
  sqlc.embed(tickets),
  sqlc.embed(boards)
FROM
  tickets
  JOIN statuses ON tickets.status_id = statuses.id
  JOIN boards ON statuses.board_id = boards.id
WHERE
  tickets.id IN (sqlc.slice('ids'))
  AND statuses.board_id = ?
  AND boards.user_id = ?;

-- name: GetTicketsByStatusID :many
SELECT
  *
FROM
  tickets
WHERE
  status_id = ?
ORDER BY
  sort_order ASC;

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

-- name: UpdateTicketSortOrderAndStatusID :exec
UPDATE
  tickets
SET
  sort_order = ?,
  status_id = sqlc.arg('status_id'),
  updated_at = CASE
    WHEN sqlc.arg('status_id') <> tickets.status_id THEN NOW()
    ELSE updated_at
  END
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