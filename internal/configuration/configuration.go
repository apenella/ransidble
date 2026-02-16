package configuration

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

const (
	// DefaultHTTPListenAddress default HTTP listen address
	DefaultHTTPListenAddress = ":8080"
	// DefaultWorkerPoolSize default worker pool size
	DefaultWorkerPoolSize = 1
	// DefaultLogLevel default log level
	DefaultLogLevel = "info"
	// DefaultProjectStorageLocalPath default local storage path
	DefaultProjectStorageLocalPath = "storage/projects"
	// DefaultProjectRepositoryLocalPath default local repository path
	DefaultProjectRepositoryLocalPath = "repository/projects"

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

	// ProjectStorageKey key for project storage configuration
	ProjectStorageKey = "storage"
	// ProjectStorageTypeKey key for project storage type configuration
	ProjectStorageTypeKey = "type"
	// ProjectStorageLocalPathKey key for project storage local path configuration
	ProjectStorageLocalPathKey = "local_path"

	// ProjectRepositoryKey key for project repository configuration
	ProjectRepositoryKey = "repository"
	// ProjectRepositoryTypeKey key for project repository type configuration
	ProjectRepositoryTypeKey = "type"
	// ProjectRepositoryLocalPathKey key for project repository local path configuration
	ProjectRepositoryLocalPathKey = "local_path"
)

// Configuration represents the configuration
type Configuration struct {
	Server ServerConfiguration `mapstructure:"server"`
}

// ServerConfiguration represents the server configuration
type ServerConfiguration struct {
	// HTTPListenAddress represents the HTTP listen address
	HTTPListenAddress string `mapstructure:"http_listen_address" validate:"required,listen_addr"`
	// WorkerPoolSize represents the worker pool size
	WorkerPoolSize int `mapstructure:"worker_pool_size" validate:"required,gt=0"`
	// LogLevel represents the log level
	LogLevel string `mapstructure:"log_level" validate:"required,oneof=debug info warn error"`
	// Project represents the project configuration
	Project ProjectConfiguration `mapstructure:"project"`
}

// ProjectConfiguration represents the project configuration
type ProjectConfiguration struct {
	ProjectStorageConfiguration    ProjectStorageConfiguration    `mapstructure:"storage"`
	ProjectRepositoryConfiguration ProjectRepositoryConfiguration `mapstructure:"repository"`
}

// ProjectStorageConfiguration represents the project storage configuration
type ProjectStorageConfiguration struct {
	// LocalStoragePath represents the local storage path
	LocalStoragePath string `mapstructure:"local_path" validate:"required_if=Type local"`
	// Type represents the type of storage (e.g., memory, local, http, registry, etc.)
	Type string `mapstructure:"type" validate:"required,oneof=local memory"`
}

// ProjectRepositoryConfiguration represents the project repository configuration
type ProjectRepositoryConfiguration struct {
	// LocalRepositoryPath represents the local repository path
	LocalRepositoryPath string `mapstructure:"local_path" validate:"required_if=Type local"`
	// Type represents the type of repository (e.g., memory, local, database, etc.)
	Type string `mapstructure:"type" validate:"required,oneof=local memory"`
}

// LoadConfig loads the configuration
func LoadConfig() (*Configuration, error) {
	var config Configuration
	var err error

	v := viper.New()

	v.BindEnv(strings.Join([]string{ServerKey, HTTPListenAddressKey}, "."))
	v.BindEnv(strings.Join([]string{ServerKey, LogLevelKey}, "."))
	v.BindEnv(strings.Join([]string{ServerKey, ProjectKey, ProjectRepositoryKey, ProjectRepositoryLocalPathKey}, "."))
	v.BindEnv(strings.Join([]string{ServerKey, ProjectKey, ProjectRepositoryKey, ProjectRepositoryTypeKey}, "."))
	v.BindEnv(strings.Join([]string{ServerKey, ProjectKey, ProjectStorageKey, ProjectStorageLocalPathKey}, "."))
	v.BindEnv(strings.Join([]string{ServerKey, ProjectKey, ProjectStorageKey, ProjectStorageTypeKey}, "."))
	v.BindEnv(strings.Join([]string{ServerKey, WorkerPoolSizeKey}, "."))

	v.SetDefault(strings.Join([]string{ServerKey, HTTPListenAddressKey}, "."), DefaultHTTPListenAddress)
	v.SetDefault(strings.Join([]string{ServerKey, LogLevelKey}, "."), DefaultLogLevel)
	v.SetDefault(strings.Join([]string{ServerKey, ProjectKey, ProjectRepositoryKey, ProjectRepositoryLocalPathKey}, "."), DefaultProjectRepositoryLocalPath)
	v.SetDefault(strings.Join([]string{ServerKey, ProjectKey, ProjectRepositoryKey, ProjectRepositoryTypeKey}, "."), "local")
	v.SetDefault(strings.Join([]string{ServerKey, ProjectKey, ProjectStorageKey, ProjectStorageLocalPathKey}, "."), DefaultProjectStorageLocalPath)
	v.SetDefault(strings.Join([]string{ServerKey, ProjectKey, ProjectStorageKey, ProjectStorageTypeKey}, "."), "local")
	v.SetDefault(strings.Join([]string{ServerKey, WorkerPoolSizeKey}, "."), DefaultWorkerPoolSize)

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

	err = config.Validate()
	if err != nil {
		return nil, fmt.Errorf("configuration validation error: %w", err)
	}

	return &config, nil
}

func (c *Configuration) Validate() error {

	validate := validator.New()
	err := validate.RegisterValidation("listen_addr", listenAddrValidation)
	if err != nil {
		return err
	}

	return validate.Struct(c)
}

func listenAddrValidation(fl validator.FieldLevel) bool {
	listenAddrRegex := regexp.MustCompile(`^((\d{1,3}\.){3}\d{1,3})?:\d{1,5}$`)

	addr := fl.Field().String()
	return listenAddrRegex.MatchString(addr)
}
