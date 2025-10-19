package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateAccountParams_Struct(t *testing.T) {
	clientID := uuid.New()
	params := CreateAccountParams{
		ClientID: clientID,
		Name:     "Test Account",
	}

	assert.Equal(t, clientID, params.ClientID)
	assert.Equal(t, "Test Account", params.Name)
}

func TestCreateAccountParams_JSONSerialization(t *testing.T) {
	clientID := uuid.New()
	params := CreateAccountParams{
		ClientID: clientID,
		Name:     "Test Account",
	}

	jsonData, err := json.Marshal(params)
	require.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	var decoded CreateAccountParams
	err = json.Unmarshal(jsonData, &decoded)
	require.NoError(t, err)
	
	assert.Equal(t, params.ClientID, decoded.ClientID)
	assert.Equal(t, params.Name, decoded.Name)
}

func TestCreateAccountParams_EmptyName(t *testing.T) {
	params := CreateAccountParams{
		ClientID: uuid.New(),
		Name:     "",
	}

	assert.Equal(t, "", params.Name)
}

func TestCreateAccountParams_SpecialCharacters(t *testing.T) {
	specialNames := []string{
		"Account with spaces",
		"Account-with-dashes",
		"Account_with_underscores",
		"账户", // Chinese
		"🚀 Rocket",
	}

	for _, name := range specialNames {
		params := CreateAccountParams{
			ClientID: uuid.New(),
			Name:     name,
		}
		assert.Equal(t, name, params.Name)
	}
}

func TestGetAccountByIDAndClientIDParams_Struct(t *testing.T) {
	id := uuid.New()
	clientID := uuid.New()
	
	params := GetAccountByIDAndClientIDParams{
		ID:       id,
		ClientID: clientID,
	}

	assert.Equal(t, id, params.ID)
	assert.Equal(t, clientID, params.ClientID)
}

func TestGetAccountByIDAndClientIDParams_JSONSerialization(t *testing.T) {
	id := uuid.New()
	clientID := uuid.New()
	
	params := GetAccountByIDAndClientIDParams{
		ID:       id,
		ClientID: clientID,
	}

	jsonData, err := json.Marshal(params)
	require.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	var decoded GetAccountByIDAndClientIDParams
	err = json.Unmarshal(jsonData, &decoded)
	require.NoError(t, err)
	
	assert.Equal(t, params.ID, decoded.ID)
	assert.Equal(t, params.ClientID, decoded.ClientID)
}

func TestGetAccountByIDAndClientIDParams_NilUUIDs(t *testing.T) {
	params := GetAccountByIDAndClientIDParams{
		ID:       uuid.Nil,
		ClientID: uuid.Nil,
	}

	assert.Equal(t, uuid.Nil, params.ID)
	assert.Equal(t, uuid.Nil, params.ClientID)
}

func TestQueries_CreateAccount_Success(t *testing.T) {
	mockDB := new(MockDBTX)
	queries := New(mockDB)
	
	ctx := context.Background()
	params := CreateAccountParams{
		ClientID: uuid.New(),
		Name:     "Test Account",
	}

	mockResult := new(MockResult)
	mockDB.On("ExecContext", ctx, createAccount, mock.Anything).Return(mockResult, nil)

	err := queries.CreateAccount(ctx, params)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestQueries_CreateAccount_Error(t *testing.T) {
	mockDB := new(MockDBTX)
	queries := New(mockDB)
	
	ctx := context.Background()
	params := CreateAccountParams{
		ClientID: uuid.New(),
		Name:     "Test Account",
	}

	expectedErr := errors.New("database error")
	mockDB.On("ExecContext", ctx, createAccount, mock.Anything).Return(nil, expectedErr)

	err := queries.CreateAccount(ctx, params)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockDB.AssertExpectations(t)
}

func TestQueries_CreateAccount_ContextCancellation(t *testing.T) {
	mockDB := new(MockDBTX)
	queries := New(mockDB)
	
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	
	params := CreateAccountParams{
		ClientID: uuid.New(),
		Name:     "Test Account",
	}

	mockDB.On("ExecContext", ctx, createAccount, mock.Anything).Return(nil, context.Canceled)

	err := queries.CreateAccount(ctx, params)

	assert.Error(t, err)
	mockDB.AssertExpectations(t)
}

func TestQueries_CreateAccount_EmptyName(t *testing.T) {
	mockDB := new(MockDBTX)
	queries := New(mockDB)
	
	ctx := context.Background()
	params := CreateAccountParams{
		ClientID: uuid.New(),
		Name:     "",
	}

	mockResult := new(MockResult)
	mockDB.On("ExecContext", ctx, createAccount, mock.Anything).Return(mockResult, nil)

	err := queries.CreateAccount(ctx, params)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestQueries_CreateAccount_WithTransaction(t *testing.T) {
	mockDB := new(MockDBTX)
	queries := New(mockDB)
	
	mockTx := &sql.Tx{}
	txQueries := queries.WithTx(mockTx)
	
	ctx := context.Background()
	params := CreateAccountParams{
		ClientID: uuid.New(),
		Name:     "Test Account",
	}

	// The transaction should be used, not the original DB
	assert.NotNil(t, txQueries)
	assert.NotEqual(t, queries.db, txQueries.db)
}

func TestQueries_GetAccountByIDAndClientID_Success(t *testing.T) {
	mockDB := new(MockDBTX)
	queries := New(mockDB)
	
	ctx := context.Background()
	id := uuid.New()
	clientID := uuid.New()
	
	params := GetAccountByIDAndClientIDParams{
		ID:       id,
		ClientID: clientID,
	}

	expectedAccount := Account{
		ID:        id,
		ClientID:  clientID,
		Name:      "Test Account",
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}

	mockRow := new(MockRow)
	mockDB.On("QueryRowContext", ctx, getAccountByIDAndClientID, mock.Anything).Return(mockRow)

	// Note: We can't fully test Scan behavior without a real database or more complex mocking
	// This tests the method call structure
	_, _ = queries.GetAccountByIDAndClientID(ctx, params)
	
	mockDB.AssertExpectations(t)
}

func TestQueries_GetAccountsByClientID_EmptyResult(t *testing.T) {
	mockDB := new(MockDBTX)
	queries := New(mockDB)
	
	ctx := context.Background()
	clientID := uuid.New()

	mockRows := new(MockRows)
	mockDB.On("QueryContext", ctx, getAccountsByClientID, mock.Anything).Return(mockRows, nil)
	mockRows.On("Close").Return(nil)
	mockRows.On("Next").Return(false)
	mockRows.On("Err").Return(nil)

	accounts, err := queries.GetAccountsByClientID(ctx, clientID)

	assert.NoError(t, err)
	assert.Empty(t, accounts)
	mockDB.AssertExpectations(t)
}

func TestQueries_GetAccountsByClientID_Error(t *testing.T) {
	mockDB := new(MockDBTX)
	queries := New(mockDB)
	
	ctx := context.Background()
	clientID := uuid.New()

	expectedErr := errors.New("query error")
	mockDB.On("QueryContext", ctx, getAccountsByClientID, mock.Anything).Return(nil, expectedErr)

	accounts, err := queries.GetAccountsByClientID(ctx, clientID)

	assert.Error(t, err)
	assert.Nil(t, accounts)
	assert.Equal(t, expectedErr, err)
	mockDB.AssertExpectations(t)
}

func TestQueries_GetAccountsByClientID_NilClientID(t *testing.T) {
	mockDB := new(MockDBTX)
	queries := New(mockDB)
	
	ctx := context.Background()
	clientID := uuid.Nil

	mockRows := new(MockRows)
	mockDB.On("QueryContext", ctx, getAccountsByClientID, mock.Anything).Return(mockRows, nil)
	mockRows.On("Close").Return(nil)
	mockRows.On("Next").Return(false)
	mockRows.On("Err").Return(nil)

	accounts, err := queries.GetAccountsByClientID(ctx, clientID)

	assert.NoError(t, err)
	assert.Empty(t, accounts)
	mockDB.AssertExpectations(t)
}

func TestCreateAccountSQL(t *testing.T) {
	expectedSQL := "-- name: CreateAccount :exec\nINSERT INTO accounts (client_id, name) VALUES ($1, $2)\n"
	assert.Equal(t, expectedSQL, createAccount)
}

func TestGetAccountByIDAndClientIDSQL(t *testing.T) {
	expectedSQL := "-- name: GetAccountByIDAndClientID :one\nSELECT id, client_id, name, created_at\nFROM accounts\nWHERE id = $1 AND client_id = $2\nLIMIT 1\n"
	assert.Equal(t, expectedSQL, getAccountByIDAndClientID)
}

func TestGetAccountsByClientIDSQL(t *testing.T) {
	expectedSQL := "-- name: GetAccountsByClientID :many\nSELECT id, client_id, name, created_at\nFROM accounts\nWHERE client_id = $1\n"
	assert.Equal(t, expectedSQL, getAccountsByClientID)
}

// Mock helpers
type MockResult struct {
	mock.Mock
}

func (m *MockResult) LastInsertId() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockResult) RowsAffected() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

type MockRow struct {
	mock.Mock
}

func (m *MockRow) Scan(dest ...interface{}) error {
	args := m.Called(dest)
	return args.Error(0)
}

type MockRows struct {
	mock.Mock
}

func (m *MockRows) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRows) Next() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockRows) Scan(dest ...interface{}) error {
	args := m.Called(dest)
	return args.Error(0)
}

func (m *MockRows) Err() error {
	args := m.Called()
	return args.Error(0)
}