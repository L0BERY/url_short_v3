package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
	"url_shortener_v3/internal/repository"

	"github.com/jackc/pgx/v5/pgconn"
)

type Service interface {
	GetURL(ctx context.Context, code string) (string, error)
	AddNewURL(ctx context.Context, url string) (string, error)
}

type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s *service) GenerateCode() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (s *service) AddNewURL(ctx context.Context, url string) (string, error) {
	if url == "" {
		return "", ErrEmptyUrl
	}

	for range 10 {
		code := s.GenerateCode()
		code, err := s.repo.AddNewURL(ctx, url, code)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "urls_code_key" {
			continue
		}
		return code, err
	}
	return "", ErrTooManyAttempts
}

// func (s *service) GetCode(url string) (string, error) {

// }

func (s *service) GetURL(ctx context.Context, code string) (string, error) {
	if code == "" {
		return "", ErrEmptyCode
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	url, err := s.repo.GetURL(ctx, code)
	if errors.Is(err, repository.ErrUrlNotFound) {
		return "", ErrUrlNotFound
	} else if err != nil {
		return "", ErrServer
	}

	return url, nil
}
