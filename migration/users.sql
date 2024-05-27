-- name: GetUser :one
SELECT
  *
FROM
  users
WHERE
  id = ?
LIMIT
  1;

-- name: FindUserByEmail :one
SELECT
  *
FROM
  users
WHERE
  email = ?
LIMIT
  1;

-- name: FindUserByID :one
SELECT
  *
FROM
  users
WHERE
  id = ?;

-- name: CreateUser :exec
INSERT INTO
  users (name, lastname, email, password, created_at)
VALUES
  (?, ?, ?, ?, NOW());

-- name: UpdateUser :exec
UPDATE
  users
SET
  name = ?,
  lastname = ?,
  email = ?,
  password = ?,
  updated_at = NOW()
WHERE
  id = ?;