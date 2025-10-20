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
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			t.Errorf("expected migration file %s does not exist", filename)
		}
	}
}

// TestMigrationFilesNotEmpty verifies migration files are not empty
func TestMigrationFilesNotEmpty(t *testing.T) {
	files := []string{
		"002_accounts.sql",
		"003_payments.sql",
		"004_payments_attempts.sql",
		"005_logs.sql",
	}

	for _, filename := range files {
		content, err := os.ReadFile(filename)
		if err != nil {
			t.Errorf("failed to read %s: %v", filename, err)
			continue
		}

		if len(content) == 0 {
			t.Errorf("migration file %s is empty", filename)
		}
	}
}

// TestAccountsMigrationStructure validates accounts table migration
func TestAccountsMigrationStructure(t *testing.T) {
	content, err := os.ReadFile("002_accounts.sql")
	if err != nil {
		t.Fatalf("failed to read accounts migration: %v", err)
	}

	sqlContent := string(content)

	// Check for table creation
	if !strings.Contains(sqlContent, "CREATE TABLE accounts") {
		t.Error("accounts migration should create accounts table")
	}

	// Check for required columns
	requiredColumns := []string{
		"id UUID PRIMARY KEY",
		"client_id UUID NOT NULL",
		"name STRING NOT NULL",
		"address_index INT DEFAULT 0",
		"created_at TIMESTAMPTZ DEFAULT now()",
	}

	for _, col := range requiredColumns {
		if !strings.Contains(sqlContent, col) {
			t.Errorf("accounts migration missing column definition: %s", col)
		}
	}

	// Check for foreign key constraint
	if !strings.Contains(sqlContent, "FOREIGN KEY (client_id) REFERENCES clients(id)") {
		t.Error("accounts migration should have foreign key to clients table")
	}

	// Check for cascade delete
	if !strings.Contains(sqlContent, "ON DELETE CASCADE") {
		t.Error("accounts migration should have ON DELETE CASCADE")
	}
}

// TestPaymentsMigrationStructure validates payments table migration
func TestPaymentsMigrationStructure(t *testing.T) {
	content, err := os.ReadFile("003_payments.sql")
	if err != nil {
		t.Fatalf("failed to read payments migration: %v", err)
	}

	sqlContent := string(content)

	// Check for table creation
	if !strings.Contains(sqlContent, "CREATE TABLE payments") {
		t.Error("payments migration should create payments table")
	}

	// Check for required columns
	requiredColumns := []string{
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
	}

	for _, col := range requiredColumns {
		if !strings.Contains(sqlContent, col) {
			t.Errorf("payments migration missing column definition: %s", col)
		}
	}

	// Check for foreign key constraints
	foreignKeys := []string{
		"REFERENCES clients(id)",
		"REFERENCES accounts(id)",
	}

	for _, fk := range foreignKeys {
		if !strings.Contains(sqlContent, fk) {
			t.Errorf("payments migration missing foreign key: %s", fk)
		}
	}

	// Check for cascade deletes
	if !strings.Contains(sqlContent, "ON DELETE CASCADE") {
		t.Error("payments migration should have ON DELETE CASCADE")
	}
}

// TestPaymentAttemptsMigrationStructure validates payment_attempts table migration
func TestPaymentAttemptsMigrationStructure(t *testing.T) {
	content, err := os.ReadFile("004_payments_attempts.sql")
	if err != nil {
		t.Fatalf("failed to read payment_attempts migration: %v", err)
	}

	sqlContent := string(content)

	// Check for table creation
	if !strings.Contains(sqlContent, "CREATE TABLE payment_attempts") {
		t.Error("payment_attempts migration should create payment_attempts table")
	}

	// Check for required columns
	requiredColumns := []string{
		"id UUID PRIMARY KEY",
		"payment_id UUID NOT NULL",
		"attempt_number INT NOT NULL",
		"generated_wallet STRING NOT NULL",
		"generated_at TIMESTAMPTZ DEFAULT now()",
	}

	for _, col := range requiredColumns {
		if !strings.Contains(sqlContent, col) {
			t.Errorf("payment_attempts migration missing column definition: %s", col)
		}
	}

	// Check for foreign key to payments
	if !strings.Contains(sqlContent, "REFERENCES payments(id)") {
		t.Error("payment_attempts migration should reference payments table")
	}

	// Check for cascade delete
	if !strings.Contains(sqlContent, "ON DELETE CASCADE") {
		t.Error("payment_attempts migration should have ON DELETE CASCADE")
	}
}

// TestLogsMigrationStructure validates logs table migration
func TestLogsMigrationStructure(t *testing.T) {
	content, err := os.ReadFile("005_logs.sql")
	if err != nil {
		t.Fatalf("failed to read logs migration: %v", err)
	}

	sqlContent := string(content)

	// Check for table creation
	if !strings.Contains(sqlContent, "CREATE TABLE logs") {
		t.Error("logs migration should create logs table")
	}

	// Check for required columns
	requiredColumns := []string{
		"id UUID PRIMARY KEY",
		"payment_id UUID",
		"event_type STRING NOT NULL",
		"message STRING",
		"raw_data JSONB",
		"created_at TIMESTAMPTZ DEFAULT now()",
	}

	for _, col := range requiredColumns {
		if !strings.Contains(sqlContent, col) {
			t.Errorf("logs migration missing column definition: %s", col)
		}
	}

	// Check for optional foreign key (logs can exist without payment)
	if !strings.Contains(sqlContent, "REFERENCES payments(id)") {
		t.Error("logs migration should reference payments table")
	}

	// Check for cascade delete
	if !strings.Contains(sqlContent, "ON DELETE CASCADE") {
		t.Error("logs migration should have ON DELETE CASCADE")
	}
}

// TestMigrationFileNamingConvention validates file naming
func TestMigrationFileNamingConvention(t *testing.T) {
	files, err := filepath.Glob("*.sql")
	if err != nil {
		t.Fatalf("failed to glob migration files: %v", err)
	}

	for _, file := range files {
		// Check if file follows XXX_name.sql pattern
		if !strings.HasSuffix(file, ".sql") {
			t.Errorf("file %s should have .sql extension", file)
		}

		// Check if name starts with digits
		basename := filepath.Base(file)
		if len(basename) < 3 || !isDigit(basename[0]) || !isDigit(basename[1]) || !isDigit(basename[2]) {
			t.Errorf("migration file %s should start with 3 digits", file)
		}
	}
}

// TestMigrationsSQLSyntax checks for common SQL syntax issues
func TestMigrationsSQLSyntax(t *testing.T) {
	files := []string{
		"002_accounts.sql",
		"003_payments.sql",
		"004_payments_attempts.sql",
		"005_logs.sql",
	}

	for _, filename := range files {
		content, err := os.ReadFile(filename)
		if err != nil {
			t.Errorf("failed to read %s: %v", filename, err)
			continue
		}

		sqlContent := string(content)

		// Check for CREATE TABLE statements
		if !strings.Contains(sqlContent, "CREATE TABLE") {
			t.Errorf("%s should contain CREATE TABLE statement", filename)
		}

		// Check that table names don't have typos in common words
		if strings.Contains(strings.ToLower(sqlContent), "tabel") {
			t.Errorf("%s contains typo: 'tabel' instead of 'table'", filename)
		}

		// Check for proper UUID usage
		if strings.Contains(sqlContent, "UUID") && !strings.Contains(sqlContent, "PRIMARY KEY") && !strings.Contains(sqlContent, "DEFAULT gen_random_uuid()") && !strings.Contains(sqlContent, "REFERENCES") && !strings.Contains(sqlContent, "NOT NULL") {
			// This is a basic check - UUID columns should have some constraint
			t.Logf("%s: UUID columns should have constraints", filename)
		}
	}
}

// TestMigrationsDefaultValues validates default values are set appropriately
func TestMigrationsDefaultValues(t *testing.T) {
	testCases := []struct {
		file     string
		contains string
	}{
		{"002_accounts.sql", "address_index INT DEFAULT 0"},
		{"003_payments.sql", "status STRING NOT NULL DEFAULT 'PENDING'"},
		{"003_payments.sql", "attempt_count INT DEFAULT 0"},
		{"004_payments_attempts.sql", "generated_at TIMESTAMPTZ DEFAULT now()"},
		{"005_logs.sql", "created_at TIMESTAMPTZ DEFAULT now()"},
	}

	for _, tc := range testCases {
		content, err := os.ReadFile(tc.file)
		if err != nil {
			t.Errorf("failed to read %s: %v", tc.file, err)
			continue
		}

		if !strings.Contains(string(content), tc.contains) {
			t.Errorf("%s should contain default value: %s", tc.file, tc.contains)
		}
	}
}

// TestMigrationsNOTNULLConstraints validates NOT NULL constraints
func TestMigrationsNOTNULLConstraints(t *testing.T) {
	testCases := []struct {
		file   string
		column string
	}{
		{"002_accounts.sql", "client_id UUID NOT NULL"},
		{"002_accounts.sql", "name STRING NOT NULL"},
		{"003_payments.sql", "client_id UUID NOT NULL"},
		{"003_payments.sql", "account_id UUID NOT NULL"},
		{"003_payments.sql", "amount DECIMAL(18,6) NOT NULL"},
		{"003_payments.sql", "unique_wallet STRING NOT NULL"},
		{"003_payments.sql", "status STRING NOT NULL"},
		{"004_payments_attempts.sql", "payment_id UUID NOT NULL"},
		{"004_payments_attempts.sql", "attempt_number INT NOT NULL"},
		{"004_payments_attempts.sql", "generated_wallet STRING NOT NULL"},
		{"005_logs.sql", "event_type STRING NOT NULL"},
	}

	for _, tc := range testCases {
		content, err := os.ReadFile(tc.file)
		if err != nil {
			t.Errorf("failed to read %s: %v", tc.file, err)
			continue
		}

		if !strings.Contains(string(content), tc.column) {
			t.Errorf("%s should have NOT NULL constraint: %s", tc.file, tc.column)
		}
	}
}

// TestMigrationsPrimaryKeys validates primary key definitions
func TestMigrationsPrimaryKeys(t *testing.T) {
	files := []string{
		"002_accounts.sql",
		"003_payments.sql",
		"004_payments_attempts.sql",
		"005_logs.sql",
	}

	for _, filename := range files {
		content, err := os.ReadFile(filename)
		if err != nil {
			t.Errorf("failed to read %s: %v", filename, err)
			continue
		}

		sqlContent := string(content)

		if !strings.Contains(sqlContent, "id UUID PRIMARY KEY") {
			t.Errorf("%s should define id as UUID PRIMARY KEY", filename)
		}

		if !strings.Contains(sqlContent, "DEFAULT gen_random_uuid()") {
			t.Errorf("%s should use gen_random_uuid() for id default", filename)
		}
	}
}

// TestPaymentsDecimalPrecision validates decimal precision for amounts
func TestPaymentsDecimalPrecision(t *testing.T) {
	content, err := os.ReadFile("003_payments.sql")
	if err != nil {
		t.Fatalf("failed to read payments migration: %v", err)
	}

	sqlContent := string(content)

	// Check for proper decimal precision (18,6 allows for large amounts with 6 decimal places)
	if !strings.Contains(sqlContent, "DECIMAL(18,6)") {
		t.Error("payments amount should be DECIMAL(18,6) for proper precision")
	}
}

// TestLogsJSONBDataType validates JSONB is used for raw_data
func TestLogsJSONBDataType(t *testing.T) {
	content, err := os.ReadFile("005_logs.sql")
	if err != nil {
		t.Fatalf("failed to read logs migration: %v", err)
	}

	sqlContent := string(content)

	if !strings.Contains(sqlContent, "raw_data JSONB") {
		t.Error("logs raw_data should use JSONB data type for efficient JSON storage")
	}
}

// TestMigrationsComments validates that migrations have helpful comments
func TestMigrationsComments(t *testing.T) {
	files := []string{
		"003_payments.sql",
		"004_payments_attempts.sql",
		"005_logs.sql",
	}

	for _, filename := range files {
		content, err := os.ReadFile(filename)
		if err != nil {
			t.Errorf("failed to read %s: %v", filename, err)
			continue
		}

		sqlContent := string(content)

		// Check for SQL comments
		if !strings.Contains(sqlContent, "--") {
			t.Logf("%s: consider adding comments for better documentation", filename)
		}
	}
}

// Helper function to check if a byte is a digit
func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

// TestMigrationOrderConsistency validates migrations are in correct order
func TestMigrationOrderConsistency(t *testing.T) {
	migrations := []struct {
		number int
		file   string
	}{
		{2, "002_accounts.sql"},
		{3, "003_payments.sql"},
		{4, "004_payments_attempts.sql"},
		{5, "005_logs.sql"},
	}

	for i := 1; i < len(migrations); i++ {
		if migrations[i].number <= migrations[i-1].number {
			t.Errorf("migration order inconsistent: %s should come after %s", 
				migrations[i].file, migrations[i-1].file)
		}
	}
}

// TestAccountsAddressIndexColumn validates the new address_index column
func TestAccountsAddressIndexColumn(t *testing.T) {
	content, err := os.ReadFile("002_accounts.sql")
	if err != nil {
		t.Fatalf("failed to read accounts migration: %v", err)
	}

	sqlContent := string(content)

	// Verify address_index exists with proper type and default
	if !strings.Contains(sqlContent, "address_index INT") {
		t.Error("accounts table should have address_index INT column")
	}

	if !strings.Contains(sqlContent, "address_index INT DEFAULT 0") {
		t.Error("address_index should have DEFAULT 0")
	}
}

// TestReferentialIntegrityChain validates the foreign key dependency chain
func TestReferentialIntegrityChain(t *testing.T) {
	// The dependency chain should be: clients -> accounts -> payments -> payment_attempts/logs
	
	// Accounts depends on clients
	accountsContent, _ := os.ReadFile("002_accounts.sql")
	if !strings.Contains(string(accountsContent), "REFERENCES clients(id)") {
		t.Error("accounts should reference clients table")
	}

	// Payments depends on both clients and accounts
	paymentsContent, _ := os.ReadFile("003_payments.sql")
	if !strings.Contains(string(paymentsContent), "REFERENCES clients(id)") {
		t.Error("payments should reference clients table")
	}
	if !strings.Contains(string(paymentsContent), "REFERENCES accounts(id)") {
		t.Error("payments should reference accounts table")
	}

	// Payment attempts depends on payments
	attemptsContent, _ := os.ReadFile("004_payments_attempts.sql")
	if !strings.Contains(string(attemptsContent), "REFERENCES payments(id)") {
		t.Error("payment_attempts should reference payments table")
	}

	// Logs depends on payments (optional)
	logsContent, _ := os.ReadFile("005_logs.sql")
	if !strings.Contains(string(logsContent), "REFERENCES payments(id)") {
		t.Error("logs should reference payments table")
	}
}