package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type repository struct {
	db DB
}

type DB interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type Repository interface {
	AddNewUrl(ctx context.Context, url, code string) error
	GetCode(ctx context.Context, url string) (string, error)
	GetUrl(ctx context.Context, code string) (string, error)
}

func NewRepository(db DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) AddNewUrl(ctx context.Context, url, code string) error {
	query := `INSERT INTO urls (url, code) VALUES ($1, $2)`

	_, err := r.db.Exec(ctx, query, url, code)
	if err != nil {
		return err
	}

	return nil
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

func (r *repository) GetUrl(ctx context.Context, code string) (string, error) {
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
