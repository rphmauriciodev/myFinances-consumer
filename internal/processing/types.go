package processing

import "context"

type Transaction struct {
	Name       string `json:"name"`
	Merchant   string `json:"merchant"`
	CardOrPass string `json:"cardOrPass"`
	Amount     string `json:"amount"`
	Date       string `json:"date"`
}
type Repository interface {
	SaveTransaction(ctx context.Context, transaction *Transaction) error
}
