package db

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestGetAccount(t *testing.T) {
	createdAccount := createTestAccount(t)

	account, err := testQueries.GetAccount(context.Background(), createdAccount.ID)
	assert.NoError(t, err)
	assert.Equal(t, createdAccount, account)
}

func TestCreateAccount(t *testing.T) {
	createTestAccount(t)
}

func TestUpdateAccount(t *testing.T) {
	createdAccount := createTestAccount(t)

	input := UpdateAccountParams{
		ID:      createdAccount.ID,
		Balance: getRandomBalance(),
	}

	err := testQueries.UpdateAccount(context.Background(), input)
	assert.NoError(t, err)

	updatedAccount, err := testQueries.GetAccount(context.Background(), createdAccount.ID)
	assert.NoError(t, err)
	assert.Equal(t, input.ID, updatedAccount.ID)
	assert.Equal(t, input.Balance, updatedAccount.Balance)
}

func TestDeleteAccount(t *testing.T) {
	createdAccount := createTestAccount(t)

	err := testQueries.DeleteAccount(context.Background(), createdAccount.ID)
	assert.NoError(t, err)

	_, err = testQueries.GetAccount(context.Background(), createdAccount.ID)
	assert.Error(t, err)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createTestAccount(t)
	}

	input := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), input)
	assert.NoError(t, err)
	assert.Len(t, accounts, 5)
	for _, account := range accounts {
		assert.NotEmpty(t, account)
	}
}

func createTestAccount(t *testing.T) Account {
	input := CreateAccountParams{
		Owner:    gofakeit.Name(),
		Balance:  getRandomBalance(),
		Currency: "EUR",
	}

	account, err := testQueries.CreateAccount(context.Background(), input)
	assert.NoError(t, err)
	assert.NotEmpty(t, account)

	assert.Equal(t, input.Owner, account.Owner)
	assert.Equal(t, input.Balance, account.Balance)
	assert.Equal(t, input.Currency, account.Currency)

	assert.NotZero(t, account.ID)
	assert.NotZero(t, account.CreatedAt)

	return account
}

func getRandomBalance() int64 {
	return int64(gofakeit.Number(0, 1000))
}
