/**
 * @description
 * This file handles configuration management for the Customer service.
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
	DatabaseURL       string `mapstructure:"DATABASE_URL"`
	RabbitMQURL       string `mapstructure:"RABBITMQ_URL"`
	AnchorAPIKey      string `mapstructure:"ANCHOR_API_KEY"`
	AnchorBaseURL     string `mapstructure:"ANCHOR_BASE_URL"`
	Port              string `mapstructure:"PORT"`
	UserCreatedQueue  string `mapstructure:"USER_CREATED_QUEUE"`
	UserCreatedEx     string `mapstructure:"USER_CREATED_EX"`
	UserCreatedRK     string `mapstructure:"USER_CREATED_RK"`
	ConsumerTag       string `mapstructure:"CONSUMER_TAG"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig() (config Config, err error) {
	viper.AddConfigPath("./")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	// Set default values for robust startup
	viper.SetDefault("PORT", "8081") // Use a different default port from auth
	viper.SetDefault("ANCHOR_BASE_URL", "https://api.sandbox.getanchor.co")
	viper.SetDefault("USER_CREATED_EX", "user_events")
	viper.SetDefault("USER_CREATED_RK", "user.created")
	viper.SetDefault("USER_CREATED_QUEUE", "customer_service_user_created")
	viper.SetDefault("CONSUMER_TAG", "customer_service_consumer")

	err = viper.ReadInConfig()
	// It's okay if the config file is not found, we can rely on env vars.
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return
		}
	}

	err = viper.Unmarshal(&config)
	return
}