package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_LoadConfig_ValidYAML(t *testing.T) {
	// Create a temporary valid config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	validYAML := `
debug: true
appPort: 8080
database:
  user: testuser
  password: testpass
  host: localhost
  database: testdb
  maxConnections: 10
`
	err := os.WriteFile(configPath, []byte(validYAML), 0644)
	require.NoError(t, err)

	// Test loading the config
	var cfg Config
	err = cfg.LoadConfig(configPath)

	assert.NoError(t, err)
	assert.True(t, cfg.Debug)
	assert.Equal(t, 8080, cfg.AppPort)
	assert.Equal(t, "testuser", cfg.DatabaseConfig.User)
	assert.Equal(t, "testpass", cfg.DatabaseConfig.Password)
	assert.Equal(t, "localhost", cfg.DatabaseConfig.Host)
	assert.Equal(t, "testdb", cfg.DatabaseConfig.Database)
	assert.Equal(t, 10, cfg.DatabaseConfig.MaxConnections)
}

func TestConfig_LoadConfig_MinimalConfig(t *testing.T) {
	// Test with minimal valid config (all optional fields)
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "minimal.yaml")

	minimalYAML := `
debug: false
appPort: 0
database:
  user: ""
  password: ""
  host: ""
  database: ""
  maxConnections: 0
`
	err := os.WriteFile(configPath, []byte(minimalYAML), 0644)
	require.NoError(t, err)

	var cfg Config
	err = cfg.LoadConfig(configPath)

	assert.NoError(t, err)
	assert.False(t, cfg.Debug)
	assert.Equal(t, 0, cfg.AppPort)
	assert.Equal(t, "", cfg.DatabaseConfig.User)
	assert.Equal(t, "", cfg.DatabaseConfig.Password)
}

func TestConfig_LoadConfig_DebugDisabled(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	yaml := `
debug: false
appPort: 3000
database:
  user: produser
  password: prodpass
  host: prod.example.com
  database: proddb
  maxConnections: 50
`
	err := os.WriteFile(configPath, []byte(yaml), 0644)
	require.NoError(t, err)

	var cfg Config
	err = cfg.LoadConfig(configPath)

	assert.NoError(t, err)
	assert.False(t, cfg.Debug)
	assert.Equal(t, 3000, cfg.AppPort)
}

func TestConfig_LoadConfig_HighPortNumber(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	yaml := `
debug: true
appPort: 65535
database:
  user: user
  password: pass
  host: host
  database: db
  maxConnections: 100
`
	err := os.WriteFile(configPath, []byte(yaml), 0644)
	require.NoError(t, err)

	var cfg Config
	err = cfg.LoadConfig(configPath)

	assert.NoError(t, err)
	assert.Equal(t, 65535, cfg.AppPort)
}

func TestConfig_LoadConfig_LargeMaxConnections(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	yaml := `
debug: false
appPort: 8000
database:
  user: user
  password: pass
  host: host
  database: db
  maxConnections: 10000
`
	err := os.WriteFile(configPath, []byte(yaml), 0644)
	require.NoError(t, err)

	var cfg Config
	err = cfg.LoadConfig(configPath)

	assert.NoError(t, err)
	assert.Equal(t, 10000, cfg.DatabaseConfig.MaxConnections)
}

func TestConfig_LoadConfig_SpecialCharactersInPassword(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	yaml := `
debug: true
appPort: 9000
database:
  user: admin
  password: "p@ss!w0rd#$%^&*()"
  host: db.example.com
  database: mydb
  maxConnections: 25
`
	err := os.WriteFile(configPath, []byte(yaml), 0644)
	require.NoError(t, err)

	var cfg Config
	err = cfg.LoadConfig(configPath)

	assert.NoError(t, err)
	assert.Equal(t, "p@ss!w0rd#$%^&*()", cfg.DatabaseConfig.Password)
}

func TestConfig_LoadConfig_IPv6Host(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	yaml := `
debug: true
appPort: 8080
database:
  user: user
  password: pass
  host: "::1"
  database: db
  maxConnections: 5
`
	err := os.WriteFile(configPath, []byte(yaml), 0644)
	require.NoError(t, err)

	var cfg Config
	err = cfg.LoadConfig(configPath)

	assert.NoError(t, err)
	assert.Equal(t, "::1", cfg.DatabaseConfig.Host)
}

func TestConfig_LoadConfig_FileNotFound(t *testing.T) {
	var cfg Config
	err := cfg.LoadConfig("/nonexistent/path/config.yaml")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read config")
}

func TestConfig_LoadConfig_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid.yaml")

	invalidYAML := `
debug: true
appPort: not_a_number
database:
  user: user
`
	err := os.WriteFile(configPath, []byte(invalidYAML), 0644)
	require.NoError(t, err)

	var cfg Config
	err = cfg.LoadConfig(configPath)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse config")
}

func TestConfig_LoadConfig_MalformedYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "malformed.yaml")

	malformedYAML := `
debug: true
  appPort: 8080
	database:
user: test
`
	err := os.WriteFile(configPath, []byte(malformedYAML), 0644)
	require.NoError(t, err)

	var cfg Config
	err = cfg.LoadConfig(configPath)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse config")
}

func TestConfig_LoadConfig_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "empty.yaml")

	err := os.WriteFile(configPath, []byte(""), 0644)
	require.NoError(t, err)

	var cfg Config
	err = cfg.LoadConfig(configPath)

	// Empty YAML should parse successfully with zero values
	assert.NoError(t, err)
	assert.False(t, cfg.Debug)
	assert.Equal(t, 0, cfg.AppPort)
	assert.Equal(t, "", cfg.DatabaseConfig.User)
}

func TestConfig_LoadConfig_PartialConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "partial.yaml")

	partialYAML := `
debug: true
appPort: 5000
`
	err := os.WriteFile(configPath, []byte(partialYAML), 0644)
	require.NoError(t, err)

	var cfg Config
	err = cfg.LoadConfig(configPath)

	assert.NoError(t, err)
	assert.True(t, cfg.Debug)
	assert.Equal(t, 5000, cfg.AppPort)
	// Database config should have zero values
	assert.Equal(t, "", cfg.DatabaseConfig.User)
	assert.Equal(t, 0, cfg.DatabaseConfig.MaxConnections)
}

func TestConfig_LoadConfig_NoPermission(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("Skipping permission test when running as root")
	}

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "noperm.yaml")

	err := os.WriteFile(configPath, []byte("debug: true"), 0000)
	require.NoError(t, err)

	var cfg Config
	err = cfg.LoadConfig(configPath)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read config")
}

func TestConfig_LoadConfig_Directory(t *testing.T) {
	tmpDir := t.TempDir()

	var cfg Config
	err := cfg.LoadConfig(tmpDir)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read config")
}

func TestConfig_LoadConfig_MultipleLoads(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	firstYAML := `
debug: true
appPort: 8080
database:
  user: firstuser
  password: firstpass
  host: localhost
  database: firstdb
  maxConnections: 10
`
	err := os.WriteFile(configPath, []byte(firstYAML), 0644)
	require.NoError(t, err)

	var cfg Config
	err = cfg.LoadConfig(configPath)
	assert.NoError(t, err)
	assert.Equal(t, "firstuser", cfg.DatabaseConfig.User)
	assert.Equal(t, 8080, cfg.AppPort)

	// Load again with different content
	secondYAML := `
debug: false
appPort: 9090
database:
  user: seconduser
  password: secondpass
  host: remotehost
  database: seconddb
  maxConnections: 20
`
	err = os.WriteFile(configPath, []byte(secondYAML), 0644)
	require.NoError(t, err)

	err = cfg.LoadConfig(configPath)
	assert.NoError(t, err)
	assert.Equal(t, "seconduser", cfg.DatabaseConfig.User)
	assert.Equal(t, 9090, cfg.AppPort)
	assert.False(t, cfg.Debug)
}

func TestConfig_LoadConfig_NegativePort(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	yaml := `
debug: true
appPort: -1
database:
  user: user
  password: pass
  host: host
  database: db
  maxConnections: 10
`
	err := os.WriteFile(configPath, []byte(yaml), 0644)
	require.NoError(t, err)

	var cfg Config
	err = cfg.LoadConfig(configPath)

	assert.NoError(t, err)
	assert.Equal(t, -1, cfg.AppPort)
}

func TestConfig_LoadConfig_NegativeMaxConnections(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	yaml := `
debug: true
appPort: 8080
database:
  user: user
  password: pass
  host: host
  database: db
  maxConnections: -5
`
	err := os.WriteFile(configPath, []byte(yaml), 0644)
	require.NoError(t, err)

	var cfg Config
	err = cfg.LoadConfig(configPath)

	assert.NoError(t, err)
	assert.Equal(t, -5, cfg.DatabaseConfig.MaxConnections)
}

func TestConfig_LoadConfig_ZeroValues(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	yaml := `
debug: false
appPort: 0
database:
  user: ""
  password: ""
  host: ""
  database: ""
  maxConnections: 0
`
	err := os.WriteFile(configPath, []byte(yaml), 0644)
	require.NoError(t, err)

	var cfg Config
	err = cfg.LoadConfig(configPath)

	assert.NoError(t, err)
	assert.False(t, cfg.Debug)
	assert.Equal(t, 0, cfg.AppPort)
	assert.Equal(t, "", cfg.DatabaseConfig.User)
	assert.Equal(t, "", cfg.DatabaseConfig.Password)
	assert.Equal(t, "", cfg.DatabaseConfig.Host)
	assert.Equal(t, "", cfg.DatabaseConfig.Database)
	assert.Equal(t, 0, cfg.DatabaseConfig.MaxConnections)
}

func TestDatabaseConfig_Fields(t *testing.T) {
	dbConfig := DatabaseConfig{
		User:           "testuser",
		Password:       "testpass",
		Host:           "testhost",
		Database:       "testdb",
		MaxConnections: 100,
	}

	assert.Equal(t, "testuser", dbConfig.User)
	assert.Equal(t, "testpass", dbConfig.Password)
	assert.Equal(t, "testhost", dbConfig.Host)
	assert.Equal(t, "testdb", dbConfig.Database)
	assert.Equal(t, 100, dbConfig.MaxConnections)
}

func TestConfig_StructFields(t *testing.T) {
	cfg := Config{
		Debug:   true,
		AppPort: 8080,
		DatabaseConfig: DatabaseConfig{
			User:           "user",
			Password:       "pass",
			Host:           "host",
			Database:       "db",
			MaxConnections: 50,
		},
	}

	assert.True(t, cfg.Debug)
	assert.Equal(t, 8080, cfg.AppPort)
	assert.Equal(t, "user", cfg.DatabaseConfig.User)
	assert.Equal(t, "pass", cfg.DatabaseConfig.Password)
	assert.Equal(t, "host", cfg.DatabaseConfig.Host)
	assert.Equal(t, "db", cfg.DatabaseConfig.Database)
	assert.Equal(t, 50, cfg.DatabaseConfig.MaxConnections)
}