-- name: GetTicketsByStatusId :many
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
    created_at
  )
VALUES
  (?, ?, ?, ?, NOW());

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