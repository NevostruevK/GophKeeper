package postgres_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/NevostruevK/GophKeeper/internal/storage"
	"github.com/NevostruevK/GophKeeper/internal/storage/postgres"
	"github.com/NevostruevK/GophKeeper/internal/tools/crypto"
	"github.com/NevostruevK/GophKeeper/internal/tools/cut"
	_ "github.com/jackc/pgx/v5/stdlib"

	"time"

	"github.com/testcontainers/testcontainers-go"
	container "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testStorage storage.Storage

type postgreSQLContainer struct {
	*container.PostgresContainer
	mappedPort string
	host       string
}

func TestMain(m *testing.M) {
	const (
		cryptoKey   = "secretKeyForDataEncryptionForTheGophKeeper"
		cryptoNonce = "nonceForGophKeeper"
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	container, err := runPostgreSQLContainer()
	if err != nil {
		log.Fatalf("failed to setup postgres tests: %v", err)
	}
	crypto, err := crypto.NewCrypto([]byte(cut.Cut(cryptoKey, 32)), []byte(cut.Cut(cryptoNonce, 12)))
	if err != nil {
		log.Fatalf("failed to init crypto: %v", err)
	}

	testStorage, err = postgres.NewStorage(ctx, container.GetDSN(), crypto)
	if err != nil {
		log.Fatalf("failed to create storage: %v", err)
	}

	code := m.Run()
	container.Terminate(ctx)
	os.Exit(code)
}

const (
	dbName     = "keeper"
	dbUser     = "user"
	dbPassword = "password"
)

func (c postgreSQLContainer) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, c.host, c.mappedPort, dbName)
}

func runPostgreSQLContainer() (*postgreSQLContainer, error) {
	ctx := context.Background()

	postgresContainer, err := container.RunContainer(ctx,
		testcontainers.WithImage("docker.io/postgres:15.2-alpine"),
		container.WithDatabase(dbName),
		container.WithUsername(dbUser),
		container.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	host, err := postgresContainer.Host(ctx)
	if err != nil {
		return nil, err
	}

	mappedPort, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		return nil, err
	}

	return &postgreSQLContainer{
		PostgresContainer: postgresContainer,
		mappedPort:        mappedPort.Port(),
		host:              host,
	}, nil
}
