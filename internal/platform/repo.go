package platform

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rphmauriciodev/myFinances-consumer/internal/processing"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(p *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		pool: p,
	}
}

func (r *PostgresRepository) SaveTransaction(ctx context.Context, transaction *processing.Transaction) error {
	query := `
		INSERT INTO transactions (name, merchant, card_or_pass, amount_value, transaction_date)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.pool.Exec(
		ctx,
		query,
		transaction.Name,
		transaction.Merchant,
		transaction.CardOrPass,
		transaction.AmountValue,
		transaction.TransactionDate,
	)

	if err != nil {
		return err
	}
	return nil
}
