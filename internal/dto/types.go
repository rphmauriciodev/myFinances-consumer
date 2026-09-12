package dto

import "encoding/json"

type Transaction struct {
	Name       string `json:"name"`
	Merchant   string `json:"merchant"`
	CardOrPass string `json:"cardOrPass"`
	Amount     string `json:"amount"`
	Date       string `json:"date"`
}

func ParseTransaction(message string) (*Transaction, error) {
	transaction := &Transaction{}
	err := json.Unmarshal([]byte(message), transaction)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}
