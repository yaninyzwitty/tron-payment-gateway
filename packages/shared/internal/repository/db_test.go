package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	mockDB := new(MockDBTX)

	queries := New(mockDB)

	assert.NotNil(t, queries)
	assert.NotNil(t, queries.db)
	assert.Equal(t, mockDB, queries.db)
}

func TestNew_NilDB(t *testing.T) {
	queries := New(nil)

	assert.NotNil(t, queries)
	assert.Nil(t, queries.db)
}

func TestQueries_WithTx(t *testing.T) {
	mockDB := new(MockDBTX)
	queries := New(mockDB)

	// Create a mock transaction
	mockTx := &MockTx{}

	txQueries := queries.WithTx(mockTx)

	assert.NotNil(t, txQueries)
	assert.NotNil(t, txQueries.db)
	assert.Equal(t, mockTx, txQueries.db)

	// Original queries should still have the original DB
	assert.Equal(t, mockDB, queries.db)
}

func TestQueries_WithTx_NilTransaction(t *testing.T) {
	mockDB := new(MockDBTX)
	queries := New(mockDB)

	txQueries := queries.WithTx(nil)

	assert.NotNil(t, txQueries)
	assert.Nil(t, txQueries.db)
}

func TestQueries_WithTx_ChainedCalls(t *testing.T) {
	mockDB := new(MockDBTX)
	queries := New(mockDB)

	tx1 := &MockTx{}
	tx2 := &MockTx{}

	txQueries1 := queries.WithTx(tx1)
	txQueries2 := txQueries1.WithTx(tx2)

	assert.Equal(t, tx1, txQueries1.db)
	assert.Equal(t, tx2, txQueries2.db)
	assert.Equal(t, mockDB, queries.db)
}

func TestQueries_WithTx_MultipleTransactions(t *testing.T) {
	mockDB := new(MockDBTX)
	queries := New(mockDB)

	tx1 := &MockTx{}
	tx2 := &MockTx{}

	txQueries1 := queries.WithTx(tx1)
	txQueries2 := queries.WithTx(tx2)

	// Both should have different transactions
	assert.NotEqual(t, txQueries1.db, txQueries2.db)
	assert.Equal(t, tx1, txQueries1.db)
	assert.Equal(t, tx2, txQueries2.db)
}

func TestDBTX_Interface(t *testing.T) {
	// Test that MockDBTX implements DBTX interface
	var _ DBTX = (*MockDBTX)(nil)
}

func TestQueries_DBAccess(t *testing.T) {
	mockDB := new(MockDBTX)
	queries := New(mockDB)

	// Verify that the db field is accessible through methods
	assert.NotNil(t, queries)

	// Test that we can create a new queries with transaction
	tx := &MockTx{}
	txQueries := queries.WithTx(tx)
	assert.NotNil(t, txQueries)
}

func TestNew_MultipleInstances(t *testing.T) {
	mockDB1 := new(MockDBTX)
	mockDB2 := new(MockDBTX)

	queries1 := New(mockDB1)
	queries2 := New(mockDB2)

	assert.NotEqual(t, queries1, queries2)
	assert.Equal(t, mockDB1, queries1.db)
	assert.Equal(t, mockDB2, queries2.db)
}

func TestQueries_IsolationBetweenInstances(t *testing.T) {
	mockDB := new(MockDBTX)
	queries1 := New(mockDB)
	queries2 := New(mockDB)

	tx := &MockTx{}
	queries1WithTx := queries1.WithTx(tx)

	// queries2 should still have the original DB
	assert.Equal(t, mockDB, queries2.db)
	assert.Equal(t, tx, queries1WithTx.db)
}
