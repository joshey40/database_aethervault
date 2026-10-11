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
WHERE username = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM Users 
WHERE user_id = $1 LIMIT 1;

-- name: DeleteUser :exec
DELETE FROM Users
WHERE user_id = $1;

-- name: UpdateUserPW :exec
UPDATE Users
  set pwhash = $2
WHERE user_id = $1;

-- name: ListDevices :many
SELECT * FROM Devices
ORDER BY device_id;

-- name: CreateDevice :one
INSERT INTO Devices (
  user_id, last_seen, login_token, device_name
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetAllDevicesFromUser :many
SELECT * FROM Devices
WHERE user_id = $1;

-- name: GetDeviceByID :one
SELECT * FROM Devices
WHERE device_id = $1 LIMIT 1;

-- name: CheckLoginTokenForUser :one
SELECT * From Devices
WHERE user_id = $1 AND login_token = $2 LIMIT 1;