package configs

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// ConfigLoader is a generic configuration loader that provides a fluent API
// for loading configuration files with defaults, validation, and merging.
type ConfigLoader[T any] struct {
	config    *T
	viper     *viper.Viper
	validator func(*T) error
}

// New creates a new ConfigLoader for the given config struct.
func New[T any](cfg *T) *ConfigLoader[T] {
	return &ConfigLoader[T]{
		config: cfg,
		viper:  viper.New(),
	}
}

// WithViper sets a custom viper instance for the loader.
func (cl *ConfigLoader[T]) WithViper(v *viper.Viper) *ConfigLoader[T] {
	cl.viper = v
	return cl
}

// WithValidation sets a validation callback that will be called after the config is loaded.
// The callback should return an error if the configuration is invalid.
//
// Example:
//
//	cfg := &MyConfig{}
//	err := New(cfg).WithValidation(func(cfg *MyConfig) error {
//	    if cfg.Port <= 0 {
//	        return fmt.Errorf("port must be positive")
//	    }
//	    if cfg.Database.Host == "" {
//	        return fmt.Errorf("database host is required")
//	    }
//	    return nil
//	}).Load(AppEnvironmentDev, "./configs")
func (cl *ConfigLoader[T]) WithValidation(validator func(*T) error) *ConfigLoader[T] {
	cl.validator = validator
	return cl
}

// Load loads the configuration from files and environment variables.
// Simplified flow:
// 1. Load config.example.yaml (base)
// 2. Merge config.yaml (overrides)
// 3. Merge config.{env}.yaml (environment-specific)
// 4. Merge environment variables
// 5. Unmarshal into config struct
// 6. Run validation callback if provided
func (cl *ConfigLoader[T]) Load(appEnv AppEnvironment, configPath string) error {
	// Setup viper
	if err := cl.setupViper(configPath); err != nil {
		return fmt.Errorf("failed to setup viper: %w", err)
	}

	// Setup environment variable binding early
	cl.viper.AutomaticEnv()
	cl.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Load base configuration (config.example.yaml)
	if err := cl.loadBaseConfig(); err != nil {
		return fmt.Errorf("failed to load base config: %w", err)
	}

	// Merge environment-specific configuration
	if err := cl.mergeEnvConfig(appEnv); err != nil {
		return fmt.Errorf("failed to merge environment config: %w", err)
	}

	// Unmarshal into the config struct
	if err := cl.viper.Unmarshal(cl.config); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Run validation callback if provided
	if cl.validator != nil {
		if err := cl.validator(cl.config); err != nil {
			return fmt.Errorf("config validation failed: %w", err)
		}
	}

	return nil
}

// setupViper initializes viper with the configuration path.
func (cl *ConfigLoader[T]) setupViper(configPath string) error {
	cl.viper.SetConfigName("config")
	cl.viper.SetConfigType("yaml")
	cl.viper.AddConfigPath(configPath)
	return nil
}

// loadBaseConfig loads the base configuration files.
func (cl *ConfigLoader[T]) loadBaseConfig() error {
	// Try to load config.example.yaml as base defaults
	cl.viper.SetConfigName("config.example")
	if err := cl.viper.ReadInConfig(); err == nil {
		// If config.example.yaml exists, try to merge config.yaml on top
		cl.viper.SetConfigName("config")
		_ = cl.viper.MergeInConfig() // ignore error if config.yaml doesn't exist
		return nil
	}

	// If config.example.yaml doesn't exist, try config.yaml instead
	cl.viper.SetConfigName("config")
	if err := cl.viper.ReadInConfig(); err != nil {
		// If neither exists, that's okay - we'll use defaults and env vars
		return nil
	}

	return nil
}

// mergeEnvConfig merges environment-specific configuration.
func (cl *ConfigLoader[T]) mergeEnvConfig(appEnv AppEnvironment) error {
	// Set the environment-specific config name
	envConfigName := fmt.Sprintf("config.%s", string(appEnv))
	cl.viper.SetConfigName(envConfigName)

	// Try to merge environment-specific config
	if err := cl.viper.MergeInConfig(); err != nil {
		// It's okay if environment-specific config doesn't exist
		return nil
	}

	return nil
}
