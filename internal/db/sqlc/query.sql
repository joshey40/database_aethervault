-- name: CreateDB :exec
CREATE TABLE IF NOT EXISTS Users (
  user_id uuid DEFAULT uuidv7() PRIMARY KEY,
  username text UNIQUE NOT NULL,
  pwhash text NOT NULL,
  role int NOT NULL
);
CREATE TABLE IF NOT EXISTS Devices (
	device_uuid uuid DEFAULT uuidv7() PRIMARY KEY,
	user_id uuid NOT NULL,
	last_seen timestamp with time zone NOT NULL,
	login_token text NOT NULL,
	device_name text NOT NULL,
  CONSTRAINT Devices_fk0 FOREIGN KEY (user_id) REFERENCES Users(uuid)
);

-- name: ListUsers :many
SELECT * FROM Users
ORDER BY user_id;

-- name: CreateUser :one
INSERT INTO Users (
  username, pwhash, role
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetUserByName :one
SELECT * FROM Users 
WHERE username = $1 Limit 1;

-- name: GetUserByUUID :one
SELECT * FROM Users 
WHERE user_id = $1 Limit 1;

-- name: DeleteUser :exec
DELETE FROM Users
WHERE user_id = $1;

-- name: UpdateUserPW :exec
UPDATE Users
  set pwhash = $2
WHERE user_id = $1;
