package main

import (
	"fmt"
	"log"

	"logger.walkaba.net/cmd/server"
	"logger.walkaba.net/internal/config"
)

func main() {
	// Carrega configurações
	err := config.LoadEnv()
	if err != nil {
		log.Printf("Aviso: Não foi possível carregar o arquivo .env: %v\n", err)
	}

	// Inicia o servidor HTTP
	port := 8080
	fmt.Printf("Iniciando servidor na porta %d...\n", port)
	serverErr := server.StartServer(port)
	if serverErr != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v\n", serverErr)
	}
}
