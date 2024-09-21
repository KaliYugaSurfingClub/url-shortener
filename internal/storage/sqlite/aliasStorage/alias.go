package aliasStorage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"link_shortener/internal/storage"
	"time"
)

type AliasStorage struct {
	db *sql.DB
}

func New(db *sql.DB) (*AliasStorage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS url(
			id INTEGER PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			original TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	  	);
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &AliasStorage{db}, nil
}

func (s *AliasStorage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetURL"

	stmt, err := s.db.Prepare("SELECT original FROM url WHERE alias=?")
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

func (s *AliasStorage) SaveAlias(original string, alias string, timeToGenerate time.Duration) (string, error) {
	const op = "storage.sqlite.SaveAlias"

	stmt, err := s.db.Prepare("INSERT INTO url(original, alias) VALUES(?, ?)")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	shouldGenerate := alias == ""
	if shouldGenerate {
		alias = randomAlias()
	}

	for startTime := time.Now(); time.Since(startTime) < timeToGenerate; {
		_, err = stmt.Exec(original, alias)

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

func (s *AliasStorage) DeleteAlias(alias string) error {
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
