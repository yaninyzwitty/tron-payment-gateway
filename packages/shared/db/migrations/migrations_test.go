package migrations

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestMigrationFilesExist verifies all expected migration files exist
func TestMigrationFilesExist(t *testing.T) {
	expectedFiles := []string{
		"002_accounts.sql",
		"003_payments.sql",
		"004_payments_attempts.sql",
		"005_logs.sql",
	}

	for _, filename := range expectedFiles {
		path := filepath.Join(".", filename)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Expected migration file does not exist: %s", filename)
		}
	}
}

// TestAccountsMigrationContent validates 002_accounts.sql contains expected schema
func TestAccountsMigrationContent(t *testing.T) {
	content, err := os.ReadFile("002_accounts.sql")
	if err != nil {
		t.Fatalf("Failed to read accounts migration: %v", err)
	}

	contentStr := string(content)

	// Check for required columns
	requiredElements := []string{
		"CREATE TABLE accounts",
		"id UUID PRIMARY KEY",
		"client_id UUID NOT NULL",
		"name STRING NOT NULL",
		"address_index INT DEFAULT 0",
		"created_at TIMESTAMPTZ DEFAULT now()",
		"FOREIGN KEY (client_id) REFERENCES clients(id)",
	}

	for _, element := range requiredElements {
		if !strings.Contains(contentStr, element) {
			t.Errorf("Accounts migration missing expected element: %s", element)
		}
	}
}

// TestPaymentsMigrationContent validates 003_payments.sql contains expected schema
func TestPaymentsMigrationContent(t *testing.T) {
	content, err := os.ReadFile("003_payments.sql")
	if err != nil {
		t.Fatalf("Failed to read payments migration: %v", err)
	}

	contentStr := string(content)

	requiredElements := []string{
		"CREATE TABLE payments",
		"id UUID PRIMARY KEY",
		"client_id UUID NOT NULL",
		"account_id UUID NOT NULL",
		"amount DECIMAL(18,6) NOT NULL",
		"unique_wallet STRING NOT NULL",
		"status STRING NOT NULL DEFAULT 'PENDING'",
		"expires_at TIMESTAMPTZ NOT NULL",
		"confirmed_at TIMESTAMPTZ",
		"attempt_count INT DEFAULT 0",
		"created_at TIMESTAMPTZ DEFAULT now()",
		"REFERENCES clients(id)",
		"REFERENCES accounts(id)",
	}

	for _, element := range requiredElements {
		if !strings.Contains(contentStr, element) {
			t.Errorf("Payments migration missing expected element: %s", element)
		}
	}
}

// TestPaymentAttemptsMigrationContent validates 004_payments_attempts.sql
func TestPaymentAttemptsMigrationContent(t *testing.T) {
	content, err := os.ReadFile("004_payments_attempts.sql")
	if err != nil {
		t.Fatalf("Failed to read payment_attempts migration: %v", err)
	}

	contentStr := string(content)

	requiredElements := []string{
		"CREATE TABLE payment_attempts",
		"id UUID PRIMARY KEY",
		"payment_id UUID NOT NULL",
		"attempt_number INT NOT NULL",
		"generated_wallet STRING NOT NULL",
		"generated_at TIMESTAMPTZ DEFAULT now()",
		"REFERENCES payments(id)",
	}

	for _, element := range requiredElements {
		if !strings.Contains(contentStr, element) {
			t.Errorf("Payment attempts migration missing expected element: %s", element)
		}
	}
}

// TestLogsMigrationContent validates 005_logs.sql
func TestLogsMigrationContent(t *testing.T) {
	content, err := os.ReadFile("005_logs.sql")
	if err != nil {
		t.Fatalf("Failed to read logs migration: %v", err)
	}

	contentStr := string(content)

	requiredElements := []string{
		"CREATE TABLE logs",
		"id UUID PRIMARY KEY",
		"payment_id UUID",
		"event_type STRING NOT NULL",
		"message STRING",
		"raw_data JSONB",
		"created_at TIMESTAMPTZ DEFAULT now()",
		"REFERENCES payments(id)",
	}

	for _, element := range requiredElements {
		if !strings.Contains(contentStr, element) {
			t.Errorf("Logs migration missing expected element: %s", element)
		}
	}
}

// TestMigrationsHaveProperStructure validates basic SQL structure
func TestMigrationsHaveProperStructure(t *testing.T) {
	migrations := []string{
		"002_accounts.sql",
		"003_payments.sql",
		"004_payments_attempts.sql",
		"005_logs.sql",
	}

	for _, migration := range migrations {
		content, err := os.ReadFile(migration)
		if err != nil {
			t.Errorf("Failed to read migration %s: %v", migration, err)
			continue
		}

		contentStr := string(content)

		// Check for CREATE TABLE statement
		if !strings.Contains(contentStr, "CREATE TABLE") {
			t.Errorf("Migration %s missing CREATE TABLE statement", migration)
		}

		// Check for proper SQL syntax elements
		if !strings.Contains(contentStr, "(") || !strings.Contains(contentStr, ")") {
			t.Errorf("Migration %s has invalid SQL structure", migration)
		}

		// Check for semicolons (proper SQL termination)
		if !strings.Contains(contentStr, ";") {
			t.Errorf("Migration %s missing semicolons", migration)
		}
	}
}

// TestMigrationsNoSQLInjectionPatterns checks for suspicious patterns
func TestMigrationsNoSQLInjectionPatterns(t *testing.T) {
	migrations := []string{
		"002_accounts.sql",
		"003_payments.sql",
		"004_payments_attempts.sql",
		"005_logs.sql",
	}

	suspiciousPatterns := []string{
		"DROP DATABASE",
		"DROP SCHEMA",
		"TRUNCATE",
		"DELETE FROM clients",
		"'; --",
		"UNION SELECT",
	}

	for _, migration := range migrations {
		content, err := os.ReadFile(migration)
		if err != nil {
			t.Errorf("Failed to read migration %s: %v", migration, err)
			continue
		}

		contentStr := strings.ToUpper(string(content))

		for _, pattern := range suspiciousPatterns {
			if strings.Contains(contentStr, strings.ToUpper(pattern)) {
				t.Errorf("Migration %s contains suspicious pattern: %s", migration, pattern)
			}
		}
	}
}

// TestMigrationFileNaming validates migration file naming convention
func TestMigrationFileNaming(t *testing.T) {
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("Failed to read migrations directory: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		name := entry.Name()

		// Check naming pattern: NNN_name.sql
		if len(name) < 5 {
			continue
		}

		// First 3 characters should be digits
		prefix := name[:3]
		for _, c := range prefix {
			if c < '0' || c > '9' {
				t.Errorf("Migration file %s doesn't follow naming convention (NNN_name.sql)", name)
				break
			}
		}

		// Fourth character should be underscore
		if name[3] != '_' {
			t.Errorf("Migration file %s doesn't follow naming convention (NNN_name.sql)", name)
		}
	}
}

// TestAddressIndexDefaultValue validates address_index has correct default
func TestAddressIndexDefaultValue(t *testing.T) {
	content, err := os.ReadFile("002_accounts.sql")
	if err != nil {
		t.Fatalf("Failed to read accounts migration: %v", err)
	}

	contentStr := string(content)

	// Verify address_index has DEFAULT 0
	if !strings.Contains(contentStr, "address_index INT DEFAULT 0") &&
		!strings.Contains(contentStr, "address_index INT DEFAULT  0") {
		t.Error("address_index column missing or has incorrect default value")
	}
}

// TestPaymentStatusDefault validates payment status default value
func TestPaymentStatusDefault(t *testing.T) {
	content, err := os.ReadFile("003_payments.sql")
	if err != nil {
		t.Fatalf("Failed to read payments migration: %v", err)
	}

	contentStr := string(content)

	if !strings.Contains(contentStr, "DEFAULT 'PENDING'") {
		t.Error("Payment status missing or has incorrect default value")
	}
}

// TestForeignKeyConstraints validates all foreign key relationships
func TestForeignKeyConstraints(t *testing.T) {
	testCases := []struct {
		file              string
		expectedFK        string
		referencedTable   string
		onDeleteBehavior  string
	}{
		{"002_accounts.sql", "client_id", "clients(id)", "ON DELETE CASCADE"},
		{"003_payments.sql", "client_id", "clients(id)", "ON DELETE CASCADE"},
		{"003_payments.sql", "account_id", "accounts(id)", "ON DELETE CASCADE"},
		{"004_payments_attempts.sql", "payment_id", "payments(id)", "ON DELETE CASCADE"},
		{"005_logs.sql", "payment_id", "payments(id)", "ON DELETE CASCADE"},
	}

	for _, tc := range testCases {
		content, err := os.ReadFile(tc.file)
		if err != nil {
			t.Errorf("Failed to read %s: %v", tc.file, err)
			continue
		}

		contentStr := string(content)

		if !strings.Contains(contentStr, tc.referencedTable) {
			t.Errorf("%s missing foreign key reference to %s", tc.file, tc.referencedTable)
		}

		if tc.onDeleteBehavior != "" && !strings.Contains(contentStr, tc.onDeleteBehavior) {
			t.Errorf("%s missing or has incorrect ON DELETE behavior", tc.file)
		}
	}
}

// TestDecimalPrecision validates amount field precision in payments
func TestDecimalPrecision(t *testing.T) {
	content, err := os.ReadFile("003_payments.sql")
	if err != nil {
		t.Fatalf("Failed to read payments migration: %v", err)
	}

	contentStr := string(content)

	if !strings.Contains(contentStr, "DECIMAL(18,6)") {
		t.Error("Payment amount field missing or has incorrect precision")
	}
}

// TestJSONBColumn validates JSONB column in logs
func TestJSONBColumn(t *testing.T) {
	content, err := os.ReadFile("005_logs.sql")
	if err != nil {
		t.Fatalf("Failed to read logs migration: %v", err)
	}

	contentStr := string(content)

	if !strings.Contains(contentStr, "raw_data JSONB") {
		t.Error("Logs table missing JSONB raw_data column")
	}
}

// TestTimestampColumns validates timestamp columns have proper defaults
func TestTimestampColumns(t *testing.T) {
	testCases := []struct {
		file   string
		column string
	}{
		{"002_accounts.sql", "created_at"},
		{"003_payments.sql", "created_at"},
		{"004_payments_attempts.sql", "generated_at"},
		{"005_logs.sql", "created_at"},
	}

	for _, tc := range testCases {
		content, err := os.ReadFile(tc.file)
		if err != nil {
			t.Errorf("Failed to read %s: %v", tc.file, err)
			continue
		}

		contentStr := string(content)

		if !strings.Contains(contentStr, tc.column+" TIMESTAMPTZ DEFAULT now()") &&
			!strings.Contains(contentStr, tc.column+" TIMESTAMPTZ DEFAULT  now()") {
			t.Errorf("%s missing proper timestamp default for %s", tc.file, tc.column)
		}
	}
}