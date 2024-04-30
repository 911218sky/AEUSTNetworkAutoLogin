package config

import (
	"time"

	"github.com/go-ini/ini"

	"AEUSTNetworkAutoLogin/src/utils"
)

type Config struct {
	Ping         string
	Interval     time.Duration
	Username     string
	Password     string
	LoginLogPath string
	ErrorLogPath string
}

// LoadConfig loads the application configuration from the specified file.
func LoadConfig(configFilePath string) (*Config, error) {
	cfg, err := ini.Load(configFilePath)
	if err != nil {
		return nil, err
	}

	section := cfg.Section("Settings")
	ping := section.Key("PING").String()
	intervalStr := section.Key("INTERVAL").String()
	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		return nil, err
	}
	username := section.Key("USERNAME").String()
	password := section.Key("PASSWORD").String()
	loginLogPath := section.Key("LOGIN_LOGFILE_PATH").String()
	errorLogPath := section.Key("ERROR_LOGFILE_PATH").String()

	return &Config{
		Ping:         ping,
		Interval:     interval,
		Username:     username,
		Password:     password,
		LoginLogPath: loginLogPath,
		ErrorLogPath: errorLogPath,
	}, nil
}

// CreateDefaultConfig creates a default
func CreateDefaultConfig() (*Config, error) {
	username := utils.PromptUserInput("Please enter your student account number: ")
	password := utils.PromptUserInput("Please enter your student password: ")
	defaultConfig := &Config{
		Ping:         "8.8.8.8",
		Interval:     time.Second * 1,
		Username:     username,
		Password:     password,
		LoginLogPath: "login.log",
		ErrorLogPath: "error.log",
	}

	cfg := ini.Empty()
	LOGIN_LOGFILE_PATH, err := utils.CreateLogFile(defaultConfig.LoginLogPath)
	if err != nil {
		return nil, err
	}

	ERROR_LOGFILE_PATH, err := utils.CreateLogFile(defaultConfig.ErrorLogPath)
	if err != nil {
		return nil, err
	}

	section, _ := cfg.NewSection("Settings")
	section.NewKey("PING", defaultConfig.Ping)
	section.NewKey("INTERVAL", defaultConfig.Interval.String())
	section.NewKey("USERNAME", defaultConfig.Username)
	section.NewKey("PASSWORD", defaultConfig.Password)
	section.NewKey("LOGIN_LOGFILE_PATH", LOGIN_LOGFILE_PATH)
	section.NewKey("ERROR_LOGFILE_PATH", ERROR_LOGFILE_PATH)

	err = cfg.SaveTo("config.ini")
	if err != nil {
		return nil, err
	}

	return defaultConfig, nil
}
