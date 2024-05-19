package configuration

import (
	"strings"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/spf13/viper"
)

const (
	// DefaultHTTPListenAddress default HTTP listen address
	DefaultHTTPListenAddress = ":8080"
	// DefaultWorkerPoolSize default worker pool size
	DefaultWorkerPoolSize = 1
	// DefaultLogLevel default log level
	DefaultLogLevel = "info"
	// DefaultStorageType default storage type
	DefaultStorageType = entity.ProjectTypeLocal
	// DefaultLocalStoragePath default local storage path
	DefaultLocalStoragePath = "projects"

	// ServerKey key for server configuration
	ServerKey = "server"
	// HTTPListenAddressKey key for HTTP listen address configuration
	HTTPListenAddressKey = "http_listen_address"
	// WorkerPoolSizeKey key for worker pool size configuration
	WorkerPoolSizeKey = "worker_pool_size"
	// LogLevelKey key for log level configuration
	LogLevelKey = "log_level"

	// ProjectKey key for project configuration
	ProjectKey = "project"
	// StorageTypeKey key for storage type configuration
	StorageTypeKey = "storage_type"
	// LocalStoragePathKey key for local storage path configuration
	LocalStoragePathKey = "local_storage_path"
)

// Configuration represents the configuration
type Configuration struct {
	Server ServerConfiguration `mapstructure:"server"`
}

type ServerConfiguration struct {
	// HTTPListenAddress represents the HTTP listen address
	HTTPListenAddress string `mapstructure:"http_listen_address"`
	// WorkerPoolSize represents the worker pool size
	WorkerPoolSize int `mapstructure:"worker_pool_size"`
	// LogLevel represents the log level
	LogLevel string `mapstructure:"log_level"`
	// Project represents the project configuration
	Project ProjectConfiguration `mapstructure:"project"`
}

type ProjectConfiguration struct {
	// StorageType represents the storage type
	StorageType string `mapstructure:"storage_type"`
	// LocalStoragePath represents the local storage path
	LocalStoragePath string `mapstructure:"local_storage_path"`
}

func LoadConfig() (*Configuration, error) {
	var config Configuration
	var err error

	v := viper.New()

	v.SetDefault(strings.Join([]string{ServerKey, HTTPListenAddressKey}, "."), DefaultHTTPListenAddress)
	v.SetDefault(strings.Join([]string{ServerKey, WorkerPoolSizeKey}, "."), DefaultWorkerPoolSize)
	v.SetDefault(strings.Join([]string{ServerKey, LogLevelKey}, "."), DefaultLogLevel)

	v.SetDefault(strings.Join([]string{ServerKey, ProjectKey, StorageTypeKey}, "."), DefaultStorageType)
	v.SetDefault(strings.Join([]string{ServerKey, ProjectKey, LocalStoragePathKey}, "."), DefaultLocalStoragePath)

	replacer := strings.NewReplacer(".", "_")
	v.SetEnvKeyReplacer(replacer)
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

	// fmt.Printf(">>> Configuration: %+v\n", config)

	// // Initialize config variable
	// config.Server = ServerConfiguration{
	// 	HTTPListenAddress: v.GetString(strings.Join([]string{ServerKey, HTTPListenAddressKey}, ".")),
	// 	WorkerPoolSize:    v.GetInt(strings.Join([]string{ServerKey, WorkerPoolSizeKey}, ".")),
	// 	LogLevel:          v.GetString(strings.Join([]string{ServerKey, LogLevelKey}, ".")),
	// }

	return &config, nil
}
