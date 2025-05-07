package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nicksonjean/logger.walkaba.net/pkg/config"
	"github.com/nicksonjean/logger.walkaba.net/pkg/logger"
	"github.com/nicksonjean/logger.walkaba.net/pkg/utils"
)

func LoggerMiddlewareFiber(channel, appName, tagName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		correlationID := c.Get("x-correlation-id")
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
			return c.Status(fiber.StatusInternalServerError).SendString("Error initializing logger")
		}

		requestLogger.SetCorrelationID(correlationID)

		c.Locals(logger.CorrelationIDKey, correlationID)
		c.Locals(logger.LoggerKey, requestLogger)

		requestLogger.Info("Request received", map[string]string{
			"method": c.Method(),
			"path":   c.Path(),
			"host":   c.Hostname(),
		})

		c.Set("x-correlation-id", correlationID)

		return c.Next()
	}
}

func GetLoggerFromFiberCtx(c *fiber.Ctx) *logger.CustomLogger {
	if logger, ok := c.Locals(logger.LoggerKey).(*logger.CustomLogger); ok {
		return logger
	}

	channel, appName, tagName := config.GetLoggerConfig()
	defaultLogger, _ := logger.NewCustomLogger(channel, appName, tagName)
	return defaultLogger
}

func GetCorrelationIDFromFiberCtx(c *fiber.Ctx) string {
	if id, ok := c.Locals(logger.CorrelationIDKey).(string); ok {
		return id
	}
	return ""
}
