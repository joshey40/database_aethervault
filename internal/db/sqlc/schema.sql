CREATE TABLE IF NOT EXISTS users (
  uuid uuid DEFAULT uuidv7() PRIMARY KEY,
  username text UNIQUE NOT NULL,
  pwhash text NOT NULL,
  role int NOT NULL
);