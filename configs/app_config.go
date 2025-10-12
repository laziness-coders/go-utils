package configs

import (
	"time"
)

// AppConfig represents the main application configuration.
type AppConfig struct {
	// App settings
	AppName     string         `mapstructure:"APP_NAME"`
	AppEnv      AppEnvironment `mapstructure:"APP_ENV"`
	AppLogPath  string         `mapstructure:"APP_LOG_PATH"`
	AppLogLevel string         `mapstructure:"APP_LOG_LEVEL"`

	// Databases
	Database PostgresConfig `mapstructure:"DATABASE"`
	MySQL    *MySQLConfig   `mapstructure:"MYSQL"`
	Redis    RedisConfig    `mapstructure:"REDIS"`
	MongoDB  *MongoDBConfig `mapstructure:"MONGODB"`

	// Messaging
	Telegram TelegramConfig `mapstructure:"TELEGRAM"`
	Email    EmailConfig    `mapstructure:"EMAIL"`

	// Server settings
	ServerPort    int    `mapstructure:"SERVER_PORT"`
	ServerHost    string `mapstructure:"SERVER_HOST"`
	ServerTimeout int    `mapstructure:"SERVER_TIMEOUT"`

	// JWT settings
	JWTSecretKey     string        `mapstructure:"JWT_SECRET_KEY"`
	JWTExpiration    time.Duration `mapstructure:"JWT_EXPIRATION"`
	JWTRefreshExpiry time.Duration `mapstructure:"JWT_REFRESH_EXPIRY"`
}
