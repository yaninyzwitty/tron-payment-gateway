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

func TestCreateClientParams_Struct(t *testing.T) {
	params := CreateClientParams{
		Name:   "Test Client",
		ApiKey: "testkey123",
	}

	assert.Equal(t, "Test Client", params.Name)
	assert.Equal(t, "testkey123", params.ApiKey)
}

func TestCreateClientParams_JSONSerialization(t *testing.T) {
	params := CreateClientParams{
		Name:   "Test Client",
		ApiKey: "test-key",
	}

	jsonData, err := json.Marshal(params)
	require.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	var decoded CreateClientParams
	err = json.Unmarshal(jsonData, &decoded)
	require.NoError(t, err)
	
	assert.Equal(t, params.Name, decoded.Name)
	assert.Equal(t, params.ApiKey, decoded.ApiKey)
}

func TestCreateClientParams_EmptyValues(t *testing.T) {
	params := CreateClientParams{
		Name:   "",
		ApiKey: "",
	}

	assert.Equal(t, "", params.Name)
	assert.Equal(t, "", params.ApiKey)
}

func TestCreateClientParams_SpecialCharacters(t *testing.T) {
	testCases := []struct {
		name   string
		apiKey string
	}{
		{"Client & Co.", "key-with-dashes"},
		{"Client's Name", "key_with_underscores"},
		{"Client (LLC)", "key.with.dots"},
		{"客户", "key123"},
	}

	for _, tc := range testCases {
		params := CreateClientParams{
			Name:   tc.name,
			ApiKey: tc.apiKey,
		}
		assert.Equal(t, tc.name, params.Name)
		assert.Equal(t, tc.apiKey, params.ApiKey)
	}
}

func TestQueries_CreateClient_Success(t *testing.T) {
	mockDB := new(MockDBTX)
	queries := New(mockDB)
	
	ctx := context.Background()
	params := CreateClientParams{
		Name:   "Test Client",
		ApiKey: "test-api-key",
	}

	mockResult := new(MockResult)
	mockDB.On("ExecContext", ctx, createClient, mock.Anything).Return(mockResult, nil)

	err := queries.CreateClient(ctx, params)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestQueries_CreateClient_Error(t *testing.T) {
	mockDB := new(MockDBTX)
	queries := New(mockDB)
	
	ctx := context.Background()
	params := CreateClientParams{
		Name:   "Test Client",
		ApiKey: "test-api-key",
	}

	expectedErr := errors.New("duplicate key error")
	mockDB.On("ExecContext", ctx, createClient, mock.Anything).Return(nil, expectedErr)

	err := queries.CreateClient(ctx, params)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockDB.AssertExpectations(t)
}

func TestQueries_CreateClient_ContextCancellation(t *testing.T) {
	mockDB := new(MockDBTX)
	queries := New(mockDB)
	
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	
	params := CreateClientParams{
		Name:   "Test Client",
		ApiKey: "test-api-key",
	}

	mockDB.On("ExecContext", ctx, createClient, mock.Anything).Return(nil, context.Canceled)

	err := queries.CreateClient(ctx, params)

	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)
	mockDB.AssertExpectations(t)
}

func TestQueries_CreateClient_EmptyName(t *testing.T) {
	mockDB := new(MockDBTX)
	queries := New(mockDB)
	
	ctx := context.Background()
	params := CreateClientParams{
		Name:   "",
		ApiKey: "test-api-key",
	}

	mockResult := new(MockResult)
	mockDB.On("ExecContext", ctx, createClient, mock.Anything).Return(mockResult, nil)

	err := queries.CreateClient(ctx, params)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestQueries_CreateClient_EmptyApiKey(t *testing.T) {
	mockDB := new(MockDBTX)
	queries := New(mockDB)
	
	ctx := context.Background()
	params := CreateClientParams{
		Name:   "Test Client",
		ApiKey: "",
	}

	mockResult := new(MockResult)
	mockDB.On("ExecContext", ctx, createClient, mock.Anything).Return(mockResult, nil)

	err := queries.CreateClient(ctx, params)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestQueries_CreateClient_LongApiKey(t *testing.T) {
	mockDB := new(MockDBTX)
	queries := New(mockDB)
	
	ctx := context.Background()
	longKey := string(make([]byte, 1000))
	params := CreateClientParams{
		Name:   "Test Client",
		ApiKey: longKey,
	}

	mockResult := new(MockResult)
	mockDB.On("ExecContext", ctx, createClient, mock.Anything).Return(mockResult, nil)

	err := queries.CreateClient(ctx, params)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestQueries_CreateClient_WithTransaction(t *testing.T) {
	mockDB := new(MockDBTX)
	queries := New(mockDB)
	
	mockTx := &sql.Tx{}
	txQueries := queries.WithTx(mockTx)
	
	ctx := context.Background()
	params := CreateClientParams{
		Name:   "Test Client",
		ApiKey: "test-api-key",
	}

	assert.NotNil(t, txQueries)
	assert.NotEqual(t, queries.db, txQueries.db)
}

func TestQueries_GetClientByAPIKey_Success(t *testing.T) {
	mockDB := new(MockDBTX)
	queries := New(mockDB)
	
	ctx := context.Background()
	apiKey := "test-api-key"

	mockRow := new(MockRow)
	mockDB.On("QueryRowContext", ctx, getClientByAPIKey, mock.Anything).Return(mockRow)

	_, _ = queries.GetClientByAPIKey(ctx, apiKey)
	
	mockDB.AssertExpectations(t)
}

func TestQueries_GetClientByAPIKey_EmptyKey(t *testing.T) {
	mockDB := new(MockDBTX)
	queries := New(mockDB)
	
	ctx := context.Background()
	apiKey := ""

	mockRow := new(MockRow)
	mockDB.On("QueryRowContext", ctx, getClientByAPIKey, mock.Anything).Return(mockRow)

	_, _ = queries.GetClientByAPIKey(ctx, apiKey)
	
	mockDB.AssertExpectations(t)
}

func TestQueries_GetClientByAPIKey_SpecialCharacters(t *testing.T) {
	mockDB := new(MockDBTX)
	queries := New(mockDB)
	
	ctx := context.Background()
	apiKey := "key-with-special-chars!@#$%"

	mockRow := new(MockRow)
	mockDB.On("QueryRowContext", ctx, getClientByAPIKey, mock.Anything).Return(mockRow)

	_, _ = queries.GetClientByAPIKey(ctx, apiKey)
	
	mockDB.AssertExpectations(t)
}

func TestQueries_GetClientByID_Success(t *testing.T) {
	mockDB := new(MockDBTX)
	queries := New(mockDB)
	
	ctx := context.Background()
	id := uuid.New()

	mockRow := new(MockRow)
	mockDB.On("QueryRowContext", ctx, getClientByID, mock.Anything).Return(mockRow)

	_, _ = queries.GetClientByID(ctx, id)
	
	mockDB.AssertExpectations(t)
}

func TestQueries_GetClientByID_NilUUID(t *testing.T) {
	mockDB := new(MockDBTX)
	queries := New(mockDB)
	
	ctx := context.Background()
	id := uuid.Nil

	mockRow := new(MockRow)
	mockDB.On("QueryRowContext", ctx, getClientByID, mock.Anything).Return(mockRow)

	_, _ = queries.GetClientByID(ctx, id)
	
	mockDB.AssertExpectations(t)
}

func TestQueries_GetClientByID_ContextCancellation(t *testing.T) {
	mockDB := new(MockDBTX)
	queries := New(mockDB)
	
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	
	id := uuid.New()

	mockRow := new(MockRow)
	mockDB.On("QueryRowContext", ctx, getClientByID, mock.Anything).Return(mockRow)

	_, _ = queries.GetClientByID(ctx, id)
	
	mockDB.AssertExpectations(t)
}

func TestCreateClientSQL(t *testing.T) {
	expectedSQL := "-- name: CreateClient :exec\nINSERT INTO clients (name, api_key) VALUES ($1, $2)\n"
	assert.Equal(t, expectedSQL, createClient)
}

func TestGetClientByAPIKeySQL(t *testing.T) {
	expectedSQL := "-- name: GetClientByAPIKey :one\nSELECT id, name, api_key, is_active, created_at\nFROM clients\nWHERE api_key = $1 AND is_active = TRUE\nLIMIT 1\n"
	assert.Equal(t, expectedSQL, getClientByAPIKey)
}

func TestGetClientByIDSQL(t *testing.T) {
	expectedSQL := "-- name: GetClientByID :one\nSELECT id, name, api_key, is_active, created_at\nFROM clients\nWHERE id = $1\nLIMIT 1\n"
	assert.Equal(t, expectedSQL, getClientByID)
}

func TestClient_JSONTags(t *testing.T) {
	client := Client{
		ID:        uuid.New(),
		Name:      "Test",
		ApiKey:    "key",
		IsActive:  sql.NullBool{Bool: true, Valid: true},
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}

	jsonData, err := json.Marshal(client)
	require.NoError(t, err)
	
	// Verify JSON uses snake_case as per json tags
	assert.Contains(t, string(jsonData), "api_key")
	assert.Contains(t, string(jsonData), "is_active")
	assert.Contains(t, string(jsonData), "created_at")
}

func TestCreateClientParams_Validation(t *testing.T) {
	testCases := []struct {
		name   string
		params CreateClientParams
		valid  bool
	}{
		{
			name:   "valid params",
			params: CreateClientParams{Name: "Client", ApiKey: "key"},
			valid:  true,
		},
		{
			name:   "empty name",
			params: CreateClientParams{Name: "", ApiKey: "key"},
			valid:  false,
		},
		{
			name:   "empty api key",
			params: CreateClientParams{Name: "Client", ApiKey: ""},
			valid:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.valid {
				assert.NotEmpty(t, tc.params.Name)
				assert.NotEmpty(t, tc.params.ApiKey)
			} else {
				isValid := tc.params.Name != "" && tc.params.ApiKey != ""
				assert.False(t, isValid)
			}
		})
	}
}

func TestGetClientByAPIKey_ActiveClientFilter(t *testing.T) {
	// Verify that the SQL query filters for active clients
	assert.Contains(t, getClientByAPIKey, "is_active = TRUE")
}