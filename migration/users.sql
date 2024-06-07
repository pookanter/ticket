-- name: GetUserByID :one
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

-- name: GetLastInsertUser :one
SELECT
  *
FROM
  users
WHERE
  id = (
    SELECT
      LAST_INSERT_ID()
    FROM
      users AS u
    LIMIT
      1
  );

-- name: GetLastInsertUserID :one
SELECT
  LAST_INSERT_ID()
FROM
  users
LIMIT
  1;