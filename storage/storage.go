package storage

import (
	"database/sql"
	"encoding/json"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type QuoteStorage interface {
	GetRandomQuote() (string, error)
	SyncFromJSON(filename string) error
}

type SQLiteStorage struct {
	db *sql.DB
}

func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS quotes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			text TEXT NOT NULL UNIQUE
		)
	`)
	if err != nil {
		return nil, err
	}

	return &SQLiteStorage{db: db}, nil
}

func (s *SQLiteStorage) GetRandomQuote() (string, error) {
	var quote string
	err := s.db.QueryRow("SELECT text FROM quotes ORDER BY RANDOM() LIMIT 1").Scan(&quote)
	return quote, err
}

func (s *SQLiteStorage) SyncFromJSON(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var quotes []string
	if err := json.Unmarshal(data, &quotes); err != nil {
		return err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT OR IGNORE INTO quotes (text) VALUES (?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, quote := range quotes {
		if _, err := stmt.Exec(quote); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}