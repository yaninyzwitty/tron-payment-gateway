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

// Tests for Account.AddressIndex field (new field added)
func TestAccount_WithAddressIndex(t *testing.T) {
	id := uuid.New()
	clientID := uuid.New()
	now := time.Now()
	addressIndex := int32(5)

	account := Account{
		ID:           id,
		ClientID:     clientID,
		Name:         "Test Account",
		AddressIndex: &addressIndex,
		CreatedAt:    pgtype.Timestamptz{Time: now, Valid: true},
	}

	assert.Equal(t, id, account.ID)
	assert.Equal(t, clientID, account.ClientID)
	assert.Equal(t, "Test Account", account.Name)
	assert.NotNil(t, account.AddressIndex)
	assert.Equal(t, int32(5), *account.AddressIndex)
	assert.True(t, account.CreatedAt.Valid)
}

func TestAccount_NullAddressIndex(t *testing.T) {
	account := Account{
		ID:           uuid.New(),
		ClientID:     uuid.New(),
		Name:         "Account",
		AddressIndex: nil,
		CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.Nil(t, account.AddressIndex)
}

func TestAccount_ZeroAddressIndex(t *testing.T) {
	addressIndex := int32(0)
	account := Account{
		ID:           uuid.New(),
		ClientID:     uuid.New(),
		Name:         "Account",
		AddressIndex: &addressIndex,
		CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.NotNil(t, account.AddressIndex)
	assert.Equal(t, int32(0), *account.AddressIndex)
}

func TestAccount_NegativeAddressIndex(t *testing.T) {
	addressIndex := int32(-1)
	account := Account{
		ID:           uuid.New(),
		ClientID:     uuid.New(),
		Name:         "Account",
		AddressIndex: &addressIndex,
		CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.NotNil(t, account.AddressIndex)
	assert.Equal(t, int32(-1), *account.AddressIndex)
}

func TestAccount_LargeAddressIndex(t *testing.T) {
	addressIndex := int32(2147483647) // max int32
	account := Account{
		ID:           uuid.New(),
		ClientID:     uuid.New(),
		Name:         "Account",
		AddressIndex: &addressIndex,
		CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.NotNil(t, account.AddressIndex)
	assert.Equal(t, int32(2147483647), *account.AddressIndex)
}

func TestAccount_JSONSerializationWithAddressIndex(t *testing.T) {
	id := uuid.New()
	clientID := uuid.New()
	now := time.Now()
	addressIndex := int32(10)

	account := Account{
		ID:           id,
		ClientID:     clientID,
		Name:         "Test Account",
		AddressIndex: &addressIndex,
		CreatedAt:    pgtype.Timestamptz{Time: now, Valid: true},
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
	assert.NotNil(t, decoded.AddressIndex)
	assert.Equal(t, *account.AddressIndex, *decoded.AddressIndex)
}

// Tests for Log model
func TestLog_Struct(t *testing.T) {
	id := uuid.New()
	paymentID := uuid.New()
	now := time.Now()
	message := "Payment initiated"
	rawData := []byte(`{"transaction": "tx123"}`)

	log := Log{
		ID:        id,
		PaymentID: pgtype.UUID{Bytes: paymentID, Valid: true},
		EventType: "payment.initiated",
		Message:   &message,
		RawData:   rawData,
		CreatedAt: pgtype.Timestamptz{Time: now, Valid: true},
	}

	assert.Equal(t, id, log.ID)
	assert.True(t, log.PaymentID.Valid)
	assert.Equal(t, paymentID, log.PaymentID.Bytes)
	assert.Equal(t, "payment.initiated", log.EventType)
	assert.NotNil(t, log.Message)
	assert.Equal(t, "Payment initiated", *log.Message)
	assert.Equal(t, rawData, log.RawData)
	assert.True(t, log.CreatedAt.Valid)
}

func TestLog_ZeroValues(t *testing.T) {
	var log Log

	assert.Equal(t, uuid.Nil, log.ID)
	assert.False(t, log.PaymentID.Valid)
	assert.Equal(t, "", log.EventType)
	assert.Nil(t, log.Message)
	assert.Nil(t, log.RawData)
	assert.False(t, log.CreatedAt.Valid)
}

func TestLog_NullPaymentID(t *testing.T) {
	log := Log{
		ID:        uuid.New(),
		PaymentID: pgtype.UUID{Valid: false},
		EventType: "system.event",
		Message:   stringPtr("System message"),
		RawData:   []byte(`{}`),
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.False(t, log.PaymentID.Valid)
}

func TestLog_NullMessage(t *testing.T) {
	log := Log{
		ID:        uuid.New(),
		PaymentID: pgtype.UUID{Bytes: uuid.New(), Valid: true},
		EventType: "payment.completed",
		Message:   nil,
		RawData:   []byte(`{}`),
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.Nil(t, log.Message)
}

func TestLog_EmptyRawData(t *testing.T) {
	log := Log{
		ID:        uuid.New(),
		PaymentID: pgtype.UUID{Bytes: uuid.New(), Valid: true},
		EventType: "payment.failed",
		Message:   stringPtr("Payment failed"),
		RawData:   []byte{},
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.NotNil(t, log.RawData)
	assert.Empty(t, log.RawData)
}

func TestLog_NilRawData(t *testing.T) {
	log := Log{
		ID:        uuid.New(),
		PaymentID: pgtype.UUID{Bytes: uuid.New(), Valid: true},
		EventType: "payment.pending",
		Message:   stringPtr("Awaiting confirmation"),
		RawData:   nil,
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.Nil(t, log.RawData)
}

func TestLog_JSONSerialization(t *testing.T) {
	id := uuid.New()
	paymentID := uuid.New()
	now := time.Now()
	message := "Test message"
	rawData := []byte(`{"key": "value"}`)

	log := Log{
		ID:        id,
		PaymentID: pgtype.UUID{Bytes: paymentID, Valid: true},
		EventType: "test.event",
		Message:   &message,
		RawData:   rawData,
		CreatedAt: pgtype.Timestamptz{Time: now, Valid: true},
	}

	jsonData, err := json.Marshal(log)
	require.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	var decoded Log
	err = json.Unmarshal(jsonData, &decoded)
	require.NoError(t, err)

	assert.Equal(t, log.ID, decoded.ID)
	assert.Equal(t, log.EventType, decoded.EventType)
}

func TestLog_DifferentEventTypes(t *testing.T) {
	eventTypes := []string{
		"payment.created",
		"payment.initiated",
		"payment.confirmed",
		"payment.completed",
		"payment.failed",
		"payment.expired",
		"wallet.generated",
		"system.error",
	}

	for _, eventType := range eventTypes {
		log := Log{
			ID:        uuid.New(),
			PaymentID: pgtype.UUID{Bytes: uuid.New(), Valid: true},
			EventType: eventType,
			Message:   stringPtr("Test"),
			RawData:   []byte(`{}`),
			CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		}

		assert.Equal(t, eventType, log.EventType)
	}
}

func TestLog_LargeRawData(t *testing.T) {
	largeData := make([]byte, 10000)
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	log := Log{
		ID:        uuid.New(),
		PaymentID: pgtype.UUID{Bytes: uuid.New(), Valid: true},
		EventType: "test.large",
		Message:   stringPtr("Large data test"),
		RawData:   largeData,
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.Equal(t, largeData, log.RawData)
	assert.Len(t, log.RawData, 10000)
}

// Tests for Payment model
func TestPayment_Struct(t *testing.T) {
	id := uuid.New()
	clientID := uuid.New()
	accountID := uuid.New()
	now := time.Now()
	attemptCount := int32(1)

	payment := Payment{
		ID:           id,
		ClientID:     clientID,
		AccountID:    accountID,
		Amount:       pgtype.Numeric{Int: big.NewInt(100000), Exp: -2, Valid: true},
		UniqueWallet: "TWallet123",
		Status:       "pending",
		ExpiresAt:    pgtype.Timestamptz{Time: now.Add(1 * time.Hour), Valid: true},
		ConfirmedAt:  pgtype.Timestamptz{Valid: false},
		AttemptCount: &attemptCount,
		CreatedAt:    pgtype.Timestamptz{Time: now, Valid: true},
	}

	assert.Equal(t, id, payment.ID)
	assert.Equal(t, clientID, payment.ClientID)
	assert.Equal(t, accountID, payment.AccountID)
	assert.True(t, payment.Amount.Valid)
	assert.Equal(t, "TWallet123", payment.UniqueWallet)
	assert.Equal(t, "pending", payment.Status)
	assert.True(t, payment.ExpiresAt.Valid)
	assert.False(t, payment.ConfirmedAt.Valid)
	assert.NotNil(t, payment.AttemptCount)
	assert.Equal(t, int32(1), *payment.AttemptCount)
}

func TestPayment_ZeroValues(t *testing.T) {
	var payment Payment

	assert.Equal(t, uuid.Nil, payment.ID)
	assert.Equal(t, uuid.Nil, payment.ClientID)
	assert.Equal(t, uuid.Nil, payment.AccountID)
	assert.False(t, payment.Amount.Valid)
	assert.Equal(t, "", payment.UniqueWallet)
	assert.Equal(t, "", payment.Status)
	assert.False(t, payment.ExpiresAt.Valid)
	assert.False(t, payment.ConfirmedAt.Valid)
	assert.Nil(t, payment.AttemptCount)
}

func TestPayment_DifferentStatuses(t *testing.T) {
	statuses := []string{
		"pending",
		"confirmed",
		"completed",
		"expired",
		"failed",
		"cancelled",
	}

	for _, status := range statuses {
		payment := Payment{
			ID:           uuid.New(),
			ClientID:     uuid.New(),
			AccountID:    uuid.New(),
			Amount:       pgtype.Numeric{Int: big.NewInt(100000), Exp: -2, Valid: true},
			UniqueWallet: "TWallet",
			Status:       status,
			ExpiresAt:    pgtype.Timestamptz{Time: time.Now().Add(1 * time.Hour), Valid: true},
			ConfirmedAt:  pgtype.Timestamptz{Valid: false},
			AttemptCount: int32Ptr(0),
			CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
		}

		assert.Equal(t, status, payment.Status)
	}
}

func TestPayment_NullAttemptCount(t *testing.T) {
	payment := Payment{
		ID:           uuid.New(),
		ClientID:     uuid.New(),
		AccountID:    uuid.New(),
		Amount:       pgtype.Numeric{Int: big.NewInt(100), Exp: 0, Valid: true},
		UniqueWallet: "TWallet",
		Status:       "pending",
		ExpiresAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
		ConfirmedAt:  pgtype.Timestamptz{Valid: false},
		AttemptCount: nil,
		CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.Nil(t, payment.AttemptCount)
}

func TestPayment_ZeroAttemptCount(t *testing.T) {
	attemptCount := int32(0)
	payment := Payment{
		ID:           uuid.New(),
		ClientID:     uuid.New(),
		AccountID:    uuid.New(),
		Amount:       pgtype.Numeric{Int: big.NewInt(100), Exp: 0, Valid: true},
		UniqueWallet: "TWallet",
		Status:       "pending",
		ExpiresAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
		ConfirmedAt:  pgtype.Timestamptz{Valid: false},
		AttemptCount: &attemptCount,
		CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.NotNil(t, payment.AttemptCount)
	assert.Equal(t, int32(0), *payment.AttemptCount)
}

func TestPayment_HighAttemptCount(t *testing.T) {
	attemptCount := int32(99)
	payment := Payment{
		ID:           uuid.New(),
		ClientID:     uuid.New(),
		AccountID:    uuid.New(),
		Amount:       pgtype.Numeric{Int: big.NewInt(100), Exp: 0, Valid: true},
		UniqueWallet: "TWallet",
		Status:       "failed",
		ExpiresAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
		ConfirmedAt:  pgtype.Timestamptz{Valid: false},
		AttemptCount: &attemptCount,
		CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.NotNil(t, payment.AttemptCount)
	assert.Equal(t, int32(99), *payment.AttemptCount)
}

func TestPayment_ConfirmedPayment(t *testing.T) {
	now := time.Now()
	payment := Payment{
		ID:           uuid.New(),
		ClientID:     uuid.New(),
		AccountID:    uuid.New(),
		Amount:       pgtype.Numeric{Int: big.NewInt(100), Exp: 0, Valid: true},
		UniqueWallet: "TWallet",
		Status:       "confirmed",
		ExpiresAt:    pgtype.Timestamptz{Time: now.Add(1 * time.Hour), Valid: true},
		ConfirmedAt:  pgtype.Timestamptz{Time: now, Valid: true},
		AttemptCount: int32Ptr(1),
		CreatedAt:    pgtype.Timestamptz{Time: now, Valid: true},
	}

	assert.True(t, payment.ConfirmedAt.Valid)
	assert.Equal(t, "confirmed", payment.Status)
}

func TestPayment_ExpiredPayment(t *testing.T) {
	now := time.Now()
	payment := Payment{
		ID:           uuid.New(),
		ClientID:     uuid.New(),
		AccountID:    uuid.New(),
		Amount:       pgtype.Numeric{Int: big.NewInt(100), Exp: 0, Valid: true},
		UniqueWallet: "TWallet",
		Status:       "expired",
		ExpiresAt:    pgtype.Timestamptz{Time: now.Add(-1 * time.Hour), Valid: true},
		ConfirmedAt:  pgtype.Timestamptz{Valid: false},
		AttemptCount: int32Ptr(0),
		CreatedAt:    pgtype.Timestamptz{Time: now.Add(-2 * time.Hour), Valid: true},
	}

	assert.Equal(t, "expired", payment.Status)
	assert.False(t, payment.ConfirmedAt.Valid)
}

func TestPayment_JSONSerialization(t *testing.T) {
	id := uuid.New()
	clientID := uuid.New()
	accountID := uuid.New()
	now := time.Now()
	attemptCount := int32(2)

	payment := Payment{
		ID:           id,
		ClientID:     clientID,
		AccountID:    accountID,
		Amount:       pgtype.Numeric{Int: big.NewInt(100000), Exp: -2, Valid: true},
		UniqueWallet: "TWallet123",
		Status:       "pending",
		ExpiresAt:    pgtype.Timestamptz{Time: now.Add(1 * time.Hour), Valid: true},
		ConfirmedAt:  pgtype.Timestamptz{Valid: false},
		AttemptCount: &attemptCount,
		CreatedAt:    pgtype.Timestamptz{Time: now, Valid: true},
	}

	jsonData, err := json.Marshal(payment)
	require.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	var decoded Payment
	err = json.Unmarshal(jsonData, &decoded)
	require.NoError(t, err)

	assert.Equal(t, payment.ID, decoded.ID)
	assert.Equal(t, payment.ClientID, decoded.ClientID)
	assert.Equal(t, payment.AccountID, decoded.AccountID)
	assert.Equal(t, payment.UniqueWallet, decoded.UniqueWallet)
	assert.Equal(t, payment.Status, decoded.Status)
}

func TestPayment_DifferentAmounts(t *testing.T) {
	testCases := []struct {
		name     string
		amount   *big.Int
		exp      int32
		expected string
	}{
		{"small amount", big.NewInt(100), -2, "1.00"},
		{"large amount", big.NewInt(1000000), -2, "10000.00"},
		{"zero amount", big.NewInt(0), -2, "0.00"},
		{"fractional", big.NewInt(12345), -2, "123.45"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			payment := Payment{
				ID:           uuid.New(),
				ClientID:     uuid.New(),
				AccountID:    uuid.New(),
				Amount:       pgtype.Numeric{Int: tc.amount, Exp: tc.exp, Valid: true},
				UniqueWallet: "TWallet",
				Status:       "pending",
				ExpiresAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
				ConfirmedAt:  pgtype.Timestamptz{Valid: false},
				AttemptCount: int32Ptr(0),
				CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
			}

			assert.True(t, payment.Amount.Valid)
			assert.Equal(t, tc.amount, payment.Amount.Int)
			assert.Equal(t, tc.exp, payment.Amount.Exp)
		})
	}
}

func TestPayment_EmptyUniqueWallet(t *testing.T) {
	payment := Payment{
		ID:           uuid.New(),
		ClientID:     uuid.New(),
		AccountID:    uuid.New(),
		Amount:       pgtype.Numeric{Int: big.NewInt(100), Exp: 0, Valid: true},
		UniqueWallet: "",
		Status:       "pending",
		ExpiresAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
		ConfirmedAt:  pgtype.Timestamptz{Valid: false},
		AttemptCount: int32Ptr(0),
		CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.Equal(t, "", payment.UniqueWallet)
}

// Tests for PaymentAttempt model
func TestPaymentAttempt_Struct(t *testing.T) {
	id := uuid.New()
	paymentID := uuid.New()
	now := time.Now()

	attempt := PaymentAttempt{
		ID:              id,
		PaymentID:       paymentID,
		AttemptNumber:   1,
		GeneratedWallet: "TWallet123",
		GeneratedAt:     pgtype.Timestamptz{Time: now, Valid: true},
	}

	assert.Equal(t, id, attempt.ID)
	assert.Equal(t, paymentID, attempt.PaymentID)
	assert.Equal(t, int32(1), attempt.AttemptNumber)
	assert.Equal(t, "TWallet123", attempt.GeneratedWallet)
	assert.True(t, attempt.GeneratedAt.Valid)
	assert.Equal(t, now, attempt.GeneratedAt.Time)
}

func TestPaymentAttempt_ZeroValues(t *testing.T) {
	var attempt PaymentAttempt

	assert.Equal(t, uuid.Nil, attempt.ID)
	assert.Equal(t, uuid.Nil, attempt.PaymentID)
	assert.Equal(t, int32(0), attempt.AttemptNumber)
	assert.Equal(t, "", attempt.GeneratedWallet)
	assert.False(t, attempt.GeneratedAt.Valid)
}

func TestPaymentAttempt_FirstAttempt(t *testing.T) {
	attempt := PaymentAttempt{
		ID:              uuid.New(),
		PaymentID:       uuid.New(),
		AttemptNumber:   1,
		GeneratedWallet: "TWallet1",
		GeneratedAt:     pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.Equal(t, int32(1), attempt.AttemptNumber)
}

func TestPaymentAttempt_MultipleAttempts(t *testing.T) {
	paymentID := uuid.New()

	attempts := []PaymentAttempt{
		{
			ID:              uuid.New(),
			PaymentID:       paymentID,
			AttemptNumber:   1,
			GeneratedWallet: "TWallet1",
			GeneratedAt:     pgtype.Timestamptz{Time: time.Now(), Valid: true},
		},
		{
			ID:              uuid.New(),
			PaymentID:       paymentID,
			AttemptNumber:   2,
			GeneratedWallet: "TWallet2",
			GeneratedAt:     pgtype.Timestamptz{Time: time.Now(), Valid: true},
		},
		{
			ID:              uuid.New(),
			PaymentID:       paymentID,
			AttemptNumber:   3,
			GeneratedWallet: "TWallet3",
			GeneratedAt:     pgtype.Timestamptz{Time: time.Now(), Valid: true},
		},
	}

	for i, attempt := range attempts {
		assert.Equal(t, paymentID, attempt.PaymentID)
		assert.Equal(t, int32(i+1), attempt.AttemptNumber)
		assert.NotEqual(t, attempts[0].ID, attempts[1].ID)
	}
}

func TestPaymentAttempt_HighAttemptNumber(t *testing.T) {
	attempt := PaymentAttempt{
		ID:              uuid.New(),
		PaymentID:       uuid.New(),
		AttemptNumber:   100,
		GeneratedWallet: "TWallet100",
		GeneratedAt:     pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.Equal(t, int32(100), attempt.AttemptNumber)
}

func TestPaymentAttempt_JSONSerialization(t *testing.T) {
	id := uuid.New()
	paymentID := uuid.New()
	now := time.Now()

	attempt := PaymentAttempt{
		ID:              id,
		PaymentID:       paymentID,
		AttemptNumber:   5,
		GeneratedWallet: "TWallet123",
		GeneratedAt:     pgtype.Timestamptz{Time: now, Valid: true},
	}

	jsonData, err := json.Marshal(attempt)
	require.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	var decoded PaymentAttempt
	err = json.Unmarshal(jsonData, &decoded)
	require.NoError(t, err)

	assert.Equal(t, attempt.ID, decoded.ID)
	assert.Equal(t, attempt.PaymentID, decoded.PaymentID)
	assert.Equal(t, attempt.AttemptNumber, decoded.AttemptNumber)
	assert.Equal(t, attempt.GeneratedWallet, decoded.GeneratedWallet)
}

func TestPaymentAttempt_EmptyGeneratedWallet(t *testing.T) {
	attempt := PaymentAttempt{
		ID:              uuid.New(),
		PaymentID:       uuid.New(),
		AttemptNumber:   1,
		GeneratedWallet: "",
		GeneratedAt:     pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.Equal(t, "", attempt.GeneratedWallet)
}

func TestPaymentAttempt_NullGeneratedAt(t *testing.T) {
	attempt := PaymentAttempt{
		ID:              uuid.New(),
		PaymentID:       uuid.New(),
		AttemptNumber:   1,
		GeneratedWallet: "TWallet",
		GeneratedAt:     pgtype.Timestamptz{Valid: false},
	}

	assert.False(t, attempt.GeneratedAt.Valid)
}

func TestPaymentAttempt_DifferentWalletFormats(t *testing.T) {
	walletFormats := []string{
		"TWallet123",
		"TXabcdef123456",
		"T" + string(make([]byte, 33)),
		"TShortWallet",
		"TLongWalletAddressWithManyCharacters",
	}

	for i, wallet := range walletFormats {
		attempt := PaymentAttempt{
			ID:              uuid.New(),
			PaymentID:       uuid.New(),
			AttemptNumber:   int32(i + 1),
			GeneratedWallet: wallet,
			GeneratedAt:     pgtype.Timestamptz{Time: time.Now(), Valid: true},
		}

		assert.Equal(t, wallet, attempt.GeneratedWallet)
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func int32Ptr(i int32) *int32 {
	return &i
}