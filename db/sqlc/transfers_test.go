package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueries_GetTransfer(t *testing.T) {
	FromAccount := createTestAccount(t)
	ToAccount := createTestAccount(t)

	createdTransfer := createTestTransfer(t, FromAccount.ID, ToAccount.ID)

	result, err := testQueries.GetTransfer(context.Background(), createdTransfer.ID)
	require.NoError(t, err)
	require.Equal(t, createdTransfer, result)
}

func TestQueries_CreateTransfer(t *testing.T) {
	FromAccount := createTestAccount(t)
	ToAccount := createTestAccount(t)

	createTestTransfer(t, FromAccount.ID, ToAccount.ID)
}

func TestQueries_ListTransfers(t *testing.T) {
	FromAccount := createTestAccount(t)
	ToAccount := createTestAccount(t)

	for i := 0; i < 10; i++ {
		createTestTransfer(t, FromAccount.ID, ToAccount.ID)
	}

	input := ListTransfersParams{
		FromAccountID: FromAccount.ID,
		ToAccountID:   ToAccount.ID,
		Limit:         5,
		Offset:        5,
	}

	result, err := testQueries.ListTransfers(context.Background(), input)
	require.NoError(t, err)
	require.Equal(t, 5, len(result))
	for _, transfer := range result {
		require.NotEmpty(t, transfer.ID)
	}
}

func createTestTransfer(t *testing.T, fromAccountId, ToAccountId int64) Transfer {
	input := CreateTransferParams{
		FromAccountID: fromAccountId,
		ToAccountID:   ToAccountId,
		Amount:        getRandomAmount(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), input)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, transfer.FromAccountID, input.FromAccountID)
	require.Equal(t, transfer.ToAccountID, input.ToAccountID)
	require.Equal(t, transfer.Amount, input.Amount)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}
