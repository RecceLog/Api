package database

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// InitDBConnection sets up a basic connection pool to interact with db.
func InitDBConnection(ctx context.Context, connStr string) (*pgxpool.Pool, error) {

	// Create configuration for connection pool
	config, err := configDBConnection(connStr)
	if err != nil {
		return nil, err
	}

	// Create connection pool
	pool, err := pgxpool.NewWithConfig(ctx, config)

	return pool, err
}

func configDBConnection(connStr string) (*pgxpool.Config, error) {

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		slog.Error("Error parsing configuration for database connection pool", "error", err.Error())
		return nil, err
	}

	config.MinConns = 1
	config.MaxConns = 10
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = time.Minute * 30
	config.HealthCheckPeriod = time.Minute
	config.ConnConfig.ConnectTimeout = time.Second * 5

	config.PrepareConn = func(ctx context.Context, conn *pgx.Conn) (bool, error) {
		err := conn.Ping(ctx)
		if err != nil {
			slog.Error("Error pinging database", "error", err.Error())
			return false, err
		}
		return true, nil
	}

	config.AfterRelease = func(c *pgx.Conn) bool {
		return true
	}

	config.BeforeClose = func(c *pgx.Conn) {

	}

	return config, nil
}
