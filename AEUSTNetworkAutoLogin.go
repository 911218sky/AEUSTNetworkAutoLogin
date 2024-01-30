package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"path/filepath"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-ini/ini"
	"github.com/go-ping/ping"
	"github.com/go-resty/resty/v2"
)

var configFilePath string

type Config struct {
	Ping         string
	Interval     time.Duration
	Username     string
	Password     string
	LoginLogPath string
	ErrorLogPath string
}

func readConfig() (*Config, error) {
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

func promptUserInput(question string) string {
	fmt.Print(question)
	var input string
	fmt.Scanln(&input)
	return input
}

func configExists() bool {
	_, err := os.Stat(configFilePath)
	return err == nil
}

func createLogFilePath(filename string) string {
	executablePath, _ := os.Executable()
	dir := filepath.Dir(executablePath)
	return filepath.Join(dir, filename)
}

func createDefaultConfig() (*Config, error) {
	username := promptUserInput("Please enter your student account number: ")
	password := promptUserInput("Please enter your student password: ")
	defaultConfig := &Config{
		Ping:         "8.8.8.8",
		Interval:     time.Second * 1,
		Username:     username,
		Password:     password,
		LoginLogPath: "login.log",
		ErrorLogPath: "error.log",
	}

	cfg := ini.Empty()
	section, _ := cfg.NewSection("Settings")
	section.NewKey("PING", defaultConfig.Ping)
	section.NewKey("INTERVAL", defaultConfig.Interval.String())
	section.NewKey("USERNAME", defaultConfig.Username)
	section.NewKey("PASSWORD", defaultConfig.Password)
	section.NewKey("LOGIN_LOGFILE_PATH", createLogFilePath(defaultConfig.LoginLogPath))
	section.NewKey("ERROR_LOGFILE_PATH", createLogFilePath(defaultConfig.ErrorLogPath))

	err := cfg.SaveTo("config.ini")
	if err != nil {
		return nil, err
	}

	return defaultConfig, nil
}

func getFormattedDateTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func appendToLogFile(logFilePath, logEntry string) error {
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(logEntry)
	return err
}

func logSuccess(config *Config, message string) {
	logEntry := fmt.Sprintf("[%s] : %s %s\n", getFormattedDateTime(), config.Username, message)
	fmt.Print(logEntry)
	appendToLogFile(config.LoginLogPath, logEntry)
}

func logError(config *Config, err error) {
	logEntry := fmt.Sprintf("[%s] ERROR: %v\n", getFormattedDateTime(), err)
	fmt.Print(logEntry)
	appendToLogFile(config.ErrorLogPath, logEntry)
}

func pingHost(host string) bool {
	pinger, err := ping.NewPinger(host)
	if err != nil {
		logError(nil, err)
		return false
	}
	pinger.Count = 1
	pinger.Timeout = time.Second * 2
	pinger.SetPrivileged(true)
	pinger.Run()
	stats := pinger.Statistics()
	return stats.PacketsRecv > 0
}

func mainProcess(config *Config) {
	client := resty.New()
	pingResult := pingHost(config.Ping)
	if !pingResult {
		resp, err := client.R().
			Get("http://www.gstatic.com/generate_204")

		if err != nil {
			logError(config, err)
			return
		}

		body := resp.String()

		if body == "" {
			return
		}

		regex := regexp.MustCompile(`window\.location\s*=\s*"([^"]+)";`)
		match := regex.FindStringSubmatch(body)
		if match == nil || len(match) < 2 {
			return
		}
		loginUrl := match[1]

		resp, err = client.R().
			Get(loginUrl)

		body = resp.String()

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
		if err != nil {
			logError(config, err)
			return
		}

		magic := doc.Find("input[type=hidden]:nth-child(1)").AttrOr("value", "")

		if magic == "" {
			return
		}

		resp, err = client.R().
			SetFormData(map[string]string{
				"magic":    magic,
				"4Tredir":  "http://edge-http.microsoft.com/captiveportal/generate_204",
				"username": config.Username,
				"password": config.Password,
			}).
			SetHeaders(map[string]string{
				"Host":       "fg.aeust.edu.tw:1442",
				"Origin":     "https://fg.aeust.edu.tw:1442",
				"Referer":    loginUrl,
				"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
			}).
			Post("https://fg.aeust.edu.tw:1442/")
		if err != nil {
			logError(config, err)
			return
		}
		if resp.StatusCode() == 200 {
			logSuccess(config, "Login successful "+config.Username)
		}
	}
}

func runMain(config *Config) {
	for {
		mainProcess(config)
		time.Sleep(config.Interval)
	}
}

func initialize() {
	executablePath, _ := os.Executable()
	dir := filepath.Dir(executablePath)
	configFilePath = filepath.Join(dir, "config.ini")

	var config *Config
	var err error
	if !configExists() {
		config, err = createDefaultConfig()
		if err != nil {
			logError(nil, err)
			return
		}
	} else {
		config, err = readConfig()
		if err != nil {
			logError(nil, err)
			return
		}
	}

	fmt.Printf("%s configuration loaded successfully\n", config.Username)
	runMain(config)
}

func main() {
	initialize()
}
