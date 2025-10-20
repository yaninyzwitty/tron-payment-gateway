package db

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func readMigration(t *testing.T, name string) string {
	t.Helper()
	p := filepath.Join("migrations", name)
	data, err := os.ReadFile(p)
	require.NoErrorf(t, err, "failed to read %s", p)
	return string(data)
}

func TestMigrations_AccountsContainsAddressIndex(t *testing.T) {
	s := readMigration(t, "002_accounts.sql")
	require.Contains(t, s, "CREATE TABLE accounts")
	require.Contains(t, s, "address_index")
}

func TestMigrations_PaymentsSchema(t *testing.T) {
	s := readMigration(t, "003_payments.sql")
	require.True(t, strings.Contains(s, "CREATE TABLE payments"), "payments table missing")
	for _, want := range []string{"unique_wallet", "status", "expires_at", "attempt_count"} {
		require.Containsf(t, s, want, "expected %q in payments schema", want)
	}
}

func TestMigrations_PaymentAttemptsSchema(t *testing.T) {
	s := readMigration(t, "004_payments_attempts.sql")
	require.Contains(t, s, "CREATE TABLE payment_attempts")
	for _, want := range []string{"attempt_number", "generated_wallet", "generated_at"} {
		require.Contains(t, s, want)
	}
}

func TestMigrations_LogsSchema(t *testing.T) {
	s := readMigration(t, "005_logs.sql")
	require.Contains(t, s, "CREATE TABLE logs")
	for _, want := range []string{"event_type", "raw_data", "created_at"} {
		require.Contains(t, s, want)
	}
}