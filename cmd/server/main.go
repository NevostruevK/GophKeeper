package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NevostruevK/GophKeeper/internal/api/grpc/server"
	"github.com/NevostruevK/GophKeeper/internal/api/grpc/server/auth"
	"github.com/NevostruevK/GophKeeper/internal/api/grpc/server/keeper"
	"github.com/NevostruevK/GophKeeper/internal/config"
	"github.com/NevostruevK/GophKeeper/internal/config/duration"

	"github.com/NevostruevK/GophKeeper/internal/storage/postgres"
)

const (
	address       = "127.0.0.1:8080"
	DSN           = "user=postgres sslmode=disable"
	tokenKey      = "test_secret_key"
	tokenDuration = time.Hour
)

func main() {
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	log.Println("Start server")
	cfg := config.NewConfig(
		config.WithAddress(address),
		config.WithDSN(DSN),
		config.WithTokenKey(tokenKey),
		config.WithTokenDuration(duration.NewDuration(tokenDuration)),
		config.WithEnableTLS(true),
	)

	storage, err := postgres.NewStorage(context.Background(), cfg.DSN)
	if err != nil {
		log.Fatalf("failed to init storage: %v", err)
	}

	keeperServer := keeper.NewKeeperServer(storage)
	jwtManager := auth.NewJWTManager(cfg.TokenKey, cfg.TokenDuration.Duration)
	options, err := server.NewServerOptions(jwtManager, cfg.EnableTLS)
	if err != nil {
		log.Fatalf("failed to initial server %v", err)
	}
	authServer := auth.NewAuthServer(storage, jwtManager)
	s := server.NewServer(authServer, keeperServer, options)
	go s.Start(cfg.Address)
	<-gracefulShutdown
	s.GracefulStop()
	err = storage.DeleteAll(context.Background())
	if err != nil {
		log.Println("failed to clean storage: ", err)
	}
	storage.Close()
}
