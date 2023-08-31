package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NevostruevK/GophKeeper/internal/api/grpc/client"
	"github.com/NevostruevK/GophKeeper/internal/config"
	"github.com/NevostruevK/GophKeeper/internal/config/duration"
	"github.com/NevostruevK/GophKeeper/internal/models"
)

const (
	address   = "127.0.0.1:8080"
	enableTLS = true
)

var tokenDuration = duration.NewDuration(time.Hour)

func testClient(c *client.Client) {
	ctx := context.Background()
	u := models.User{Login: "some_login", Password: "some_password"}
	_, err := c.Auth.Register(ctx, u.Login, u.Password)
	if err != nil {
		log.Println(err)
	}
	titles := []string{"title1", "title2", "title3"}
	for _, t := range titles {
		r := models.NewRecord(models.TEXT, t, []byte("data for "+t))
		log.Println("Add Record")
		id, err := c.Keeper.AddRecord(ctx, r)
		if err != nil {
			log.Println(err)
		}
		log.Println(id)
	}
	log.Println("Add type of FILE ")
	r := models.NewRecord(models.FILE, "some title", []byte("some data"))
	_, err = c.Keeper.AddRecord(ctx, r)
	if err != nil {
		log.Println(err)
	}

	log.Println("Get all specs: ")
	specs, err := c.Keeper.GetSpecs(ctx)
	if err != nil {
		log.Println(err)
	}
	log.Println(specs)

	log.Println("Get specs for TEXT:")
	specs, err = c.Keeper.GetSpecsOfType(ctx, models.TEXT)
	if err != nil {
		log.Println(err)
	}
	log.Println(specs)
	log.Println("Get specs for FILE:")
	specs, err = c.Keeper.GetSpecsOfType(ctx, models.FILE)
	if err != nil {
		log.Println(err)
	}
	log.Println(specs)
	log.Println("Get Data:")
	data, err := c.Keeper.GetData(ctx, models.DataSpec{ID: specs[0].ID, DataSize: specs[0].DataSize})
	if err != nil {
		log.Println(err)
	}
	log.Println(string(data))
}

func main() {
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	cfg := config.NewConfig(
		config.WithAddress(address),
		config.WithTokenDuration(tokenDuration),
		config.WithEnableTLS(enableTLS),
	)

	client, err := client.NewClient(cfg.Address, cfg.EnableTLS)
	if err != nil {
		log.Fatalf("initial gRPC client failed with error: %v", err)
	}
	testClient(client)
	<-gracefulShutdown
	if err = client.Close(); err != nil {
		log.Println(err)
	}
}
