package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueries_GetEntry(t *testing.T) {
	account := createTestAccount(t)
	createdEntry := createTestEntry(t, account.ID)

	result, err := testQueries.GetEntry(context.Background(), createdEntry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, result)
	require.Equal(t, createdEntry, result)
}

func TestQueries_CreateEntry(t *testing.T) {
	account := createTestAccount(t)
	createTestEntry(t, account.ID)
}

func TestQueries_ListEntries(t *testing.T) {
	account := createTestAccount(t)
	for i := 0; i < 10; i++ {
		createTestEntry(t, account.ID)
	}

	input := ListEntriesParams{
		Limit:     5,
		Offset:    5,
		AccountID: account.ID,
	}

	result, err := testQueries.ListEntries(context.Background(), input)
	require.NoError(t, err)
	require.Len(t, result, 5)
	for _, entry := range result {
		require.NotEmpty(t, entry)
	}
}

func createTestEntry(t *testing.T, accountId int64) Entry {
	input := CreateEntryParams{
		AccountID: accountId,
		Amount:    getRandomAmount(),
	}

	result, err := testQueries.CreateEntry(context.Background(), input)
	require.NoError(t, err)
	require.NotEmpty(t, result)
	require.Equal(t, input.Amount, result.Amount)
	require.Equal(t, input.AccountID, result.AccountID)
	require.NotZero(t, result.ID)
	require.NotZero(t, result.CreatedAt)

	return result
}
