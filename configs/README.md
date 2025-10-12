# Configs Package

A generic, flexible configuration loader for Go applications using Viper with support for defaults and environment-specific overrides.

## Features

- **Generic Configuration Loading**: Use any struct with the fluent API
- **Layered Configuration**: Load from multiple sources with proper precedence
- **Environment-Specific Configs**: Support for different environments (dev, prod, test)
- **Environment Variable Override**: Automatic environment variable mapping
- **Fluent API**: Clean, chainable interface for configuration loading
- **Optional Default Configs**: Pre-built configs available in `configs/defaults` package

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/vnteam/go-utils/configs"
)

type MyConfig struct {
    Name string `mapstructure:"NAME"`
    Port int    `mapstructure:"PORT"`
    Host string `mapstructure:"HOST"`
}

func main() {
    cfg := &MyConfig{}
    
    err := configs.New(cfg).Load(configs.AppEnvironmentDev, "./configs")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("App: %s, Port: %d, Host: %s\n", cfg.Name, cfg.Port, cfg.Host)
}
```

## Configuration Loading Flow

The config loader follows a simple, predictable flow:

1. **Load config.example.yaml** (base defaults) - if it exists
2. **Merge config.yaml** (overrides) - if it exists  
3. **Merge config.{env}.yaml** (environment-specific) - if it exists
4. **Merge environment variables** (highest priority)

**Note**: If both `config.yaml` and `config.example.yaml` exist, they are merged with `config.yaml` taking precedence over `config.example.yaml`.

## Usage Examples

### Basic Configuration

```go
type AppConfig struct {
    Name string `mapstructure:"APP_NAME"`
    Port int    `mapstructure:"PORT"`
    Host string `mapstructure:"HOST"`
}

cfg := &AppConfig{}
err := configs.New(cfg).Load(configs.AppEnvironmentDev, "./configs")
```

### With Custom Viper Instance

```go
v := viper.New()
v.SetConfigName("my-config")
v.SetConfigType("yaml")

cfg := &MyConfig{}
err := configs.New(cfg).WithViper(v).Load(configs.AppEnvironmentDev, "./configs")
```

### Environment-Specific Configuration

```yaml
# config.example.yaml (base defaults)
app_name: "my-app"
port: 8080
host: "localhost"

# config.prod.yaml (production overrides)
app_name: "my-app-prod"
port: 80
host: "prod.example.com"
```

## Using Default Configurations

The configs package includes pre-built database and messaging configurations:

```go
import "github.com/vnteam/go-utils/configs"

type AppConfig struct {
    AppName string                    `mapstructure:"APP_NAME"`
    AppEnv  configs.AppEnvironment   `mapstructure:"APP_ENV"`
    Database configs.PostgresConfig  `mapstructure:"DATABASE"`
    Redis   configs.RedisConfig      `mapstructure:"REDIS"`
    Telegram configs.TelegramConfig  `mapstructure:"TELEGRAM"`
}

cfg := &AppConfig{}
err := configs.New(cfg).Load(configs.AppEnvironmentDev, "./configs")
```

## Configuration File Structure

```
configs/
├── config.example.yaml    # Base defaults (required)
├── config.yaml           # General overrides (optional)
├── config.dev.yaml       # Development overrides (optional)
├── config.prod.yaml      # Production overrides (optional)
└── config.test.yaml      # Test overrides (optional)
```

## Environment Variables

Environment variables automatically override config file values:

```bash
export APP_NAME="my-app"
export PORT=9000
export DATABASE_HOST="prod-db.example.com"
```

## API Reference

### ConfigLoader

```go
type ConfigLoader[T any] struct {
    config *T
    viper  *viper.Viper
}
```

### Methods

- `New[T any](cfg *T) *ConfigLoader[T]` - Create a new config loader
- `WithViper(v *viper.Viper) *ConfigLoader[T]` - Use a custom Viper instance
- `Load(appEnv AppEnvironment, configPath string) error` - Load configuration

### AppEnvironment

```go
type AppEnvironment string

const (
    AppEnvironmentDev  AppEnvironment = "dev"
    AppEnvironmentProd AppEnvironment = "prod"
    AppEnvironmentTest AppEnvironment = "test"
)
```

## Default Configurations

The configs package provides pre-built configuration structs:

### Database Configurations

- `PostgresConfig` - PostgreSQL database configuration
- `MySQLConfig` - MySQL database configuration  
- `RedisConfig` - Redis configuration
- `MongoDBConfig` - MongoDB configuration

### Messaging Configurations

- `TelegramConfig` - Telegram bot configuration
- `EmailConfig` - Email/SMTP configuration

### Example Usage

```go
import "github.com/vnteam/go-utils/configs"

// Use individual configs
postgres := &configs.PostgresConfig{}
redis := &configs.RedisConfig{}

// Or use the complete AppConfig
appConfig := &configs.AppConfig{}
```

## Best Practices

1. **Use config.example.yaml for defaults** - This provides a clear template for all required configuration
2. **Environment-specific overrides** - Use `config.{env}.yaml` for environment-specific settings
3. **Environment variables for secrets** - Use env vars for sensitive data like passwords and API keys
4. **Keep configs simple** - Avoid deeply nested structures when possible
5. **Use the defaults package sparingly** - Only import what you need

## Migration from Old Config System

If you're migrating from the old config system:

1. **Use generic loader** - `configs.New(cfg).Load(env, path)`
2. **Update struct tags** - Remove `validate` and `default` tags, keep only `mapstructure`
3. **Use config files for defaults** - Use `config.example.yaml` instead of struct tags

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

MIT License - see LICENSE file for details.