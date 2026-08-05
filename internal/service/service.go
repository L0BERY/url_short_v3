package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
	"url_shortener_v3/internal/repository"
)

type Service interface {
}

type service struct {
	repo repository.PostgresRepository
}

func NewService(repo repository.PostgresRepository) Service {
	return &service{repo: repo}
}

func (s *service) GenearteCode() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// func (s *service) GetCode(url string) (string, error) {

// }

func (s *service) GetUrl(ctx context.Context, code string) (string, error) {
	if code == "" {
		return "", ErrEmptyCode
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	url, err := s.repo.GetUrl(ctx, code)
	if errors.Is(err, repository.ErrUrlNotFound) {
		return "", ErrUrlNotFound
	} else if err != nil {
		return "", ErrServer
	}

	return url, nil
}
