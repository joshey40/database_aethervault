-- +goose Up
CREATE TABLE Devices (
	device_id text PRIMARY KEY,
	user_id text NOT NULL,
	last_seen integer NOT NULL,
	device_name text NOT NULL,
  CONSTRAINT Devices_fk0 FOREIGN KEY (user_id) REFERENCES Users(user_id)
);

-- +goose Down
DROP TABLE Devices;
