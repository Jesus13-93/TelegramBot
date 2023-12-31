package sqlite

import (
	"TelegramBot/storage"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

// New creates new SQLite storage
func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}

	return &Storage{db: db}, nil
}

// Save saves page to storage
func (s *Storage) Save(ctx context.Context, p *storage.Page) error {
	q := `Insert into pages (url, user_name,id_user) values (?,?,?)`

	if _, err := s.db.ExecContext(ctx, q, p.URL, p.UserName, p.IdUser); err != nil {
		return fmt.Errorf("can't save page: %w", err)
	}
	return nil
}

// PickRandom picks random page from storage
func (s *Storage) PickRandom(ctx context.Context, userName string, idUser int) (*storage.Page, error) {
	q := `SELECT url from pages where id_user = ? order by RANDOM() LIMIT 1`

	var url string

	err := s.db.QueryRowContext(ctx, q, idUser).Scan(&url)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNoSavedPages
	}
	if err != nil {
		return nil, fmt.Errorf("can't pick random page: %w", err)
	}

	return &storage.Page{
		URL:      url,
		UserName: userName,
		IdUser:   idUser,
	}, nil
}

// Remove removes page from storage
func (s *Storage) Remove(ctx context.Context, page *storage.Page) error {
	q := `Delete from pages where url = ? and id_user = ?`

	if _, err := s.db.ExecContext(ctx, q, page.URL, page.IdUser); err != nil {
		return fmt.Errorf("can't remove page: %w", err)
	}

	return nil
}

// IsExists checks if page exists in storage
func (s *Storage) IsExists(ctx context.Context, page *storage.Page) (bool, error) {
	q := `select COUNT(*) from pages where url = ? and id_user = ?`

	var count int

	if err := s.db.QueryRowContext(ctx, q, page.URL, page.IdUser).Scan(&count); err != nil {
		return false, fmt.Errorf("can't check if page exists: %w", err)
	}

	return count > 0, nil
}

func (s *Storage) Init(ctx context.Context) error {
	q := `Create table if not exists pages (url text, user_name text, id_user INT)`

	_, err := s.db.ExecContext(ctx, q)
	if err != nil {
		return fmt.Errorf("can't create the table: %w", err)
	}

	return nil
}
