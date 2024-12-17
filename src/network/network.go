package network

import (
	"context"
	"errors"
	"fmt"
	"net"
	"path/filepath"
	"regexp"
	"runtime"
	"time"

	"github.com/go-ping/ping"
	"github.com/go-resty/resty/v2"

	"AEUSTNetworkAutoLogin/src/config"
	"AEUSTNetworkAutoLogin/src/logger"
	"AEUSTNetworkAutoLogin/src/utils"
)

const (
	generateURL    = "http://www.gstatic.com/generate_204"
	defaultTimeout = 5 * time.Second
	pingTimeout    = 2 * time.Second
	pingCount      = 1
	userAgent      = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36"
	logoutFileName = "logout"
)

var (
	// ErrInvalidCredentials indicates invalid username or password
	ErrInvalidCredentials = errors.New("invalid username or password")
	// ErrNoMatchingURL indicates no matching URL found
	ErrNoMatchingURL = errors.New("no matching URL found")
	// ErrLogoutURLNotFound indicates logout URL not found
	ErrLogoutURLNotFound = errors.New("logout URL not found")

	// Regular expression for parsing redirect URL
	redirectURLRegex = regexp.MustCompile(`window\.location="?(https://fg\.aeust\.edu\.tw:(\d+)/fgtauth\?([^"&]+))"?;`)
	// Regular expression for parsing logout URL
	logoutURLRegex = regexp.MustCompile(`/keepalive\?([^"]+)`)
)

// PingHost checks if the specified host is reachable via ping
func PingHost(ctx context.Context, cfg *config.Config) bool {
	pinger, err := ping.NewPinger(cfg.Ping)
	if err != nil {
		logger.LogError(err, cfg.ErrorLogPath)
		return false
	}

	if runtime.GOOS == "windows" {
		pinger.SetPrivileged(true)
	}

	pinger.Count = pingCount
	pinger.Timeout = pingTimeout

	// Use context for timeout control
	done := make(chan bool)
	go func() {
		pinger.Run()
		done <- true
	}()

	select {
	case <-ctx.Done():
		pinger.Stop()
		return false
	case <-done:
		stats := pinger.Statistics()
		return stats.PacketsRecv > 0
	}
}

// PerformLogin attempts to log in using the provided credentials and configuration
func PerformLogin(cfg *config.Config, client *resty.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	if PingHost(ctx, cfg) {
		return nil
	}

	// Get redirect URL
	resp, err := client.R().Get(generateURL)
	if err != nil {
		return fmt.Errorf("failed to access generate page: %w", err)
	}

	body := resp.String()
	if body == "" {
		return nil
	}

	// Parse redirect URL
	match := redirectURLRegex.FindStringSubmatch(body)
	if len(match) < 4 {
		return ErrNoMatchingURL
	}

	fgtauthURL := match[1]
	port := match[2]
	magicValue := match[3]

	// Access authentication page
	if _, err := client.R().Get(fgtauthURL); err != nil {
		return fmt.Errorf("failed to access authentication page: %w", err)
	}

	// Submit login form
	resp, err = client.SetTimeout(defaultTimeout).
		R().
		SetFormData(map[string]string{
			"magic":    magicValue,
			"4Tredir":  generateURL,
			"username": cfg.Username,
			"password": cfg.Password,
			"submit":   "Submit",
		}).
		SetHeaders(map[string]string{
			"Host":       fmt.Sprintf("fg.aeust.edu.tw:%s", port),
			"Origin":     fmt.Sprintf("https://fg.aeust.edu.tw:%s", port),
			"Referer":    fgtauthURL,
			"User-Agent": userAgent,
		}).
		Post(fmt.Sprintf("https://fg.aeust.edu.tw:%s/", port))

	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return ErrInvalidCredentials
		}
		return fmt.Errorf("login request failed: %w", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("login failed with status code: %d", resp.StatusCode())
	}

	// Parse and save logout URL
	match = logoutURLRegex.FindStringSubmatch(resp.String())
	if len(match) < 2 {
		return ErrLogoutURLNotFound
	}

	logoutURL := fmt.Sprintf("https://fg.aeust.edu.tw:%s/logout?%s", port, match[1])
	if err := utils.WriteBinaryFile(filepath.Join(cfg.TempPath, logoutFileName), logoutURL); err != nil {
		return fmt.Errorf("failed to save logout URL: %w", err)
	}

	if err := logger.LogSuccess(cfg.Username, "Login successful", cfg.LoginLogPath); err != nil {
		return fmt.Errorf("failed to log login success: %w", err)
	}

	return nil
}

// Logout performs network logout operation
func Logout(cfg *config.Config) error {
	// Read logout URL
	logoutURL, err := utils.ReadBinaryFile(filepath.Join(cfg.TempPath, logoutFileName))
	if err != nil {
		return fmt.Errorf("failed to read logout URL: %w", err)
	}

	// Execute logout request
	client := resty.New().SetTimeout(defaultTimeout)
	resp, err := client.R().Get(logoutURL)
	if err != nil {
		return fmt.Errorf("logout request failed: %w", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("logout failed with status code: %d", resp.StatusCode())
	}

	if err := logger.LogSuccess(cfg.Username, "Logout successful", cfg.LoginLogPath); err != nil {
		return fmt.Errorf("failed to log logout success: %w", err)
	}

	return nil
}
