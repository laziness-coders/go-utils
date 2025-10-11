package configs

import (
	"os"
	"testing"
)

// SimpleConfig represents a simple configuration for testing
type SimpleConfig struct {
	AppName string `mapstructure:"APP_NAME"`
	AppEnv  string `mapstructure:"APP_ENV"`
	Port    int    `mapstructure:"PORT"`
	Host    string `mapstructure:"HOST"`
}

func TestFullConfigLoadingFlow(t *testing.T) {
	// Use testdata directory
	testDataDir := "../testdata/configs"

	// Test loading configuration
	cfg := &SimpleConfig{}
	err := New(cfg).Load(AppEnvironmentDev, testDataDir)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Verify configuration was loaded from testdata (config.yaml overrides config.example.yaml)
	if cfg.AppName != "override-app" {
		t.Errorf("Expected AppName to be 'override-app', got '%s'", cfg.AppName)
	}
	if cfg.AppEnv != "dev" {
		t.Errorf("Expected AppEnv to be 'dev', got '%s'", cfg.AppEnv)
	}
	if cfg.Port != 8080 {
		t.Errorf("Expected Port to be 8080, got %d", cfg.Port)
	}
	if cfg.Host != "0.0.0.0" {
		t.Errorf("Expected Host to be '0.0.0.0', got '%s'", cfg.Host)
	}
}

func TestConfigWithExampleFile(t *testing.T) {
	// Use testdata directory
	testDataDir := "../testdata/configs"

	// Test loading configuration
	cfg := &SimpleConfig{}
	err := New(cfg).Load(AppEnvironmentDev, testDataDir)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Verify configuration was loaded from example file (config.yaml overrides config.example.yaml)
	if cfg.AppName != "override-app" {
		t.Errorf("Expected AppName to be 'override-app', got '%s'", cfg.AppName)
	}
	if cfg.AppEnv != "dev" {
		t.Errorf("Expected AppEnv to be 'dev', got '%s'", cfg.AppEnv)
	}
}

func TestConfigWithEnvironmentOverride(t *testing.T) {
	// Use testdata directory
	testDataDir := "../testdata/configs"

	// Set environment variables
	os.Setenv("APP_NAME", "env-override-app")
	os.Setenv("PORT", "9000")
	defer func() {
		os.Unsetenv("APP_NAME")
		os.Unsetenv("PORT")
	}()

	// Test loading configuration
	cfg := &SimpleConfig{}
	err := New(cfg).Load(AppEnvironmentDev, testDataDir)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Verify environment variables override config file values
	if cfg.AppName != "env-override-app" {
		t.Errorf("Expected AppName to be 'env-override-app', got '%s'", cfg.AppName)
	}
	if cfg.Port != 9000 {
		t.Errorf("Expected Port to be 9000, got %d", cfg.Port)
	}
	// These should still come from the config file
	if cfg.AppEnv != "dev" {
		t.Errorf("Expected AppEnv to be 'dev', got '%s'", cfg.AppEnv)
	}
}

func TestConfigWithEnvironmentSpecificFile(t *testing.T) {
	// Use testdata directory
	testDataDir := "../testdata/configs"

	// Test loading configuration for production environment
	cfg := &SimpleConfig{}
	err := New(cfg).Load(AppEnvironmentProd, testDataDir)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Verify production config overrides example config
	if cfg.AppName != "prod-app" {
		t.Errorf("Expected AppName to be 'prod-app', got '%s'", cfg.AppName)
	}
	if cfg.Port != 80 {
		t.Errorf("Expected Port to be 80, got %d", cfg.Port)
	}
	if cfg.Host != "prod.example.com" {
		t.Errorf("Expected Host to be 'prod.example.com', got '%s'", cfg.Host)
	}
	// This should come from the prod config
	if cfg.AppEnv != "prod" {
		t.Errorf("Expected AppEnv to be 'prod', got '%s'", cfg.AppEnv)
	}
}

func TestConfigMergingWithBothFiles(t *testing.T) {
	// Use testdata directory
	testDataDir := "../testdata/configs"

	// Test loading configuration
	cfg := &SimpleConfig{}
	err := New(cfg).Load(AppEnvironmentDev, testDataDir)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Verify that config.yaml overrides config.example.yaml
	if cfg.AppName != "override-app" {
		t.Errorf("Expected AppName to be 'override-app', got '%s'", cfg.AppName)
	}
	if cfg.Port != 8080 {
		t.Errorf("Expected Port to be 8080, got %d", cfg.Port)
	}
	if cfg.Host != "0.0.0.0" {
		t.Errorf("Expected Host to be '0.0.0.0', got '%s'", cfg.Host)
	}
	// This should still come from the example config (not overridden)
	if cfg.AppEnv != "dev" {
		t.Errorf("Expected AppEnv to be 'dev', got '%s'", cfg.AppEnv)
	}
}
