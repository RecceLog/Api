package database

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB(ctx context.Context) (*pgxpool.Pool, error) {
	pool, err := initDBConnection(ctx)
	if err != nil {
		return nil, err
	}

	for attempt := 1; attempt < 4; attempt++ {
		err = pool.Ping(ctx)
		if err != nil {
			slog.Warn(
				fmt.Sprintf("Attempt %d to connect to database failed, retrying in 5 seconds...", attempt),
				"error message", err.Error())
		} else {
			break
		}
	}

	return pool, err
}

func initDBConnection(ctx context.Context) (*pgxpool.Pool, error) {

	connString, exists := os.LookupEnv("DB_CONN_STRING")
	if !exists {
		return nil, fmt.Errorf("DB_CONN_STRING environment variable not found")
	}

	// Create configuration for connection pool
	config, err := configDBConnection(connString)
	if err != nil {
		return nil, fmt.Errorf("error creating pool configuration: %s", err)
	}

	// Create connection pool
	return pgxpool.NewWithConfig(ctx, config)
}

func configDBConnection(connStr string) (*pgxpool.Config, error) {

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	config.MinConns = 1
	config.MaxConns = 10
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = time.Minute * 30
	config.HealthCheckPeriod = time.Minute
	config.ConnConfig.ConnectTimeout = time.Second * 5

	/*config.PrepareConn = func(ctx context.Context, conn *pgx.Conn) (bool, error) {
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

	}*/

	return config, nil
}
