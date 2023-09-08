package tui

import (
	"github.com/rivo/tview"
)

type menu struct {
	flex *tview.Flex
	*tview.List
	options []string
}

func newMenu(
	options []string,
	switchToStore func(),
	switchToLoad func(),
	switchToUser func(),
) *menu {
	flex := tview.NewFlex()
	m := &menu{
		flex:    flex,
		List:    tview.NewList().ShowSecondaryText(false),
		options: options,
	}
	for index, option := range m.options {
		m.AddItem(option, " ", rune(49+index), nil)
	}
	m.SetSelectedFunc(func(index int, name string, second_name string, shortcut rune) {
		switch name {
		case menuUser:
			//			pages.SwitchToPage(pageUser)
			switchToUser()
		case menuLoad:
			switchToLoad()
			//			pages.SwitchToPage(pageLoadForm)
		case menuStore:
			switchToStore()
		}
	})
	flex.
		AddItem(m, 0, 1, true).
		AddItem(messager, 0, 1, false)
	return m
}