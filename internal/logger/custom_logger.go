package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/nicksonjean/logger.walkaba.net/internal/config"
	"github.com/nicksonjean/logger.walkaba.net/internal/domain/constants"
	"github.com/nicksonjean/logger.walkaba.net/internal/domain/models"
	"github.com/nicksonjean/logger.walkaba.net/pkg/utils"
)

type CustomLogger struct {
	zapLogger     *zap.Logger
	channel       string
	appName       string
	tagName       string
	correlationID string
}

func NewCustomLogger(channel, appName, tagName string) (*CustomLogger, error) {
	if channel == "" || appName == "" || tagName == "" {
		envChannel, envAppName, envTagName := config.GetLoggerConfig()

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
		correlationID: utils.GenerateUUID(),
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

func (l *CustomLogger) createLogEntry(message string, level int, extra interface{}, exc *models.Exception) models.CustomLog {
	now := time.Now().UTC().Format(time.RFC3339Nano)[:23] + "Z"

	return models.CustomLog{
		Message: message,
		Context: models.Context{
			CorrelationID: l.correlationID,
			RequestID:     utils.GenerateUUID(),
			AppName:       l.appName,
			TagName:       l.tagName,
			Exception:     exc,
		},
		Level:     level,
		LevelName: constants.LevelNames[level],
		Channel:   l.channel,
		Datetime:  now,
		Extra:     extra,
	}
}

func (l *CustomLogger) writeLog(entry models.CustomLog) {
	jsonData, err := json.Marshal(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao converter log para JSON: %v\n", err)
		return
	}
	fmt.Println(string(jsonData))
}

func (l *CustomLogger) Debug(message string, extra interface{}) {
	entry := l.createLogEntry(message, constants.DebugLevel, extra, nil)
	l.writeLog(entry)
	l.zapLogger.Debug(message)
}

func (l *CustomLogger) Info(message string, extra interface{}) {
	entry := l.createLogEntry(message, constants.InfoLevel, extra, nil)
	l.writeLog(entry)
	l.zapLogger.Info(message)
}

func (l *CustomLogger) Notice(message string, extra interface{}) {
	entry := l.createLogEntry(message, constants.NoticeLevel, extra, nil)
	l.writeLog(entry)
	l.zapLogger.Info(message)
}

func (l *CustomLogger) Warning(message string, extra interface{}) {
	entry := l.createLogEntry(message, constants.WarningLevel, extra, nil)
	l.writeLog(entry)
	l.zapLogger.Warn(message)
}

func (l *CustomLogger) Error(message string, extra interface{}) {
	exc := &models.Exception{
		Message: message,
		File:    GetStacktrace(),
	}
	entry := l.createLogEntry(message, constants.ErrorLevel, extra, exc)
	l.writeLog(entry)
	l.zapLogger.Error(message)
}

func (l *CustomLogger) Critical(message string, extra interface{}) {
	exc := &models.Exception{
		Message: message,
		File:    GetStacktrace(),
	}
	entry := l.createLogEntry(message, constants.CriticalLevel, extra, exc)
	l.writeLog(entry)
	l.zapLogger.Error(message, zap.String("level", "critical"))
}

func (l *CustomLogger) Alert(message string, extra interface{}) {
	exc := &models.Exception{
		Message: message,
		File:    GetStacktrace(),
	}
	entry := l.createLogEntry(message, constants.AlertLevel, extra, exc)
	l.writeLog(entry)
	l.zapLogger.Error(message, zap.String("level", "alert"))
}

func (l *CustomLogger) Emergency(message string, extra interface{}) {
	exc := &models.Exception{
		Message: message,
		File:    GetStacktrace(),
	}
	entry := l.createLogEntry(message, constants.EmergencyLevel, extra, exc)
	l.writeLog(entry)
	l.zapLogger.Error(message, zap.String("level", "emergency"))
}
