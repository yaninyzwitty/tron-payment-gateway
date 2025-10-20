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

// ============================================================================
// Tests for new AddressIndex field in Account model
// ============================================================================

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
	assert.Equal(t, now, account.CreatedAt.Time)
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

func TestAccount_LargeAddressIndex(t *testing.T) {
	addressIndex := int32(2147483647) // Max int32
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

// ============================================================================
// Tests for Payment model
// ============================================================================

func TestPayment_Struct(t *testing.T) {
	id := uuid.New()
	clientID := uuid.New()
	accountID := uuid.New()
	now := time.Now()
	attemptCount := int32(0)

	payment := Payment{
		ID:           id,
		ClientID:     clientID,
		AccountID:    accountID,
		Amount:       pgtype.Numeric{Int: big.NewInt(1000000), Exp: -6, Valid: true},
		UniqueWallet: "TXYZabc123",
		Status:       "PENDING",
		ExpiresAt:    pgtype.Timestamptz{Time: now.Add(time.Hour), Valid: true},
		ConfirmedAt:  pgtype.Timestamptz{Valid: false},
		AttemptCount: &attemptCount,
		CreatedAt:    pgtype.Timestamptz{Time: now, Valid: true},
	}

	assert.Equal(t, id, payment.ID)
	assert.Equal(t, clientID, payment.ClientID)
	assert.Equal(t, accountID, payment.AccountID)
	assert.True(t, payment.Amount.Valid)
	assert.Equal(t, "TXYZabc123", payment.UniqueWallet)
	assert.Equal(t, "PENDING", payment.Status)
	assert.True(t, payment.ExpiresAt.Valid)
	assert.False(t, payment.ConfirmedAt.Valid)
	assert.NotNil(t, payment.AttemptCount)
	assert.Equal(t, int32(0), *payment.AttemptCount)
	assert.True(t, payment.CreatedAt.Valid)
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
	assert.False(t, payment.CreatedAt.Valid)
}

func TestPayment_StatusValues(t *testing.T) {
	statuses := []string{"PENDING", "CONFIRMED", "EXPIRED"}

	for _, status := range statuses {
		payment := Payment{
			ID:           uuid.New(),
			ClientID:     uuid.New(),
			AccountID:    uuid.New(),
			Amount:       pgtype.Numeric{Int: big.NewInt(1000000), Exp: -6, Valid: true},
			UniqueWallet: "TXYZabc123",
			Status:       status,
			ExpiresAt:    pgtype.Timestamptz{Time: time.Now().Add(time.Hour), Valid: true},
			CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
		}

		assert.Equal(t, status, payment.Status)
	}
}

func TestPayment_ConfirmedPayment(t *testing.T) {
	now := time.Now()
	confirmedTime := now.Add(time.Minute * 30)

	payment := Payment{
		ID:           uuid.New(),
		ClientID:     uuid.New(),
		AccountID:    uuid.New(),
		Amount:       pgtype.Numeric{Int: big.NewInt(1000000), Exp: -6, Valid: true},
		UniqueWallet: "TXYZabc123",
		Status:       "CONFIRMED",
		ExpiresAt:    pgtype.Timestamptz{Time: now.Add(time.Hour), Valid: true},
		ConfirmedAt:  pgtype.Timestamptz{Time: confirmedTime, Valid: true},
		CreatedAt:    pgtype.Timestamptz{Time: now, Valid: true},
	}

	assert.Equal(t, "CONFIRMED", payment.Status)
	assert.True(t, payment.ConfirmedAt.Valid)
	assert.Equal(t, confirmedTime, payment.ConfirmedAt.Time)
}

func TestPayment_ExpiredPayment(t *testing.T) {
	now := time.Now()
	pastTime := now.Add(-time.Hour)

	payment := Payment{
		ID:           uuid.New(),
		ClientID:     uuid.New(),
		AccountID:    uuid.New(),
		Amount:       pgtype.Numeric{Int: big.NewInt(1000000), Exp: -6, Valid: true},
		UniqueWallet: "TXYZabc123",
		Status:       "EXPIRED",
		ExpiresAt:    pgtype.Timestamptz{Time: pastTime, Valid: true},
		ConfirmedAt:  pgtype.Timestamptz{Valid: false},
		CreatedAt:    pgtype.Timestamptz{Time: now.Add(-time.Hour * 2), Valid: true},
	}

	assert.Equal(t, "EXPIRED", payment.Status)
	assert.True(t, payment.ExpiresAt.Valid)
	assert.True(t, payment.ExpiresAt.Time.Before(now))
	assert.False(t, payment.ConfirmedAt.Valid)
}

func TestPayment_AttemptCount(t *testing.T) {
	testCases := []int32{0, 1, 5, 10, 100}

	for _, count := range testCases {
		attemptCount := count
		payment := Payment{
			ID:           uuid.New(),
			ClientID:     uuid.New(),
			AccountID:    uuid.New(),
			Amount:       pgtype.Numeric{Int: big.NewInt(1000000), Exp: -6, Valid: true},
			UniqueWallet: "TXYZabc123",
			Status:       "PENDING",
			ExpiresAt:    pgtype.Timestamptz{Time: time.Now().Add(time.Hour), Valid: true},
			AttemptCount: &attemptCount,
			CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
		}

		assert.NotNil(t, payment.AttemptCount)
		assert.Equal(t, count, *payment.AttemptCount)
	}
}

func TestPayment_NullAttemptCount(t *testing.T) {
	payment := Payment{
		ID:           uuid.New(),
		ClientID:     uuid.New(),
		AccountID:    uuid.New(),
		Amount:       pgtype.Numeric{Int: big.NewInt(1000000), Exp: -6, Valid: true},
		UniqueWallet: "TXYZabc123",
		Status:       "PENDING",
		ExpiresAt:    pgtype.Timestamptz{Time: time.Now().Add(time.Hour), Valid: true},
		AttemptCount: nil,
		CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.Nil(t, payment.AttemptCount)
}

func TestPayment_Amount(t *testing.T) {
	testCases := []struct {
		name     string
		intValue int64
		exp      int32
	}{
		{"1 TRX", 1000000, -6},
		{"0.1 TRX", 100000, -6},
		{"10 TRX", 10000000, -6},
		{"100.5 TRX", 100500000, -6},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			payment := Payment{
				ID:           uuid.New(),
				ClientID:     uuid.New(),
				AccountID:    uuid.New(),
				Amount:       pgtype.Numeric{Int: big.NewInt(tc.intValue), Exp: tc.exp, Valid: true},
				UniqueWallet: "TXYZabc123",
				Status:       "PENDING",
				ExpiresAt:    pgtype.Timestamptz{Time: time.Now().Add(time.Hour), Valid: true},
				CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
			}

			assert.True(t, payment.Amount.Valid)
			assert.Equal(t, tc.intValue, payment.Amount.Int.Int64())
			assert.Equal(t, tc.exp, payment.Amount.Exp)
		})
	}
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
		Amount:       pgtype.Numeric{Int: big.NewInt(1000000), Exp: -6, Valid: true},
		UniqueWallet: "TXYZabc123",
		Status:       "PENDING",
		ExpiresAt:    pgtype.Timestamptz{Time: now.Add(time.Hour), Valid: true},
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

func TestPayment_EmptyWallet(t *testing.T) {
	payment := Payment{
		ID:           uuid.New(),
		ClientID:     uuid.New(),
		AccountID:    uuid.New(),
		Amount:       pgtype.Numeric{Int: big.NewInt(1000000), Exp: -6, Valid: true},
		UniqueWallet: "",
		Status:       "PENDING",
		ExpiresAt:    pgtype.Timestamptz{Time: time.Now().Add(time.Hour), Valid: true},
		CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.Equal(t, "", payment.UniqueWallet)
}

func TestPayment_LongWalletAddress(t *testing.T) {
	longWallet := "T" + string(make([]byte, 100))
	payment := Payment{
		ID:           uuid.New(),
		ClientID:     uuid.New(),
		AccountID:    uuid.New(),
		Amount:       pgtype.Numeric{Int: big.NewInt(1000000), Exp: -6, Valid: true},
		UniqueWallet: longWallet,
		Status:       "PENDING",
		ExpiresAt:    pgtype.Timestamptz{Time: time.Now().Add(time.Hour), Valid: true},
		CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.Equal(t, longWallet, payment.UniqueWallet)
}

// ============================================================================
// Tests for PaymentAttempt model
// ============================================================================

func TestPaymentAttempt_Struct(t *testing.T) {
	id := uuid.New()
	paymentID := uuid.New()
	now := time.Now()

	attempt := PaymentAttempt{
		ID:              id,
		PaymentID:       paymentID,
		AttemptNumber:   1,
		GeneratedWallet: "TXYZabc123",
		GeneratedAt:     pgtype.Timestamptz{Time: now, Valid: true},
	}

	assert.Equal(t, id, attempt.ID)
	assert.Equal(t, paymentID, attempt.PaymentID)
	assert.Equal(t, int32(1), attempt.AttemptNumber)
	assert.Equal(t, "TXYZabc123", attempt.GeneratedWallet)
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

func TestPaymentAttempt_MultipleAttempts(t *testing.T) {
	paymentID := uuid.New()

	for i := int32(1); i <= 5; i++ {
		attempt := PaymentAttempt{
			ID:              uuid.New(),
			PaymentID:       paymentID,
			AttemptNumber:   i,
			GeneratedWallet: "TXYZabc" + string(rune('0'+i)),
			GeneratedAt:     pgtype.Timestamptz{Time: time.Now(), Valid: true},
		}

		assert.Equal(t, i, attempt.AttemptNumber)
		assert.Equal(t, paymentID, attempt.PaymentID)
	}
}

func TestPaymentAttempt_JSONSerialization(t *testing.T) {
	id := uuid.New()
	paymentID := uuid.New()
	now := time.Now()

	attempt := PaymentAttempt{
		ID:              id,
		PaymentID:       paymentID,
		AttemptNumber:   3,
		GeneratedWallet: "TXYZabc123",
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

func TestPaymentAttempt_NullGeneratedAt(t *testing.T) {
	attempt := PaymentAttempt{
		ID:              uuid.New(),
		PaymentID:       uuid.New(),
		AttemptNumber:   1,
		GeneratedWallet: "TXYZabc123",
		GeneratedAt:     pgtype.Timestamptz{Valid: false},
	}

	assert.False(t, attempt.GeneratedAt.Valid)
}

func TestPaymentAttempt_EmptyWallet(t *testing.T) {
	attempt := PaymentAttempt{
		ID:              uuid.New(),
		PaymentID:       uuid.New(),
		AttemptNumber:   1,
		GeneratedWallet: "",
		GeneratedAt:     pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.Equal(t, "", attempt.GeneratedWallet)
}

func TestPaymentAttempt_LargeAttemptNumber(t *testing.T) {
	attempt := PaymentAttempt{
		ID:              uuid.New(),
		PaymentID:       uuid.New(),
		AttemptNumber:   2147483647, // Max int32
		GeneratedWallet: "TXYZabc123",
		GeneratedAt:     pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.Equal(t, int32(2147483647), attempt.AttemptNumber)
}

// ============================================================================
// Tests for Log model
// ============================================================================

func TestLog_Struct(t *testing.T) {
	id := uuid.New()
	paymentID := uuid.New()
	now := time.Now()
	message := "Test log message"
	rawData := []byte(`{"key": "value"}`)

	log := Log{
		ID:        id,
		PaymentID: pgtype.UUID{Bytes: paymentID, Valid: true},
		EventType: "ADDRESS_GENERATED",
		Message:   &message,
		RawData:   rawData,
		CreatedAt: pgtype.Timestamptz{Time: now, Valid: true},
	}

	assert.Equal(t, id, log.ID)
	assert.True(t, log.PaymentID.Valid)
	assert.Equal(t, paymentID, log.PaymentID.Bytes)
	assert.Equal(t, "ADDRESS_GENERATED", log.EventType)
	assert.NotNil(t, log.Message)
	assert.Equal(t, "Test log message", *log.Message)
	assert.Equal(t, rawData, log.RawData)
	assert.True(t, log.CreatedAt.Valid)
	assert.Equal(t, now, log.CreatedAt.Time)
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

func TestLog_EventTypes(t *testing.T) {
	eventTypes := []string{
		"ADDRESS_GENERATED",
		"TX_CONFIRMED",
		"WEBHOOK_SENT",
		"ERROR",
		"PAYMENT_EXPIRED",
	}

	for _, eventType := range eventTypes {
		log := Log{
			ID:        uuid.New(),
			PaymentID: pgtype.UUID{Bytes: uuid.New(), Valid: true},
			EventType: eventType,
			CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		}

		assert.Equal(t, eventType, log.EventType)
	}
}

func TestLog_NullPaymentID(t *testing.T) {
	log := Log{
		ID:        uuid.New(),
		PaymentID: pgtype.UUID{Valid: false},
		EventType: "SYSTEM_EVENT",
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.False(t, log.PaymentID.Valid)
}

func TestLog_NullMessage(t *testing.T) {
	log := Log{
		ID:        uuid.New(),
		PaymentID: pgtype.UUID{Bytes: uuid.New(), Valid: true},
		EventType: "TX_CONFIRMED",
		Message:   nil,
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.Nil(t, log.Message)
}

func TestLog_EmptyMessage(t *testing.T) {
	message := ""
	log := Log{
		ID:        uuid.New(),
		PaymentID: pgtype.UUID{Bytes: uuid.New(), Valid: true},
		EventType: "TX_CONFIRMED",
		Message:   &message,
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.NotNil(t, log.Message)
	assert.Equal(t, "", *log.Message)
}

func TestLog_LongMessage(t *testing.T) {
	longMessage := string(make([]byte, 10000))
	log := Log{
		ID:        uuid.New(),
		PaymentID: pgtype.UUID{Bytes: uuid.New(), Valid: true},
		EventType: "ERROR",
		Message:   &longMessage,
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.NotNil(t, log.Message)
	assert.Len(t, *log.Message, 10000)
}

func TestLog_NullRawData(t *testing.T) {
	log := Log{
		ID:        uuid.New(),
		PaymentID: pgtype.UUID{Bytes: uuid.New(), Valid: true},
		EventType: "ADDRESS_GENERATED",
		RawData:   nil,
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.Nil(t, log.RawData)
}

func TestLog_EmptyRawData(t *testing.T) {
	log := Log{
		ID:        uuid.New(),
		PaymentID: pgtype.UUID{Bytes: uuid.New(), Valid: true},
		EventType: "ADDRESS_GENERATED",
		RawData:   []byte{},
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.NotNil(t, log.RawData)
	assert.Empty(t, log.RawData)
}

func TestLog_JSONRawData(t *testing.T) {
	rawData := []byte(`{"transaction_id": "abc123", "amount": 1000000}`)
	log := Log{
		ID:        uuid.New(),
		PaymentID: pgtype.UUID{Bytes: uuid.New(), Valid: true},
		EventType: "TX_CONFIRMED",
		RawData:   rawData,
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.NotNil(t, log.RawData)
	assert.Equal(t, rawData, log.RawData)

	// Verify it's valid JSON
	var jsonMap map[string]interface{}
	err := json.Unmarshal(log.RawData, &jsonMap)
	require.NoError(t, err)
	assert.Equal(t, "abc123", jsonMap["transaction_id"])
	assert.Equal(t, float64(1000000), jsonMap["amount"])
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
		EventType: "ADDRESS_GENERATED",
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

func TestLog_NullCreatedAt(t *testing.T) {
	log := Log{
		ID:        uuid.New(),
		PaymentID: pgtype.UUID{Bytes: uuid.New(), Valid: true},
		EventType: "ERROR",
		CreatedAt: pgtype.Timestamptz{Valid: false},
	}

	assert.False(t, log.CreatedAt.Valid)
}

// ============================================================================
// Integration tests across models
// ============================================================================

func TestPaymentWithAttempts_Relationship(t *testing.T) {
	paymentID := uuid.New()
	clientID := uuid.New()
	accountID := uuid.New()

	payment := Payment{
		ID:           paymentID,
		ClientID:     clientID,
		AccountID:    accountID,
		Amount:       pgtype.Numeric{Int: big.NewInt(1000000), Exp: -6, Valid: true},
		UniqueWallet: "TXYZabc123",
		Status:       "PENDING",
		ExpiresAt:    pgtype.Timestamptz{Time: time.Now().Add(time.Hour), Valid: true},
		CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	attempt1 := PaymentAttempt{
		ID:              uuid.New(),
		PaymentID:       paymentID,
		AttemptNumber:   1,
		GeneratedWallet: "TXYZabc123",
		GeneratedAt:     pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	attempt2 := PaymentAttempt{
		ID:              uuid.New(),
		PaymentID:       paymentID,
		AttemptNumber:   2,
		GeneratedWallet: "TXYZabc456",
		GeneratedAt:     pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.Equal(t, payment.ID, attempt1.PaymentID)
	assert.Equal(t, payment.ID, attempt2.PaymentID)
	assert.NotEqual(t, attempt1.AttemptNumber, attempt2.AttemptNumber)
}

func TestPaymentWithLogs_Relationship(t *testing.T) {
	paymentID := uuid.New()

	payment := Payment{
		ID:           paymentID,
		ClientID:     uuid.New(),
		AccountID:    uuid.New(),
		Amount:       pgtype.Numeric{Int: big.NewInt(1000000), Exp: -6, Valid: true},
		UniqueWallet: "TXYZabc123",
		Status:       "CONFIRMED",
		ExpiresAt:    pgtype.Timestamptz{Time: time.Now().Add(time.Hour), Valid: true},
		ConfirmedAt:  pgtype.Timestamptz{Time: time.Now(), Valid: true},
		CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	log1 := Log{
		ID:        uuid.New(),
		PaymentID: pgtype.UUID{Bytes: paymentID, Valid: true},
		EventType: "ADDRESS_GENERATED",
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	log2 := Log{
		ID:        uuid.New(),
		PaymentID: pgtype.UUID{Bytes: paymentID, Valid: true},
		EventType: "TX_CONFIRMED",
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.Equal(t, payment.ID, log1.PaymentID.Bytes)
	assert.Equal(t, payment.ID, log2.PaymentID.Bytes)
	assert.True(t, log1.PaymentID.Valid)
	assert.True(t, log2.PaymentID.Valid)
}

func TestAccountWithPayment_Relationship(t *testing.T) {
	clientID := uuid.New()
	accountID := uuid.New()

	account := Account{
		ID:        accountID,
		ClientID:  clientID,
		Name:      "Test Account",
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	payment := Payment{
		ID:           uuid.New(),
		ClientID:     clientID,
		AccountID:    accountID,
		Amount:       pgtype.Numeric{Int: big.NewInt(1000000), Exp: -6, Valid: true},
		UniqueWallet: "TXYZabc123",
		Status:       "PENDING",
		ExpiresAt:    pgtype.Timestamptz{Time: time.Now().Add(time.Hour), Valid: true},
		CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	assert.Equal(t, account.ID, payment.AccountID)
	assert.Equal(t, account.ClientID, payment.ClientID)
}

// Helper function to create an int32 pointer
func int32Ptr(i int32) *int32 {
	return &i
}

// Helper function to create a string pointer  
func stringPtr(s string) *string {
	return &s
}