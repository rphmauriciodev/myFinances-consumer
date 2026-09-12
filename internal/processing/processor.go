package processing

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/rphmauriciodev/myFinances-consumer/internal/dto"
)

type Processor struct {
	repo Repository
}

func NewProcessor(repo Repository) *Processor {
	return &Processor{
		repo: repo,
	}
}

func (p *Processor) ProcessTransactions(ctx context.Context, transaction dto.Transaction) error {
	transactionEntity, err := parseTransaction(transaction)
	if err != nil {
		return err
	}

	if err := p.repo.SaveTransaction(ctx, transactionEntity); err != nil {
		return err
	}
	return nil
}

func parseTransaction(dto dto.Transaction) (*Transaction, error) {

	amountValue, err := parseAmount(dto.Amount)
	if err != nil {
		return nil, err
	}

	transactionDate, err := parseDate(dto.Date)
	if err != nil {
		return nil, err
	}

	transaction := &Transaction{
		Name:            dto.Name,
		Merchant:        dto.Merchant,
		CardOrPass:      dto.CardOrPass,
		AmountValue:     amountValue,
		TransactionDate: transactionDate,
	}

	return transaction, nil
}

func parseAmount(amount string) (int64, error) {
	amount = strings.ReplaceAll(amount, "R$", "")
	amount = strings.ReplaceAll(amount, "\u00a0", "")
	amount = strings.ReplaceAll(amount, " ", "")
	amount = strings.ReplaceAll(amount, ".", "")
	amount = strings.ReplaceAll(amount, ",", "")

	value, err := strconv.ParseInt(amount, 10, 64)
	if err != nil {
		return 0, err
	}

	return value, nil
}

func parseDate(date string) (time.Time, error) {
	return time.Parse("2 Jan 2006 at 15:04", date)
}
