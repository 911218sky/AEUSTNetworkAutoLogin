package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"AEUSTNetworkAutoLogin/src/utils"

	"github.com/go-ini/ini"
)

// Config defines the application configuration structure
type Config struct {
	Ping         string        // Target address for network connectivity test
	Interval     time.Duration // Time interval for login checks
	Username     string        // Username for authentication
	Password     string        // Password for authentication
	LoginLogPath string        // Path to the login log file
	ErrorLogPath string        // Path to the error log file
	TempPath     string        // Path to temporary files directory
}

const (
	defaultPing     = "8.8.8.8"
	defaultInterval = time.Second
)

// LoadConfig loads the application configuration from the specified file
func LoadConfig(configFilePath string) (*Config, error) {
	cfg, err := ini.Load(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config file: %w", err)
	}

	section := cfg.Section("Settings")
	interval, err := time.ParseDuration(section.Key("INTERVAL").String())
	if err != nil {
		return nil, fmt.Errorf("failed to parse time interval: %w", err)
	}

	config := &Config{
		Ping:         section.Key("PING").String(),
		Interval:     interval,
		Username:     section.Key("USERNAME").String(),
		Password:     section.Key("PASSWORD").String(),
		LoginLogPath: section.Key("LOGIN_LOGFILE_PATH").String(),
		ErrorLogPath: section.Key("ERROR_LOGFILE_PATH").String(),
		TempPath:     section.Key("TEMP_PATH").String(),
	}

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

// validateConfig validates the configuration values
func validateConfig(cfg *Config) error {
	if cfg.Username == "" {
		return fmt.Errorf("username cannot be empty")
	}
	if cfg.Password == "" {
		return fmt.Errorf("password cannot be empty")
	}
	if cfg.Interval < time.Second {
		return fmt.Errorf("interval cannot be less than 1 second")
	}
	return nil
}

// CreateDefaultConfig creates a default configuration file
func CreateDefaultConfig() error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	config := &Config{
		Ping:         defaultPing,
		Interval:     defaultInterval,
		Username:     utils.PromptUserInput("Please enter your student account number: "),
		Password:     utils.PromptUserInput("Please enter your password: "),
		LoginLogPath: filepath.Join(dir, "logs", "login.log"),
		ErrorLogPath: filepath.Join(dir, "logs", "error.log"),
		TempPath:     filepath.Join(dir, "temp"),
	}

	// Create required directories
	if err := createRequiredDirectories(config); err != nil {
		return err
	}

	// Create log files
	if err := createLogFiles(config); err != nil {
		return err
	}

	// Save configuration to file
	return saveConfigToFile(config)
}

// createRequiredDirectories creates necessary directories
func createRequiredDirectories(cfg *Config) error {
	dirs := []string{
		filepath.Dir(cfg.LoginLogPath),
		filepath.Dir(cfg.ErrorLogPath),
		cfg.TempPath,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	return nil
}

// createLogFiles creates log files
func createLogFiles(cfg *Config) error {
	loginLogPath, err := utils.CreateLogFile(cfg.LoginLogPath)
	if err != nil {
		return fmt.Errorf("failed to create login log file: %w", err)
	}
	cfg.LoginLogPath = loginLogPath

	errorLogPath, err := utils.CreateLogFile(cfg.ErrorLogPath)
	if err != nil {
		return fmt.Errorf("failed to create error log file: %w", err)
	}
	cfg.ErrorLogPath = errorLogPath

	return nil
}

// saveConfigToFile saves the configuration to a file
func saveConfigToFile(cfg *Config) error {
	iniCfg := ini.Empty()
	section, err := iniCfg.NewSection("Settings")
	if err != nil {
		return fmt.Errorf("failed to create config section: %w", err)
	}

	// Set configuration items
	section.NewKey("PING", cfg.Ping)
	section.NewKey("INTERVAL", cfg.Interval.String())
	section.NewKey("USERNAME", cfg.Username)
	section.NewKey("PASSWORD", cfg.Password)
	section.NewKey("LOGIN_LOGFILE_PATH", cfg.LoginLogPath)
	section.NewKey("ERROR_LOGFILE_PATH", cfg.ErrorLogPath)
	section.NewKey("TEMP_PATH", cfg.TempPath)

	if err := iniCfg.SaveTo("config.ini"); err != nil {
		return fmt.Errorf("failed to save config file: %w", err)
	}

	return nil
}
