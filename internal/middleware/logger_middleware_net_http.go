package middleware

import (
	"context"
	"net/http"

	"github.com/nicksonjean/logger.walkaba.net/internal/config"
	"github.com/nicksonjean/logger.walkaba.net/internal/logger"
	"github.com/nicksonjean/logger.walkaba.net/pkg/utils"
)

func LoggerMiddlewareNetHttp(channel, appName, tagName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			correlationID := r.Header.Get("x-correlation-id")
			if correlationID == "" {
				correlationID = utils.GenerateUUID()
			}

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

			requestLogger, err := logger.NewCustomLogger(channel, appName, tagName)
			if err != nil {
				http.Error(w, "Error initializing logger", http.StatusInternalServerError)
				return
			}

			requestLogger.SetCorrelationID(correlationID)

			ctx := context.WithValue(r.Context(), logger.CorrelationIDKey, correlationID)
			ctx = context.WithValue(ctx, logger.LoggerKey, requestLogger)

			requestLogger.Info("Request received", map[string]string{
				"method": r.Method,
				"path":   r.URL.Path,
				"host":   r.Host,
			})

			w.Header().Set("X-Correlation-ID", correlationID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetLoggerFromContext(ctx context.Context) *logger.CustomLogger {
	if logger, ok := ctx.Value(logger.LoggerKey).(*logger.CustomLogger); ok {
		return logger
	}

	channel, appName, tagName := config.GetLoggerConfig()
	defaultLogger, _ := logger.NewCustomLogger(channel, appName, tagName)
	return defaultLogger
}

func GetCorrelationIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(logger.CorrelationIDKey).(string); ok {
		return id
	}
	return ""
}
