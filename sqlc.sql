-- name: InsertUser :execresult
INSERT INTO
  user (user_name, email, password)
VALUES
  (?, ?, ?);

-- name: GetLastInsertID :one
SELECT
  LAST_INSERT_ID();

-- name: GetUserByEmail :one
SELECT
  user_id,
  user_name,
  email
FROM
  user
WHERE
  email = ?;

-- name: ExistsUser :one
SELECT
  EXISTS(
    SELECT
      1
    FROM
      user
    WHERE
      email = ?
  );
