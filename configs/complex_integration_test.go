package configs

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestComplexConfig_MultipleTelegramConfigs tests loading multiple Telegram configurations
func TestComplexConfig_MultipleTelegramConfigs(t *testing.T) {
	// Create a temporary directory for config files
	tempDir := t.TempDir()

	// Create config.example.yaml with multiple Telegram configs
	configContent := `
telegrams:
  - bot_token: "bot_token_1"
    channel_id: 123456789
    message_thread_id: 1
    enabled: true
  - bot_token: "bot_token_2"
    channel_id: 987654321
    message_thread_id: 2
    enabled: false
  - bot_token: "bot_token_3"
    channel_id: 555555555
    message_thread_id: 3
    enabled: true
`

	configFile := filepath.Join(tempDir, "config.example.yaml")
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Use a wrapper struct to match the YAML structure
	type TestConfig struct {
		Telegrams TelegramConfigs `mapstructure:"telegrams"`
	}

	testCfg := &TestConfig{}
	err := New(testCfg).Load(AppEnvironmentDev, tempDir)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Check that we have 3 Telegram configs
	if len(testCfg.Telegrams) != 3 {
		t.Fatalf("Expected 3 Telegram configs, got %d", len(testCfg.Telegrams))
	}

	// Verify first Telegram config
	if testCfg.Telegrams[0].BotToken != "bot_token_1" {
		t.Errorf("Expected first BotToken to be 'bot_token_1', got '%s'", testCfg.Telegrams[0].BotToken)
	}
	if testCfg.Telegrams[0].ChannelID != 123456789 {
		t.Errorf("Expected first ChannelID to be 123456789, got %d", testCfg.Telegrams[0].ChannelID)
	}
	if testCfg.Telegrams[0].MessageThreadID != 1 {
		t.Errorf("Expected first MessageThreadID to be 1, got %d", testCfg.Telegrams[0].MessageThreadID)
	}
	if !testCfg.Telegrams[0].Enabled {
		t.Error("Expected first Telegram config to be enabled")
	}

	// Verify second Telegram config
	if testCfg.Telegrams[1].BotToken != "bot_token_2" {
		t.Errorf("Expected second BotToken to be 'bot_token_2', got '%s'", testCfg.Telegrams[1].BotToken)
	}
	if testCfg.Telegrams[1].ChannelID != 987654321 {
		t.Errorf("Expected second ChannelID to be 987654321, got %d", testCfg.Telegrams[1].ChannelID)
	}
	if testCfg.Telegrams[1].MessageThreadID != 2 {
		t.Errorf("Expected second MessageThreadID to be 2, got %d", testCfg.Telegrams[1].MessageThreadID)
	}
	if testCfg.Telegrams[1].Enabled {
		t.Error("Expected second Telegram config to be disabled")
	}

	// Verify third Telegram config
	if testCfg.Telegrams[2].BotToken != "bot_token_3" {
		t.Errorf("Expected third BotToken to be 'bot_token_3', got '%s'", testCfg.Telegrams[2].BotToken)
	}
	if testCfg.Telegrams[2].ChannelID != 555555555 {
		t.Errorf("Expected third ChannelID to be 555555555, got %d", testCfg.Telegrams[2].ChannelID)
	}
	if testCfg.Telegrams[2].MessageThreadID != 3 {
		t.Errorf("Expected third MessageThreadID to be 3, got %d", testCfg.Telegrams[2].MessageThreadID)
	}
	if !testCfg.Telegrams[2].Enabled {
		t.Error("Expected third Telegram config to be enabled")
	}
}

// TestComplexConfig_MultipleMySQL tests loading multiple MySQL database configurations
func TestComplexConfig_MultipleMySQL(t *testing.T) {
	// Create a temporary directory for config files
	tempDir := t.TempDir()

	// Create config.example.yaml with multiple MySQL configs
	configContent := `
mysql_a:
  host: "mysql-a.example.com"
  port: 3306
  user: "user_a"
  password: "password_a"
  database: "database_a"
  charset: "utf8mb4"
  parse_time: true
  loc: "Local"
  max_open_conns: 100
  max_idle_conns: 10
  conn_max_lifetime: "1h"
  conn_max_idle_time: "30m"
  enabled: true

mysql_b:
  host: "mysql-b.example.com"
  port: 3307
  user: "user_b"
  password: "password_b"
  database: "database_b"
  charset: "utf8mb4"
  parse_time: false
  loc: "UTC"
  max_open_conns: 50
  max_idle_conns: 5
  conn_max_lifetime: "30m"
  conn_max_idle_time: "15m"
  enabled: false
`

	configFile := filepath.Join(tempDir, "config.example.yaml")
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Use a wrapper struct to match the YAML structure
	type TestConfig struct {
		MySQLA MySQLConfig `mapstructure:"mysql_a"`
		MySQLB MySQLConfig `mapstructure:"mysql_b"`
	}

	testCfg := &TestConfig{}
	err := New(testCfg).Load(AppEnvironmentDev, tempDir)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Verify MySQL A config
	if testCfg.MySQLA.Host != "mysql-a.example.com" {
		t.Errorf("Expected MySQLA Host to be 'mysql-a.example.com', got '%s'", testCfg.MySQLA.Host)
	}
	if testCfg.MySQLA.Port != 3306 {
		t.Errorf("Expected MySQLA Port to be 3306, got %d", testCfg.MySQLA.Port)
	}
	if testCfg.MySQLA.User != "user_a" {
		t.Errorf("Expected MySQLA User to be 'user_a', got '%s'", testCfg.MySQLA.User)
	}
	if testCfg.MySQLA.Password != "password_a" {
		t.Errorf("Expected MySQLA Password to be 'password_a', got '%s'", testCfg.MySQLA.Password)
	}
	if testCfg.MySQLA.Database != "database_a" {
		t.Errorf("Expected MySQLA Database to be 'database_a', got '%s'", testCfg.MySQLA.Database)
	}
	if testCfg.MySQLA.Charset != "utf8mb4" {
		t.Errorf("Expected MySQLA Charset to be 'utf8mb4', got '%s'", testCfg.MySQLA.Charset)
	}
	if !testCfg.MySQLA.ParseTime {
		t.Error("Expected MySQLA ParseTime to be true")
	}
	if testCfg.MySQLA.Loc != "Local" {
		t.Errorf("Expected MySQLA Loc to be 'Local', got '%s'", testCfg.MySQLA.Loc)
	}
	if testCfg.MySQLA.MaxOpenConns != 100 {
		t.Errorf("Expected MySQLA MaxOpenConns to be 100, got %d", testCfg.MySQLA.MaxOpenConns)
	}
	if testCfg.MySQLA.MaxIdleConns != 10 {
		t.Errorf("Expected MySQLA MaxIdleConns to be 10, got %d", testCfg.MySQLA.MaxIdleConns)
	}
	if testCfg.MySQLA.ConnMaxLifetime != time.Hour {
		t.Errorf("Expected MySQLA ConnMaxLifetime to be 1h, got %v", testCfg.MySQLA.ConnMaxLifetime)
	}
	if testCfg.MySQLA.ConnMaxIdleTime != 30*time.Minute {
		t.Errorf("Expected MySQLA ConnMaxIdleTime to be 30m, got %v", testCfg.MySQLA.ConnMaxIdleTime)
	}
	if !testCfg.MySQLA.Enabled {
		t.Error("Expected MySQLA to be enabled")
	}

	// Verify MySQL A DSN
	expectedDSNA := "user_a:password_a@tcp(mysql-a.example.com:3306)/database_a?charset=utf8mb4&parseTime=true&loc=Local"
	if testCfg.MySQLA.GetDSN() != expectedDSNA {
		t.Errorf("Expected MySQLA DSN to be '%s', got '%s'", expectedDSNA, testCfg.MySQLA.GetDSN())
	}

	// Verify MySQL B config
	if testCfg.MySQLB.Host != "mysql-b.example.com" {
		t.Errorf("Expected MySQLB Host to be 'mysql-b.example.com', got '%s'", testCfg.MySQLB.Host)
	}
	if testCfg.MySQLB.Port != 3307 {
		t.Errorf("Expected MySQLB Port to be 3307, got %d", testCfg.MySQLB.Port)
	}
	if testCfg.MySQLB.User != "user_b" {
		t.Errorf("Expected MySQLB User to be 'user_b', got '%s'", testCfg.MySQLB.User)
	}
	if testCfg.MySQLB.Password != "password_b" {
		t.Errorf("Expected MySQLB Password to be 'password_b', got '%s'", testCfg.MySQLB.Password)
	}
	if testCfg.MySQLB.Database != "database_b" {
		t.Errorf("Expected MySQLB Database to be 'database_b', got '%s'", testCfg.MySQLB.Database)
	}
	if testCfg.MySQLB.Charset != "utf8mb4" {
		t.Errorf("Expected MySQLB Charset to be 'utf8mb4', got '%s'", testCfg.MySQLB.Charset)
	}
	if testCfg.MySQLB.ParseTime {
		t.Error("Expected MySQLB ParseTime to be false")
	}
	if testCfg.MySQLB.Loc != "UTC" {
		t.Errorf("Expected MySQLB Loc to be 'UTC', got '%s'", testCfg.MySQLB.Loc)
	}
	if testCfg.MySQLB.MaxOpenConns != 50 {
		t.Errorf("Expected MySQLB MaxOpenConns to be 50, got %d", testCfg.MySQLB.MaxOpenConns)
	}
	if testCfg.MySQLB.MaxIdleConns != 5 {
		t.Errorf("Expected MySQLB MaxIdleConns to be 5, got %d", testCfg.MySQLB.MaxIdleConns)
	}
	if testCfg.MySQLB.ConnMaxLifetime != 30*time.Minute {
		t.Errorf("Expected MySQLB ConnMaxLifetime to be 30m, got %v", testCfg.MySQLB.ConnMaxLifetime)
	}
	if testCfg.MySQLB.ConnMaxIdleTime != 15*time.Minute {
		t.Errorf("Expected MySQLB ConnMaxIdleTime to be 15m, got %v", testCfg.MySQLB.ConnMaxIdleTime)
	}
	if testCfg.MySQLB.Enabled {
		t.Error("Expected MySQLB to be disabled")
	}

	// Verify MySQL B DSN
	expectedDSNB := "user_b:password_b@tcp(mysql-b.example.com:3307)/database_b?charset=utf8mb4&parseTime=false&loc=UTC"
	if testCfg.MySQLB.GetDSN() != expectedDSNB {
		t.Errorf("Expected MySQLB DSN to be '%s', got '%s'", expectedDSNB, testCfg.MySQLB.GetDSN())
	}
}

// TestComplexConfig_CombinedMultipleConfigs tests loading multiple complex configurations together
func TestComplexConfig_CombinedMultipleConfigs(t *testing.T) {
	// Create a temporary directory for config files
	tempDir := t.TempDir()

	// Create config.example.yaml with multiple Telegram configs and multiple MySQL databases
	configContent := `
# Multiple Telegram configurations
telegrams:
  - bot_token: "production_bot"
    channel_id: 111111111
    message_thread_id: 10
    enabled: true
  - bot_token: "staging_bot"
    channel_id: 222222222
    message_thread_id: 20
    enabled: true
  - bot_token: "development_bot"
    channel_id: 333333333
    message_thread_id: 30
    enabled: false

# Primary MySQL database
mysql_primary:
  host: "primary-db.example.com"
  port: 3306
  user: "primary_user"
  password: "primary_pass"
  database: "primary_db"
  charset: "utf8mb4"
  parse_time: true
  loc: "Local"
  max_open_conns: 200
  max_idle_conns: 20
  conn_max_lifetime: "2h"
  conn_max_idle_time: "1h"
  enabled: true

# Secondary MySQL database (read replica)
mysql_secondary:
  host: "secondary-db.example.com"
  port: 3306
  user: "secondary_user"
  password: "secondary_pass"
  database: "secondary_db"
  charset: "utf8mb4"
  parse_time: true
  loc: "Local"
  max_open_conns: 150
  max_idle_conns: 15
  conn_max_lifetime: "1h30m"
  conn_max_idle_time: "45m"
  enabled: true

# Analytics MySQL database
mysql_analytics:
  host: "analytics-db.example.com"
  port: 3306
  user: "analytics_user"
  password: "analytics_pass"
  database: "analytics_db"
  charset: "utf8mb4"
  parse_time: true
  loc: "UTC"
  max_open_conns: 75
  max_idle_conns: 10
  conn_max_lifetime: "1h"
  conn_max_idle_time: "30m"
  enabled: false
`

	configFile := filepath.Join(tempDir, "config.example.yaml")
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Use a wrapper struct to match the YAML structure
	type TestConfig struct {
		Telegrams        TelegramConfigs `mapstructure:"telegrams"`
		MySQLPrimary     MySQLConfig     `mapstructure:"mysql_primary"`
		MySQLSecondary   MySQLConfig     `mapstructure:"mysql_secondary"`
		MySQLAnalytics   MySQLConfig     `mapstructure:"mysql_analytics"`
	}

	testCfg := &TestConfig{}
	err := New(testCfg).Load(AppEnvironmentDev, tempDir)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Verify Telegram configs
	if len(testCfg.Telegrams) != 3 {
		t.Fatalf("Expected 3 Telegram configs, got %d", len(testCfg.Telegrams))
	}

	// Spot check Telegram configs
	if testCfg.Telegrams[0].BotToken != "production_bot" {
		t.Errorf("Expected first BotToken to be 'production_bot', got '%s'", testCfg.Telegrams[0].BotToken)
	}
	if testCfg.Telegrams[1].ChannelID != 222222222 {
		t.Errorf("Expected second ChannelID to be 222222222, got %d", testCfg.Telegrams[1].ChannelID)
	}
	if testCfg.Telegrams[2].Enabled {
		t.Error("Expected third Telegram config to be disabled")
	}

	// Verify Primary MySQL
	if testCfg.MySQLPrimary.Host != "primary-db.example.com" {
		t.Errorf("Expected Primary Host to be 'primary-db.example.com', got '%s'", testCfg.MySQLPrimary.Host)
	}
	if testCfg.MySQLPrimary.Database != "primary_db" {
		t.Errorf("Expected Primary Database to be 'primary_db', got '%s'", testCfg.MySQLPrimary.Database)
	}
	if testCfg.MySQLPrimary.MaxOpenConns != 200 {
		t.Errorf("Expected Primary MaxOpenConns to be 200, got %d", testCfg.MySQLPrimary.MaxOpenConns)
	}
	if !testCfg.MySQLPrimary.Enabled {
		t.Error("Expected Primary MySQL to be enabled")
	}

	// Verify Secondary MySQL
	if testCfg.MySQLSecondary.Host != "secondary-db.example.com" {
		t.Errorf("Expected Secondary Host to be 'secondary-db.example.com', got '%s'", testCfg.MySQLSecondary.Host)
	}
	if testCfg.MySQLSecondary.Database != "secondary_db" {
		t.Errorf("Expected Secondary Database to be 'secondary_db', got '%s'", testCfg.MySQLSecondary.Database)
	}
	if testCfg.MySQLSecondary.MaxOpenConns != 150 {
		t.Errorf("Expected Secondary MaxOpenConns to be 150, got %d", testCfg.MySQLSecondary.MaxOpenConns)
	}
	if !testCfg.MySQLSecondary.Enabled {
		t.Error("Expected Secondary MySQL to be enabled")
	}

	// Verify Analytics MySQL
	if testCfg.MySQLAnalytics.Host != "analytics-db.example.com" {
		t.Errorf("Expected Analytics Host to be 'analytics-db.example.com', got '%s'", testCfg.MySQLAnalytics.Host)
	}
	if testCfg.MySQLAnalytics.Database != "analytics_db" {
		t.Errorf("Expected Analytics Database to be 'analytics_db', got '%s'", testCfg.MySQLAnalytics.Database)
	}
	if testCfg.MySQLAnalytics.MaxOpenConns != 75 {
		t.Errorf("Expected Analytics MaxOpenConns to be 75, got %d", testCfg.MySQLAnalytics.MaxOpenConns)
	}
	if testCfg.MySQLAnalytics.Loc != "UTC" {
		t.Errorf("Expected Analytics Loc to be 'UTC', got '%s'", testCfg.MySQLAnalytics.Loc)
	}
	if testCfg.MySQLAnalytics.Enabled {
		t.Error("Expected Analytics MySQL to be disabled")
	}

	// Verify DSN generation for each database
	primaryDSN := testCfg.MySQLPrimary.GetDSN()
	if primaryDSN == "" {
		t.Error("Expected Primary DSN to be non-empty")
	}
	secondaryDSN := testCfg.MySQLSecondary.GetDSN()
	if secondaryDSN == "" {
		t.Error("Expected Secondary DSN to be non-empty")
	}
	analyticsDSN := testCfg.MySQLAnalytics.GetDSN()
	if analyticsDSN == "" {
		t.Error("Expected Analytics DSN to be non-empty")
	}

	// Verify all DSNs are different
	if primaryDSN == secondaryDSN || primaryDSN == analyticsDSN || secondaryDSN == analyticsDSN {
		t.Error("Expected all DSNs to be different")
	}
}

// TestComplexConfig_WithEnvironmentOverride tests environment variable overrides for complex configs
func TestComplexConfig_WithEnvironmentOverride(t *testing.T) {
	// Create a temporary directory for config files
	tempDir := t.TempDir()

	// Create config.example.yaml with multiple MySQL configs
	configContent := `
mysql_a:
  host: "mysql-a.example.com"
  port: 3306
  user: "user_a"
  password: "password_a"
  database: "database_a"
  charset: "utf8mb4"
  parse_time: true
  loc: "Local"
  max_open_conns: 100
  max_idle_conns: 10
  conn_max_lifetime: "1h"
  conn_max_idle_time: "30m"
  enabled: true

mysql_b:
  host: "mysql-b.example.com"
  port: 3307
  user: "user_b"
  password: "password_b"
  database: "database_b"
  charset: "utf8mb4"
  parse_time: false
  loc: "UTC"
  max_open_conns: 50
  max_idle_conns: 5
  conn_max_lifetime: "30m"
  conn_max_idle_time: "15m"
  enabled: false
`

	configFile := filepath.Join(tempDir, "config.example.yaml")
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Set environment variables to override config values
	err := os.Setenv("MYSQL_A_HOST", "env-override-host-a")
	if err != nil {
		t.Fatalf("Failed to set MYSQL_A_HOST env: %v", err)
	}
	err = os.Setenv("MYSQL_B_PORT", "9999")
	if err != nil {
		t.Fatalf("Failed to set MYSQL_B_PORT env: %v", err)
	}

	defer func() {
		_ = os.Unsetenv("MYSQL_A_HOST")
		_ = os.Unsetenv("MYSQL_B_PORT")
	}()

	// Use a wrapper struct to match the YAML structure
	type TestConfig struct {
		MySQLA MySQLConfig `mapstructure:"mysql_a"`
		MySQLB MySQLConfig `mapstructure:"mysql_b"`
	}

	testCfg := &TestConfig{}
	err = New(testCfg).Load(AppEnvironmentDev, tempDir)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Verify MySQL A host was overridden by environment variable
	if testCfg.MySQLA.Host != "env-override-host-a" {
		t.Errorf("Expected MySQLA Host to be 'env-override-host-a', got '%s'", testCfg.MySQLA.Host)
	}

	// Verify MySQL A other fields remain from config file
	if testCfg.MySQLA.Port != 3306 {
		t.Errorf("Expected MySQLA Port to be 3306, got %d", testCfg.MySQLA.Port)
	}

	// Verify MySQL B port was overridden by environment variable
	if testCfg.MySQLB.Port != 9999 {
		t.Errorf("Expected MySQLB Port to be 9999, got %d", testCfg.MySQLB.Port)
	}

	// Verify MySQL B other fields remain from config file
	if testCfg.MySQLB.Host != "mysql-b.example.com" {
		t.Errorf("Expected MySQLB Host to be 'mysql-b.example.com', got '%s'", testCfg.MySQLB.Host)
	}
}
