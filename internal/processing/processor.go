package processing

import (
	"context"
	"encoding/json"
)

type Processor struct {
	repo Repository
}

func NewProcessor(repo Repository) *Processor {
	return &Processor{
		repo: repo,
	}
}

func (p *Processor) ProcessTransactions(ctx context.Context, transaction *Transaction) error {
	if err := p.repo.SaveTransaction(ctx, transaction); err != nil {
		return err
	}
	return nil
}

func ParseTransaction(messageBody string) (*Transaction, error) {

	var transaction Transaction

	err := json.Unmarshal([]byte(messageBody), &transaction)
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}
