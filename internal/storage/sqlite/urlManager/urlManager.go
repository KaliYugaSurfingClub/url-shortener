package urlManager

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"link_shortener/internal/storage"
	"time"
)

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) (*Storage, error) {
	const op = "storage.sqlite.urlManager.New"

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
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db}, nil
}

func (s *Storage) GetOriginalByAlias(ctx context.Context, alias string) (int, string, error) {
	const op = "storage.sqlite.urlManager.GetOriginalByAlias"

	var original string
	var id int

	err := s.db.QueryRowContext(ctx, `SELECT id, original FROM link WHERE alias=?`, alias).Scan(&id, &original)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, "", fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
	}
	if err != nil {
		return 0, "", fmt.Errorf("%s: %w", op, err)
	}

	return id, original, nil
}

func (s *Storage) UpdateLastAccess(ctx context.Context, ID int, timestamp time.Time) error {
	op := "storage.sqlite.urlManager.UpdateLastAccess"

	_, err := s.db.ExecContext(ctx, `UPDATE link SET last_access=? WHERE id=?`, timestamp, ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) Save(linkID int) error {
	const op = "storage.sqlite.clicks.CreateClick"

	if _, err := s.db.Exec(`INSERT INTO clicks(link_id) VALUES (?)`, linkID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
