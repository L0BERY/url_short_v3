package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
)

func newMockRepo(t *testing.T) (Repository, pgxmock.PgxPoolIface) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("не удалось создать pgxmock pool: %v", err)
	}
	t.Cleanup(mock.Close)
	t.Cleanup(func() {
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("остались невыполненные ожидания к БД: %v", err)
		}
	})

	return NewRepository(mock), mock
}

func TestAddNewURL(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo, mock := newMockRepo(t)

		const (
			url  = "https://example.com"
			code = "abc123"
		)

		mock.ExpectQuery("INSERT INTO urls").
			WithArgs(url, code).
			WillReturnRows(pgxmock.NewRows([]string{"code"}).AddRow(code))

		got, err := repo.AddNewURL(ctx, url, code)
		if err != nil {
			t.Fatalf("неожиданная ошибка: %v", err)
		}
		if got != code {
			t.Errorf("код = %q, ожидали %q", got, code)
		}
	})

	t.Run("URL conflict - returning existing code", func(t *testing.T) {
		repo, mock := newMockRepo(t)

		const (
			url          = "https://example.com"
			newCode      = "new111"
			existingCode = "old999"
		)

		mock.ExpectQuery("INSERT INTO urls").
			WithArgs(url, newCode).
			WillReturnRows(pgxmock.NewRows([]string{"code"}).AddRow(existingCode))

		got, err := repo.AddNewURL(ctx, url, newCode)
		if err != nil {
			t.Fatalf("неожиданная ошибка: %v", err)
		}
		if got != existingCode {
			t.Errorf("код = %q, ожидали %q", got, existingCode)
		}
	})

	t.Run("ошибка базы данных пробрасывается наружу", func(t *testing.T) {
		repo, mock := newMockRepo(t)

		dbErr := errors.New("connection refused")

		mock.ExpectQuery("INSERT INTO urls").
			WithArgs("https://example.com", "abc123").
			WillReturnError(dbErr)

		_, err := repo.AddNewURL(ctx, "https://example.com", "abc123")
		if !errors.Is(err, dbErr) {
			t.Fatalf("err = %v, ожидали %v", err, dbErr)
		}
	})

}

func TestRepository_GetCode(t *testing.T) {
	ctx := context.Background()

	t.Run("код найден", func(t *testing.T) {
		repo, mock := newMockRepo(t)

		const (
			url  = "https://example.com"
			code = "abc123"
		)

		mock.ExpectQuery("SELECT code FROM urls").
			WithArgs(url).
			WillReturnRows(pgxmock.NewRows([]string{"code"}).AddRow(code))

		got, err := repo.GetCode(ctx, url)
		if err != nil {
			t.Fatalf("неожиданная ошибка: %v", err)
		}
		if got != code {
			t.Errorf("код = %q, ожидали %q", got, code)
		}
	})

	t.Run("url не найден — маппится в ErrCodeNotFound", func(t *testing.T) {
		repo, mock := newMockRepo(t)

		mock.ExpectQuery("SELECT code FROM urls").
			WithArgs("https://missing.example.com").
			WillReturnError(pgx.ErrNoRows)

		_, err := repo.GetCode(ctx, "https://missing.example.com")
		if !errors.Is(err, ErrCodeNotFound) {
			t.Fatalf("err = %v, ожидали ErrCodeNotFound", err)
		}
	})

	t.Run("прочая ошибка БД не превращается в ErrCodeNotFound", func(t *testing.T) {
		repo, mock := newMockRepo(t)

		dbErr := errors.New("connection refused")

		mock.ExpectQuery("SELECT code FROM urls").
			WithArgs("https://example.com").
			WillReturnError(dbErr)

		_, err := repo.GetCode(ctx, "https://example.com")
		if !errors.Is(err, dbErr) {
			t.Fatalf("err = %v, ожидали %v", err, dbErr)
		}
		if errors.Is(err, ErrCodeNotFound) {
			t.Fatalf("произвольная ошибка БД не должна маскироваться под ErrCodeNotFound")
		}
	})
}

func TestRepository_GetURL(t *testing.T) {
	ctx := context.Background()

	t.Run("url найден", func(t *testing.T) {
		repo, mock := newMockRepo(t)

		const (
			code = "abc123"
			url  = "https://example.com"
		)

		mock.ExpectQuery("SELECT url FROM urls").
			WithArgs(code).
			WillReturnRows(pgxmock.NewRows([]string{"url"}).AddRow(url))

		got, err := repo.GetURL(ctx, code)
		if err != nil {
			t.Fatalf("неожиданная ошибка: %v", err)
		}
		if got != url {
			t.Errorf("url = %q, ожидали %q", got, url)
		}
	})

	t.Run("код не найден — маппится в ErrUrlNotFound", func(t *testing.T) {
		repo, mock := newMockRepo(t)

		mock.ExpectQuery("SELECT url FROM urls").
			WithArgs("zzzzzz").
			WillReturnError(pgx.ErrNoRows)

		_, err := repo.GetURL(ctx, "zzzzzz")
		if !errors.Is(err, ErrUrlNotFound) {
			t.Fatalf("err = %v, ожидали ErrUrlNotFound", err)
		}
	})
}
