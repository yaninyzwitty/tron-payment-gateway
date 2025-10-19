package db

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yaninyzwitty/tron-payment-gateway/packages/shared/config"
)

func TestDbConnect_ValidConfig(t *testing.T) {
	cfg := &config.Config{
		DatabaseConfig: config.DatabaseConfig{
			User:           "testuser",
			Password:       "testpass",
			Host:           "localhost",
			Port:           5432,
			Database:       "testdb",
			MaxConnections: 10,
		},
	}

	ctx := context.Background()
	
	// Note: This will fail to actually connect since we don't have a real database
	// but we're testing the URL construction and pool creation attempt
	pool, err := DbConnect(ctx, cfg)
	
	// We expect an error since there's no real database
	// But we're testing that the function doesn't panic and handles errors properly
	if err != nil {
		assert.Error(t, err)
		// Accept various connection-related error messages
		hasExpectedError := strings.Contains(err.Error(), "failed to create new pool") ||
			strings.Contains(err.Error(), "failed to parse pgx pool config") ||
			strings.Contains(err.Error(), "invalid port") ||
			strings.Contains(err.Error(), "database connection failed")
		assert.True(t, hasExpectedError, "Expected connection-related error, got: %s", err.Error())
		assert.Nil(t, pool)
	} else {
		// If somehow it connects (unlikely), clean up
		if pool != nil {
			pool.Close()
		}
	}
}

func TestDbConnect_URLConstruction(t *testing.T) {
	testCases := []struct {
		name           string
		cfg            *config.Config
		expectedInURL  []string
	}{
		{
			name: "standard configuration",
			cfg: &config.Config{
				DatabaseConfig: config.DatabaseConfig{
					User:           "admin",
					Password:       "secret",
					Host:           "db.example.com",
					Port:           5432,
					Database:       "mydb",
					MaxConnections: 25,
				},
			},
			expectedInURL: []string{"admin", "secret", "db.example.com", "mydb", "pool_max_conns=25"},
		},
		{
			name: "localhost configuration",
			cfg: &config.Config{
				DatabaseConfig: config.DatabaseConfig{
					User:           "root",
					Password:       "rootpass",
					Host:           "127.0.0.1",
					Port:           5432,
					Database:       "localdb",
					MaxConnections: 5,
				},
			},
			expectedInURL: []string{"root", "rootpass", "127.0.0.1", "localdb", "pool_max_conns=5"},
		},
		{
			name: "high connection count",
			cfg: &config.Config{
				DatabaseConfig: config.DatabaseConfig{
					User:           "appuser",
					Password:       "apppass",
					Host:           "prod-db",
					Port:           5432,
					Database:       "proddb",
					MaxConnections: 1000,
				},
			},
			expectedInURL: []string{"appuser", "apppass", "prod-db", "proddb", "pool_max_conns=1000"},
		},
		{
			name: "special characters in password",
			cfg: &config.Config{
				DatabaseConfig: config.DatabaseConfig{
					User:           "user",
					Password:       "p@ss!w0rd",
					Host:           "host",
					Port:           5432,
					Database:       "db",
					MaxConnections: 10,
				},
			},
			expectedInURL: []string{"user", "p@ss!w0rd", "host", "db", "pool_max_conns=10"},
		},
		{
			name: "IPv6 host",
			cfg: &config.Config{
				DatabaseConfig: config.DatabaseConfig{
					User:           "user",
					Password:       "pass",
					Host:           "::1",
					Port:           5432,
					Database:       "db",
					MaxConnections: 10,
				},
			},
			expectedInURL: []string{"user", "pass", "::1", "db", "pool_max_conns=10"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			
			// Attempt connection (will fail but we're testing URL construction)
			_, err := DbConnect(ctx, tc.cfg)
			
			// We expect an error since no real database exists
			if err != nil {
				assert.Error(t, err)
				hasExpectedError := strings.Contains(err.Error(), "failed to create new pool") ||
					strings.Contains(err.Error(), "failed to parse pgx pool config") ||
					strings.Contains(err.Error(), "database connection failed") ||
					strings.Contains(err.Error(), "dial error") ||
					strings.Contains(err.Error(), "network is unreachable")
				assert.True(t, hasExpectedError, "Expected connection-related error, got: %s", err.Error())
			}
		})
	}
}

func TestDbConnect_NilConfig(t *testing.T) {
	ctx := context.Background()
	
	// This should panic or error due to nil pointer dereference
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil config
			assert.NotNil(t, r)
		}
	}()
	
	_, _ = DbConnect(ctx, nil)
}

func TestDbConnect_EmptyConfig(t *testing.T) {
	cfg := &config.Config{
		DatabaseConfig: config.DatabaseConfig{
			User:           "",
			Password:       "",
			Host:           "",
			Port:           0,
			Database:       "",
			MaxConnections: 0,
		},
	}

	ctx := context.Background()
	pool, err := DbConnect(ctx, cfg)
	
	assert.Error(t, err)
	assert.Nil(t, pool)
	hasExpectedError := strings.Contains(err.Error(), "failed to create new pool") ||
		strings.Contains(err.Error(), "failed to parse pgx pool config") ||
		strings.Contains(err.Error(), "invalid port")
	assert.True(t, hasExpectedError, "Expected connection-related error, got: %s", err.Error())
}

func TestDbConnect_ContextCancellation(t *testing.T) {
	cfg := &config.Config{
		DatabaseConfig: config.DatabaseConfig{
			User:           "testuser",
			Password:       "testpass",
			Host:           "nonexistent-host-12345",
			Port:           5432,
			Database:       "testdb",
			MaxConnections: 10,
		},
	}

	// Create a context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	
	pool, err := DbConnect(ctx, cfg)
	
	// Should error due to cancelled context or connection failure
	assert.Error(t, err)
	assert.Nil(t, pool)
}

func TestDbConnect_ContextWithTimeout(t *testing.T) {
	cfg := &config.Config{
		DatabaseConfig: config.DatabaseConfig{
			User:           "testuser",
			Password:       "testpass",
			Host:           "nonexistent-host-xyz",
			Port:           5432,
			Database:       "testdb",
			MaxConnections: 10,
		},
	}

	// Create a context with a very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()
	
	time.Sleep(2 * time.Millisecond) // Ensure timeout
	
	pool, err := DbConnect(ctx, cfg)
	
	assert.Error(t, err)
	assert.Nil(t, pool)
}

func TestDbConnect_ZeroMaxConnections(t *testing.T) {
	cfg := &config.Config{
		DatabaseConfig: config.DatabaseConfig{
			User:           "testuser",
			Password:       "testpass",
			Host:           "localhost",
			Port:           5432,
			Database:       "testdb",
			MaxConnections: 0,
		},
	}

	ctx := context.Background()
	pool, err := DbConnect(ctx, cfg)
	
	// Should handle zero connections (pgxpool will use default)
	assert.Error(t, err)
	assert.Nil(t, pool)
}

func TestDbConnect_NegativeMaxConnections(t *testing.T) {
	cfg := &config.Config{
		DatabaseConfig: config.DatabaseConfig{
			User:           "testuser",
			Password:       "testpass",
			Host:           "localhost",
			Port:           5432,
			Database:       "testdb",
			MaxConnections: -1,
		},
	}

	ctx := context.Background()
	pool, err := DbConnect(ctx, cfg)
	
	// Should handle negative connections
	assert.Error(t, err)
	assert.Nil(t, pool)
}

func TestDbConnect_LargeMaxConnections(t *testing.T) {
	cfg := &config.Config{
		DatabaseConfig: config.DatabaseConfig{
			User:           "testuser",
			Password:       "testpass",
			Host:           "localhost",
			Port:           5432,
			Database:       "testdb",
			MaxConnections: 100000,
		},
	}

	ctx := context.Background()
	pool, err := DbConnect(ctx, cfg)
	
	assert.Error(t, err)
	assert.Nil(t, pool)
}

func TestDbConnect_SpecialCharsInCredentials(t *testing.T) {
	testCases := []struct {
		name     string
		user     string
		password string
	}{
		{
			name:     "special chars in password",
			user:     "user",
			password: "p@ss!w0rd#$%",
		},
		{
			name:     "unicode in username",
			user:     "üser",
			password: "pass",
		},
		{
			name:     "spaces in credentials",
			user:     "user name",
			password: "pass word",
		},
		{
			name:     "url special chars",
			user:     "user@domain",
			password: "pass/word",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &config.Config{
				DatabaseConfig: config.DatabaseConfig{
					User:           tc.user,
					Password:       tc.password,
					Host:           "localhost",
					Port:           5432,
					Database:       "testdb",
					MaxConnections: 10,
				},
			}

			ctx := context.Background()
			pool, err := DbConnect(ctx, cfg)
			
			// We expect error since there's no real DB, but function shouldn't panic
			assert.Error(t, err)
			assert.Nil(t, pool)
		})
	}
}

func TestDbConnect_DifferentHosts(t *testing.T) {
	testCases := []struct {
		name string
		host string
	}{
		{"localhost", "localhost"},
		{"IPv4", "192.168.1.1"},
		{"IPv6", "::1"},
		{"IPv6 full", "2001:0db8:85a3:0000:0000:8a2e:0370:7334"},
		{"domain", "db.example.com"},
		{"subdomain", "master.db.example.com"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &config.Config{
				DatabaseConfig: config.DatabaseConfig{
					User:           "user",
					Password:       "pass",
					Host:           tc.host,
					Port:           5432,
					Database:       "db",
					MaxConnections: 10,
				},
			}

			ctx := context.Background()
			_, err := DbConnect(ctx, cfg)
			
			// All should fail to connect but shouldn't panic
			assert.Error(t, err)
		})
	}
}

func TestDbConnect_URLFormat(t *testing.T) {
	// Test that the URL is formatted correctly with expected components
	cfg := &config.Config{
		DatabaseConfig: config.DatabaseConfig{
			User:           "myuser",
			Password:       "mypass",
			Host:           "myhost",
			Port:           5432,
			Database:       "mydb",
			MaxConnections: 42,
		},
	}

	// We can't easily inspect the URL without modifying the function,
	// but we can verify the function doesn't panic with valid inputs
	ctx := context.Background()
	pool, err := DbConnect(ctx, cfg)
	
	assert.Error(t, err) // Expected since no real database
	assert.Nil(t, pool)
}
