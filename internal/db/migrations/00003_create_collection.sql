-- +goose Up
CREATE TABLE IF NOT EXISTS `CollectionEntries` (
	`collection_entry_id` TEXT NOT NULL UNIQUE,
	`scryfall_id` TEXT NOT NULL,
	`is_proxy` INTEGER NOT NULL,
	`is_altered` INTEGER NOT NULL,
	`is_signed` INTEGER NOT NULL,
	`quality` INTEGER NOT NULL,
	`user_id` TEXT NOT NULL,
	`finish` INTEGER NOT NULL,
	`folder_id` TEXT NOT NULL,
	`version_at` INTEGER NOT NULL,
	`is_owned` INTEGER NOT NULL,
	PRIMARY KEY (`collection_entry_id`),
	FOREIGN KEY (`user_id`) REFERENCES `Users`(`user_id`),
	FOREIGN KEY (`folder_id`) REFERENCES `Folders`(`folder_id`)
);
CREATE TABLE IF NOT EXISTS `Folders` (
	`folder_id` TEXT NOT NULL UNIQUE,
	`folder_name` TEXT NOT NULL UNIQUE,
	`parent_folder` TEXT,
	`is_private` INTEGER NOT NULL,
	`color` TEXT NOT NULL,
	`is_deck` INTEGER NOT NULL,
	`format` TEXT,
	`description` TEXT NOT NULL,
	`bracket` INTEGER,
	`user_id` TEXT NOT NULL,
	`version_at` INTEGER NOT NULL,
	PRIMARY KEY (`folder_id`),
	FOREIGN KEY (`parent_folder`) REFERENCES `Folders`(`folder_id`),
	FOREIGN KEY (`user_id`) REFERENCES `Users`(`user_id`)
);

-- +goose Down
DROP TABLE Folders;
DROP TABLE CollectionEntries;
