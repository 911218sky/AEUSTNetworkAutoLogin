package network

import (
	"errors"
	"fmt"
	"net"
	"path"
	"regexp"
	"runtime"
	"time"

	"github.com/go-ping/ping"
	"github.com/go-resty/resty/v2"

	"AEUSTNetworkAutoLogin/src/config"
	"AEUSTNetworkAutoLogin/src/logger"
	"AEUSTNetworkAutoLogin/src/utils"
)

// PingHost checks if the specified host is reachable via ping.
func PingHost(cfg *config.Config) bool {
	pinger, err := ping.NewPinger(cfg.Ping)
	if err != nil {
		logger.LogError(err, cfg.ErrorLogPath)
		return false
	}
	// If the OS is Windows, we need to use privileged mode
	if runtime.GOOS == "windows" {
		pinger.SetPrivileged(true)
	}
	pinger.Count = 1
	pinger.Timeout = time.Second * 2
	pinger.Run()
	stats := pinger.Statistics()
	return stats.PacketsRecv > 0
}

// PerformLogin attempts to log in using the provided credentials and configuration.
func PerformLogin(cfg *config.Config, client *resty.Client) error {
	pingResult := PingHost(cfg)
	if pingResult {
		return nil
	}

	const generateUrl = "http://www.gstatic.com/generate_204"

	resp, err := client.R().
		Get(generateUrl)

	if err != nil {
		logger.LogError(err, cfg.ErrorLogPath)
		return err
	}

	body := resp.String()
	if body == "" {
		return nil
	}

	regex := regexp.MustCompile(`window\.location="?(https://fg\.aeust\.edu\.tw:(\d+)/fgtauth\?([^"&]+))"?;`)
	match := regex.FindStringSubmatch(body)
	var fgtauthUrl, port, magicValue string
	if len(match) >= 4 {
		fgtauthUrl = match[1]
		port = match[2]
		magicValue = match[3]
	} else {
		return errors.New("no matching URL, port, and magic value found")
	}

	client.R().Get(fgtauthUrl)

	resp, err = client.SetTimeout(time.Second).
		R().
		SetFormData(map[string]string{
			"magic":    magicValue,
			"4Tredir":  generateUrl,
			"username": cfg.Username,
			"password": cfg.Password,
			"submit":   "確認",
		}).
		SetHeaders(map[string]string{
			"Host":       "fg.aeust.edu.tw:" + port,
			"Origin":     "https://fg.aeust.edu.tw:" + port,
			"Referer":    fgtauthUrl,
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36",
		}).
		Post("https://fg.aeust.edu.tw:" + port + "/")

	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			logger.LogError(fmt.Errorf("invalid username or password, please check your credentials"), cfg.ErrorLogPath)
			return nil
		}
		logger.LogError(err, cfg.ErrorLogPath)
		return err
	}

	if resp.StatusCode() == 200 {
		regex := regexp.MustCompile(`/keepalive\?([^"]+)`)
		match := regex.FindStringSubmatch(resp.String())
		if match == nil || len(match) < 2 {
			return fmt.Errorf("logout URL not found")
		}
		logoutUrl := fmt.Sprintf("https://fg.aeust.edu.tw:%s/logout?%s", port, match[1])
		utils.WriteCustomBinaryFile(path.Join(cfg.TempPath, "logout"), logoutUrl)
		logger.LogSuccess(cfg.Username, "Login successful", cfg.LoginLogPath)
	} else {
		logger.LogError(fmt.Errorf("login failed with status code: %d", resp.StatusCode()), cfg.ErrorLogPath)
		return fmt.Errorf("login failed with status code: %d", resp.StatusCode())
	}

	return nil
}

// Logout logs out of the network.
func Logout(cfg *config.Config) error {
	logoutUrl, err := utils.ReadCustomBinaryFile(path.Join(cfg.TempPath, "logout"))
	if err != nil {
		logger.LogError(fmt.Errorf("logout URL not found in temp file"), cfg.ErrorLogPath)
		return err
	}

	client := resty.New()
	resp, err := client.R().Get(logoutUrl)
	if err != nil {
		logger.LogError(err, cfg.ErrorLogPath)
		return err
	}
	if resp.StatusCode() == 200 {
		logger.LogSuccess(cfg.Username, "Logout successful", cfg.LoginLogPath)
	} else {
		logger.LogError(fmt.Errorf("logout failed with status code: %d", resp.StatusCode()), cfg.ErrorLogPath)
		return fmt.Errorf("logout failed with status code: %d", resp.StatusCode())
	}
	return nil
}
