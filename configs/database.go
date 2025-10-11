package configs

import (
	"fmt"
	"time"
)

// PostgresConfig represents PostgreSQL database configuration.
type PostgresConfig struct {
	Host            string        `mapstructure:"HOST"`
	Port            int           `mapstructure:"PORT"`
	User            string        `mapstructure:"USER"`
	Password        string        `mapstructure:"PASSWORD"`
	Database        string        `mapstructure:"DATABASE"`
	SSLMode         string        `mapstructure:"SSL_MODE"`
	Timezone        string        `mapstructure:"TIMEZONE"`
	MaxOpenConns    int           `mapstructure:"MAX_OPEN_CONNS"`
	MaxIdleConns    int           `mapstructure:"MAX_IDLE_CONNS"`
	ConnMaxLifetime time.Duration `mapstructure:"CONN_MAX_LIFETIME"`
	ConnMaxIdleTime time.Duration `mapstructure:"CONN_MAX_IDLE_TIME"`
	Enabled         bool          `mapstructure:"ENABLED"`
}

// GetDSN returns the PostgreSQL DSN string.
func (c *PostgresConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode, c.Timezone)
}

// MySQLConfig represents MySQL database configuration.
type MySQLConfig struct {
	Host            string        `mapstructure:"HOST"`
	Port            int           `mapstructure:"PORT"`
	User            string        `mapstructure:"USER"`
	Password        string        `mapstructure:"PASSWORD"`
	Database        string        `mapstructure:"DATABASE"`
	Charset         string        `mapstructure:"CHARSET"`
	ParseTime       bool          `mapstructure:"PARSE_TIME"`
	Loc             string        `mapstructure:"LOC"`
	MaxOpenConns    int           `mapstructure:"MAX_OPEN_CONNS"`
	MaxIdleConns    int           `mapstructure:"MAX_IDLE_CONNS"`
	ConnMaxLifetime time.Duration `mapstructure:"CONN_MAX_LIFETIME"`
	ConnMaxIdleTime time.Duration `mapstructure:"CONN_MAX_IDLE_TIME"`
	Enabled         bool          `mapstructure:"ENABLED"`
}

// GetDSN returns the MySQL DSN string.
func (c *MySQLConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		c.User, c.Password, c.Host, c.Port, c.Database, c.Charset, c.ParseTime, c.Loc)
}

// RedisConfig represents Redis configuration.
type RedisConfig struct {
	Host         string        `mapstructure:"HOST"`
	Port         int           `mapstructure:"PORT"`
	DB           int           `mapstructure:"DB"`
	Password     string        `mapstructure:"PASSWORD"`
	PoolSize     int           `mapstructure:"POOL_SIZE"`
	MaxRetries   int           `mapstructure:"MAX_RETRIES"`
	DialTimeout  time.Duration `mapstructure:"DIAL_TIMEOUT"`
	ReadTimeout  time.Duration `mapstructure:"READ_TIMEOUT"`
	WriteTimeout time.Duration `mapstructure:"WRITE_TIMEOUT"`
	Enabled      bool          `mapstructure:"ENABLED"`
}

// GetAddr returns the Redis address string.
func (c *RedisConfig) GetAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// MongoDBConfig represents MongoDB configuration.
type MongoDBConfig struct {
	Host        string        `mapstructure:"HOST"`
	Port        int           `mapstructure:"PORT"`
	Database    string        `mapstructure:"DATABASE"`
	AuthSource  string        `mapstructure:"AUTH_SOURCE"`
	ReplicaSet  string        `mapstructure:"REPLICA_SET"`
	MaxPoolSize uint64        `mapstructure:"MAX_POOL_SIZE"`
	MinPoolSize uint64        `mapstructure:"MIN_POOL_SIZE"`
	Timeout     time.Duration `mapstructure:"TIMEOUT"`
	Enabled     bool          `mapstructure:"ENABLED"`
}

// GetURI returns the MongoDB connection URI.
func (c *MongoDBConfig) GetURI() string {
	uri := fmt.Sprintf("mongodb://%s:%d/%s", c.Host, c.Port, c.Database)

	if c.AuthSource != "" {
		uri += "?authSource=" + c.AuthSource
	}

	if c.ReplicaSet != "" {
		if c.AuthSource != "" {
			uri += "&replicaSet=" + c.ReplicaSet
		} else {
			uri += "?replicaSet=" + c.ReplicaSet
		}
	}

	return uri
}
