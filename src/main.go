package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
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
	if !utils.CheckFileExists(cfg.LoginLogPath) {
		_, err = utils.CreateLogFile(cfg.LoginLogPath)
		if err != nil {
			fmt.Printf("Failed to create login log file: %v\n", err)
			os.Exit(1)
		}
	}

	if !utils.CheckFileExists(cfg.ErrorLogPath) {
		_, err = utils.CreateLogFile(cfg.ErrorLogPath)
		if err != nil {
			fmt.Printf("Failed to create error log file: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Printf("Configuration loaded. Username: %s\n", cfg.Username)

	// Handle SIGINT and SIGTERM to perform cleanup before exiting.
	logoutFilePath := filepath.Join(cfg.TempPath, "logout")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		network.Logout(cfg)
		err := os.Remove(logoutFilePath)
		if err != nil {
			fmt.Printf("Failed to remove temp file: %v\n", err)
		}
		os.Exit(0)
	}()

	if utils.CheckFileExists(logoutFilePath) {
		err := network.Logout(cfg)
		if err != nil {
			os.Exit(1)
		}
	}

	ticker := time.NewTicker(cfg.Interval)
	for {
		err = network.PerformLogin(cfg)
		if err != nil {
			logger.LogError(err, cfg.ErrorLogPath)
		}
		<-ticker.C
	}
}
