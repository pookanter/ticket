-- name: GetTicketsByStatusID :many
SELECT
  *
FROM
  tickets
WHERE
  status_id = ?;

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
  title = ?,
  description = ?,
  contact = ?,
  updated_at = NOW()
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