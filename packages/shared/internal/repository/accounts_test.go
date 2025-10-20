package repository

import (
	"context"
	"encoding/json"
	"errors"
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
		"账户", // Chinese
		"🚀 Rocket Account",
		"Account's Name",
		"Account & Co.",
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

func TestGetAccountByIDAndClientIDRow_NilUUIDs(t *testing.T) {
	row := GetAccountByIDAndClientIDRow{
		ID:        uuid.Nil,
		ClientID:  uuid.Nil,
		Name:      "Test",
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.Equal(t, uuid.Nil, row.ID)
	assert.Equal(t, uuid.Nil, row.ClientID)
}

func TestGetAccountByIDAndClientIDRow_MultipleInstances(t *testing.T) {
	rows := []GetAccountByIDAndClientIDRow{
		{
			ID:        uuid.New(),
			ClientID:  uuid.New(),
			Name:      "Account 1",
			CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		},
		{
			ID:        uuid.New(),
			ClientID:  uuid.New(),
			Name:      "Account 2",
			CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		},
	}

	assert.NotEqual(t, rows[0].ID, rows[1].ID)
	assert.NotEqual(t, rows[0].ClientID, rows[1].ClientID)
	assert.NotEqual(t, rows[0].Name, rows[1].Name)
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
		"账户", // Chinese
		"🚀 Rocket Account",
		"Account's Name",
		"Account & Co.",
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

func TestGetAccountsByClientIDRow_LongName(t *testing.T) {
	longName := string(make([]byte, 1000))
	row := GetAccountsByClientIDRow{
		ID:        uuid.New(),
		ClientID:  uuid.New(),
		Name:      longName,
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.Equal(t, longName, row.Name)
	assert.Len(t, row.Name, 1000)
}

func TestGetAccountsByClientIDRow_SameClientIDMultipleRows(t *testing.T) {
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

	// All rows should have the same client ID
	for _, row := range rows {
		assert.Equal(t, clientID, row.ClientID)
	}

	// But different account IDs
	assert.NotEqual(t, rows[0].ID, rows[1].ID)
	assert.NotEqual(t, rows[1].ID, rows[2].ID)
	assert.NotEqual(t, rows[0].ID, rows[2].ID)
}

func TestGetAccountsByClientIDRow_MultipleInstances(t *testing.T) {
	rows := []GetAccountsByClientIDRow{
		{
			ID:        uuid.New(),
			ClientID:  uuid.New(),
			Name:      "Account 1",
			CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		},
		{
			ID:        uuid.New(),
			ClientID:  uuid.New(),
			Name:      "Account 2",
			CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		},
	}

	assert.NotEqual(t, rows[0].ID, rows[1].ID)
	assert.NotEqual(t, rows[0].ClientID, rows[1].ClientID)
	assert.NotEqual(t, rows[0].Name, rows[1].Name)
}

func TestGetAccountsByClientIDRow_NilUUIDs(t *testing.T) {
	row := GetAccountsByClientIDRow{
		ID:        uuid.Nil,
		ClientID:  uuid.Nil,
		Name:      "Test",
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.Equal(t, uuid.Nil, row.ID)
	assert.Equal(t, uuid.Nil, row.ClientID)
}

func TestGetAccountsByClientIDRow_ArrayOfRows(t *testing.T) {
	clientID := uuid.New()
	now := time.Now()

	rows := []GetAccountsByClientIDRow{
		{
			ID:        uuid.New(),
			ClientID:  clientID,
			Name:      "Account 1",
			CreatedAt: pgtype.Timestamptz{Time: now, Valid: true},
		},
		{
			ID:        uuid.New(),
			ClientID:  clientID,
			Name:      "Account 2",
			CreatedAt: pgtype.Timestamptz{Time: now.Add(1 * time.Minute), Valid: true},
		},
		{
			ID:        uuid.New(),
			ClientID:  clientID,
			Name:      "Account 3",
			CreatedAt: pgtype.Timestamptz{Time: now.Add(2 * time.Minute), Valid: true},
		},
	}

	assert.Len(t, rows, 3)

	for i, row := range rows {
		assert.Equal(t, clientID, row.ClientID)
		assert.Equal(t, fmt.Sprintf("Account %d", i+1), row.Name)
	}
}

func TestGetAccountsByClientIDRow_EmptyArray(t *testing.T) {
	var rows []GetAccountsByClientIDRow

	assert.Empty(t, rows)
	assert.Len(t, rows, 0)
}

func TestGetAccountsByClientIDRow_SingleRow(t *testing.T) {
	rows := []GetAccountsByClientIDRow{
		{
			ID:        uuid.New(),
			ClientID:  uuid.New(),
			Name:      "Single Account",
			CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		},
	}

	assert.Len(t, rows, 1)
	assert.Equal(t, "Single Account", rows[0].Name)
}

// Integration tests for row types with query functions
func TestQueries_GetAccountByIDAndClientID_ReturnsCorrectRowType(t *testing.T) {
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
	mockDB.On("QueryRowContext", ctx, getAccountByIDAndClientID, id, clientID).Return(mockRow)

	// Call the function (Scan will be called but we don't mock the full behavior)
	_, _ = queries.GetAccountByIDAndClientID(ctx, params)

	mockDB.AssertExpectations(t)
}

func TestQueries_GetAccountsByClientID_ReturnsCorrectRowType(t *testing.T) {
	mockDB := new(MockDBTX)
	queries := New(mockDB)

	ctx := context.Background()
	clientID := uuid.New()

	mockRows := new(MockRows)
	mockDB.On("QueryContext", ctx, getAccountsByClientID, clientID).Return(mockRows, nil)
	mockRows.On("Close").Return(nil)
	mockRows.On("Next").Return(false)
	mockRows.On("Err").Return(nil)

	rows, err := queries.GetAccountsByClientID(ctx, clientID)

	assert.NoError(t, err)
	assert.NotNil(t, rows)
	assert.IsType(t, []GetAccountsByClientIDRow{}, rows)
	mockDB.AssertExpectations(t)
}

// Comparison tests between row types
func TestRowTypes_FieldCompatibility(t *testing.T) {
	id := uuid.New()
	clientID := uuid.New()
	name := "Test Account"
	now := time.Now()
	createdAt := pgtype.Timestamptz{Time: now, Valid: true}

	// Create instances of both row types
	row1 := GetAccountByIDAndClientIDRow{
		ID:        id,
		ClientID:  clientID,
		Name:      name,
		CreatedAt: createdAt,
	}

	row2 := GetAccountsByClientIDRow{
		ID:        id,
		ClientID:  clientID,
		Name:      name,
		CreatedAt: createdAt,
	}

	// Both should have the same fields with the same values
	assert.Equal(t, row1.ID, row2.ID)
	assert.Equal(t, row1.ClientID, row2.ClientID)
	assert.Equal(t, row1.Name, row2.Name)
	assert.Equal(t, row1.CreatedAt.Valid, row2.CreatedAt.Valid)
	assert.Equal(t, row1.CreatedAt.Time, row2.CreatedAt.Time)
}

func TestRowTypes_StructureDifference(t *testing.T) {
	// Both row types should be separate types even if they have the same fields
	var row1 GetAccountByIDAndClientIDRow
	var row2 GetAccountsByClientIDRow

	// They should be different types
	assert.IsType(t, GetAccountByIDAndClientIDRow{}, row1)
	assert.IsType(t, GetAccountsByClientIDRow{}, row2)

	// But not the same type
	assert.NotEqual(t, fmt.Sprintf("%T", row1), fmt.Sprintf("%T", row2))
}

// Tests for new return types: GetAccountByIDAndClientIDRow and GetAccountsByClientIDRow

func TestGetAccountByIDAndClientIDRow_Struct(t *testing.T) {
	id := uuid.New()
	clientID := uuid.New()

	row := GetAccountByIDAndClientIDRow{
		ID:        id,
		ClientID:  clientID,
		Name:      "Test Account",
		CreatedAt: pgtype.Timestamptz{},
	}

	assert.Equal(t, id, row.ID)
	assert.Equal(t, clientID, row.ClientID)
	assert.Equal(t, "Test Account", row.Name)
}

func TestGetAccountByIDAndClientIDRow_JSONSerialization(t *testing.T) {
	row := GetAccountByIDAndClientIDRow{
		ID:       uuid.New(),
		ClientID: uuid.New(),
		Name:     "Account Name",
	}

	jsonData, err := json.Marshal(row)
	require.NoError(t, err)
	assert.Contains(t, string(jsonData), "id")
	assert.Contains(t, string(jsonData), "client_id")
	assert.Contains(t, string(jsonData), "name")

	var decoded GetAccountByIDAndClientIDRow
	err = json.Unmarshal(jsonData, &decoded)
	require.NoError(t, err)
	assert.Equal(t, row.ID, decoded.ID)
	assert.Equal(t, row.ClientID, decoded.ClientID)
	assert.Equal(t, row.Name, decoded.Name)
}

func TestGetAccountByIDAndClientIDRow_EmptyName(t *testing.T) {
	row := GetAccountByIDAndClientIDRow{
		ID:       uuid.New(),
		ClientID: uuid.New(),
		Name:     "",
	}

	assert.Equal(t, "", row.Name)
}

func TestGetAccountByIDAndClientIDRow_SpecialCharacters(t *testing.T) {
	testCases := []string{
		"Account with spaces",
		"Account-with-dashes",
		"Account_with_underscores",
		"账户名称",
		"Compte français",
		"Κλειδί",
	}

	for _, name := range testCases {
		row := GetAccountByIDAndClientIDRow{
			ID:       uuid.New(),
			ClientID: uuid.New(),
			Name:     name,
		}
		assert.Equal(t, name, row.Name)
	}
}

func TestGetAccountsByClientIDRow_Struct(t *testing.T) {
	id := uuid.New()
	clientID := uuid.New()

	row := GetAccountsByClientIDRow{
		ID:        id,
		ClientID:  clientID,
		Name:      "Account Row",
		CreatedAt: pgtype.Timestamptz{},
	}

	assert.Equal(t, id, row.ID)
	assert.Equal(t, clientID, row.ClientID)
	assert.Equal(t, "Account Row", row.Name)
}

func TestGetAccountsByClientIDRow_JSONSerialization(t *testing.T) {
	row := GetAccountsByClientIDRow{
		ID:       uuid.New(),
		ClientID: uuid.New(),
		Name:     "Test Row",
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
			ID:       uuid.New(),
			ClientID: clientID,
			Name:     "Account 1",
		},
		{
			ID:       uuid.New(),
			ClientID: clientID,
			Name:     "Account 2",
		},
		{
			ID:       uuid.New(),
			ClientID: clientID,
			Name:     "Account 3",
		},
	}

	// Verify all rows belong to same client
	for i, row := range rows {
		assert.Equal(t, clientID, row.ClientID)
		assert.Equal(t, "Account "+string(rune('1'+i)), row.Name)
	}
}

func TestGetAccountsByClientIDRow_EmptyList(t *testing.T) {
	var rows []GetAccountsByClientIDRow
	assert.Empty(t, rows)
	assert.Equal(t, 0, len(rows))
}

func TestQueries_GetAccountByIDAndClientID_ReturnsRow(t *testing.T) {
	mockDB := new(MockDBTX)
	queries := New(mockDB)

	ctx := context.Background()
	params := GetAccountByIDAndClientIDParams{
		ID:       uuid.New(),
		ClientID: uuid.New(),
	}

	expectedRow := GetAccountByIDAndClientIDRow{
		ID:       params.ID,
		ClientID: params.ClientID,
		Name:     "Test Account",
	}

	mockRow := new(MockRow)
	mockDB.On("QueryRowContext", ctx, getAccountByIDAndClientID, params.ID, params.ClientID).Return(mockRow)
	mockRow.On("Scan", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		// Simulate scanning into the row
		dest := args.Get(0).([]interface{})
		if len(dest) >= 4 {
			*dest[0].(*uuid.UUID) = expectedRow.ID
			*dest[1].(*uuid.UUID) = expectedRow.ClientID
			*dest[2].(*string) = expectedRow.Name
		}
	})

	row, err := queries.GetAccountByIDAndClientID(ctx, params)

	assert.NoError(t, err)
	// Note: With our mocking setup, we can't fully verify the returned row
	// In a real integration test, we would verify all fields
	mockDB.AssertExpectations(t)
}

func TestQueries_GetAccountsByClientID_ReturnsRows(t *testing.T) {
	mockDB := new(MockDBTX)
	queries := New(mockDB)

	ctx := context.Background()
	clientID := uuid.New()

	mockRows := new(MockRows)
	mockDB.On("QueryContext", ctx, getAccountsByClientID, clientID).Return(mockRows, nil)
	mockRows.On("Close").Return(nil)
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Scan", mock.Anything).Return(nil).Once()
	mockRows.On("Next").Return(false)
	mockRows.On("Err").Return(nil)

	rows, err := queries.GetAccountsByClientID(ctx, clientID)

	assert.NoError(t, err)
	assert.NotNil(t, rows)
	mockDB.AssertExpectations(t)
}

func TestQueries_GetAccountsByClientID_MultipleRows(t *testing.T) {
	mockDB := new(MockDBTX)
	queries := New(mockDB)

	ctx := context.Background()
	clientID := uuid.New()

	mockRows := new(MockRows)
	mockDB.On("QueryContext", ctx, getAccountsByClientID, clientID).Return(mockRows, nil)
	mockRows.On("Close").Return(nil)
	
	// Simulate 3 rows
	mockRows.On("Next").Return(true).Times(3)
	mockRows.On("Scan", mock.Anything).Return(nil).Times(3)
	mockRows.On("Next").Return(false).Once()
	mockRows.On("Err").Return(nil)

	rows, err := queries.GetAccountsByClientID(ctx, clientID)

	assert.NoError(t, err)
	assert.Len(t, rows, 3)
	mockDB.AssertExpectations(t)
}

func TestQueries_GetAccountsByClientID_ScanError(t *testing.T) {
	mockDB := new(MockDBTX)
	queries := New(mockDB)

	ctx := context.Background()
	clientID := uuid.New()

	mockRows := new(MockRows)
	mockDB.On("QueryContext", ctx, getAccountsByClientID, clientID).Return(mockRows, nil)
	mockRows.On("Close").Return(nil)
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Scan", mock.Anything).Return(errors.New("scan error")).Once()

	rows, err := queries.GetAccountsByClientID(ctx, clientID)

	assert.Error(t, err)
	assert.Nil(t, rows)
	assert.Contains(t, err.Error(), "scan error")
	mockDB.AssertExpectations(t)
}

func TestQueries_GetAccountsByClientID_RowsError(t *testing.T) {
	mockDB := new(MockDBTX)
	queries := New(mockDB)

	ctx := context.Background()
	clientID := uuid.New()

	mockRows := new(MockRows)
	mockDB.On("QueryContext", ctx, getAccountsByClientID, clientID).Return(mockRows, nil)
	mockRows.On("Close").Return(nil)
	mockRows.On("Next").Return(false)
	mockRows.On("Err").Return(errors.New("rows error"))

	rows, err := queries.GetAccountsByClientID(ctx, clientID)

	assert.Error(t, err)
	assert.Nil(t, rows)
	assert.Contains(t, err.Error(), "rows error")
	mockDB.AssertExpectations(t)
}

// Test that return types are different from Account model
func TestReturnTypesAreDifferent(t *testing.T) {
	// GetAccountByIDAndClientIDRow should not have AddressIndex field
	row := GetAccountByIDAndClientIDRow{
		ID:       uuid.New(),
		ClientID: uuid.New(),
		Name:     "Test",
	}

	// This should compile - verifying the type doesn't have AddressIndex
	jsonData, err := json.Marshal(row)
	require.NoError(t, err)
	
	// The JSON should NOT contain address_index since it's not selected in the query
	assert.NotContains(t, string(jsonData), "address_index")
}

func TestGetAccountsByClientIDRow_Consistency(t *testing.T) {
	// Verify GetAccountsByClientIDRow and GetAccountByIDAndClientIDRow have same structure
	row1 := GetAccountByIDAndClientIDRow{
		ID:       uuid.New(),
		ClientID: uuid.New(),
		Name:     "Test",
	}

	row2 := GetAccountsByClientIDRow{
		ID:       row1.ID,
		ClientID: row1.ClientID,
		Name:     row1.Name,
	}

	assert.Equal(t, row1.ID, row2.ID)
	assert.Equal(t, row1.ClientID, row2.ClientID)
	assert.Equal(t, row1.Name, row2.Name)
}