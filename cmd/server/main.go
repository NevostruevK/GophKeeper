// Инициализация сервера.
// Для задания dsn поддерживается флаг -d.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NevostruevK/GophKeeper/internal/api/ftp"
	"github.com/NevostruevK/GophKeeper/internal/api/grpc/server"
	"github.com/NevostruevK/GophKeeper/internal/api/grpc/server/auth"
	"github.com/NevostruevK/GophKeeper/internal/api/grpc/server/keeper"
	"github.com/NevostruevK/GophKeeper/internal/config"
	"github.com/NevostruevK/GophKeeper/internal/tools/crypto"
	"github.com/NevostruevK/GophKeeper/internal/tools/cut"

	"github.com/NevostruevK/GophKeeper/internal/storage/postgres"
)

const shutDownTimeOut = time.Second * 3

const (
	address       = "127.0.0.1:8080"
	ftpAddress    = "127.0.0.1:8082"
	ftpDir        = "./../../build"
	dsn           = "user=postgres sslmode=disable"
	tokenKey      = "secretKeyForUserIdentification"
	tokenDuration = time.Hour
	cryptoKey     = "secretKeyForDataEncryptionForTheGophKeeper"
	cryptoNonce   = "nonceForGophKeeper"
)

func main() {
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	dbDSN := flag.String("d", dsn, "dsn for postgres")
	flag.Parse()

	cfg := config.Config{
		Address:       address,
		DSN:           *dbDSN,
		TokenKey:      tokenKey,
		EnableTLS:     true,
		TokenDuration: tokenDuration,
		FtpAddress:    ftpAddress,
		FtpDir:        ftpDir,
		CryptoKey:     cut.Cut(cryptoKey, 32),
		CryptoNonce:   cut.Cut(cryptoNonce, 12),
	}
	fmt.Println(cfg.DSN)

	crypto, err := crypto.NewCrypto([]byte(cfg.CryptoKey), []byte(cfg.CryptoNonce))
	if err != nil {
		log.Fatalf("failed to init crypto: %v", err)
	}

	storage, err := postgres.NewStorage(context.Background(), cfg.DSN, crypto)
	if err != nil {
		log.Fatalf("failed to init storage: %v", err)
	}

	keeperServer := keeper.NewKeeperServer(storage)
	jwtManager := auth.NewJWTManager(cfg.TokenKey, cfg.TokenDuration)
	options, err := server.NewServerOptions(jwtManager, cfg.EnableTLS)
	if err != nil {
		log.Fatalf("failed to initial server %v", err)
	}
	authServer := auth.NewAuthServer(storage, jwtManager)
	s := server.NewServer(authServer, keeperServer, options)
	go func() {
		if err = s.Start(cfg.Address); err != nil {
			log.Fatalf("failed to start gRPC server %v", err)
		}
	}()
	fs := ftp.NewServer(cfg.FtpAddress, cfg.FtpDir)
	go func() {
		if err = fs.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("failed to start ftp server %v", err)
		}
	}()

	<-gracefulShutdown
	s.GracefulStop()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutDownTimeOut)
	defer cancel()
	if err = fs.Shutdown(shutdownCtx); err != nil {
		log.Printf("ERROR : Server Shutdown error %v", err)
	} else {
		log.Printf("Server Shutdown ")
	}
	storage.Close()
}
