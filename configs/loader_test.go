package configs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestConfig represents a simple test configuration.
type TestConfig struct {
	Name       string        `mapstructure:"NAME" validate:"required" default:"test-app"`
	Port       int           `mapstructure:"PORT" validate:"required,port" default:"8080"`
	Timeout    time.Duration `mapstructure:"TIMEOUT" default:"5s"`
	Debug      bool          `mapstructure:"DEBUG" default:"false"`
	Database   TestDBConfig  `mapstructure:"DATABASE"`
	OptionalDB *TestDBConfig `mapstructure:"OPTIONAL_DB"`
}

// TestDBConfig represents a test database configuration.
type TestDBConfig struct {
	Host     string `mapstructure:"HOST" validate:"required" default:"localhost"`
	Port     int    `mapstructure:"PORT" validate:"required,port" default:"5432"`
	Database string `mapstructure:"DATABASE" validate:"required" default:"testdb"`
	Enabled  bool   `mapstructure:"ENABLED" default:"true"`
}

func TestNew(t *testing.T) {
	cfg := &TestConfig{}
	loader := New(cfg)

	if loader == nil {
		t.Fatal("New() returned nil")
	}

	if loader.config != cfg {
		t.Error("Config pointer not set correctly")
	}
}

func TestLoadWithConfigExample(t *testing.T) {
	cfg := &TestConfig{}

	// Create a temporary directory for config files
	tempDir := t.TempDir()

	// Create config.example.yaml
	configContent := `
name: "example-app"
port: 9000
timeout: "10s"
debug: true
database:
  host: "example-db.example.com"
  port: 5432
  database: "example_db"
  enabled: true
`

	configFile := filepath.Join(tempDir, "config.example.yaml")
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	err := New(cfg).Load(AppEnvironmentDev, tempDir)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Check that config values were loaded
	if cfg.Name != "example-app" {
		t.Errorf("Expected Name to be 'example-app', got '%s'", cfg.Name)
	}

	if cfg.Port != 9000 {
		t.Errorf("Expected Port to be 9000, got %d", cfg.Port)
	}

	if cfg.Timeout != 10*time.Second {
		t.Errorf("Expected Timeout to be 10s, got %v", cfg.Timeout)
	}

	if cfg.Debug != true {
		t.Errorf("Expected Debug to be true, got %v", cfg.Debug)
	}

	// Check nested struct values
	if cfg.Database.Host != "example-db.example.com" {
		t.Errorf("Expected Database.Host to be 'example-db.example.com', got '%s'", cfg.Database.Host)
	}

	if cfg.Database.Port != 5432 {
		t.Errorf("Expected Database.Port to be 5432, got %d", cfg.Database.Port)
	}
}

func TestLoadWithConfigFile(t *testing.T) {
	cfg := &TestConfig{}

	// Create a temporary directory for config files
	tempDir := t.TempDir()

	// Create a config.yaml file
	configContent := `
name: "custom-app"
port: 9000
timeout: "10s"
debug: true
database:
  host: "db.example.com"
  port: 3306
  database: "customdb"
`

	configFile := filepath.Join(tempDir, "config.yaml")
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	err := New(cfg).Load(AppEnvironmentDev, tempDir)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Check that config file values were loaded
	if cfg.Name != "custom-app" {
		t.Errorf("Expected Name to be 'custom-app', got '%s'", cfg.Name)
	}

	if cfg.Port != 9000 {
		t.Errorf("Expected Port to be 9000, got %d", cfg.Port)
	}

	if cfg.Timeout != 10*time.Second {
		t.Errorf("Expected Timeout to be 10s, got %v", cfg.Timeout)
	}

	if cfg.Debug != true {
		t.Errorf("Expected Debug to be true, got %v", cfg.Debug)
	}

	// Check nested struct values
	if cfg.Database.Host != "db.example.com" {
		t.Errorf("Expected Database.Host to be 'db.example.com', got '%s'", cfg.Database.Host)
	}

	if cfg.Database.Port != 3306 {
		t.Errorf("Expected Database.Port to be 3306, got %d", cfg.Database.Port)
	}

	if cfg.Database.Database != "customdb" {
		t.Errorf("Expected Database.Database to be 'customdb', got '%s'", cfg.Database.Database)
	}
}

func TestLoadWithEnvironmentSpecificConfig(t *testing.T) {
	cfg := &TestConfig{}

	// Create a temporary directory for config files
	tempDir := t.TempDir()

	// Create base config.yaml
	baseConfig := `
name: "base-app"
port: 8080
database:
  host: "localhost"
  port: 5432
`

	configFile := filepath.Join(tempDir, "config.yaml")
	if err := os.WriteFile(configFile, []byte(baseConfig), 0644); err != nil {
		t.Fatalf("Failed to write base config file: %v", err)
	}

	// Create environment-specific config
	envConfig := `
name: "prod-app"
port: 9000
database:
  host: "prod-db.example.com"
  port: 5432
`

	envConfigFile := filepath.Join(tempDir, "config.prod.yaml")
	if err := os.WriteFile(envConfigFile, []byte(envConfig), 0644); err != nil {
		t.Fatalf("Failed to write env config file: %v", err)
	}

	err := New(cfg).Load(AppEnvironmentProd, tempDir)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Check that environment-specific values were loaded
	if cfg.Name != "prod-app" {
		t.Errorf("Expected Name to be 'prod-app', got '%s'", cfg.Name)
	}

	if cfg.Port != 9000 {
		t.Errorf("Expected Port to be 9000, got %d", cfg.Port)
	}

	if cfg.Database.Host != "prod-db.example.com" {
		t.Errorf("Expected Database.Host to be 'prod-db.example.com', got '%s'", cfg.Database.Host)
	}
}

func TestLoadWithoutValidation(t *testing.T) {
	cfg := &TestConfig{}

	// Create a temporary directory for config files
	tempDir := t.TempDir()

	// This should not fail because validation is disabled by default
	err := New(cfg).Load(AppEnvironmentDev, tempDir)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Check that config was loaded (even with missing required fields)
	// Since validation is disabled, it should work
}

func TestLoadWithInvalidConfig(t *testing.T) {
	cfg := &TestConfig{}

	// Create a temporary directory for config files
	tempDir := t.TempDir()

	// Create a config with invalid port
	configContent := `
name: ""
port: 99999
`

	configFile := filepath.Join(tempDir, "config.yaml")
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// This should succeed since validation is disabled by default
	err := New(cfg).Load(AppEnvironmentDev, tempDir)
	if err != nil {
		t.Errorf("Load() failed even without validation: %v", err)
	}
}

func TestPointerFieldDefaults(t *testing.T) {
	cfg := &TestConfig{}

	// Create a temporary directory for config files
	tempDir := t.TempDir()

	err := New(cfg).Load(AppEnvironmentDev, tempDir)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// OptionalDB should be nil since it's a pointer and no config was provided
	if cfg.OptionalDB != nil {
		t.Error("Expected OptionalDB to be nil")
	}
}

func TestFluentAPI(t *testing.T) {
	cfg := &TestConfig{}

	// Test fluent API chaining
	loader := New(cfg)
	if loader == nil {
		t.Fatal("Fluent API returned nil")
	}

	// Test that we can still call Load after chaining
	tempDir := t.TempDir()
	err := loader.Load(AppEnvironmentDev, tempDir)
	if err != nil {
		t.Fatalf("Load() failed after fluent chaining: %v", err)
	}
}

func TestLoadWithValidationSuccess(t *testing.T) {
	cfg := &TestConfig{}

	// Create a temporary directory for config files
	tempDir := t.TempDir()

	// Create config.yaml with valid data
	configContent := `
name: "valid-app"
port: 8080
timeout: "5s"
debug: true
database:
  host: "localhost"
  port: 5432
  database: "testdb"
  enabled: true
`

	configFile := filepath.Join(tempDir, "config.yaml")
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Define validation function
	validator := func(cfg *TestConfig) error {
		if cfg.Name == "" {
			return fmt.Errorf("name is required")
		}
		if cfg.Port <= 0 || cfg.Port > 65535 {
			return fmt.Errorf("port must be between 1 and 65535")
		}
		if cfg.Database.Host == "" {
			return fmt.Errorf("database host is required")
		}
		return nil
	}

	err := New(cfg).WithValidation(validator).Load(AppEnvironmentDev, tempDir)
	if err != nil {
		t.Fatalf("Load() with validation failed: %v", err)
	}

	// Verify config was loaded correctly
	if cfg.Name != "valid-app" {
		t.Errorf("Expected Name to be 'valid-app', got '%s'", cfg.Name)
	}
	if cfg.Port != 8080 {
		t.Errorf("Expected Port to be 8080, got %d", cfg.Port)
	}
}

func TestLoadWithValidationFailure(t *testing.T) {
	cfg := &TestConfig{}

	// Create a temporary directory for config files
	tempDir := t.TempDir()

	// Create config.yaml with invalid data
	configContent := `
name: ""
port: 0
database:
  host: ""
`

	configFile := filepath.Join(tempDir, "config.yaml")
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Define validation function that will fail
	validator := func(cfg *TestConfig) error {
		if cfg.Name == "" {
			return fmt.Errorf("name is required")
		}
		if cfg.Port <= 0 {
			return fmt.Errorf("port must be greater than 0")
		}
		if cfg.Database.Host == "" {
			return fmt.Errorf("database host is required")
		}
		return nil
	}

	err := New(cfg).WithValidation(validator).Load(AppEnvironmentDev, tempDir)
	if err == nil {
		t.Fatal("Expected Load() to fail with validation error, but it succeeded")
	}

	// Check that the error message contains validation information
	if !strings.Contains(err.Error(), "config validation failed") {
		t.Errorf("Expected error to contain 'config validation failed', got: %v", err)
	}
}

func TestLoadWithValidationChaining(t *testing.T) {
	cfg := &TestConfig{}

	// Create a temporary directory for config files
	tempDir := t.TempDir()

	// Create config.yaml
	configContent := `
name: "chained-app"
port: 3000
database:
  host: "localhost"
  port: 5432
`

	configFile := filepath.Join(tempDir, "config.yaml")
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Test fluent API chaining with validation
	validator := func(cfg *TestConfig) error {
		if cfg.Name == "" {
			return fmt.Errorf("name is required")
		}
		return nil
	}

	loader := New(cfg).WithValidation(validator)
	if loader == nil {
		t.Fatal("Fluent API with validation returned nil")
	}

	err := loader.Load(AppEnvironmentDev, tempDir)
	if err != nil {
		t.Fatalf("Load() with chained validation failed: %v", err)
	}

	// Verify config was loaded
	if cfg.Name != "chained-app" {
		t.Errorf("Expected Name to be 'chained-app', got '%s'", cfg.Name)
	}
}

func TestLoadWithoutValidationCallback(t *testing.T) {
	cfg := &TestConfig{}

	// Create a temporary directory for config files
	tempDir := t.TempDir()

	// Create config.yaml
	configContent := `
name: "no-validation-app"
port: 4000
`

	configFile := filepath.Join(tempDir, "config.yaml")
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Load without validation callback - should succeed
	err := New(cfg).Load(AppEnvironmentDev, tempDir)
	if err != nil {
		t.Fatalf("Load() without validation failed: %v", err)
	}

	// Verify config was loaded
	if cfg.Name != "no-validation-app" {
		t.Errorf("Expected Name to be 'no-validation-app', got '%s'", cfg.Name)
	}
}

func TestValidationWithComplexConfig(t *testing.T) {
	cfg := &TestConfig{}

	// Create a temporary directory for config files
	tempDir := t.TempDir()

	// Create config.yaml with complex nested structure
	configContent := `
name: "complex-app"
port: 5000
timeout: "30s"
debug: false
database:
  host: "db.example.com"
  port: 3306
  database: "complexdb"
  enabled: true
optional_db:
  host: "optional.example.com"
  port: 5432
  database: "optionaldb"
  enabled: false
`

	configFile := filepath.Join(tempDir, "config.yaml")
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Define comprehensive validation function
	validator := func(cfg *TestConfig) error {
		// Validate main config
		if cfg.Name == "" {
			return fmt.Errorf("name is required")
		}
		if cfg.Port <= 0 || cfg.Port > 65535 {
			return fmt.Errorf("port must be between 1 and 65535")
		}
		if cfg.Timeout <= 0 {
			return fmt.Errorf("timeout must be positive")
		}

		// Validate database config
		if cfg.Database.Host == "" {
			return fmt.Errorf("database host is required")
		}
		if cfg.Database.Port <= 0 || cfg.Database.Port > 65535 {
			return fmt.Errorf("database port must be between 1 and 65535")
		}
		if cfg.Database.Database == "" {
			return fmt.Errorf("database name is required")
		}

		// Validate optional database if present
		if cfg.OptionalDB != nil {
			if cfg.OptionalDB.Host == "" {
				return fmt.Errorf("optional database host cannot be empty if database is specified")
			}
			if cfg.OptionalDB.Port <= 0 || cfg.OptionalDB.Port > 65535 {
				return fmt.Errorf("optional database port must be between 1 and 65535")
			}
		}

		return nil
	}

	err := New(cfg).WithValidation(validator).Load(AppEnvironmentDev, tempDir)
	if err != nil {
		t.Fatalf("Load() with complex validation failed: %v", err)
	}

	// Verify all config was loaded correctly
	if cfg.Name != "complex-app" {
		t.Errorf("Expected Name to be 'complex-app', got '%s'", cfg.Name)
	}
	if cfg.Port != 5000 {
		t.Errorf("Expected Port to be 5000, got %d", cfg.Port)
	}
	if cfg.Database.Host != "db.example.com" {
		t.Errorf("Expected Database.Host to be 'db.example.com', got '%s'", cfg.Database.Host)
	}
	if cfg.OptionalDB == nil {
		t.Error("Expected OptionalDB to be set")
	} else if cfg.OptionalDB.Host != "optional.example.com" {
		t.Errorf("Expected OptionalDB.Host to be 'optional.example.com', got '%s'", cfg.OptionalDB.Host)
	}
}
