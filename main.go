package main

import (
	"log"

	"logger.walkaba.net/cmd/server"
	"logger.walkaba.net/internal/config"
)

func main() {
	err := config.LoadEnv()
	if err != nil {
		log.Printf("Aviso: Não foi possível carregar o arquivo .env: %v\n", err)
	}

	middleware, host, port := config.GetStartServerConfig()

	var serverErr error
	if middleware == "fiber" {
		serverErr = server.StartServerFiber(port)
	} else {
		serverErr = server.StartServerNetHttp(host, port)
	}
	if serverErr != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v\n", serverErr)
	}
}
