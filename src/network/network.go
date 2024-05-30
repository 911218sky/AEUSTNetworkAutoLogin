package network

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/go-ping/ping"
	"github.com/go-resty/resty/v2"

	"AEUSTNetworkAutoLogin/src/config"
	"AEUSTNetworkAutoLogin/src/logger"
)

// PingHost checks if the specified host is reachable via ping.
func PingHost(host string, cfg *config.Config) bool {
	pinger, err := ping.NewPinger(host)
	if err != nil {
		logger.LogError(err, cfg.ErrorLogPath)
		return false
	}
	pinger.Count = 1
	pinger.Timeout = time.Second * 2
	pinger.Run()
	stats := pinger.Statistics()
	return stats.PacketsRecv > 0
}

// PerformLogin attempts to log in using the provided credentials and configuration.
func PerformLogin(cfg *config.Config) error {
	client := resty.New()
	pingResult := PingHost(cfg.Ping, cfg)

	if pingResult {
		return nil
	}

	resp, err := client.R().
		Get("http://www.gstatic.com/generate_204")

	if err != nil {
		logger.LogError(err, cfg.ErrorLogPath)
		return err
	}

	body := resp.String()

	if body == "" {
		return nil
	}

	regex := regexp.MustCompile(`\?([^&]+)";`)
	match := regex.FindStringSubmatch(body)
	if match == nil || len(match) < 2 {
		return errors.New("login URL not found")
	}
	magic := match[1]

	if magic == "" {
		return errors.New("magic value not found")
	}

	fgtauthUrl := fmt.Sprintf("https://fg.aeust.edu.tw:1442/fgtauth?%s", magic)
	client.R().Get(fgtauthUrl)

	resp, err = client.SetTimeout(time.Second).
		R().
		SetFormData(map[string]string{
			"magic":    magic,
			"4Tredir":  "http://www.gstatic.com/generate_204",
			"username": cfg.Username,
			"password": cfg.Password,
			"submit":   "確認",
		}).
		SetHeaders(map[string]string{
			"Host":       "fg.aeust.edu.tw:1442",
			"Origin":     "https://fg.aeust.edu.tw:1442",
			"Referer":    fgtauthUrl,
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36",
		}).
		Post("https://fg.aeust.edu.tw:1442/")

	if err != nil {
		logger.LogError(err, cfg.ErrorLogPath)
		return err
	}

	if resp.StatusCode() == 200 {
		logger.LogSuccess(cfg.Username, "Login successful", cfg.LoginLogPath)
	} else {
		logger.LogError(fmt.Errorf("login failed with status code: %d", resp.StatusCode()), cfg.ErrorLogPath)
		return fmt.Errorf("login failed with status code: %d", resp.StatusCode())
	}

	return nil
}
