/**
 * @description
 * This file handles configuration management for the Auth service.
 * It uses the Viper library to read configuration from environment variables
 * and a local .env file, making the service easily configurable across
 * different environments (development, staging, production).
 *
 * @dependencies
 * - "github.com/spf13/viper": A popular library for handling application configuration.
 */
package config

import "github.com/spf13/viper"

// Config stores all configuration for the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	DatabaseURL     string `mapstructure:"DATABASE_URL"`
	RabbitMQURL     string `mapstructure:"RABBITMQ_URL"`
	ClerkSecretKey  string `mapstructure:"CLERK_SECRET_KEY"`
	Port            string `mapstructure:"PORT"`
	UserCreatedEx   string `mapstructure:"USER_CREATED_EX"`
	UserCreatedRK   string `mapstructure:"USER_CREATED_RK"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig() (config Config, err error) {
	viper.AddConfigPath("./")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	// Set default values
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("USER_CREATED_EX", "user_events")
	viper.SetDefault("USER_CREATED_RK", "user.created")


	err = viper.ReadInConfig()
	if err != nil {
		// It's okay if the config file is not found, we can rely on env vars.
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return
		}
	}

	err = viper.Unmarshal(&config)
	return
}