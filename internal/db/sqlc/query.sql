-- name: CreateUserTable :exec
CREATE TABLE IF NOT EXISTS users (
  uuid uuid DEFAULT uuidv7() PRIMARY KEY,
  username text UNIQUE NOT NULL,
  pwhash text NOT NULL
);

-- name: ListUsers :many
SELECT * FROM users
ORDER BY uuid;

-- name: CreateUser :one
INSERT INTO users (
  username, pwhash
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetUserByName :one
SELECT * FROM users 
WHERE username = $1 Limit 1;

-- name: GetUserByUUID :one
SELECT * FROM users 
WHERE uuid = $1 Limit 1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE uuid = $1;

-- name: UpdateUserPW :exec
UPDATE users
  set pwhash = $2
WHERE uuid = $1;
