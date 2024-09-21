package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

func InitTables(db *sql.DB) error {
	const op = "storage.sqlite.InitTables"

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS link(
			id INTEGER PRIMARY KEY,
			original TEXT NOT NULL,
			alias TEXT NOT NULL UNIQUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
			last_access TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL	
	  	);
		CREATE TABLE IF NOT EXISTS clicks(
		    id INTEGER PRIMARY KEY,
			link_id INTEGER,
		    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
			FOREIGN KEY (link_id) REFERENCES link(id)
		);
	`)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
