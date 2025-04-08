package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewStorage creates and returns connections pull to PostgreSQL
func NewStorage(ctx context.Context, dsn string, maxAttempts int,
	attemptTimeout time.Duration) (*pgxpool.Pool, error) {
	var pool *pgxpool.Pool
	var err error

	// Parse config for pool
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to parse postgres DSN: %w", err)
	}

	// Try to connect several times
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		// Create pool
		pool, err = pgxpool.NewWithConfig(ctx, poolConfig)
		if err == nil {
			// Check connection
			err = pool.Ping(ctx)
			if err == nil {
				return pool, nil // Successfully connection
			}
		}
		fmt.Printf("Attempt %d: Failed to conntect to postgres: %v\n", attempt, err)
		if attempt < maxAttempts {
			time.Sleep(attemptTimeout)
		}
	}

	// If all attempts failed
	return nil, fmt.Errorf("failed to connect to postgres after %d attempts: %w",
		maxAttempts, err)
}
