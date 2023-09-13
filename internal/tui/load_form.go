package tui

import (
	"context"
	"fmt"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type loadForm struct {
	flex     *tview.Flex
	choice   *tview.DropDown
	list     *tview.List
	title    *tview.TextView
	helpInfo *tview.TextView
	specs    []models.Spec
}

func newLoadForm() *loadForm {
	flex := tview.NewFlex()
	title := tview.NewTextView()
	list := tview.NewList().ShowSecondaryText(true)
	helpInfo := tview.NewTextView().
		SetText(" Press Esc to go to the menu\n Press Enter on an entry from the list to get data")

	lf := &loadForm{
		flex:     flex,
		list:     list,
		title:    title,
		helpInfo: helpInfo,
	}

	lf.setChoice()
	list.SetChangedFunc(func(index int, name string, second_name string, shortcut rune) {
		title.Clear()
		title.SetText(lf.specs[index].Title)
	})
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			index := list.GetCurrentItem()
			entry, err := srv.GetData(context.Background(), lf.specs[index])
			if err != nil {
				messager.setError(err.Error())
				return nil
			}
			title.SetText(entry.Show())
			return nil
		}
		return event
	})

	lf.updateFlex()
	return lf
}

func (lf *loadForm) updateFlex() {
	lf.flex.AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(lf.choice, 0, 5, true).
		AddItem(lf.list, 0, 15, false).
		AddItem(lf.helpInfo, 0, 2, false), 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(messager.flex, 0, 1, false).
			AddItem(lf.title, 0, 1, false), 0, 1, false)
	lf.flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			pages.SwitchToPage(pageMenu)
			return nil
		}
		return event
	})
}

func (lf *loadForm) setChoice() {
	choice := tview.NewDropDown().SetLabel("Type").SetFieldWidth(5)
	choice.SetOptions(types, func(typ string, index int) {
		if index >= 0 {
			specs, err := srv.LoadSpecs(context.Background(), models.StringToMType(typ))
			if err != nil {
				messager.setError(err.Error())
				return
			}
			lf.specs = specs
			if len(specs) == 0 {
				messager.setMessage("there are no entries")
				pages.SwitchToPage(pageMenu)
				return
			}
			messager.setMessage(fmt.Sprintf("found %d specs", len(specs)))
			lf.list.Clear()
			for index, spec := range specs {
				lf.list.AddItem(spec.Title, string(spec.Type), rune(49+index), nil)
			}
			app.SetFocus(lf.list)
		}
	})
	lf.choice = choice
}

func (lf *loadForm) getSwitchFromMenuFunc() func() {
	return func() {
		pages.SwitchToPage(pageLoadForm)
		lf.choice.SetCurrentOption(-1)
		lf.title.Clear()
		app.SetFocus(lf.choice)
	}
}
