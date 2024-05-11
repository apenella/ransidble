package configuration

import (
	"github.com/spf13/viper"
)

const (
	// DefaultHTTPListenAddress default HTTP listen address
	DefaultHTTPListenAddress = ":8080"
	// DefaultWorkerPoolSize default worker pool size
	DefaultWorkerPoolSize = 1

	// HTTPListenAddressKey key for HTTP listen address configuration
	HTTPListenAddressKey = "http_listen_address"
	// WorkerPoolSizeKey key for worker pool size configuration
	WorkerPoolSizeKey = "worker_pool_size"
)

// Configuration represents the configuration
type Configuration struct {
	// HTTPListenAddress represents the HTTP listen address
	HTTPListenAddress string `mapstructure:"http_listen_address"`
	// WorkerPoolSize represents the worker pool size
	WorkerPoolSize int `mapstructure:"worker_pool_size"`
}

func LoadConfig() (*Configuration, error) {
	var config Configuration
	var err error

	v := viper.New()

	v.SetDefault(HTTPListenAddressKey, DefaultHTTPListenAddress)
	v.SetDefault(WorkerPoolSizeKey, DefaultWorkerPoolSize)

	v.SetConfigName("ransidble") // Name of the config file (config.yaml, config.json, etc.)
	v.AddConfigPath(".")         // Search the current directory for the config file
	v.SetEnvPrefix("ransidble")  // Prefix for environment variables (e.g., RANSIDBLE_HTTP_LISTEN_ADDRESS)
	v.AutomaticEnv()             // Automatically read environment variables
	v.SetConfigType("yaml")      // Config file type (can be json, yaml, etc.)

	err = v.ReadInConfig()
	if err != nil {
		if _, isConfigFileNotFoundError := err.(viper.ConfigFileNotFoundError); !isConfigFileNotFoundError {
			return nil, err
		}
	}

	// Unmarshal configuration into struct
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	// Initialize config variable
	config.HTTPListenAddress = v.GetString("http_listen_address")
	config.WorkerPoolSize = v.GetInt("worker_pool_size")

	return &config, nil
}
