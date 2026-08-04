package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresConfig struct {
	DSN                   string
	MaxConns              int32
	MinConns              int32
	MaxConnLifetime       time.Duration
	MaxConnLifetimeJitter time.Duration
	MaxConnIdleTime       time.Duration
	HealthCheckPeriod     time.Duration
	StatementTimeoutMS    int
	ApplicationName       string
	ConnectTimeout        time.Duration
}

func DefaultPostgresConfig(dsn string) PostgresConfig {
	return PostgresConfig{
		DSN:                   dsn,
		MaxConns:              25,
		MinConns:              4,
		MaxConnLifetime:       time.Hour,
		MaxConnLifetimeJitter: 10 * time.Minute,
		MaxConnIdleTime:       15 * time.Minute,
		HealthCheckPeriod:     time.Minute,
		StatementTimeoutMS:    5000,
		ApplicationName:       "url-shortener-v3",
		ConnectTimeout:        5 * time.Second,
	}
}

func ConnectPostgres(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	return ConnectPostgresWithConfig(ctx, DefaultPostgresConfig(dsn))
}

func ConnectPostgresWithConfig(ctx context.Context, config PostgresConfig) (*pgxpool.Pool, error) {
	if config.DSN == "" {
		return nil, errors.New("postgres dsn is empty")
	}

	ctx, cancel := context.WithTimeout(ctx, config.ConnectTimeout)
	defer cancel()

	cfg, err := pgxpool.ParseConfig(config.DSN)
	if err != nil {
		return nil, fmt.Errorf("parse dsn string: %w", err)
	}

	cfg.MinConns = config.MinConns
	cfg.MaxConns = config.MaxConns
	cfg.MaxConnLifetime = config.MaxConnLifetime
	cfg.MaxConnLifetimeJitter = config.MaxConnLifetimeJitter
	cfg.MaxConnIdleTime = config.MaxConnIdleTime
	cfg.HealthCheckPeriod = config.HealthCheckPeriod

	if cfg.ConnConfig.RuntimeParams == nil {
		cfg.ConnConfig.RuntimeParams = make(map[string]string)
	}

	cfg.ConnConfig.RuntimeParams["statement_timeout"] = fmt.Sprintf("%d", config.StatementTimeoutMS)
	cfg.ConnConfig.RuntimeParams["application_name"] = config.ApplicationName

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return pool, nil
}
