package repository

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccount_Struct(t *testing.T) {
	id := uuid.New()
	clientID := uuid.New()
	now := time.Now()

	account := Account{
		ID:        id,
		ClientID:  clientID,
		Name:      "Test Account",
		CreatedAt: pgtype.Timestamptz{Time: now, Valid: true},
	}

	assert.Equal(t, id, account.ID)
	assert.Equal(t, clientID, account.ClientID)
	assert.Equal(t, "Test Account", account.Name)
	assert.True(t, account.CreatedAt.Valid)
	assert.Equal(t, now, account.CreatedAt.Time)
}

func TestAccount_ZeroValues(t *testing.T) {
	var account Account

	assert.Equal(t, uuid.Nil, account.ID)
	assert.Equal(t, uuid.Nil, account.ClientID)
	assert.Equal(t, "", account.Name)
	assert.False(t, account.CreatedAt.Valid)
}

func TestAccount_NullCreatedAt(t *testing.T) {
	account := Account{
		ID:        uuid.New(),
		ClientID:  uuid.New(),
		Name:      "Account",
		CreatedAt: pgtype.Timestamptz{Valid: false},
	}

	assert.False(t, account.CreatedAt.Valid)
}

func TestAccount_JSONSerialization(t *testing.T) {
	id := uuid.New()
	clientID := uuid.New()
	now := time.Now()

	account := Account{
		ID:        id,
		ClientID:  clientID,
		Name:      "Test Account",
		CreatedAt: pgtype.Timestamptz{Time: now, Valid: true},
	}

	jsonData, err := json.Marshal(account)
	require.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	var decoded Account
	err = json.Unmarshal(jsonData, &decoded)
	require.NoError(t, err)

	assert.Equal(t, account.ID, decoded.ID)
	assert.Equal(t, account.ClientID, decoded.ClientID)
	assert.Equal(t, account.Name, decoded.Name)
}

func TestAccount_EmptyName(t *testing.T) {
	account := Account{
		ID:        uuid.New(),
		ClientID:  uuid.New(),
		Name:      "",
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.Equal(t, "", account.Name)
}

func TestAccount_LongName(t *testing.T) {
	longName := string(make([]byte, 1000))
	account := Account{
		ID:        uuid.New(),
		ClientID:  uuid.New(),
		Name:      longName,
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.Equal(t, longName, account.Name)
	assert.Len(t, account.Name, 1000)
}

func TestAccount_SpecialCharactersInName(t *testing.T) {
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
		account := Account{
			ID:        uuid.New(),
			ClientID:  uuid.New(),
			Name:      name,
			CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		}

		assert.Equal(t, name, account.Name)
	}
}

func TestClient_Struct(t *testing.T) {
	id := uuid.New()
	now := time.Now()

	client := Client{
		ID:        id,
		Name:      "Test Client",
		ApiKey:    "mock-api-key",
		IsActive:  boolPtr(true),
		CreatedAt: pgtype.Timestamptz{Time: now, Valid: true},
	}

	assert.Equal(t, id, client.ID)
	assert.Equal(t, "Test Client", client.Name)
	assert.Equal(t, "mock-api-key", client.ApiKey)
	assert.NotNil(t, client.IsActive)
	assert.True(t, *client.IsActive)
	assert.True(t, client.CreatedAt.Valid)
	assert.Equal(t, now, client.CreatedAt.Time)
}

func TestClient_ZeroValues(t *testing.T) {
	var client Client

	assert.Equal(t, uuid.Nil, client.ID)
	assert.Equal(t, "", client.Name)
	assert.Equal(t, "", client.ApiKey)
	assert.Nil(t, client.IsActive)
	assert.False(t, client.CreatedAt.Valid)
}

func TestClient_InactiveClient(t *testing.T) {
	client := Client{
		ID:        uuid.New(),
		Name:      "Inactive Client",
		ApiKey:    "inactive-key",
		IsActive:  boolPtr(false),
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.NotNil(t, client.IsActive)
	assert.False(t, *client.IsActive)
}

func TestClient_NullIsActive(t *testing.T) {
	client := Client{
		ID:        uuid.New(),
		Name:      "Client",
		ApiKey:    "key",
		IsActive:  nil,
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.Nil(t, client.IsActive)
}

func TestClient_NullCreatedAt(t *testing.T) {
	client := Client{
		ID:        uuid.New(),
		Name:      "Client",
		ApiKey:    "key",
		IsActive:  boolPtr(true),
		CreatedAt: pgtype.Timestamptz{Valid: false},
	}

	assert.False(t, client.CreatedAt.Valid)
}

func TestClient_JSONSerialization(t *testing.T) {
	id := uuid.New()
	now := time.Now()

	client := Client{
		ID:        id,
		Name:      "Test Client",
		ApiKey:    "test-key",
		IsActive:  boolPtr(true),
		CreatedAt: pgtype.Timestamptz{Time: now, Valid: true},
	}

	jsonData, err := json.Marshal(client)
	require.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	var decoded Client
	err = json.Unmarshal(jsonData, &decoded)
	require.NoError(t, err)

	assert.Equal(t, client.ID, decoded.ID)
	assert.Equal(t, client.Name, decoded.Name)
	assert.Equal(t, client.ApiKey, decoded.ApiKey)
}

func TestClient_EmptyName(t *testing.T) {
	client := Client{
		ID:        uuid.New(),
		Name:      "",
		ApiKey:    "key",
		IsActive:  boolPtr(true),
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.Equal(t, "", client.Name)
}

func TestClient_EmptyApiKey(t *testing.T) {
	client := Client{
		ID:        uuid.New(),
		Name:      "Client",
		ApiKey:    "",
		IsActive:  boolPtr(true),
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.Equal(t, "", client.ApiKey)
}

func TestClient_LongApiKey(t *testing.T) {
	longKey := string(make([]byte, 500))
	client := Client{
		ID:        uuid.New(),
		Name:      "Client",
		ApiKey:    longKey,
		IsActive:  boolPtr(true),
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.Equal(t, longKey, client.ApiKey)
	assert.Len(t, client.ApiKey, 500)
}

func TestClient_SpecialCharactersInApiKey(t *testing.T) {
	specialKeys := []string{
		"key-with-dashes",
		"key_with_underscores",
		"key.with.dots",
		"key$with$special",
		"mock-key-uppercase",
		"mock-key",
	}

	for _, key := range specialKeys {
		client := Client{
			ID:        uuid.New(),
			Name:      "Client",
			ApiKey:    key,
			IsActive:  boolPtr(true),
			CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		}

		assert.Equal(t, key, client.ApiKey)
	}
}

func TestClient_SpecialCharactersInName(t *testing.T) {
	specialNames := []string{
		"Client with spaces",
		"Client-with-dashes",
		"Client_with_underscores",
		"Client & Co.",
		"Client's Name",
		"客户名", // Chinese
	}

	for _, name := range specialNames {
		client := Client{
			ID:        uuid.New(),
			Name:      name,
			ApiKey:    "key",
			IsActive:  boolPtr(true),
			CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		}

		assert.Equal(t, name, client.Name)
	}
}

func TestAccount_MultipleInstances(t *testing.T) {
	accounts := []Account{
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

	assert.NotEqual(t, accounts[0].ID, accounts[1].ID)
	assert.NotEqual(t, accounts[0].ClientID, accounts[1].ClientID)
	assert.NotEqual(t, accounts[0].Name, accounts[1].Name)
}

func TestClient_MultipleInstances(t *testing.T) {
	clients := []Client{
		{
			ID:        uuid.New(),
			Name:      "Client 1",
			ApiKey:    "key1",
			IsActive:  boolPtr(true),
			CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		},
		{
			ID:        uuid.New(),
			Name:      "Client 2",
			ApiKey:    "key2",
			IsActive:  boolPtr(false),
			CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		},
	}

	assert.NotEqual(t, clients[0].ID, clients[1].ID)
	assert.NotEqual(t, clients[0].Name, clients[1].Name)
	assert.NotEqual(t, clients[0].ApiKey, clients[1].ApiKey)
	assert.NotEqual(t, *clients[0].IsActive, *clients[1].IsActive)
}

func TestAccount_SameClientID(t *testing.T) {
	clientID := uuid.New()

	account1 := Account{
		ID:        uuid.New(),
		ClientID:  clientID,
		Name:      "Account 1",
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	account2 := Account{
		ID:        uuid.New(),
		ClientID:  clientID,
		Name:      "Account 2",
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.Equal(t, account1.ClientID, account2.ClientID)
	assert.NotEqual(t, account1.ID, account2.ID)
}

func TestUUIDNil(t *testing.T) {
	var account Account
	var client Client

	assert.Equal(t, uuid.Nil, account.ID)
	assert.Equal(t, uuid.Nil, account.ClientID)
	assert.Equal(t, uuid.Nil, client.ID)
}

func TestSQLNullTypes(t *testing.T) {
	// Test pgtype.Timestamptz with nil time
	nullTime := pgtype.Timestamptz{Valid: false}
	assert.False(t, nullTime.Valid)

	// Test pgtype.Timestamptz with valid time
	now := time.Now()
	validTime := pgtype.Timestamptz{Time: now, Valid: true}
	assert.True(t, validTime.Valid)
	assert.Equal(t, now, validTime.Time)

	// Test bool pointer with false
	falseBool := boolPtr(false)
	assert.NotNil(t, falseBool)
	assert.False(t, *falseBool)

	// Test bool pointer with true
	trueBool := boolPtr(true)
	assert.NotNil(t, trueBool)
	assert.True(t, *trueBool)
}

func TestAccount_AddressIndex_JSON(t *testing.T) {
	id := uuid.New()
	clientID := uuid.New()
	idx := int32(7)
	account := Account{
		ID:           id,
		ClientID:     clientID,
		Name:         "Indexed Account",
		AddressIndex: &idx,
		CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}
	b, err := json.Marshal(account)
	require.NoError(t, err)
	assert.Contains(t, string(b), "address_index")
}

func TestPayment_Struct_JSON(t *testing.T) {
	p := Payment{
		ID:           uuid.New(),
		ClientID:     uuid.New(),
		AccountID:    uuid.New(),
		UniqueWallet: "TXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
		Status:       "PENDING",
		ExpiresAt:    pgtype.Timestamptz{Time: time.Now().Add(5 * time.Minute), Valid: true},
		CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}
	b, err := json.Marshal(p)
	require.NoError(t, err)
	assert.Contains(t, string(b), "unique_wallet")
	assert.Contains(t, string(b), "status")
	assert.Contains(t, string(b), "expires_at")
}

func TestPaymentAttempt_Struct_JSON(t *testing.T) {
	pa := PaymentAttempt{
		ID:              uuid.New(),
		PaymentID:       uuid.New(),
		AttemptNumber:   2,
		GeneratedWallet: "TYYYYYYYYYYYYYYYYYYYYYYYYYYYYY",
		GeneratedAt:     pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}
	b, err := json.Marshal(pa)
	require.NoError(t, err)
	assert.Contains(t, string(b), "attempt_number")
	assert.Contains(t, string(b), "generated_wallet")
}

func TestLog_Struct_JSON(t *testing.T) {
	msg := "event occurred"
	l := Log{
		ID:        uuid.New(),
		EventType: "ADDRESS_GENERATED",
		Message:   &msg,
		RawData:   []byte(`{"key":"value"}`),
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}
	b, err := json.Marshal(l)
	require.NoError(t, err)
	assert.Contains(t, string(b), "event_type")
	assert.Contains(t, string(b), "raw_data")
}