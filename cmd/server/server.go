package server

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/nicksonjean/logger.walkaba.net/pkg/config"
	"github.com/nicksonjean/logger.walkaba.net/pkg/middleware"
	"github.com/nicksonjean/logger.walkaba.net/pkg/utils"
)

func StartServerNetHttp(host string, port int) error {
	mux := utils.NewCountingServeMux()

	channel, appName, tagName := config.GetLoggerConfig()
	logMiddleware := middleware.LoggerMiddlewareNetHttp(channel, appName, tagName)

	mux.HandleFunc("/api/logs", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"message": "Log endpoint funcionando!"}`))
			return
		}
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	})

	mux.HandleFunc("/api/test", func(w http.ResponseWriter, r *http.Request) {
		logger := middleware.GetLoggerFromContext(r.Context())
		logger.Info("Teste de endpoint com logger do contexto", map[string]string{
			"correlation_id": middleware.GetCorrelationIDFromContext(r.Context()),
			"path":           r.URL.Path,
		})

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Teste com logger do contexto realizado com sucesso!"}`))
	})

	utils.PrintServerBanner(host, port, mux.HandlersCount())

	handler := logMiddleware(mux)

	http.ListenAndServe(":"+strconv.Itoa(port), handler)

	return nil
}

func StartServerFiber(port int) error {
	app := fiber.New()

	channel, appName, tagName := config.GetLoggerConfig()
	app.Use(middleware.LoggerMiddlewareFiber(channel, appName, tagName))

	app.Get("/api/logs", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"message": "Log endpoint funcionando!",
		})
	})

	app.Get("/api/test", func(ctx *fiber.Ctx) error {
		logger := middleware.GetLoggerFromFiberCtx(ctx)
		logger.Info("Teste de endpoint com logger do contexto", map[string]string{
			"correlation_id": middleware.GetCorrelationIDFromContext(ctx.Context()),
			"path":           ctx.Path(),
		})

		return ctx.JSON(fiber.Map{
			"message": "Teste com logger do contexto realizado com sucesso!",
		})
	})

	app.Listen(":" + strconv.Itoa(port))

	return nil
}
