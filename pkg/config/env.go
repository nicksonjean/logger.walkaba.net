package config

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

const (
	DefaultChannel    = "development"
	DefaultAppName    = "logger"
	DefaultTagName    = "latest"
	DefaultHost       = "127.0.0.1"
	DefaultPort       = 8080
	DefaultMiddleware = "net/http"
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

func GetEnvString(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	if value, exists := envVars[key]; exists && value != "" {
		return value
	}

	return defaultValue
}

func GetEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}

	if value, exists := envVars[key]; exists && value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}

	return defaultValue
}

func GetLoggerConfig() (string, string, string) {
	channel := GetEnvString("CHANNEL", DefaultChannel)
	if channel != "production" && channel != "development" {
		channel = DefaultChannel
	}

	appName := GetEnvString("APPNAME", DefaultAppName)

	tagName := GetEnvString("TAGNAME", DefaultTagName)

	return channel, appName, tagName
}

func GetStartServerConfig() (string, string, int) {
	middleware := GetEnvString("MIDDLEWARE", DefaultMiddleware)
	if middleware != "net/http" && middleware != "fiber" {
		middleware = DefaultMiddleware
	}

	host := GetEnvString("HOST", DefaultHost)

	port := GetEnvInt("PORT", DefaultPort)

	return middleware, host, port
}
