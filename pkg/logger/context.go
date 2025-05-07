package logger

type contextKey string

const (
	CorrelationIDKey contextKey = "correlation_id"
	LoggerKey        contextKey = "logger"
)
