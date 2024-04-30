package network

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
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

	regex := regexp.MustCompile(`window\.location\s*=\s*"([^"]+)";`)
	match := regex.FindStringSubmatch(body)
	if match == nil || len(match) < 2 {
		return errors.New("login URL not found")
	}
	loginUrl := match[1]

	resp, err = client.R().Get(loginUrl)

	if err != nil {
		logger.LogError(err, cfg.ErrorLogPath)
		return err
	}

	body = resp.String()

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		logger.LogError(err, cfg.ErrorLogPath)
		return err
	}

	magic := doc.Find("input[type=hidden]:nth-child(1)").AttrOr("value", "")

	if magic == "" {
		return errors.New("magic value not found")
	}

	resp, err = client.R().
		SetFormData(map[string]string{
			"magic":    magic,
			"4Tredir":  "http://edge-http.microsoft.com/captiveportal/generate_204",
			"username": cfg.Username,
			"password": cfg.Password,
		}).
		SetHeaders(map[string]string{
			"Host":       "fg.aeust.edu.tw:1442",
			"Origin":     "https://fg.aeust.edu.tw:1442",
			"Referer":    loginUrl,
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
		}).
		Post("https://fg.aeust.edu.tw:1442/") // 替换为实际的登录 Post URL

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
