package tui

import (
	"github.com/NevostruevK/GophKeeper/internal/service"
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
