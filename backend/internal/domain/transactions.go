package domain

import (
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/models"
	"github.com/shopspring/decimal"
)

type TransactionState string

const (
	TransactionStatePending   = TransactionState("pending")
	TransactionStateConfirmed = TransactionState("confirmed")
	TransactionStateFailed    = TransactionState("failed")
)

type Transaction struct {
	Id        string
	State     TransactionState
	Timestamp int64
	PayAmount *decimal.Decimal
}

func TransactionToModel(t *Transaction) *models.Transaction {
	state := models.TransactionState(t.State)
	return &models.Transaction{
		ID:        &t.Id,
		State:     &state,
		Timestamp: &t.Timestamp,
	}
}

func TransactionFromModel(t *models.Transaction) *Transaction {
	return &Transaction{
		Id:        *t.ID,
		State:     TransactionState(*t.State),
		Timestamp: *t.Timestamp,
	}
}
