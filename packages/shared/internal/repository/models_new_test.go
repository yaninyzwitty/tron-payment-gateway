package repository

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Tests for Account model with AddressIndex field

func TestAccount_WithAddressIndex(t *testing.T) {
	id := uuid.New()
	clientID := uuid.New()
	addressIndex := int32(42)

	account := Account{
		ID:           id,
		ClientID:     clientID,
		Name:         "Test Account",
		AddressIndex: &addressIndex,
		CreatedAt:    pgtype.Timestamptz{},
	}

	assert.Equal(t, id, account.ID)
	assert.Equal(t, clientID, account.ClientID)
	assert.Equal(t, "Test Account", account.Name)
	assert.NotNil(t, account.AddressIndex)
	assert.Equal(t, int32(42), *account.AddressIndex)
}

func TestAccount_NilAddressIndex(t *testing.T) {
	account := Account{
		ID:           uuid.New(),
		ClientID:     uuid.New(),
		Name:         "Test Account",
		AddressIndex: nil,
		CreatedAt:    pgtype.Timestamptz{},
	}

	assert.Nil(t, account.AddressIndex)
}

func TestAccount_JSONSerializationWithAddressIndex(t *testing.T) {
	addressIndex := int32(100)
	account := Account{
		ID:           uuid.New(),
		ClientID:     uuid.New(),
		Name:         "Wallet Account",
		AddressIndex: &addressIndex,
	}

	jsonData, err := json.Marshal(account)
	require.NoError(t, err)
	assert.Contains(t, string(jsonData), "address_index")

	var decoded Account
	err = json.Unmarshal(jsonData, &decoded)
	require.NoError(t, err)
	assert.NotNil(t, decoded.AddressIndex)
	assert.Equal(t, int32(100), *decoded.AddressIndex)
}

func TestAccount_ZeroAddressIndex(t *testing.T) {
	addressIndex := int32(0)
	account := Account{
		ID:           uuid.New(),
		ClientID:     uuid.New(),
		Name:         "Account",
		AddressIndex: &addressIndex,
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
	}

	assert.NotNil(t, account.AddressIndex)
	assert.Equal(t, int32(-1), *account.AddressIndex)
}

func TestAccount_LargeAddressIndex(t *testing.T) {
	addressIndex := int32(2147483647) // Max int32
	account := Account{
		ID:           uuid.New(),
		ClientID:     uuid.New(),
		Name:         "Account",
		AddressIndex: &addressIndex,
	}

	assert.NotNil(t, account.AddressIndex)
	assert.Equal(t, int32(2147483647), *account.AddressIndex)
}

// Tests for Payment model

func TestPayment_Struct(t *testing.T) {
	id := uuid.New()
	clientID := uuid.New()
	accountID := uuid.New()
	attemptCount := int32(0)

	payment := Payment{
		ID:           id,
		ClientID:     clientID,
		AccountID:    accountID,
		Amount:       pgtype.Numeric{},
		UniqueWallet: "TXXmULCEzRo6JfxUP1LYfmzbiKUezUvUNj",
		Status:       "PENDING",
		ExpiresAt:    pgtype.Timestamptz{},
		ConfirmedAt:  pgtype.Timestamptz{},
		AttemptCount: &attemptCount,
		CreatedAt:    pgtype.Timestamptz{},
	}

	assert.Equal(t, id, payment.ID)
	assert.Equal(t, clientID, payment.ClientID)
	assert.Equal(t, accountID, payment.AccountID)
	assert.Equal(t, "TXXmULCEzRo6JfxUP1LYfmzbiKUezUvUNj", payment.UniqueWallet)
	assert.Equal(t, "PENDING", payment.Status)
	assert.NotNil(t, payment.AttemptCount)
	assert.Equal(t, int32(0), *payment.AttemptCount)
}

func TestPayment_JSONSerialization(t *testing.T) {
	attemptCount := int32(3)
	payment := Payment{
		ID:           uuid.New(),
		ClientID:     uuid.New(),
		AccountID:    uuid.New(),
		UniqueWallet: "TTest123",
		Status:       "CONFIRMED",
		AttemptCount: &attemptCount,
	}

	jsonData, err := json.Marshal(payment)
	require.NoError(t, err)
	assert.Contains(t, string(jsonData), "unique_wallet")
	assert.Contains(t, string(jsonData), "attempt_count")

	var decoded Payment
	err = json.Unmarshal(jsonData, &decoded)
	require.NoError(t, err)
	assert.Equal(t, "TTest123", decoded.UniqueWallet)
	assert.Equal(t, "CONFIRMED", decoded.Status)
	assert.NotNil(t, decoded.AttemptCount)
	assert.Equal(t, int32(3), *decoded.AttemptCount)
}

func TestPayment_StatusValues(t *testing.T) {
	testCases := []struct {
		status string
	}{
		{"PENDING"},
		{"CONFIRMED"},
		{"EXPIRED"},
		{"FAILED"},
	}

	for _, tc := range testCases {
		t.Run(tc.status, func(t *testing.T) {
			payment := Payment{
				ID:           uuid.New(),
				ClientID:     uuid.New(),
				AccountID:    uuid.New(),
				UniqueWallet: "TAddress",
				Status:       tc.status,
			}

			assert.Equal(t, tc.status, payment.Status)
		})
	}
}

func TestPayment_NilAttemptCount(t *testing.T) {
	payment := Payment{
		ID:           uuid.New(),
		ClientID:     uuid.New(),
		AccountID:    uuid.New(),
		UniqueWallet: "TAddress",
		Status:       "PENDING",
		AttemptCount: nil,
	}

	assert.Nil(t, payment.AttemptCount)
}

func TestPayment_MultipleAttempts(t *testing.T) {
	attemptCount := int32(5)
	payment := Payment{
		ID:           uuid.New(),
		ClientID:     uuid.New(),
		AccountID:    uuid.New(),
		UniqueWallet: "TAddress",
		Status:       "PENDING",
		AttemptCount: &attemptCount,
	}

	assert.NotNil(t, payment.AttemptCount)
	assert.Equal(t, int32(5), *payment.AttemptCount)
}

func TestPayment_ZeroAttemptCount(t *testing.T) {
	attemptCount := int32(0)
	payment := Payment{
		ID:           uuid.New(),
		ClientID:     uuid.New(),
		AccountID:    uuid.New(),
		UniqueWallet: "TAddress",
		Status:       "PENDING",
		AttemptCount: &attemptCount,
	}

	assert.NotNil(t, payment.AttemptCount)
	assert.Equal(t, int32(0), *payment.AttemptCount)
}

// Tests for PaymentAttempt model

func TestPaymentAttempt_Struct(t *testing.T) {
	id := uuid.New()
	paymentID := uuid.New()

	attempt := PaymentAttempt{
		ID:              id,
		PaymentID:       paymentID,
		AttemptNumber:   1,
		GeneratedWallet: "TGeneratedAddress123",
		GeneratedAt:     pgtype.Timestamptz{},
	}

	assert.Equal(t, id, attempt.ID)
	assert.Equal(t, paymentID, attempt.PaymentID)
	assert.Equal(t, int32(1), attempt.AttemptNumber)
	assert.Equal(t, "TGeneratedAddress123", attempt.GeneratedWallet)
}

func TestPaymentAttempt_JSONSerialization(t *testing.T) {
	attempt := PaymentAttempt{
		ID:              uuid.New(),
		PaymentID:       uuid.New(),
		AttemptNumber:   2,
		GeneratedWallet: "TWalletAddress",
	}

	jsonData, err := json.Marshal(attempt)
	require.NoError(t, err)
	assert.Contains(t, string(jsonData), "payment_id")
	assert.Contains(t, string(jsonData), "attempt_number")
	assert.Contains(t, string(jsonData), "generated_wallet")

	var decoded PaymentAttempt
	err = json.Unmarshal(jsonData, &decoded)
	require.NoError(t, err)
	assert.Equal(t, int32(2), decoded.AttemptNumber)
	assert.Equal(t, "TWalletAddress", decoded.GeneratedWallet)
}

func TestPaymentAttempt_MultipleAttempts(t *testing.T) {
	paymentID := uuid.New()

	attempts := []PaymentAttempt{
		{
			ID:              uuid.New(),
			PaymentID:       paymentID,
			AttemptNumber:   1,
			GeneratedWallet: "TAddress1",
		},
		{
			ID:              uuid.New(),
			PaymentID:       paymentID,
			AttemptNumber:   2,
			GeneratedWallet: "TAddress2",
		},
		{
			ID:              uuid.New(),
			PaymentID:       paymentID,
			AttemptNumber:   3,
			GeneratedWallet: "TAddress3",
		},
	}

	for i, attempt := range attempts {
		assert.Equal(t, paymentID, attempt.PaymentID)
		assert.Equal(t, int32(i+1), attempt.AttemptNumber)
		assert.NotEmpty(t, attempt.GeneratedWallet)
	}
}

func TestPaymentAttempt_ZeroAttemptNumber(t *testing.T) {
	attempt := PaymentAttempt{
		ID:              uuid.New(),
		PaymentID:       uuid.New(),
		AttemptNumber:   0,
		GeneratedWallet: "TAddress",
	}

	assert.Equal(t, int32(0), attempt.AttemptNumber)
}

func TestPaymentAttempt_NegativeAttemptNumber(t *testing.T) {
	attempt := PaymentAttempt{
		ID:              uuid.New(),
		PaymentID:       uuid.New(),
		AttemptNumber:   -1,
		GeneratedWallet: "TAddress",
	}

	assert.Equal(t, int32(-1), attempt.AttemptNumber)
}

// Tests for Log model

func TestLog_Struct(t *testing.T) {
	id := uuid.New()
	paymentID := pgtype.UUID{Bytes: uuid.New(), Valid: true}
	message := "Transaction confirmed"

	log := Log{
		ID:        id,
		PaymentID: paymentID,
		EventType: "TX_CONFIRMED",
		Message:   &message,
		RawData:   []byte(`{"tx_id": "abc123"}`),
		CreatedAt: pgtype.Timestamptz{},
	}

	assert.Equal(t, id, log.ID)
	assert.True(t, log.PaymentID.Valid)
	assert.Equal(t, "TX_CONFIRMED", log.EventType)
	assert.NotNil(t, log.Message)
	assert.Equal(t, "Transaction confirmed", *log.Message)
	assert.NotEmpty(t, log.RawData)
}

func TestLog_JSONSerialization(t *testing.T) {
	paymentID := pgtype.UUID{Bytes: uuid.New(), Valid: true}
	message := "Webhook sent"

	log := Log{
		ID:        uuid.New(),
		PaymentID: paymentID,
		EventType: "WEBHOOK_SENT",
		Message:   &message,
		RawData:   []byte(`{"status": "ok"}`),
	}

	jsonData, err := json.Marshal(log)
	require.NoError(t, err)
	assert.Contains(t, string(jsonData), "event_type")
	assert.Contains(t, string(jsonData), "WEBHOOK_SENT")

	var decoded Log
	err = json.Unmarshal(jsonData, &decoded)
	require.NoError(t, err)
	assert.Equal(t, "WEBHOOK_SENT", decoded.EventType)
}

func TestLog_EventTypes(t *testing.T) {
	testCases := []struct {
		eventType string
	}{
		{"ADDRESS_GENERATED"},
		{"TX_CONFIRMED"},
		{"WEBHOOK_SENT"},
		{"ERROR_OCCURRED"},
		{"PAYMENT_EXPIRED"},
	}

	for _, tc := range testCases {
		t.Run(tc.eventType, func(t *testing.T) {
			log := Log{
				ID:        uuid.New(),
				EventType: tc.eventType,
			}

			assert.Equal(t, tc.eventType, log.EventType)
		})
	}
}

func TestLog_NilMessage(t *testing.T) {
	log := Log{
		ID:        uuid.New(),
		EventType: "TEST_EVENT",
		Message:   nil,
	}

	assert.Nil(t, log.Message)
}

func TestLog_EmptyMessage(t *testing.T) {
	message := ""
	log := Log{
		ID:        uuid.New(),
		EventType: "TEST_EVENT",
		Message:   &message,
	}

	assert.NotNil(t, log.Message)
	assert.Equal(t, "", *log.Message)
}

func TestLog_NilPaymentID(t *testing.T) {
	log := Log{
		ID:        uuid.New(),
		PaymentID: pgtype.UUID{Valid: false},
		EventType: "SYSTEM_EVENT",
	}

	assert.False(t, log.PaymentID.Valid)
}

func TestLog_RawDataJSONFormat(t *testing.T) {
	rawData := []byte(`{"key": "value", "number": 123, "nested": {"field": "data"}}`)

	log := Log{
		ID:        uuid.New(),
		EventType: "DATA_EVENT",
		RawData:   rawData,
	}

	// Verify raw data is valid JSON
	var parsed map[string]interface{}
	err := json.Unmarshal(log.RawData, &parsed)
	require.NoError(t, err)
	assert.Equal(t, "value", parsed["key"])
	assert.Equal(t, float64(123), parsed["number"])
}

func TestLog_EmptyRawData(t *testing.T) {
	log := Log{
		ID:        uuid.New(),
		EventType: "TEST_EVENT",
		RawData:   []byte{},
	}

	assert.NotNil(t, log.RawData)
	assert.Equal(t, 0, len(log.RawData))
}

func TestLog_NilRawData(t *testing.T) {
	log := Log{
		ID:        uuid.New(),
		EventType: "TEST_EVENT",
		RawData:   nil,
	}

	assert.Nil(t, log.RawData)
}

func TestLog_LargeRawData(t *testing.T) {
	// Create large JSON payload
	largeData := make(map[string]string)
	for i := 0; i < 1000; i++ {
		largeData[string(rune('a'+i%26))+string(rune(i))] = "value"
	}

	jsonData, err := json.Marshal(largeData)
	require.NoError(t, err)

	log := Log{
		ID:        uuid.New(),
		EventType: "LARGE_DATA_EVENT",
		RawData:   jsonData,
	}

	assert.NotEmpty(t, log.RawData)
	assert.True(t, len(log.RawData) > 1000)
}

// Integration tests for model relationships

func TestPaymentWithAttemptsRelationship(t *testing.T) {
	paymentID := uuid.New()

	payment := Payment{
		ID:           paymentID,
		ClientID:     uuid.New(),
		AccountID:    uuid.New(),
		UniqueWallet: "TOriginalWallet",
		Status:       "PENDING",
	}

	attempts := []PaymentAttempt{
		{
			ID:              uuid.New(),
			PaymentID:       paymentID,
			AttemptNumber:   1,
			GeneratedWallet: "TAttempt1",
		},
		{
			ID:              uuid.New(),
			PaymentID:       paymentID,
			AttemptNumber:   2,
			GeneratedWallet: "TAttempt2",
		},
	}

	// Verify relationship
	for _, attempt := range attempts {
		assert.Equal(t, payment.ID, attempt.PaymentID)
	}
}

func TestPaymentWithLogsRelationship(t *testing.T) {
	paymentID := uuid.New()

	payment := Payment{
		ID:           paymentID,
		ClientID:     uuid.New(),
		AccountID:    uuid.New(),
		UniqueWallet: "TWallet",
		Status:       "CONFIRMED",
	}

	logs := []Log{
		{
			ID:        uuid.New(),
			PaymentID: pgtype.UUID{Bytes: paymentID, Valid: true},
			EventType: "ADDRESS_GENERATED",
		},
		{
			ID:        uuid.New(),
			PaymentID: pgtype.UUID{Bytes: paymentID, Valid: true},
			EventType: "TX_CONFIRMED",
		},
	}

	// Verify relationship
	for _, log := range logs {
		assert.True(t, log.PaymentID.Valid)
		assert.Equal(t, payment.ID, log.PaymentID.Bytes)
	}
}

func TestAccountWithPaymentRelationship(t *testing.T) {
	accountID := uuid.New()
	clientID := uuid.New()
	addressIndex := int32(0)

	account := Account{
		ID:           accountID,
		ClientID:     clientID,
		Name:         "Account 1",
		AddressIndex: &addressIndex,
	}

	payment := Payment{
		ID:           uuid.New(),
		ClientID:     clientID,
		AccountID:    accountID,
		UniqueWallet: "TWallet",
		Status:       "PENDING",
	}

	assert.Equal(t, account.ID, payment.AccountID)
	assert.Equal(t, account.ClientID, payment.ClientID)
}