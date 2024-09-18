package visit

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Storage struct {
	db       *sql.DB
	lifetime time.Duration
}

func New(storagePath string, lifetime time.Duration) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	//todo migrations
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS visit(
		    id INTEGER PRIMARY KEY,
		    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		    alias_id INTEGER,
			FOREIGN KEY (alias_id) REFERENCES url(id)
		);
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db, lifetime}, nil
}

func (s *Storage) aliasIsActual(aliasId int64) (bool, error) {
	const op = "storage.sqlite.aliasIsActual"

	stmt, err := s.db.Prepare("SELECT created_at FROM visit WHERE alias_id=? and CURRENT_TIMESTAMP-created_at<?")
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	err = stmt.QueryRow(aliasId, s.lifetime).Scan()
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return true, nil
}

func (s *Storage) SaveVisit(aliasId int64) error {
	const op = "storage.sqlite.GetURL"

	stmt, err := s.db.Prepare("INSERT INTO visit(alias_id) VALUES(?)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(aliasId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetUnusedAliases() ([]int64, error) {
	const op = "storage.sqlite.aliasIsActual"

	res := make([]int64, 0)

	stmt, err := s.db.Prepare("SELECT id FROM visit WHERE CURRENT_TIMESTAMP-created_at>?")
	if err != nil {
		return res, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := stmt.Query(s.lifetime)
	if errors.Is(err, sql.ErrNoRows) {
		return res, nil
	}
	if err != nil {
		return res, fmt.Errorf("%s: %w", op, err)
	}

	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			return res, fmt.Errorf("%s: %w", op, err)
		}
		res = append(res, id)
	}

	return res, nil
}
