package tui

import (
	"context"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/NevostruevK/GophKeeper/internal/service"
	"github.com/rivo/tview"
)

const messagesLimit = 5

const (
	pageUser      = "User"
	pageMenu      = "Menu"
	pageStoreForm = "Store"
	pagePickFile  = "File"
	pageInputText = "Text"
	pageLoadForm  = "Load"
)
const (
	menuLoad  = "Load data"
	menuStore = "Store data"
	menuUser  = "New user"
)

var app = tview.NewApplication()
var pages = tview.NewPages()
var messager = newMessageTextView(messagesLimit)
var srv = &service.Service{}

var types = []string{
	"ALL",
	string(models.PAIR),
	string(models.TEXT),
	string(models.FILE),
	string(models.CARD),
}

func Run(ctx context.Context, service *service.Service) error {
	srv = service
	storeForm := newStoreForm()
	filePicker := newFilePicker(storeForm.getPickFileFunc())
	inputText := newInputText(storeForm.getPickTextFunc())
	loadForm := newLoadForm()
	userForm := newUserForm()
	menu := newMenu(
		[]string{menuLoad, menuStore, menuUser},
		storeForm.getSwitchFromMenuFunc(),
		loadForm.getSwitchFromMenuFunc(),
		userForm.getSwitchFromMenuFunc(),
	)
	srv.Login(context.Background(), &models.User{Login: "konstantin", Password: "konstantin"})

	pages.AddPage(pageUser, userForm.flex, true, false)
	pages.AddPage(pageMenu, menu.flex, true, true)
	pages.AddPage(pageStoreForm, storeForm.flex, true, false)
	pages.AddPage(pagePickFile, filePicker.flex, true, false)
	pages.AddPage(pageInputText, inputText.grid, true, false)
	pages.AddPage(pageInputText, inputText.grid, true, false)
	pages.AddPage(pageLoadForm, loadForm.flex, true, false)

	return app.SetRoot(pages, true).EnableMouse(true).Run()
}