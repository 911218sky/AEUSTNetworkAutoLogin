package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"AEUSTNetworkAutoLogin/src/config"
	"AEUSTNetworkAutoLogin/src/logger"
	"AEUSTNetworkAutoLogin/src/network"
	"AEUSTNetworkAutoLogin/src/utils"
)

func main() {
	dir, _ := os.Getwd()
	configFilePath := filepath.Join(dir, "config.ini")

	isExists := utils.CheckFileExists(configFilePath)
	if !isExists {
		config.CreateDefaultConfig()
	}

	// Load the configuration.
	cfg, err := config.LoadConfig(configFilePath)
	if err != nil {
		fmt.Printf("Failed to read config: %v\n", err)
		os.Exit(1) // Exit if configuration cannot be loaded
	}

	// Check if the log paths are valid and create log files if they don't exist.
	_, err = utils.CreateLogFile(cfg.LoginLogPath)
	if err != nil {
		fmt.Printf("Failed to create login log file: %v\n", err)
		os.Exit(1)
	}
	_, err = utils.CreateLogFile(cfg.ErrorLogPath)
	if err != nil {
		fmt.Printf("Failed to create error log file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Configuration loaded. Username: %s\n", cfg.Username)

	for {
		err = network.PerformLogin(cfg)
		if err != nil {
			logger.LogError(err, cfg.ErrorLogPath)
		}
		time.Sleep(cfg.Interval)
	}
}
