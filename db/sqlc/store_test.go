package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStore_TransferTx(t *testing.T) {
	fromAccount := createTestAccount(t)
	toAccount := createTestAccount(t)

	t.Logf("Balances before all transactions: fromAccount=%d, toAccount=%d", fromAccount.Balance, toAccount.Balance)

	// run n concurrent transfers
	n := 5
	amount := int64(10)
	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		input := CreateTransferParams{
			FromAccountID: fromAccount.ID,
			ToAccountID:   toAccount.ID,
			Amount:        amount,
		}

		go func() {
			result, err := testStore.TransferTx(context.Background(), input)

			errs <- err
			results <- result
		}()
	}

	existed := map[int]bool{}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result.Transfer)

		// check transfer
		require.NotEmpty(t, result.Transfer)
		require.Equal(t, fromAccount.ID, result.Transfer.FromAccountID)
		require.Equal(t, toAccount.ID, result.Transfer.ToAccountID)
		require.Equal(t, amount, result.Transfer.Amount)
		require.NotZero(t, result.Transfer.ID)
		require.NotZero(t, result.Transfer.CreatedAt)

		_, err = testStore.GetTransfer(context.Background(), result.Transfer.ID)
		require.NoError(t, err)

		// check from entry
		require.NotEmpty(t, result.FromEntry)
		require.Equal(t, fromAccount.ID, result.FromEntry.AccountID)
		require.Equal(t, -amount, result.FromEntry.Amount)
		require.NotZero(t, result.FromEntry.ID)
		require.NotZero(t, result.FromEntry.CreatedAt)

		_, err = testStore.GetEntry(context.Background(), result.FromEntry.ID)
		require.NoError(t, err)

		// check to entry
		require.NotEmpty(t, result.ToEntry)
		require.Equal(t, toAccount.ID, result.ToEntry.AccountID)
		require.Equal(t, amount, result.ToEntry.Amount)
		require.NotZero(t, result.ToEntry.ID)
		require.NotZero(t, result.ToEntry.CreatedAt)

		_, err = testStore.GetEntry(context.Background(), result.ToEntry.ID)
		require.NoError(t, err)

		// check from account
		require.NotEmpty(t, result.FromAccount)
		require.Equal(t, fromAccount.ID, result.FromAccount.ID)

		// check to account
		require.NotEmpty(t, result.FromAccount)
		require.Equal(t, fromAccount.ID, result.FromAccount.ID)

		// check account balances
		t.Logf("Balances after a transaction: fromAccount=%d, toAccount=%d", result.FromAccount.Balance, result.ToAccount.Balance)
		diffFromAccountBalance := fromAccount.Balance - result.FromAccount.Balance
		diffToAccountBalance := result.ToAccount.Balance - toAccount.Balance
		require.Equal(t, diffFromAccountBalance, diffToAccountBalance)
		require.True(t, diffFromAccountBalance > 0)
		require.True(t, diffFromAccountBalance%amount == 0)

		k := int(diffFromAccountBalance / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated balances
	updatedFromAccount, err := testStore.GetAccount(context.Background(), fromAccount.ID)
	require.NoError(t, err)

	updatedToAccount, err := testStore.GetAccount(context.Background(), toAccount.ID)
	require.NoError(t, err)

	t.Logf("Balances after all transactions: fromAccount=%d, toAccount=%d", updatedFromAccount.Balance, updatedFromAccount.Balance)

	require.Equal(t, fromAccount.Balance-int64(n)*amount, updatedFromAccount.Balance)
	require.Equal(t, toAccount.Balance+int64(n)*amount, updatedToAccount.Balance)
}

func TestStore_TransferTxDeadlock(t *testing.T) {
	account1 := createTestAccount(t)
	account2 := createTestAccount(t)

	t.Logf("Balances before all transactions: fromAccount=%d, toAccount=%d", account1.Balance, account2.Balance)

	// run n concurrent transfers
	n := 10
	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountId := account1.ID
		toAccountId := account2.ID

		if i%2 == 1 {
			fromAccountId = account2.ID
			toAccountId = account1.ID
		}

		input := CreateTransferParams{
			FromAccountID: fromAccountId,
			ToAccountID:   toAccountId,
			Amount:        amount,
		}

		go func() {
			_, err := testStore.TransferTx(context.Background(), input)
			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// check the final updated balances
	updatedFromAccount, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedToAccount, err := testStore.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	t.Logf("Balances after all transactions: fromAccount=%d, toAccount=%d", updatedFromAccount.Balance, updatedToAccount.Balance)

	require.Equal(t, account1.Balance, updatedFromAccount.Balance)
	require.Equal(t, account2.Balance, updatedToAccount.Balance)
}
