package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"AEUSTNetworkAutoLogin/src/config"
	"AEUSTNetworkAutoLogin/src/logger"
	"AEUSTNetworkAutoLogin/src/network"
	"AEUSTNetworkAutoLogin/src/utils"

	"github.com/go-resty/resty/v2"
)

func initConfig() (*config.Config, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	configFilePath := filepath.Join(dir, "config.ini")
	if !utils.CheckFileExists(configFilePath) {
		if err := config.CreateDefaultConfig(); err != nil {
			return nil, err
		}
	}

	return config.LoadConfig(configFilePath)
}

func initLogFiles(cfg *config.Config) error {
	for _, path := range []string{cfg.LoginLogPath, cfg.ErrorLogPath} {
		if !utils.CheckFileExists(path) {
			if _, err := utils.CreateLogFile(path); err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	// Initialize logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Load configuration
	cfg, err := initConfig()
	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	// Initialize log files
	if err := initLogFiles(cfg); err != nil {
		log.Fatalf("Failed to initialize log files: %v", err)
	}

	log.Printf("Configuration loaded. Username: %s", cfg.Username)
	logoutFilePath := filepath.Join(cfg.TempPath, "logout")

	// Create context for program lifecycle control
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		if err := network.Logout(cfg); err != nil {
			log.Printf("Failed to logout: %v", err)
		}
		if err := os.Remove(logoutFilePath); err != nil && !os.IsNotExist(err) {
			log.Printf("Failed to remove temp file: %v", err)
		}
		cancel()
	}()

	// Check if logout is needed
	if utils.CheckFileExists(logoutFilePath) {
		if err := network.Logout(cfg); err != nil {
			log.Fatalf("Failed to logout: %v", err)
		}
	}

	client := resty.New()
	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := network.PerformLogin(cfg, client); err != nil {
				logger.LogError(err, cfg.ErrorLogPath)
			}
		}
	}
}
