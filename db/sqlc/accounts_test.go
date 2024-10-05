package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
)

func TestGetAccount(t *testing.T) {
	createdAccount := createTestAccount(t)

	result, err := testQueries.GetAccount(context.Background(), createdAccount.ID)
	require.NoError(t, err)
	require.Equal(t, createdAccount, result)
}

func TestCreateAccount(t *testing.T) {
	createTestAccount(t)
}

func TestUpdateAccount(t *testing.T) {
	createdAccount := createTestAccount(t)

	input := UpdateAccountParams{
		ID:      createdAccount.ID,
		Balance: getRandomAmount(),
	}

	result, err := testQueries.UpdateAccount(context.Background(), input)
	require.NoError(t, err)
	require.NoError(t, err)
	require.Equal(t, input.ID, result.ID)
	require.Equal(t, input.Balance, result.Balance)
}

func TestDeleteAccount(t *testing.T) {
	createdAccount := createTestAccount(t)

	err := testQueries.DeleteAccount(context.Background(), createdAccount.ID)
	require.NoError(t, err)

	_, err = testQueries.GetAccount(context.Background(), createdAccount.ID)
	require.Error(t, err)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createTestAccount(t)
	}

	input := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	result, err := testQueries.ListAccounts(context.Background(), input)
	require.NoError(t, err)
	require.Len(t, result, 5)
	for _, account := range result {
		require.NotEmpty(t, account)
	}
}

func createTestAccount(t *testing.T) Account {
	input := CreateAccountParams{
		Owner:    gofakeit.Name(),
		Balance:  getRandomAmount(),
		Currency: "EUR",
	}

	result, err := testQueries.CreateAccount(context.Background(), input)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, input.Owner, result.Owner)
	require.Equal(t, input.Balance, result.Balance)
	require.Equal(t, input.Currency, result.Currency)
	require.NotZero(t, result.ID)
	require.NotZero(t, result.CreatedAt)

	return result
}

func getRandomAmount() int64 {
	return int64(gofakeit.Number(100, 1000))
}
