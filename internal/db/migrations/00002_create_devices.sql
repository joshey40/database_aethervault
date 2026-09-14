-- +goose Up
CREATE TABLE Devices (
	device_uuid uuid DEFAULT uuidv7() PRIMARY KEY,
	user_id uuid NOT NULL,
	last_seen timestamp with time zone NOT NULL,
	login_token text NOT NULL,
	device_name text NOT NULL,
  CONSTRAINT Devices_fk0 FOREIGN KEY (user_id) REFERENCES Users(user_id)
);

-- +goose Down
DROP TABLE Devices;
