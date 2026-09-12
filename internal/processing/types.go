package processing

import (
	"context"
	"time"
)

type Transaction struct {
	Name            string `json:"name"`
	Merchant        string `json:"merchant"`
	CardOrPass      string `json:"cardOrPass"`
	AmountValue     int64
	TransactionDate time.Time
}

type Repository interface {
	SaveTransaction(ctx context.Context, transaction *Transaction) error
}
