package main

import (
	"Api/internal"
	"Api/internal/infrastructure/database"
	"Api/internal/presentation"
	"context"
	"fmt"
	"log/slog"
	"os"
)

func init() {
	internal.LoadEnvVariables()
}

func main() {

	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	apiCfg := presentation.ApiConfig{
		Address: ":8080",
		Db: database.DbConfig{
			ConnString: fmt.Sprintf(
				"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
				os.Getenv("DEV_DB_HOST"),
				os.Getenv("DEV_DB_PORT"),
				os.Getenv("DEV_DB_USER"),
				os.Getenv("DEV_DB_PASSWORD"),
				os.Getenv("DEV_DB"),
			),
		},
	}

	connPool, err := database.InitDBConnection(ctx, apiCfg.Db.ConnString)
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
	}
}
