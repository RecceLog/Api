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
	internal.LoadEnvVariables()
	slog.Debug("Environment variables initialized")
}

func main() {

	ctx := context.Background()

	dbConnPool, rdb := configureServices(ctx)

	/*var connString string
	value, exists := os.LookupEnv("DB_CONN_STRING")
	if exists {
		connString = value
	} else {
		connString = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("DEV_DB_HOST"),
			os.Getenv("DEV_DB_PORT"),
			os.Getenv("DEV_DB_USER"),
			os.Getenv("DEV_DB_PASSWORD"),
			os.Getenv("DEV_DB"),
		)
	}

	apiCfg := presentation.ApiConfig{
		Address: ":8080",
		Db: database.DbConfig{
			ConnString: connString,
		},
	}

	connPool, err := database.InitDBConnection(ctx)
	if err != nil {
		slog.Error("Error connecting to database", "error", err.Error())
		return
	}
	defer connPool.Close()

	slog.Info("Connected to database", "db", apiCfg.Db.ConnString)

	api := presentation.Application{
		Config: apiCfg,
		Db:     connPool,
	}

	if err := api.Run(api.Mount()); err != nil {
		slog.Error("Server failed to start", "error", err.Error())
		os.Exit(1)
	}*/
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

	// api configuration

	return dbConnPool, rdb
}
