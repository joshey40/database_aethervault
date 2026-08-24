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