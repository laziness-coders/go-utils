package configs

import (
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

// GetEnv returns the value of an environment variable or the default value if not set.
func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
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
