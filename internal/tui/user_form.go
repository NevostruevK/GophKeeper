package tui

import (
	"context"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type userForm struct {
	flex *tview.Flex
	*tview.Form
	helpInfo *tview.TextView
	user     *models.User
}

func newUserForm() *userForm {
	helpInfo := tview.NewTextView().
		SetText(" Press Esc to go to the menu")
	uf := &userForm{Form: tview.NewForm(), helpInfo: helpInfo, user: &models.User{}}
	uf.updateForm()
	uf.flex = tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(uf, 0, 15, true).
			AddItem(uf.helpInfo, 0, 1, false), 0, 1, false).
		AddItem(messager.flex, 0, 1, false)
	uf.flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			pages.SwitchToPage(pageMenu)
			return nil
		}
		return event
	})
	return uf
}

func (uf *userForm) userRequst(isLogin bool) {
	ok, msg := uf.user.IsReadyForStorage()
	if !ok {
		messager.setWarning(msg)
		uf.SetFocus(0)
		return
	}
	var id string
	var err error
	if isLogin {
		id, err = srv.Login(context.Background(), uf.user)
	} else {
		id, err = srv.Register(context.Background(), uf.user)
	}
	if err != nil {
		messager.setError(err.Error())
		return
	}
	messager.setMessage(id)
	pages.SwitchToPage(pageMenu)
}

func (uf *userForm) updateForm() {
	uf.AddInputField("Login", "", 20, nil, func(login string) {
		uf.user.Login = login
	})

	uf.AddInputField("Password", "", 20, nil, func(password string) {
		uf.user.Password = password
	})

	uf.AddButton("Login", func() {
		uf.userRequst(true)
	})

	uf.AddButton("Register", func() {
		uf.userRequst(false)
	})
}

func (uf *userForm) getSwitchFromMenuFunc() func() {
	return func() {
		pages.SwitchToPage(pageUser)
		uf.Clear(true)
		uf.updateForm()
		uf.SetFocus(0)
		app.SetFocus(uf)
		//		app.SetFocus(uf.GetFormItem(0))
	}
}
