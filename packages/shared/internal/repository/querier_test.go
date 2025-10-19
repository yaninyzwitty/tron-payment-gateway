package repository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockQuerier is a mock implementation of Querier interface
type MockQuerier struct {
	mock.Mock
}

func (m *MockQuerier) CreateAccount(ctx context.Context, arg CreateAccountParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQuerier) CreateClient(ctx context.Context, arg CreateClientParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQuerier) GetAccountByIDAndClientID(ctx context.Context, arg GetAccountByIDAndClientIDParams) (Account, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(Account), args.Error(1)
}

func (m *MockQuerier) GetAccountsByClientID(ctx context.Context, clientID uuid.UUID) ([]Account, error) {
	args := m.Called(ctx, clientID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Account), args.Error(1)
}

func (m *MockQuerier) GetClientByAPIKey(ctx context.Context, apiKey string) (Client, error) {
	args := m.Called(ctx, apiKey)
	return args.Get(0).(Client), args.Error(1)
}

func (m *MockQuerier) GetClientByID(ctx context.Context, id uuid.UUID) (Client, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(Client), args.Error(1)
}

func TestQuerier_Interface(t *testing.T) {
	// Test that MockQuerier implements Querier interface
	var _ Querier = (*MockQuerier)(nil)
}

func TestQuerier_QueriesImplementsInterface(t *testing.T) {
	// Test that Queries struct implements Querier interface
	var _ Querier = (*Queries)(nil)
}

func TestMockQuerier_CreateAccount(t *testing.T) {
	mockQuerier := new(MockQuerier)
	ctx := context.Background()
	params := CreateAccountParams{
		ClientID: uuid.New(),
		Name:     "Test Account",
	}

	mockQuerier.On("CreateAccount", ctx, params).Return(nil)

	err := mockQuerier.CreateAccount(ctx, params)

	assert.NoError(t, err)
	mockQuerier.AssertExpectations(t)
}

func TestMockQuerier_CreateClient(t *testing.T) {
	mockQuerier := new(MockQuerier)
	ctx := context.Background()
	params := CreateClientParams{
		Name:   "Test Client",
		ApiKey: "test-key",
	}

	mockQuerier.On("CreateClient", ctx, params).Return(nil)

	err := mockQuerier.CreateClient(ctx, params)

	assert.NoError(t, err)
	mockQuerier.AssertExpectations(t)
}

func TestMockQuerier_GetAccountByIDAndClientID(t *testing.T) {
	mockQuerier := new(MockQuerier)
	ctx := context.Background()
	
	id := uuid.New()
	clientID := uuid.New()
	params := GetAccountByIDAndClientIDParams{
		ID:       id,
		ClientID: clientID,
	}

	expectedAccount := Account{
		ID:       id,
		ClientID: clientID,
		Name:     "Test Account",
	}

	mockQuerier.On("GetAccountByIDAndClientID", ctx, params).Return(expectedAccount, nil)

	account, err := mockQuerier.GetAccountByIDAndClientID(ctx, params)

	assert.NoError(t, err)
	assert.Equal(t, expectedAccount, account)
	mockQuerier.AssertExpectations(t)
}

func TestMockQuerier_GetAccountsByClientID(t *testing.T) {
	mockQuerier := new(MockQuerier)
	ctx := context.Background()
	clientID := uuid.New()

	expectedAccounts := []Account{
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
	}

	mockQuerier.On("GetAccountsByClientID", ctx, clientID).Return(expectedAccounts, nil)

	accounts, err := mockQuerier.GetAccountsByClientID(ctx, clientID)

	assert.NoError(t, err)
	assert.Equal(t, expectedAccounts, accounts)
	assert.Len(t, accounts, 2)
	mockQuerier.AssertExpectations(t)
}

func TestMockQuerier_GetClientByAPIKey(t *testing.T) {
	mockQuerier := new(MockQuerier)
	ctx := context.Background()
	apiKey := "test-api-key"

	expectedClient := Client{
		ID:     uuid.New(),
		Name:   "Test Client",
		ApiKey: apiKey,
	}

	mockQuerier.On("GetClientByAPIKey", ctx, apiKey).Return(expectedClient, nil)

	client, err := mockQuerier.GetClientByAPIKey(ctx, apiKey)

	assert.NoError(t, err)
	assert.Equal(t, expectedClient, client)
	mockQuerier.AssertExpectations(t)
}

func TestMockQuerier_GetClientByID(t *testing.T) {
	mockQuerier := new(MockQuerier)
	ctx := context.Background()
	id := uuid.New()

	expectedClient := Client{
		ID:     id,
		Name:   "Test Client",
		ApiKey: "test-key",
	}

	mockQuerier.On("GetClientByID", ctx, id).Return(expectedClient, nil)

	client, err := mockQuerier.GetClientByID(ctx, id)

	assert.NoError(t, err)
	assert.Equal(t, expectedClient, client)
	mockQuerier.AssertExpectations(t)
}

func TestMockQuerier_MultipleMethodCalls(t *testing.T) {
	mockQuerier := new(MockQuerier)
	ctx := context.Background()

	// Setup expectations for multiple calls
	clientParams := CreateClientParams{
		Name:   "Test Client",
		ApiKey: "test-key",
	}
	clientID := uuid.New()
	client := Client{
		ID:     clientID,
		Name:   clientParams.Name,
		ApiKey: clientParams.ApiKey,
	}

	accountParams := CreateAccountParams{
		ClientID: clientID,
		Name:     "Test Account",
	}

	mockQuerier.On("CreateClient", ctx, clientParams).Return(nil)
	mockQuerier.On("GetClientByAPIKey", ctx, clientParams.ApiKey).Return(client, nil)
	mockQuerier.On("CreateAccount", ctx, accountParams).Return(nil)

	// Execute
	err := mockQuerier.CreateClient(ctx, clientParams)
	assert.NoError(t, err)

	retrievedClient, err := mockQuerier.GetClientByAPIKey(ctx, clientParams.ApiKey)
	assert.NoError(t, err)
	assert.Equal(t, client, retrievedClient)

	err = mockQuerier.CreateAccount(ctx, accountParams)
	assert.NoError(t, err)

	mockQuerier.AssertExpectations(t)
}

func TestQuerier_InterfaceCoverage(t *testing.T) {
	// Verify all methods are defined in the interface
	mockDB := new(MockDBTX)
	queries := New(mockDB)

	// This test verifies compile-time that Queries implements Querier
	var querier Querier = queries
	assert.NotNil(t, querier)
}

func TestMockQuerier_GetAccountsByClientID_EmptyResult(t *testing.T) {
	mockQuerier := new(MockQuerier)
	ctx := context.Background()
	clientID := uuid.New()

	mockQuerier.On("GetAccountsByClientID", ctx, clientID).Return([]Account{}, nil)

	accounts, err := mockQuerier.GetAccountsByClientID(ctx, clientID)

	assert.NoError(t, err)
	assert.Empty(t, accounts)
	mockQuerier.AssertExpectations(t)
}

func TestMockQuerier_GetAccountsByClientID_NilResult(t *testing.T) {
	mockQuerier := new(MockQuerier)
	ctx := context.Background()
	clientID := uuid.New()

	mockQuerier.On("GetAccountsByClientID", ctx, clientID).Return(nil, nil)

	accounts, err := mockQuerier.GetAccountsByClientID(ctx, clientID)

	assert.NoError(t, err)
	assert.Nil(t, accounts)
	mockQuerier.AssertExpectations(t)
}