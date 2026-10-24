-- name: ListUsers :many
SELECT * FROM Users
ORDER BY user_id;

-- name: CreateUser :one
INSERT INTO Users (
  user_id, username, pwhash, role
) VALUES (
  sqlc.arg(user_id), sqlc.arg(username), sqlc.arg(pwhash), sqlc.arg(role)
) RETURNING *;

-- name: GetUserByName :one
SELECT * FROM Users 
WHERE username = sqlc.arg(username) LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM Users 
WHERE user_id = sqlc.arg(user_id) LIMIT 1;

-- name: DeleteUser :exec
DELETE FROM Users
WHERE user_id = sqlc.arg(user_id);

-- name: UpdateUserPW :exec
UPDATE Users
  set pwhash = sqlc.arg(pwhash)
WHERE user_id = sqlc.arg(user_id);

-- name: ListDevices :many
SELECT * FROM Devices
ORDER BY device_id;

-- name: CreateDevice :one
INSERT INTO Devices (
  device_id, user_id, last_seen, device_name
) VALUES (
  sqlc.arg(device_id), sqlc.arg(user_id), sqlc.arg(last_seen), sqlc.arg(device_name)
) RETURNING *;

-- name: GetAllDevicesFromUser :many
SELECT * FROM Devices
WHERE user_id = sqlc.arg(user_id);

-- name: GetDeviceByID :one
SELECT * FROM Devices
WHERE device_id = sqlc.arg(device_id) LIMIT 1;

-- name: DeleteDevice :exec
DELETE FROM Devices
WHERE device_id = sqlc.arg(device_id);

-- name: GetDeviceByIDAndUserID :one
SELECT * FROM Devices
WHERE device_id = sqlc.arg(device_id) and user_id = sqlc.arg(user_id) LIMIT 1;