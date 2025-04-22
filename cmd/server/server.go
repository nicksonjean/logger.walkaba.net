package server

import (
	"net/http"
	"strconv"

	"logger.walkaba.net/internal/config"
	"logger.walkaba.net/internal/middleware"
)

func StartServer(port int) error {
	mux := http.NewServeMux()

	channel, appName, tagName := config.GetLoggerConfig()
	logMiddleware := middleware.LoggerMiddleware(channel, appName, tagName)

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

	handler := logMiddleware(mux)

	return http.ListenAndServe(":"+strconv.Itoa(port), handler)
}
