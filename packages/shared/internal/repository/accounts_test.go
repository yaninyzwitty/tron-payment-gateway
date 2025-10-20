package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"
	"fmt"
	"testing"

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

	mockTx := &MockTx{}
	txQueries := queries.WithTx(mockTx)

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

	mockRow := new(MockRow)
	mockDB.On("QueryRowContext", ctx, getAccountByIDAndClientID, mock.Anything).Return(mockRow)

	// Verify params structure
	assert.Equal(t, id, params.ID)
	assert.Equal(t, clientID, params.ClientID)

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

// Tests for GetAccountByIDAndClientIDRow
func TestGetAccountByIDAndClientIDRow_Struct(t *testing.T) {
	id := uuid.New()
	clientID := uuid.New()
	now := time.Now()

	row := GetAccountByIDAndClientIDRow{
		ID:        id,
		ClientID:  clientID,
		Name:      "Test Account",
		CreatedAt: pgtype.Timestamptz{Time: now, Valid: true},
	}

	assert.Equal(t, id, row.ID)
	assert.Equal(t, clientID, row.ClientID)
	assert.Equal(t, "Test Account", row.Name)
	assert.True(t, row.CreatedAt.Valid)
	assert.Equal(t, now, row.CreatedAt.Time)
}

func TestGetAccountByIDAndClientIDRow_ZeroValues(t *testing.T) {
	var row GetAccountByIDAndClientIDRow

	assert.Equal(t, uuid.Nil, row.ID)
	assert.Equal(t, uuid.Nil, row.ClientID)
	assert.Equal(t, "", row.Name)
	assert.False(t, row.CreatedAt.Valid)
}

func TestGetAccountByIDAndClientIDRow_JSONSerialization(t *testing.T) {
	id := uuid.New()
	clientID := uuid.New()
	now := time.Now()

	row := GetAccountByIDAndClientIDRow{
		ID:        id,
		ClientID:  clientID,
		Name:      "Test Account",
		CreatedAt: pgtype.Timestamptz{Time: now, Valid: true},
	}

	jsonData, err := json.Marshal(row)
	require.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	var decoded GetAccountByIDAndClientIDRow
	err = json.Unmarshal(jsonData, &decoded)
	require.NoError(t, err)

	assert.Equal(t, row.ID, decoded.ID)
	assert.Equal(t, row.ClientID, decoded.ClientID)
	assert.Equal(t, row.Name, decoded.Name)
}

func TestGetAccountByIDAndClientIDRow_NullCreatedAt(t *testing.T) {
	row := GetAccountByIDAndClientIDRow{
		ID:        uuid.New(),
		ClientID:  uuid.New(),
		Name:      "Account",
		CreatedAt: pgtype.Timestamptz{Valid: false},
	}

	assert.False(t, row.CreatedAt.Valid)
}

func TestGetAccountByIDAndClientIDRow_EmptyName(t *testing.T) {
	row := GetAccountByIDAndClientIDRow{
		ID:        uuid.New(),
		ClientID:  uuid.New(),
		Name:      "",
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.Equal(t, "", row.Name)
}

func TestGetAccountByIDAndClientIDRow_SpecialCharactersInName(t *testing.T) {
	specialNames := []string{
		"Account with spaces",
		"Account-with-dashes",
		"Account_with_underscores",
		"Account.with.dots",
		"Account@with@symbols",
		"账户名称", // Chinese characters
		"🚀 Rocket Account",
	}

	for _, name := range specialNames {
		row := GetAccountByIDAndClientIDRow{
			ID:        uuid.New(),
			ClientID:  uuid.New(),
			Name:      name,
			CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		}

		assert.Equal(t, name, row.Name)
	}
}

func TestGetAccountByIDAndClientIDRow_LongName(t *testing.T) {
	longName := string(make([]byte, 1000))
	row := GetAccountByIDAndClientIDRow{
		ID:        uuid.New(),
		ClientID:  uuid.New(),
		Name:      longName,
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.Equal(t, longName, row.Name)
	assert.Len(t, row.Name, 1000)
}

// Tests for GetAccountsByClientIDRow
func TestGetAccountsByClientIDRow_Struct(t *testing.T) {
	id := uuid.New()
	clientID := uuid.New()
	now := time.Now()

	row := GetAccountsByClientIDRow{
		ID:        id,
		ClientID:  clientID,
		Name:      "Test Account",
		CreatedAt: pgtype.Timestamptz{Time: now, Valid: true},
	}

	assert.Equal(t, id, row.ID)
	assert.Equal(t, clientID, row.ClientID)
	assert.Equal(t, "Test Account", row.Name)
	assert.True(t, row.CreatedAt.Valid)
	assert.Equal(t, now, row.CreatedAt.Time)
}

func TestGetAccountsByClientIDRow_ZeroValues(t *testing.T) {
	var row GetAccountsByClientIDRow

	assert.Equal(t, uuid.Nil, row.ID)
	assert.Equal(t, uuid.Nil, row.ClientID)
	assert.Equal(t, "", row.Name)
	assert.False(t, row.CreatedAt.Valid)
}

func TestGetAccountsByClientIDRow_JSONSerialization(t *testing.T) {
	id := uuid.New()
	clientID := uuid.New()
	now := time.Now()

	row := GetAccountsByClientIDRow{
		ID:        id,
		ClientID:  clientID,
		Name:      "Test Account",
		CreatedAt: pgtype.Timestamptz{Time: now, Valid: true},
	}

	jsonData, err := json.Marshal(row)
	require.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	var decoded GetAccountsByClientIDRow
	err = json.Unmarshal(jsonData, &decoded)
	require.NoError(t, err)

	assert.Equal(t, row.ID, decoded.ID)
	assert.Equal(t, row.ClientID, decoded.ClientID)
	assert.Equal(t, row.Name, decoded.Name)
}

func TestGetAccountsByClientIDRow_MultipleRows(t *testing.T) {
	clientID := uuid.New()

	rows := []GetAccountsByClientIDRow{
		{
			ID:        uuid.New(),
			ClientID:  clientID,
			Name:      "Account 1",
			CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		},
		{
			ID:        uuid.New(),
			ClientID:  clientID,
			Name:      "Account 2",
			CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		},
		{
			ID:        uuid.New(),
			ClientID:  clientID,
			Name:      "Account 3",
			CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		},
	}

	// All rows should have the same clientID
	for _, row := range rows {
		assert.Equal(t, clientID, row.ClientID)
	}

	// All rows should have different IDs
	assert.NotEqual(t, rows[0].ID, rows[1].ID)
	assert.NotEqual(t, rows[0].ID, rows[2].ID)
	assert.NotEqual(t, rows[1].ID, rows[2].ID)
}

func TestGetAccountsByClientIDRow_EmptySlice(t *testing.T) {
	var rows []GetAccountsByClientIDRow

	assert.Empty(t, rows)
	assert.Len(t, rows, 0)
}

func TestGetAccountsByClientIDRow_NullCreatedAt(t *testing.T) {
	row := GetAccountsByClientIDRow{
		ID:        uuid.New(),
		ClientID:  uuid.New(),
		Name:      "Account",
		CreatedAt: pgtype.Timestamptz{Valid: false},
	}

	assert.False(t, row.CreatedAt.Valid)
}

func TestGetAccountsByClientIDRow_EmptyName(t *testing.T) {
	row := GetAccountsByClientIDRow{
		ID:        uuid.New(),
		ClientID:  uuid.New(),
		Name:      "",
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.Equal(t, "", row.Name)
}

func TestGetAccountsByClientIDRow_SpecialCharactersInName(t *testing.T) {
	specialNames := []string{
		"Account with spaces",
		"Account-with-dashes",
		"Account_with_underscores",
		"Account.with.dots",
		"Account@with@symbols",
		"账户名称", // Chinese characters
		"🚀 Rocket Account",
	}

	for _, name := range specialNames {
		row := GetAccountsByClientIDRow{
			ID:        uuid.New(),
			ClientID:  uuid.New(),
			Name:      name,
			CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		}

		assert.Equal(t, name, row.Name)
	}
}

func TestGetAccountsByClientIDRow_SortedByCreationTime(t *testing.T) {
	clientID := uuid.New()
	now := time.Now()

	rows := []GetAccountsByClientIDRow{
		{
			ID:        uuid.New(),
			ClientID:  clientID,
			Name:      "First",
			CreatedAt: pgtype.Timestamptz{Time: now, Valid: true},
		},
		{
			ID:        uuid.New(),
			ClientID:  clientID,
			Name:      "Second",
			CreatedAt: pgtype.Timestamptz{Time: now.Add(1 * time.Hour), Valid: true},
		},
		{
			ID:        uuid.New(),
			ClientID:  clientID,
			Name:      "Third",
			CreatedAt: pgtype.Timestamptz{Time: now.Add(2 * time.Hour), Valid: true},
		},
	}

	// Verify timestamps are in ascending order
	assert.True(t, rows[0].CreatedAt.Time.Before(rows[1].CreatedAt.Time))
	assert.True(t, rows[1].CreatedAt.Time.Before(rows[2].CreatedAt.Time))
}

func TestGetAccountsByClientIDRow_LargeSlice(t *testing.T) {
	clientID := uuid.New()
	var rows []GetAccountsByClientIDRow

	// Create 100 rows
	for i := 0; i < 100; i++ {
		rows = append(rows, GetAccountsByClientIDRow{
			ID:        uuid.New(),
			ClientID:  clientID,
			Name:      fmt.Sprintf("Account %d", i),
			CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		})
	}

	assert.Len(t, rows, 100)
	// Verify all have the same clientID
	for _, row := range rows {
		assert.Equal(t, clientID, row.ClientID)
	}
}

// Comparison tests between row types
func TestRowTypes_ComparisonWithAccount(t *testing.T) {
	id := uuid.New()
	clientID := uuid.New()
	name := "Test Account"
	now := time.Now()
	addressIndex := int32(5)

	account := Account{
		ID:           id,
		ClientID:     clientID,
		Name:         name,
		AddressIndex: &addressIndex,
		CreatedAt:    pgtype.Timestamptz{Time: now, Valid: true},
	}

	rowByID := GetAccountByIDAndClientIDRow{
		ID:        id,
		ClientID:  clientID,
		Name:      name,
		CreatedAt: pgtype.Timestamptz{Time: now, Valid: true},
	}

	rowByClientID := GetAccountsByClientIDRow{
		ID:        id,
		ClientID:  clientID,
		Name:      name,
		CreatedAt: pgtype.Timestamptz{Time: now, Valid: true},
	}

	// Verify common fields match
	assert.Equal(t, account.ID, rowByID.ID)
	assert.Equal(t, account.ClientID, rowByID.ClientID)
	assert.Equal(t, account.Name, rowByID.Name)

	assert.Equal(t, account.ID, rowByClientID.ID)
	assert.Equal(t, account.ClientID, rowByClientID.ClientID)
	assert.Equal(t, account.Name, rowByClientID.Name)

	// Row types should match each other
	assert.Equal(t, rowByID.ID, rowByClientID.ID)
	assert.Equal(t, rowByID.ClientID, rowByClientID.ClientID)
	assert.Equal(t, rowByID.Name, rowByClientID.Name)
}

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

// ===== Tests for updated GetAccountByIDAndClientIDRow =====

func TestGetAccountByIDAndClientIDRow_Struct(t *testing.T) {
	id := uuid.New()
	clientID := uuid.New()
	now := time.Now()

	row := GetAccountByIDAndClientIDRow{
		ID:        id,
		ClientID:  clientID,
		Name:      "Test Account",
		CreatedAt: pgtype.Timestamptz{Time: now, Valid: true},
	}

	assert.Equal(t, id, row.ID)
	assert.Equal(t, clientID, row.ClientID)
	assert.Equal(t, "Test Account", row.Name)
	assert.True(t, row.CreatedAt.Valid)
	assert.Equal(t, now, row.CreatedAt.Time)
}

func TestGetAccountByIDAndClientIDRow_ZeroValues(t *testing.T) {
	var row GetAccountByIDAndClientIDRow

	assert.Equal(t, uuid.Nil, row.ID)
	assert.Equal(t, uuid.Nil, row.ClientID)
	assert.Equal(t, "", row.Name)
	assert.False(t, row.CreatedAt.Valid)
}

func TestGetAccountByIDAndClientIDRow_JSONSerialization(t *testing.T) {
	id := uuid.New()
	clientID := uuid.New()
	now := time.Now()

	row := GetAccountByIDAndClientIDRow{
		ID:        id,
		ClientID:  clientID,
		Name:      "Account",
		CreatedAt: pgtype.Timestamptz{Time: now, Valid: true},
	}

	jsonData, err := json.Marshal(row)
	require.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	var decoded GetAccountByIDAndClientIDRow
	err = json.Unmarshal(jsonData, &decoded)
	require.NoError(t, err)

	assert.Equal(t, row.ID, decoded.ID)
	assert.Equal(t, row.ClientID, decoded.ClientID)
	assert.Equal(t, row.Name, decoded.Name)
}

// ===== Tests for updated GetAccountsByClientIDRow =====

func TestGetAccountsByClientIDRow_Struct(t *testing.T) {
	id := uuid.New()
	clientID := uuid.New()
	now := time.Now()

	row := GetAccountsByClientIDRow{
		ID:        id,
		ClientID:  clientID,
		Name:      "Test Account",
		CreatedAt: pgtype.Timestamptz{Time: now, Valid: true},
	}

	assert.Equal(t, id, row.ID)
	assert.Equal(t, clientID, row.ClientID)
	assert.Equal(t, "Test Account", row.Name)
	assert.True(t, row.CreatedAt.Valid)
	assert.Equal(t, now, row.CreatedAt.Time)
}

func TestGetAccountsByClientIDRow_ZeroValues(t *testing.T) {
	var row GetAccountsByClientIDRow

	assert.Equal(t, uuid.Nil, row.ID)
	assert.Equal(t, uuid.Nil, row.ClientID)
	assert.Equal(t, "", row.Name)
	assert.False(t, row.CreatedAt.Valid)
}

func TestGetAccountsByClientIDRow_MultipleRows(t *testing.T) {
	clientID := uuid.New()

	rows := []GetAccountsByClientIDRow{
		{
			ID:        uuid.New(),
			ClientID:  clientID,
			Name:      "Account 1",
			CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		},
		{
			ID:        uuid.New(),
			ClientID:  clientID,
			Name:      "Account 2",
			CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		},
		{
			ID:        uuid.New(),
			ClientID:  clientID,
			Name:      "Account 3",
			CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		},
	}

	for _, row := range rows {
		assert.Equal(t, clientID, row.ClientID)
		assert.NotEqual(t, uuid.Nil, row.ID)
	}

	assert.NotEqual(t, rows[0].ID, rows[1].ID)
	assert.NotEqual(t, rows[0].ID, rows[2].ID)
	assert.NotEqual(t, rows[1].ID, rows[2].ID)
}

func TestGetAccountsByClientIDRow_JSONSerialization(t *testing.T) {
	id := uuid.New()
	clientID := uuid.New()
	now := time.Now()

	row := GetAccountsByClientIDRow{
		ID:        id,
		ClientID:  clientID,
		Name:      "Account",
		CreatedAt: pgtype.Timestamptz{Time: now, Valid: true},
	}

	jsonData, err := json.Marshal(row)
	require.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	var decoded GetAccountsByClientIDRow
	err = json.Unmarshal(jsonData, &decoded)
	require.NoError(t, err)

	assert.Equal(t, row.ID, decoded.ID)
	assert.Equal(t, row.ClientID, decoded.ClientID)
	assert.Equal(t, row.Name, decoded.Name)
}

func TestGetAccountsByClientIDRow_EmptySlice(t *testing.T) {
	var rows []GetAccountsByClientIDRow

	assert.Empty(t, rows)
	assert.Len(t, rows, 0)
}

func TestGetAccountsByClientIDRow_SpecialCharacters(t *testing.T) {
	specialNames := []string{
		"Account with 🚀 emoji",
		"Account-with-dashes",
		"Account_with_underscores",
		"账户名称",
		"Àccöùñt wîth ãccênts",
	}

	for _, name := range specialNames {
		row := GetAccountsByClientIDRow{
			ID:        uuid.New(),
			ClientID:  uuid.New(),
			Name:      name,
			CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		}

		assert.Equal(t, name, row.Name)
	}
}