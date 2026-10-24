-- +goose Up
CREATE TABLE Users (
  user_id text PRIMARY KEY,
  username text UNIQUE NOT NULL,
  pwhash text NOT NULL,
  role int NOT NULL
);

-- +goose Down
DROP TABLE Users;
