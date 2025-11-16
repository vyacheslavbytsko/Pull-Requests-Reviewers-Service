package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func NewDB() (*DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if err := createTables(ctx, pool); err != nil {
		return nil, err
	}

	return &DB{Pool: pool}, nil
}

func createTables(ctx context.Context, pool *pgxpool.Pool) error {
	queries := []string{
		`CREATE TABLE teams (
			team_name TEXT PRIMARY KEY
		);`,
		`CREATE TABLE users (
			user_id TEXT PRIMARY KEY,
			username TEXT NOT NULL,
			team_name TEXT NOT NULL REFERENCES teams(team_name),
			is_active BOOLEAN NOT NULL DEFAULT TRUE
		);`,
	}

	for _, q := range queries {
		if _, err := pool.Exec(ctx, q); err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}

	return nil
}

func (db *DB) Close() {
	db.Pool.Close()
}
