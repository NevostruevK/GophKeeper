// for building:
// go build -o ./../../build/windowsClient.exe -ldflags "-X main.version=v1.0.1 -X main.buildTime=2023.09.15"
// env GOOS=linux GOARCH=amd64 go build -o ./../../build/linuxAmd64 -ldflags "-X main.version=v1.0.1 -X main.buildTime=2023.09.15"
// env GOOS=darwin GOARCH=amd64 go build -o ./../../build/darwinAmd64 -ldflags "-X main.version=v1.0.1 -X main.buildTime=2023.09.15"

// Инициализация клиента.
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NevostruevK/GophKeeper/internal/api/grpc/client"
	"github.com/NevostruevK/GophKeeper/internal/config"
	"github.com/NevostruevK/GophKeeper/internal/service"
	"github.com/NevostruevK/GophKeeper/internal/tui"
)

var (
	version   = "N/A"
	buildTime = "N/A"
)

const (
	address       = "127.0.0.1:8080"
	enableTLS     = true
	tokenDuration = time.Hour
)

func main() {
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	cfg := config.Config{
		Address:       address,
		EnableTLS:     true,
		TokenDuration: tokenDuration,
	}

	client, err := client.NewClient(cfg.Address, cfg.EnableTLS)
	if err != nil {
		log.Fatalf("initial gRPC client failed with error: %v", err)
	}
	service := service.NewService(client)

	err = tui.NewTui(service, version, buildTime).Start()
	if err != nil {
		log.Fatalf("user interface terminal failed with error: %v", err)
	}
	<-gracefulShutdown
	if err = client.Close(); err != nil {
		log.Println(err)
	}
}
