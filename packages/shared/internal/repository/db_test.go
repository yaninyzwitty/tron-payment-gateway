package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDBTX is a mock implementation of DBTX interface
type MockDBTX struct {
	mock.Mock
}

func (m *MockDBTX) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	argsMock := m.Called(ctx, query, args)
	if argsMock.Get(0) == nil {
		return nil, argsMock.Error(1)
	}
	return argsMock.Get(0).(sql.Result), argsMock.Error(1)
}

func (m *MockDBTX) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	argsMock := m.Called(ctx, query)
	if argsMock.Get(0) == nil {
		return nil, argsMock.Error(1)
	}
	return argsMock.Get(0).(*sql.Stmt), argsMock.Error(1)
}

func (m *MockDBTX) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	argsMock := m.Called(ctx, query, args)
	if argsMock.Get(0) == nil {
		return nil, argsMock.Error(1)
	}
	return argsMock.Get(0).(*sql.Rows), argsMock.Error(1)
}

func (m *MockDBTX) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	argsMock := m.Called(ctx, query, args)
	return argsMock.Get(0).(*sql.Row)
}

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
	mockTx := &sql.Tx{}
	
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
	
	tx1 := &sql.Tx{}
	tx2 := &sql.Tx{}
	
	txQueries1 := queries.WithTx(tx1)
	txQueries2 := txQueries1.WithTx(tx2)
	
	assert.Equal(t, tx1, txQueries1.db)
	assert.Equal(t, tx2, txQueries2.db)
	assert.Equal(t, mockDB, queries.db)
}

func TestQueries_WithTx_MultipleTransactions(t *testing.T) {
	mockDB := new(MockDBTX)
	queries := New(mockDB)
	
	tx1 := &sql.Tx{}
	tx2 := &sql.Tx{}
	
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
	tx := &sql.Tx{}
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
	
	tx := &sql.Tx{}
	queries1WithTx := queries1.WithTx(tx)
	
	// queries2 should still have the original DB
	assert.Equal(t, mockDB, queries2.db)
	assert.Equal(t, tx, queries1WithTx.db)
}