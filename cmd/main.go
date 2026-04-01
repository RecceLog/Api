package main

import (
	"Api/internal"
	"Api/internal/infrastructure/cache"
	"Api/internal/infrastructure/database"
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func init() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	slog.Debug("Logger initialized")
	/*internal.LoadEnvVariables()
	slog.Debug("Environment variables initialized")*/
}

func main() {

	ctx := context.Background()

	dbConnPool, rdb := configureServices(ctx)

	application := internal.Application{
		Addr:  ":8080",
		Db:    dbConnPool,
		Cache: rdb,
	}

	if err := application.Run(application.Mount()); err != nil {
		slog.Error("Server failed to start", "error", err.Error())
		os.Exit(1)
	}
}

func configureServices(ctx context.Context) (*pgxpool.Pool, *redis.Client) {
	// database configuration
	dbConnPool, err := database.ConnectDB(ctx)
	if err != nil {
		slog.Error("Error connecting to database, shutting down", "error message", err.Error())
		os.Exit(1)
	}
	slog.Debug("Connected to database")

	// caching configuration
	rdb, err := cache.ConnectRedis(ctx)
	if err != nil {
		slog.Error("Error connecting to redis, shutting down", "error message", err.Error())
		os.Exit(1)
	}
	slog.Debug("Caching initialized")

	return dbConnPool, rdb
}
