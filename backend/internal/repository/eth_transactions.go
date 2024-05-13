package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/internal/domain"
	"github.com/jackc/pgx/v5"
)

type EthTransactionsRepo struct {
}

func NewEthTransactionsRepo() EthTransactions {
	return &EthTransactionsRepo{}
}

func (r *EthTransactionsRepo) GetTransaction(ctx context.Context, transaction Transaction, id string) (*domain.Transaction, error) {
	tx, ok := transaction.(pgx.Tx)
	if !ok {
		return nil, errors.New("GetTransaction: error: type assertion failed on interface Transaction")
	}

	var ethTx domain.Transaction
	row := tx.QueryRow(ctx, "SELECT id, state, timestamp FROM transactions WHERE id = $1", id)
	if err := row.Scan(&ethTx.Id, &ethTx.State, &ethTx.Timestamp); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRows
		}
		return nil, fmt.Errorf("GetTransaction/tx.QueryRow: %w", err)
	}

	return &ethTx, nil
}

func (r *EthTransactionsRepo) InsertTransaction(ctx context.Context, transaction Transaction, ethTx *domain.Transaction) error {
	tx, ok := transaction.(pgx.Tx)
	if !ok {
		return errors.New("InsertTransaction: error: type assertion failed on interface Transaction")
	}

	if _, err := tx.Exec(ctx, "INSERT INTO transactions (id, state, timestamp) VALUES ($1, $2, $3)",
		ethTx.Id, ethTx.State, ethTx.Timestamp); err != nil {
		return fmt.Errorf("InsertTransaction/tx.Exec: %w", err)
	}

	return nil
}
