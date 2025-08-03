package postgres

import (
	"context"
	"errors"
	"fmt"

	"url-shortener/internal/storage"
	"url-shortener/internal/config"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Storage struct {
	db *pgx.Conn
}

func New(cfg config.Postgres) (*Storage, error) {
	const op = "storage.postgres.New"

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
	)

	db, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}


	_, err = db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS url(
    		id SERIAL PRIMARY KEY,
    		alias TEXT NOT NULL UNIQUE,
    		url TEXT NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.postgres.SaveURL"

	var id int64
	err := s.db.QueryRow(context.Background(),
		"INSERT INTO url(url, alias) VALUES($1, $2) RETURNING id",
		urlToSave, alias).Scan(&id)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" { // Код ошибки unique_violation в PostgreSQL
				return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
			}
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	
	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.postgres.GetURL"

	var resURL string
	err := s.db.QueryRow(context.Background(),
		"SELECT url FROM url WHERE alias = $1",
		alias).Scan(&resURL)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}
		
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resURL, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.postgres.DeleteURL"

	tag, err := s.db.Exec(context.Background(),
		"DELETE FROM url WHERE alias = $1",
		alias)

	if err != nil {
		if tag.RowsAffected() == 0 {
     	   return fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
    	}
		
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}