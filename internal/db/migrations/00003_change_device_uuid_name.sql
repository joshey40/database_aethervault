-- +goose Up
ALTER TABLE Devices
RENAME COLUMN device_uuid to device_id;

-- +goose Down
ALTER TABLE Devices
RENAME COLUMN device_id to device_uuid;

