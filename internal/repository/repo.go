package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DB interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type Repository interface {
	AddNewURL(ctx context.Context, url, code string) (string, error)
	GetCode(ctx context.Context, url string) (string, error)
	GetURL(ctx context.Context, code string) (string, error)
}

type repository struct {
	db DB
}

func NewRepository(db DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) AddNewURL(ctx context.Context, url, code string) (string, error) {
	query := `INSERT INTO urls (url, code) VALUES ($1, $2)
				ON CONFLICT (url) DO UPDATE SET url = EXCLUDED.url
				RETURNING code`

	var stored string

	if err := r.db.QueryRow(ctx, query, url, code).Scan(&stored); err != nil {
		return "", err
	}

	return stored, nil
}

func (r *repository) GetCode(ctx context.Context, url string) (string, error) {
	query := `SELECT code FROM urls WHERE url = $1`

	var code string
	err := r.db.QueryRow(ctx, query, url).Scan(&code)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrCodeNotFound
	} else if err != nil {
		return "", err
	}

	return code, nil
}

func (r *repository) GetURL(ctx context.Context, code string) (string, error) {
	query := `SELECT url FROM urls WHERE code = $1`

	var url string
	err := r.db.QueryRow(ctx, query, code).Scan(&url)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrUrlNotFound
	} else if err != nil {
		return "", err
	}

	return url, nil
}
