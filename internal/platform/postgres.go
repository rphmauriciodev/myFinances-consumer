package platform

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool() (*pgxpool.Pool, error) {
	envVar := "DATABASE_URL"

	connString := os.Getenv(envVar)
	if connString == "" {
		return nil, fmt.Errorf("variável %s não definida", envVar)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	slog.Info(fmt.Sprintf("conexão com o postgres estabelecida com sucesso usando %s", envVar))
	return pool, nil
}
