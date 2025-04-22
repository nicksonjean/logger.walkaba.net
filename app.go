package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DebugLevel     = 100
	InfoLevel      = 200
	NoticeLevel    = 250
	WarningLevel   = 300
	ErrorLevel     = 400
	CriticalLevel  = 500
	AlertLevel     = 550
	EmergencyLevel = 600
)

const (
	DefaultChannel = "development"
	DefaultAppName = "logger"
	DefaultTagName = "latest"
)

var levelNames = map[int]string{
	DebugLevel:     "DEBUG",
	InfoLevel:      "INFO",
	NoticeLevel:    "NOTICE",
	WarningLevel:   "WARNING",
	ErrorLevel:     "ERROR",
	CriticalLevel:  "CRITICAL",
	AlertLevel:     "ALERT",
	EmergencyLevel: "EMERGENCY",
}

type Context struct {
	CorrelationID string     `json:"correlation_id"`
	RequestID     string     `json:"request_id"`
	AppName       string     `json:"app_name"`
	TagName       string     `json:"tag_name,omitempty"`
	Exception     *Exception `json:"exception,omitempty"`
}

type Exception struct {
	Message string `json:"message"`
	File    string `json:"file"`
}

type CustomLog struct {
	Message   string      `json:"message"`
	Context   Context     `json:"context"`
	Level     int         `json:"level"`
	LevelName string      `json:"level_name"`
	Channel   string      `json:"channel"`
	Datetime  string      `json:"datetime"`
	Extra     interface{} `json:"extra"`
}

type CustomWriter struct {
	channel string
	appName string
	tagName string
}

func (w *CustomWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func generateUUID() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	uuid := fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		r.Uint32(), r.Uint32()&0xffff, (r.Uint32()&0x0fff)|0x4000,
		(r.Uint32()&0x3fff)|0x8000, r.Int63n(0xffffffffffff))
	return uuid
}

type contextKey string

const (
	correlationIDKey contextKey = "correlation_id"
	loggerKey        contextKey = "logger"
)

type CustomLogger struct {
	zapLogger     *zap.Logger
	channel       string
	appName       string
	tagName       string
	correlationID string
}

func LoadEnvFile(filename string) (map[string]string, error) {
	envMap := make(map[string]string)

	file, err := os.Open(filename)
	if err != nil {
		return envMap, err
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

		envMap[key] = value
	}

	if err := scanner.Err(); err != nil {
		return envMap, err
	}

	return envMap, nil
}

func GetLoggerConfig() (string, string, string) {
	envVars, err := LoadEnvFile(".env")
	if err != nil {
		fmt.Printf("Warning: Não foi possível carregar o arquivo .env: %v\n", err)
	}

	channel := envVars["CHANNEL"]
	if channel == "" {
		channel = os.Getenv("CHANNEL")
	}
	if channel != "production" && channel != "development" {
		channel = DefaultChannel
	}

	appName := envVars["APPNAME"]
	if appName == "" {
		appName = os.Getenv("APPNAME")
	}
	if appName == "" {
		appName = DefaultAppName
	}

	tagName := envVars["TAGNAME"]
	if tagName == "" {
		tagName = os.Getenv("TAGNAME")
	}
	if tagName == "" {
		tagName = DefaultTagName
	}

	return channel, appName, tagName
}

func NewCustomLogger(channel, appName, tagName string) (*CustomLogger, error) {
	if channel == "" || appName == "" {
		envChannel, envAppName, envTagName := GetLoggerConfig()

		if channel == "" {
			channel = envChannel
		}

		if appName == "" {
			appName = envAppName
		}

		if tagName == "" {
			tagName = envTagName
		}
	}

	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "datetime",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	encoder := zapcore.NewJSONEncoder(encoderConfig)

	customWriter := &CustomWriter{
		channel: channel,
		appName: appName,
		tagName: tagName,
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(customWriter),
		zap.NewAtomicLevelAt(zapcore.DebugLevel),
	)

	zapLogger := zap.New(core, zap.AddStacktrace(zapcore.ErrorLevel))

	return &CustomLogger{
		zapLogger:     zapLogger,
		channel:       channel,
		appName:       appName,
		tagName:       tagName,
		correlationID: generateUUID(),
	}, nil
}

func (l *CustomLogger) SetCorrelationID(id string) {
	l.correlationID = id
}

func (l *CustomLogger) GetCorrelationID() string {
	return l.correlationID
}

func (l *CustomLogger) SetAppName(name string) {
	l.appName = name
}

func (l *CustomLogger) GetAppName() string {
	return l.appName
}

func (l *CustomLogger) SetTagName(name string) {
	l.tagName = name
}

func (l *CustomLogger) GetTagName() string {
	return l.tagName
}

func (l *CustomLogger) SetChannel(name string) {
	l.channel = name
}

func (l *CustomLogger) GetChannel() string {
	return l.channel
}

func getStacktrace() string {
	pc := make([]uintptr, 10)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])

	frame, _ := frames.Next()
	return fmt.Sprintf("%s:%d", frame.File, frame.Line)
}

func (l *CustomLogger) createLogEntry(message string, level int, extra interface{}, exc *Exception) CustomLog {
	now := time.Now().UTC().Format(time.RFC3339Nano)[:23] + "Z"

	return CustomLog{
		Message: message,
		Context: Context{
			CorrelationID: l.correlationID,
			RequestID:     generateUUID(),
			AppName:       l.appName,
			TagName:       l.tagName,
			Exception:     exc,
		},
		Level:     level,
		LevelName: levelNames[level],
		Channel:   l.channel,
		Datetime:  now,
		Extra:     extra,
	}
}

func (l *CustomLogger) writeLog(entry CustomLog) {
	jsonData, err := json.Marshal(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao converter log para JSON: %v\n", err)
		return
	}
	fmt.Println(string(jsonData))
}

func (l *CustomLogger) Debug(message string, extra interface{}) {
	entry := l.createLogEntry(message, DebugLevel, extra, nil)
	l.writeLog(entry)
	l.zapLogger.Debug(message)
}

func (l *CustomLogger) Info(message string, extra interface{}) {
	entry := l.createLogEntry(message, InfoLevel, extra, nil)
	l.writeLog(entry)
	l.zapLogger.Info(message)
}

func (l *CustomLogger) Notice(message string, extra interface{}) {
	entry := l.createLogEntry(message, NoticeLevel, extra, nil)
	l.writeLog(entry)
	l.zapLogger.Info(message)
}

func (l *CustomLogger) Warning(message string, extra interface{}) {
	entry := l.createLogEntry(message, WarningLevel, extra, nil)
	l.writeLog(entry)
	l.zapLogger.Warn(message)
}

func (l *CustomLogger) Error(message string, extra interface{}) {
	exc := &Exception{
		Message: message,
		File:    getStacktrace(),
	}
	entry := l.createLogEntry(message, ErrorLevel, extra, exc)
	l.writeLog(entry)
	l.zapLogger.Error(message)
}

func (l *CustomLogger) Critical(message string, extra interface{}) {
	exc := &Exception{
		Message: message,
		File:    getStacktrace(),
	}
	entry := l.createLogEntry(message, CriticalLevel, extra, exc)
	l.writeLog(entry)
	l.zapLogger.DPanic(message)
}

func (l *CustomLogger) Alert(message string, extra interface{}) {
	exc := &Exception{
		Message: message,
		File:    getStacktrace(),
	}
	entry := l.createLogEntry(message, AlertLevel, extra, exc)
	l.writeLog(entry)
	l.zapLogger.Error(message, zap.String("level", "alert"))
}

func (l *CustomLogger) Emergency(message string, extra interface{}) {
	exc := &Exception{
		Message: message,
		File:    getStacktrace(),
	}
	entry := l.createLogEntry(message, EmergencyLevel, extra, exc)
	l.writeLog(entry)
	l.zapLogger.Error(message, zap.String("level", "emergency"))
}

func LoggerMiddleware(channel, appName, tagName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			correlationID := r.Header.Get("x-correlation-id")
			if correlationID == "" {
				correlationID = generateUUID()
			}

			if channel == "" || appName == "" || tagName == "" {
				envChannel, envAppName, envTagName := GetLoggerConfig()

				if channel == "" {
					channel = envChannel
				}

				if appName == "" {
					appName = envAppName
				}

				if tagName == "" {
					tagName = envTagName
				}
			}

			requestLogger, err := NewCustomLogger(channel, appName, tagName)
			if err != nil {
				http.Error(w, "Error initializing logger", http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), correlationIDKey, correlationID)

			requestLogger.SetCorrelationID(correlationID)

			ctx = context.WithValue(ctx, loggerKey, requestLogger)

			requestLogger.Info("Request received", map[string]string{
				"method": r.Method,
				"path":   r.URL.Path,
			})

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetLoggerFromContext(ctx context.Context) *CustomLogger {
	if logger, ok := ctx.Value(loggerKey).(*CustomLogger); ok {
		return logger
	}

	channel, appName, tagName := GetLoggerConfig()
	defaultLogger, _ := NewCustomLogger(channel, appName, tagName)
	return defaultLogger
}

func GetCorrelationIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(correlationIDKey).(string); ok {
		return id
	}
	return ""
}

func main() {
	channel, appName, tagName := GetLoggerConfig()

	log, err := NewCustomLogger(channel, appName, tagName)
	if err != nil {
		panic(err)
	}

	zapLogger := zap.New(zapcore.NewNopCore())
	zap.ReplaceGlobals(zapLogger)

	// log.Info("Logger inicializado", map[string]string{
	// 	"channel": log.GetChannel(),
	// 	"appName": log.GetAppName(),
	// 	"tagName": log.GetTagName(),
	// })

	// Logs de exemplo
	log.Debug("Debug message", "Level 100, Priority 0, Severity 0")
	log.Info("Info message", map[string]string{"addInfo": "Level 200, Priority 8, Severify 1"})
	log.Notice("Notice message", map[string]string{"addInfo": "Level 250, Priority 16, Severity 2"})
	log.Warning("Warning message", "Level 300, Priority 24, Severify 3")
	log.Error("Error message", map[string]string{"addInfo": "Level 400, Priority 32, Severity 4"})
	log.Critical("critical message", map[string]string{"addInfo": "Level 500, Priority 40, Severity 5"})
	log.Alert("alert message", "Level 550, Priority 48, Priority 6")
	log.Emergency("emergency message", map[string]string{"addInfo": "Level 600, Priority 56, Severity 7"})
}
