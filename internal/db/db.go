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
		`CREATE TABLE IF NOT EXISTS teams (
			team_name TEXT PRIMARY KEY
		);`,
		`CREATE TABLE IF NOT EXISTS users (
			user_id TEXT PRIMARY KEY,
			username TEXT NOT NULL,
			team_name TEXT NOT NULL REFERENCES teams(team_name),
			is_active BOOLEAN NOT NULL DEFAULT TRUE
		);`,
		`CREATE TABLE IF NOT EXISTS prs (
    		pull_request_id TEXT PRIMARY KEY,
    		pull_request_name TEXT NOT NULL,
    		author_id TEXT NOT NULL REFERENCES users(user_id),
    		status TEXT NOT NULL CHECK ( status IN ('OPEN', 'MERGED')),
    		assigned_reviewers TEXT[] NOT NULL,
    		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    		merged_at TIMESTAMPTZ
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
