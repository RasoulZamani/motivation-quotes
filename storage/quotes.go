package storage

import (
	"database/sql"
	"encoding/json"
	"math/rand"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type QuoteStore interface {
	GetRandomQuote() (string, error)
}

type SQLiteStore struct {
	db *sql.DB
}

func NewQuoteStore(source string) (QuoteStore, error) {
	if strings.HasSuffix(source, ".db") {
		db, err := sql.Open("sqlite3", source)
		if err != nil {
			return nil, err
		}
		return &SQLiteStore{db: db}, nil
	} else {
		return NewJSONStore(source), nil
	}
}

func (s *SQLiteStore) GetRandomQuote() (string, error) {
	var quote string
	err := s.db.QueryRow("SELECT text FROM quotes ORDER BY RANDOM() LIMIT 1").Scan(&quote)
	return quote, err
}

// JSON Implementation (fallback)
type JSONStore struct {
	quotes []string
}

func NewJSONStore(filename string) *JSONStore {
	data, _ := os.ReadFile(filename)
	var quotes []string
	json.Unmarshal(data, &quotes)
	return &JSONStore{quotes: quotes}
}

func (j *JSONStore) GetRandomQuote() (string, error) {
	return j.quotes[rand.Intn(len(j.quotes))], nil
}