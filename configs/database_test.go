package configs

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestPostgresConfig_GetDSN(t *testing.T) {
	cfg := &PostgresConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "password",
		Database: "testdb",
		SSLMode:  "disable",
		Timezone: "UTC",
	}

	expected := "host=localhost port=5432 user=postgres password=password dbname=testdb sslmode=disable TimeZone=UTC"
	actual := cfg.GetDSN()

	if actual != expected {
		t.Errorf("Expected DSN: %s, got: %s", expected, actual)
	}
}

func TestMySQLConfig_GetDSN(t *testing.T) {
	cfg := &MySQLConfig{
		Host:      "localhost",
		Port:      3306,
		User:      "root",
		Password:  "password",
		Database:  "testdb",
		Charset:   "utf8mb4",
		ParseTime: true,
		Loc:       "Local",
	}

	expected := "root:password@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=true&loc=Local"
	actual := cfg.GetDSN()

	if actual != expected {
		t.Errorf("Expected DSN: %s, got: %s", expected, actual)
	}
}

func TestRedisConfig_GetAddr(t *testing.T) {
	cfg := &RedisConfig{
		Host: "localhost",
		Port: 6379,
	}

	expected := "localhost:6379"
	actual := cfg.GetAddr()

	if actual != expected {
		t.Errorf("Expected address: %s, got: %s", expected, actual)
	}
}

func TestMongoDBConfig_GetURI(t *testing.T) {
	cfg := &MongoDBConfig{
		Host:       "localhost",
		Port:       27017,
		Database:   "testdb",
		AuthSource: "admin",
		ReplicaSet: "rs0",
	}

	expected := "mongodb://localhost:27017/testdb?authSource=admin&replicaSet=rs0"
	actual := cfg.GetURI()

	if actual != expected {
		t.Errorf("Expected URI: %s, got: %s", expected, actual)
	}
}

func TestPostgresConfig_WithConfigFile(t *testing.T) {
	// cfg := &PostgresConfig{}

	// Create a temporary directory for config files
	tempDir := t.TempDir()

	// Create config.example.yaml with PostgreSQL config
	configContent := `
database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "password"
  database: "testdb"
  ssl_mode: "disable"
  timezone: "UTC"
  max_open_conns: 100
  max_idle_conns: 10
  conn_max_lifetime: "1h"
  conn_max_idle_time: "30m"
  enabled: true
`

	configFile := filepath.Join(tempDir, "config.example.yaml")
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Use a wrapper struct to match the YAML structure
	type TestConfig struct {
		Database PostgresConfig `mapstructure:"database"`
	}

	testCfg := &TestConfig{}
	err := New(testCfg).Load(AppEnvironmentDev, tempDir)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Check that the config was loaded correctly
	if testCfg.Database.Host != "localhost" {
		t.Errorf("Expected Host to be 'localhost', got '%s'", testCfg.Database.Host)
	}
	if testCfg.Database.Port != 5432 {
		t.Errorf("Expected Port to be 5432, got %d", testCfg.Database.Port)
	}
	if testCfg.Database.User != "postgres" {
		t.Errorf("Expected User to be 'postgres', got '%s'", testCfg.Database.User)
	}
	if testCfg.Database.Password != "password" {
		t.Errorf("Expected Password to be 'password', got '%s'", testCfg.Database.Password)
	}
	if testCfg.Database.Database != "testdb" {
		t.Errorf("Expected Database to be 'testdb', got '%s'", testCfg.Database.Database)
	}
	if testCfg.Database.SSLMode != "disable" {
		t.Errorf("Expected SSLMode to be 'disable', got '%s'", testCfg.Database.SSLMode)
	}
	if testCfg.Database.Timezone != "UTC" {
		t.Errorf("Expected Timezone to be 'UTC', got '%s'", testCfg.Database.Timezone)
	}
	if testCfg.Database.MaxOpenConns != 100 {
		t.Errorf("Expected MaxOpenConns to be 100, got %d", testCfg.Database.MaxOpenConns)
	}
	if testCfg.Database.MaxIdleConns != 10 {
		t.Errorf("Expected MaxIdleConns to be 10, got %d", testCfg.Database.MaxIdleConns)
	}
	if testCfg.Database.ConnMaxLifetime != time.Hour {
		t.Errorf("Expected ConnMaxLifetime to be 1h, got %v", testCfg.Database.ConnMaxLifetime)
	}
	if testCfg.Database.ConnMaxIdleTime != 30*time.Minute {
		t.Errorf("Expected ConnMaxIdleTime to be 30m, got %v", testCfg.Database.ConnMaxIdleTime)
	}
	if !testCfg.Database.Enabled {
		t.Error("Expected Enabled to be true")
	}
}

func TestMySQLConfig_WithConfigFile(t *testing.T) {
	// cfg := &MySQLConfig{}

	// Create a temporary directory for config files
	tempDir := t.TempDir()

	// Create config.example.yaml with MySQL config
	configContent := `
mysql:
  host: "localhost"
  port: 3306
  user: "root"
  password: "password"
  database: "testdb"
  charset: "utf8mb4"
  parse_time: true
  loc: "Local"
  max_open_conns: 100
  max_idle_conns: 10
  conn_max_lifetime: "1h"
  conn_max_idle_time: "30m"
  enabled: true
`

	configFile := filepath.Join(tempDir, "config.example.yaml")
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Use a wrapper struct to match the YAML structure
	type TestConfig struct {
		MySQL MySQLConfig `mapstructure:"mysql"`
	}

	testCfg := &TestConfig{}
	err := New(testCfg).Load(AppEnvironmentDev, tempDir)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Check that the config was loaded correctly
	if testCfg.MySQL.Host != "localhost" {
		t.Errorf("Expected Host to be 'localhost', got '%s'", testCfg.MySQL.Host)
	}
	if testCfg.MySQL.Port != 3306 {
		t.Errorf("Expected Port to be 3306, got %d", testCfg.MySQL.Port)
	}
	if testCfg.MySQL.User != "root" {
		t.Errorf("Expected User to be 'root', got '%s'", testCfg.MySQL.User)
	}
	if testCfg.MySQL.Password != "password" {
		t.Errorf("Expected Password to be 'password', got '%s'", testCfg.MySQL.Password)
	}
	if testCfg.MySQL.Database != "testdb" {
		t.Errorf("Expected Database to be 'testdb', got '%s'", testCfg.MySQL.Database)
	}
	if testCfg.MySQL.Charset != "utf8mb4" {
		t.Errorf("Expected Charset to be 'utf8mb4', got '%s'", testCfg.MySQL.Charset)
	}
	if !testCfg.MySQL.ParseTime {
		t.Error("Expected ParseTime to be true")
	}
	if testCfg.MySQL.Loc != "Local" {
		t.Errorf("Expected Loc to be 'Local', got '%s'", testCfg.MySQL.Loc)
	}
	if testCfg.MySQL.MaxOpenConns != 100 {
		t.Errorf("Expected MaxOpenConns to be 100, got %d", testCfg.MySQL.MaxOpenConns)
	}
	if testCfg.MySQL.MaxIdleConns != 10 {
		t.Errorf("Expected MaxIdleConns to be 10, got %d", testCfg.MySQL.MaxIdleConns)
	}
	if testCfg.MySQL.ConnMaxLifetime != time.Hour {
		t.Errorf("Expected ConnMaxLifetime to be 1h, got %v", testCfg.MySQL.ConnMaxLifetime)
	}
	if testCfg.MySQL.ConnMaxIdleTime != 30*time.Minute {
		t.Errorf("Expected ConnMaxIdleTime to be 30m, got %v", testCfg.MySQL.ConnMaxIdleTime)
	}
	if !testCfg.MySQL.Enabled {
		t.Error("Expected Enabled to be true")
	}
}

func TestRedisConfig_WithConfigFile(t *testing.T) {
	// cfg := &RedisConfig{}

	// Create a temporary directory for config files
	tempDir := t.TempDir()

	// Create config.example.yaml with Redis config
	configContent := `
redis:
  host: "localhost"
  port: 6379
  db: 0
  password: "redis_password"
  pool_size: 10
  max_retries: 3
  dial_timeout: "5s"
  read_timeout: "3s"
  write_timeout: "3s"
  enabled: true
`

	configFile := filepath.Join(tempDir, "config.example.yaml")
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Use a wrapper struct to match the YAML structure
	type TestConfig struct {
		Redis RedisConfig `mapstructure:"redis"`
	}

	testCfg := &TestConfig{}
	err := New(testCfg).Load(AppEnvironmentDev, tempDir)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Check that the config was loaded correctly
	if testCfg.Redis.Host != "localhost" {
		t.Errorf("Expected Host to be 'localhost', got '%s'", testCfg.Redis.Host)
	}
	if testCfg.Redis.Port != 6379 {
		t.Errorf("Expected Port to be 6379, got %d", testCfg.Redis.Port)
	}
	if testCfg.Redis.DB != 0 {
		t.Errorf("Expected DB to be 0, got %d", testCfg.Redis.DB)
	}
	if testCfg.Redis.Password != "redis_password" {
		t.Errorf("Expected Password to be 'redis_password', got '%s'", testCfg.Redis.Password)
	}
	if testCfg.Redis.PoolSize != 10 {
		t.Errorf("Expected PoolSize to be 10, got %d", testCfg.Redis.PoolSize)
	}
	if testCfg.Redis.MaxRetries != 3 {
		t.Errorf("Expected MaxRetries to be 3, got %d", testCfg.Redis.MaxRetries)
	}
	if testCfg.Redis.DialTimeout != 5*time.Second {
		t.Errorf("Expected DialTimeout to be 5s, got %v", testCfg.Redis.DialTimeout)
	}
	if testCfg.Redis.ReadTimeout != 3*time.Second {
		t.Errorf("Expected ReadTimeout to be 3s, got %v", testCfg.Redis.ReadTimeout)
	}
	if testCfg.Redis.WriteTimeout != 3*time.Second {
		t.Errorf("Expected WriteTimeout to be 3s, got %v", testCfg.Redis.WriteTimeout)
	}
	if !testCfg.Redis.Enabled {
		t.Error("Expected Enabled to be true")
	}
}

func TestMongoDBConfig_WithConfigFile(t *testing.T) {
	// cfg := &MongoDBConfig{}

	// Create a temporary directory for config files
	tempDir := t.TempDir()

	// Create config.example.yaml with MongoDB config
	configContent := `
mongodb:
  host: "localhost"
  port: 27017
  database: "testdb"
  auth_source: "admin"
  replica_set: "rs0"
  max_pool_size: 100
  min_pool_size: 5
  timeout: "30s"
  enabled: true
`

	configFile := filepath.Join(tempDir, "config.example.yaml")
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Use a wrapper struct to match the YAML structure
	type TestConfig struct {
		MongoDB MongoDBConfig `mapstructure:"mongodb"`
	}

	testCfg := &TestConfig{}
	err := New(testCfg).Load(AppEnvironmentDev, tempDir)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Check that the config was loaded correctly
	if testCfg.MongoDB.Host != "localhost" {
		t.Errorf("Expected Host to be 'localhost', got '%s'", testCfg.MongoDB.Host)
	}
	if testCfg.MongoDB.Port != 27017 {
		t.Errorf("Expected Port to be 27017, got %d", testCfg.MongoDB.Port)
	}
	if testCfg.MongoDB.Database != "testdb" {
		t.Errorf("Expected Database to be 'testdb', got '%s'", testCfg.MongoDB.Database)
	}
	if testCfg.MongoDB.AuthSource != "admin" {
		t.Errorf("Expected AuthSource to be 'admin', got '%s'", testCfg.MongoDB.AuthSource)
	}
	if testCfg.MongoDB.ReplicaSet != "rs0" {
		t.Errorf("Expected ReplicaSet to be 'rs0', got '%s'", testCfg.MongoDB.ReplicaSet)
	}
	if testCfg.MongoDB.MaxPoolSize != 100 {
		t.Errorf("Expected MaxPoolSize to be 100, got %d", testCfg.MongoDB.MaxPoolSize)
	}
	if testCfg.MongoDB.MinPoolSize != 5 {
		t.Errorf("Expected MinPoolSize to be 5, got %d", testCfg.MongoDB.MinPoolSize)
	}
	if testCfg.MongoDB.Timeout != 30*time.Second {
		t.Errorf("Expected Timeout to be 30s, got %v", testCfg.MongoDB.Timeout)
	}
	if !testCfg.MongoDB.Enabled {
		t.Error("Expected Enabled to be true")
	}
}
