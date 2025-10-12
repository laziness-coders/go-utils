package configs

import (
	"flag"
	"os"
)

// AppEnvironment represents the application environment.
type AppEnvironment string

const (
	AppEnvironmentProd        AppEnvironment = "prod"
	AppEnvironmentDev         AppEnvironment = "dev"
	AppEnvironmentTest        AppEnvironment = "test"
	AppEnvironmentIntegration AppEnvironment = "integration"
)

func (e AppEnvironment) IsProduction() bool {
	return e == AppEnvironmentProd
}

func (e AppEnvironment) IsDevelopment() bool {
	return e == AppEnvironmentDev
}

func (e AppEnvironment) IsTest() bool {
	return e == AppEnvironmentTest
}

func (e AppEnvironment) IsIntegration() bool {
	return e == AppEnvironmentIntegration
}

// GetEnv returns the value of an environment variable or the default value if not set.
// Typical the environment variable is "APP_ENV".
func GetEnv(key, defaultValue string) AppEnvironment {
	if value, exists := os.LookupEnv(key); exists {
		return AppEnvironment(value)
	}
	return AppEnvironment(defaultValue)
}

func ParseConfigDir(defaultDir string) string {
	var configDir string
	flag.StringVar(&configDir, "config-dir", defaultDir, "Configuration directory")
	return configDir
}

// GetServerPort returns the port for the main service.
// Priority is given to the port set from the environment variable.
// Then it checks the default value.
func GetServerPort(defaultServicePort string) string {
	// Check if the port is set in the environment variable
	if port := os.Getenv("PORT"); port != "" {
		return port
	}

	// If not set, return the default service port
	return defaultServicePort
}
