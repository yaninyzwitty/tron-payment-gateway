package migrations

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Test that all migration files exist and are properly numbered
func TestMigrationFilesExist(t *testing.T) {
	expectedFiles := []string{
		"001_clients.sql",
		"002_accounts.sql",
		"003_payments.sql",
		"004_payments_attempts.sql",
		"005_logs.sql",
	}

	for _, file := range expectedFiles {
		_, err := os.Stat(file)
		if os.IsNotExist(err) {
			t.Errorf("Migration file %s does not exist", file)
		}
	}
}

// Test that migration files are properly sequenced
func TestMigrationSequencing(t *testing.T) {
	files, err := filepath.Glob("*.sql")
	if err != nil {
		t.Fatalf("Failed to list migration files: %v", err)
	}

	if len(files) < 5 {
		t.Errorf("Expected at least 5 migration files, found %d", len(files))
	}

	// Verify sequential numbering
	for i := 1; i <= 5; i++ {
		prefix := strings.Repeat("0", 3-len(string(rune('0'+i)))) + string(rune('0'+i))
		found := false
		for _, file := range files {
			if strings.HasPrefix(filepath.Base(file), prefix) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Missing migration with prefix %s", prefix)
		}
	}
}

// Test accounts migration contains required columns
func TestAccountsMigrationSchema(t *testing.T) {
	content, err := os.ReadFile("002_accounts.sql")
	if err != nil {
		t.Fatalf("Failed to read accounts migration: %v", err)
	}

	migration := string(content)

	requiredElements := []string{
		"CREATE TABLE accounts",
		"id UUID PRIMARY KEY",
		"client_id UUID NOT NULL",
		"name STRING NOT NULL",
		"address_index INT DEFAULT 0",
		"created_at TIMESTAMPTZ",
		"FOREIGN KEY",
		"REFERENCES clients(id)",
		"ON DELETE CASCADE",
	}

	for _, element := range requiredElements {
		if !strings.Contains(migration, element) {
			t.Errorf("Accounts migration missing required element: %s", element)
		}
	}
}

// Test payments migration contains required columns
func TestPaymentsMigrationSchema(t *testing.T) {
	content, err := os.ReadFile("003_payments.sql")
	if err != nil {
		t.Fatalf("Failed to read payments migration: %v", err)
	}

	migration := string(content)

	requiredElements := []string{
		"CREATE TABLE payments",
		"id UUID PRIMARY KEY",
		"client_id UUID NOT NULL",
		"account_id UUID NOT NULL",
		"amount DECIMAL",
		"unique_wallet STRING NOT NULL",
		"status STRING NOT NULL",
		"expires_at TIMESTAMPTZ NOT NULL",
		"confirmed_at TIMESTAMPTZ",
		"attempt_count INT DEFAULT 0",
		"created_at TIMESTAMPTZ",
		"REFERENCES clients(id)",
		"REFERENCES accounts(id)",
	}

	for _, element := range requiredElements {
		if !strings.Contains(migration, element) {
			t.Errorf("Payments migration missing required element: %s", element)
		}
	}
}

// Test payment_attempts migration contains required columns
func TestPaymentAttemptsMigrationSchema(t *testing.T) {
	content, err := os.ReadFile("004_payments_attempts.sql")
	if err != nil {
		t.Fatalf("Failed to read payment attempts migration: %v", err)
	}

	migration := string(content)

	requiredElements := []string{
		"CREATE TABLE payment_attempts",
		"id UUID PRIMARY KEY",
		"payment_id UUID NOT NULL",
		"attempt_number INT NOT NULL",
		"generated_wallet STRING NOT NULL",
		"generated_at TIMESTAMPTZ",
		"REFERENCES payments(id)",
		"ON DELETE CASCADE",
	}

	for _, element := range requiredElements {
		if !strings.Contains(migration, element) {
			t.Errorf("Payment attempts migration missing required element: %s", element)
		}
	}
}

// Test logs migration contains required columns
func TestLogsMigrationSchema(t *testing.T) {
	content, err := os.ReadFile("005_logs.sql")
	if err != nil {
		t.Fatalf("Failed to read logs migration: %v", err)
	}

	migration := string(content)

	requiredElements := []string{
		"CREATE TABLE logs",
		"id UUID PRIMARY KEY",
		"payment_id UUID",
		"event_type STRING NOT NULL",
		"message STRING",
		"raw_data JSONB",
		"created_at TIMESTAMPTZ",
		"REFERENCES payments(id)",
	}

	for _, element := range requiredElements {
		if !strings.Contains(migration, element) {
			t.Errorf("Logs migration missing required element: %s", element)
		}
	}
}

// Test that migrations don't contain dangerous operations
func TestMigrationsSafety(t *testing.T) {
	dangerousPatterns := []string{
		"DROP DATABASE",
		"DROP SCHEMA",
		"TRUNCATE",
		"DELETE FROM",
	}

	files, err := filepath.Glob("*.sql")
	if err != nil {
		t.Fatalf("Failed to list migration files: %v", err)
	}

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			t.Errorf("Failed to read %s: %v", file, err)
			continue
		}

		migration := strings.ToUpper(string(content))
		for _, pattern := range dangerousPatterns {
			if strings.Contains(migration, pattern) {
				t.Errorf("Migration %s contains dangerous pattern: %s", file, pattern)
			}
		}
	}
}

// Test that all migrations use UUID primary keys
func TestMigrationsUseUUID(t *testing.T) {
	files, err := filepath.Glob("*.sql")
	if err != nil {
		t.Fatalf("Failed to list migration files: %v", err)
	}

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			t.Errorf("Failed to read %s: %v", file, err)
			continue
		}

		migration := string(content)
		if strings.Contains(migration, "CREATE TABLE") {
			if !strings.Contains(migration, "id UUID PRIMARY KEY") {
				t.Errorf("Migration %s doesn't use UUID for primary key", file)
			}
		}
	}
}

// Test that foreign key relationships are properly defined
func TestMigrationsForeignKeys(t *testing.T) {
	testCases := []struct {
		file     string
		hasFK    bool
		refTable string
	}{
		{"002_accounts.sql", true, "clients"},
		{"003_payments.sql", true, "clients"},
		{"003_payments.sql", true, "accounts"},
		{"004_payments_attempts.sql", true, "payments"},
		{"005_logs.sql", true, "payments"},
	}

	for _, tc := range testCases {
		content, err := os.ReadFile(tc.file)
		if err != nil {
			t.Errorf("Failed to read %s: %v", tc.file, err)
			continue
		}

		migration := string(content)
		if tc.hasFK {
			if !strings.Contains(migration, "REFERENCES "+tc.refTable) {
				t.Errorf("Migration %s missing foreign key reference to %s", tc.file, tc.refTable)
			}
		}
	}
}

// Test that all tables have created_at timestamp
func TestMigrationsHaveTimestamps(t *testing.T) {
	files, err := filepath.Glob("*.sql")
	if err != nil {
		t.Fatalf("Failed to list migration files: %v", err)
	}

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			t.Errorf("Failed to read %s: %v", file, err)
			continue
		}

		migration := string(content)
		if strings.Contains(migration, "CREATE TABLE") {
			// Skip comment-only lines
			if !strings.Contains(migration, "--") || strings.Index(migration, "CREATE TABLE") < strings.Index(migration, "--") {
				if !strings.Contains(migration, "created_at TIMESTAMPTZ") && !strings.Contains(migration, "generated_at TIMESTAMPTZ") {
					t.Errorf("Migration %s missing timestamp column", file)
				}
			}
		}
	}
}

// Test migration file naming convention
func TestMigrationNamingConvention(t *testing.T) {
	files, err := filepath.Glob("*.sql")
	if err != nil {
		t.Fatalf("Failed to list migration files: %v", err)
	}

	for _, file := range files {
		basename := filepath.Base(file)
		
		// Should start with 3 digits
		if len(basename) < 4 || basename[0] < '0' || basename[0] > '9' {
			t.Errorf("Migration file %s doesn't follow naming convention (NNN_name.sql)", file)
		}

		// Should have underscore after number
		if len(basename) >= 4 && basename[3] != '_' {
			t.Errorf("Migration file %s doesn't have underscore after number", file)
		}

		// Should end with .sql
		if !strings.HasSuffix(basename, ".sql") {
			t.Errorf("Migration file %s doesn't have .sql extension", file)
		}
	}
}

// Test that address_index was added to accounts
func TestAccountsAddressIndexField(t *testing.T) {
	content, err := os.ReadFile("002_accounts.sql")
	if err != nil {
		t.Fatalf("Failed to read accounts migration: %v", err)
	}

	migration := string(content)

	if !strings.Contains(migration, "address_index") {
		t.Error("Accounts migration missing address_index column")
	}

	if !strings.Contains(migration, "address_index INT DEFAULT 0") {
		t.Error("address_index should be INT with DEFAULT 0")
	}
}

// Test that payment status has default value
func TestPaymentsDefaultStatus(t *testing.T) {
	content, err := os.ReadFile("003_payments.sql")
	if err != nil {
		t.Fatalf("Failed to read payments migration: %v", err)
	}

	migration := string(content)

	if !strings.Contains(migration, "status STRING NOT NULL DEFAULT 'PENDING'") {
		t.Error("Payments status should have DEFAULT 'PENDING'")
	}
}

// Test that tables are ordered correctly (dependencies)
func TestMigrationDependencyOrder(t *testing.T) {
	// Read all migration files
	migrations := make(map[string]string)
	files := []string{
		"001_clients.sql",
		"002_accounts.sql",
		"003_payments.sql",
		"004_payments_attempts.sql",
		"005_logs.sql",
	}

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("Failed to read %s: %v", file, err)
		}
		migrations[file] = string(content)
	}

	// Verify clients comes before accounts
	if !strings.Contains(migrations["002_accounts.sql"], "REFERENCES clients") {
		t.Error("Accounts should reference clients, so clients must be created first")
	}

	// Verify accounts comes before payments
	if !strings.Contains(migrations["003_payments.sql"], "REFERENCES accounts") {
		t.Error("Payments should reference accounts")
	}

	// Verify payments comes before payment_attempts
	if !strings.Contains(migrations["004_payments_attempts.sql"], "REFERENCES payments") {
		t.Error("Payment attempts should reference payments")
	}

	// Verify payments comes before logs
	if !strings.Contains(migrations["005_logs.sql"], "REFERENCES payments") {
		t.Error("Logs should reference payments")
	}
}