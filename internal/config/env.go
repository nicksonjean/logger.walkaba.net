package config

import (
	"bufio"
	"os"
	"strings"
)

const (
	DefaultChannel = "development"
	DefaultAppName = "logger"
	DefaultTagName = "latest"
)

var envVars map[string]string

func LoadEnv() error {
	envVars = make(map[string]string)

	file, err := os.Open(".env")
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		value = strings.Trim(value, `"'`)

		envVars[key] = value
	}

	return scanner.Err()
}

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	if value, exists := envVars[key]; exists && value != "" {
		return value
	}

	return defaultValue
}

func GetLoggerConfig() (string, string, string) {
	channel := GetEnv("CHANNEL", DefaultChannel)
	if channel != "production" && channel != "development" {
		channel = DefaultChannel
	}

	appName := GetEnv("APPNAME", DefaultAppName)

	tagName := GetEnv("TAGNAME", DefaultTagName)

	return channel, appName, tagName
}
