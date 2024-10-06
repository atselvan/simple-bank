package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides functions to execute db queries and transactions.
type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// TransferTxResult represents the result of the money transfer transaction.
type TransferTxResult struct {
	Transfer    Transfer `json:"Transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx performs a money transfer from one another to another.
// It creates a transfer record, add account entries and updates account balances
// within a single database transaction.
func (s *Store) TransferTx(ctx context.Context, input CreateTransferParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, input)
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: input.FromAccountID,
			Amount:    -input.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: input.ToAccountID,
			Amount:    input.Amount,
		})
		if err != nil {
			return err
		}

		// Update account balance
		if input.FromAccountID < input.ToAccountID {
			result.FromAccount, result.ToAccount, err = updateBalances(ctx, q, input.FromAccountID, -input.Amount, input.ToAccountID, input.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = updateBalances(ctx, q, input.ToAccountID, input.Amount, input.FromAccountID, -input.Amount)
		}

		return nil
	})

	return result, err
}

func updateBalances(ctx context.Context, q *Queries, accountId1, amount1, accountId2, amount2 int64) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountId1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountId2,
		Amount: amount2,
	})
	if err != nil {
		return
	}
	return
}

// execTx executes a function within a database transaction.
func (s *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}
