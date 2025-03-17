package database

import "database/sql"

// RunMigrations создаёт таблицу songs, если она отсутствует.
func RunMigrations(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS songs (
		id SERIAL PRIMARY KEY,
		group_name VARCHAR(255) NOT NULL,
		song_name VARCHAR(255) NOT NULL,
		release_date VARCHAR(50),
		lyrics TEXT,
		link TEXT
	);
	`
	_, err := db.Exec(query)
	return err
}
