package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"link_shortener/internal/lib/random"
	"link_shortener/internal/storage"
	"time"
)

type Storage struct {
	db *sql.DB
}

func randomAlias() string {
	return random.NewRandomString(1, random.AlphaNumAlp()[0:2])
}

// todo maybe refactor op and wrapping

func New(storagePath string) (*Storage, error) {
	//todo find out is it worth wrapping errors with path
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS url(
			id INTEGER PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	  	);
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db}, nil
}

func (s *Storage) SaveAlias(originalURL string, alias string, timeToGenerate time.Duration) (string, error) {
	const op = "storage.sqlite.SaveAlias"

	stmt, err := s.db.Prepare("INSERT INTO url(alias, url) VALUES(?, ?)")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	shouldGenerate := alias == ""
	if shouldGenerate {
		alias = randomAlias()
	}

	for startTime := time.Now(); time.Since(startTime) < timeToGenerate; {
		_, err = stmt.Exec(alias, originalURL)

		exists := err != nil && errors.Is(err.(sqlite3.Error).ExtendedCode, sqlite3.ErrConstraintUnique)

		switch {
		//if db already contains alias, and we got empty string as alias.
		//We generate new alias and do another try
		case exists && shouldGenerate:
			alias = randomAlias()
		//We are got not empty alias, and it is already exists in db => return error
		case exists:
			return "", fmt.Errorf("%s: %w", op, storage.ErrAliasExists)
		case err == nil:
			return alias, nil
		//internal db error
		default:
			return "", fmt.Errorf("%s: %w", op, err)
		}
	}

	return "", fmt.Errorf("%s: %w", op, storage.NotEnoughTimeToGenerate)
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetURL"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias=?")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	var url string
	err = stmt.QueryRow(alias).Scan(&url)
	if errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
	}
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return url, nil
}

func (s *Storage) DeleteOverdueAliases(deadline time.Time) (int64, error) {
	const op = "storage.sqlite.DeleteOverdueAliases"

	stmt, err := s.db.Prepare("DELETE FROM url WHERE created_at < ?")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(deadline)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return affected, nil
}

func (s *Storage) DeleteAlias(alias string) error {
	const op = "storage.sqlite.DeleteAlias"

	stmt, err := s.db.Prepare("DELETE FROM url WHERE alias=?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
