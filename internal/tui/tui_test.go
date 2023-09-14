package tui

import (
	"context"
	"time"

	"github.com/NevostruevK/GophKeeper/internal/api/grpc/client"
	"github.com/NevostruevK/GophKeeper/internal/api/grpc/server"
	"github.com/NevostruevK/GophKeeper/internal/api/grpc/server/auth"
	"github.com/NevostruevK/GophKeeper/internal/api/grpc/server/keeper"
	"github.com/NevostruevK/GophKeeper/internal/config"
	"github.com/NevostruevK/GophKeeper/internal/service"
	"github.com/NevostruevK/GophKeeper/internal/storage/memory"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	errUnauthenticated        = "[Error] rpc error: code = Unauthenticated desc = access token is invalid: invalid token: token contains an invalid number of segments\n"
	errIncorrectLoginPassword = "[Error] rpc error: code = NotFound desc = incorrect login/password\n"
	mesThereAreNoEntries      = "there are no entries\n"
)

func setSimulationScreen() tcell.SimulationScreen {
	simScreen := tcell.NewSimulationScreen("UTF-8")
	simScreen.Init()
	simScreen.SetSize(10, 10)
	app = tview.NewApplication()
	app.SetScreen(simScreen)
	return simScreen
}

func run(service *service.Service, version, builtTime string, ch chan any) error {
	setSimulationScreen()
	app.SetAfterDrawFunc(func(screen tcell.Screen) {
		ch <- struct{}{}
	})
	pages = tview.NewPages()
	messager = newMessageTextView(1)
	messager.setAbout(version, builtTime)
	srv = service
	pages.AddPage(pageUser, tview.NewFlex().Box, true, false)
	pages.AddPage(pageMenu, tview.NewFlex().Box, true, false)
	pages.AddPage(pageStoreForm, tview.NewFlex().Box, true, false)
	pages.AddPage(pagePickFile, tview.NewFlex().Box, true, false)
	pages.AddPage(pageInputText, tview.NewFlex().Box, true, false)
	pages.AddPage(pageInputText, tview.NewFlex().Box, true, false)
	pages.AddPage(pageLoadForm, tview.NewFlex().Box, true, false)

	return app.SetRoot(pages, true).Run()
}

func startService() (*service.Service, *server.Server, *client.Client, error) {
	cfg := config.Config{
		Address:       "127.0.0.1:8080",
		TokenKey:      "secretKeyForUserIdentification",
		EnableTLS:     false,
		TokenDuration: time.Hour,
	}
	dataStorage := memory.NewDataStore()
	userSoorage := memory.NewUserStore()
	keeperServer := keeper.NewKeeperServer(dataStorage)
	jwtManager := auth.NewJWTManager(cfg.TokenKey, cfg.TokenDuration)
	options, err := server.NewServerOptions(jwtManager, cfg.EnableTLS)
	if err != nil {
		return nil, nil, nil, err
	}
	authServer := auth.NewAuthServer(userSoorage, jwtManager)
	server := server.NewServer(authServer, keeperServer, options)
	go server.Start(cfg.Address)
	client, err := client.NewClient(cfg.Address, cfg.EnableTLS)
	if err != nil {
		return nil, nil, nil, err
	}
	service := service.NewService(client)
	return service, server, client, nil
}

func stopService(server *server.Server, client *client.Client) error {
	server.Shutdown(context.TODO())
	return client.Close()
}
