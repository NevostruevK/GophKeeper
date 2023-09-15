package tui

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type record struct {
	typ   models.MType
	title string
	entry models.Entry
}

type storeForm struct {
	flex     *tview.Flex
	choice   *tview.DropDown
	data     *tview.Form
	title    *tview.Form
	helpInfo *tview.TextView
	record   *record
}

func newStoreForm() *storeForm {
	flex := tview.NewFlex()
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			pages.SwitchToPage(pageMenu)
			return nil
		}
		return event
	})
	helpInfo := tview.NewTextView().
		SetText(" Press Esc to go to the menu")

	sf := &storeForm{
		data:     tview.NewForm(),
		flex:     flex,
		helpInfo: helpInfo,
		record:   &record{},
	}
	sf.setChoice()
	sf.title = tview.NewForm()

	sf.updateFlex()
	return sf
}

func (sf *storeForm) getPickFileFunc() func(string) {
	return func(path string) {
		file := &models.File{Name: path}
		sf.record.entry = file
		pages.SwitchToPage(pageStoreForm)
		sf.data.Clear(true)
		sf.updateFile(path)
		app.SetFocus(sf.title)
	}
}

func (sf *storeForm) getPickTextFunc() func(string) {
	return func(text string) {
		t := &models.Text{Text: []byte(text)}
		sf.record.entry = t
		pages.SwitchToPage(pageStoreForm)
		sf.data.Clear(true)
		sf.updateText()
		app.SetFocus(sf.title)
	}
}

func (sf *storeForm) getSwitchFromMenuFunc() func() {
	return func() {
		pages.SwitchToPage(pageStoreForm)
		sf.data.Clear(true)
		sf.title.Clear(true)
		app.SetFocus(sf.choice)
	}
}

func (sf *storeForm) updateTitleForm() {
	const emptyTitle = "empty title"
	sf.title.Clear(true)
	sf.title.AddInputField("Title", "", 20, nil, func(title string) {
		sf.record.title = title
	})
	sf.title.AddButton("Store", func() {
		if sf.record.title == "" {
			messager.setWarning(emptyTitle)
			sf.title.SetFocus(0)
			return
		}
		ok, msg := sf.record.entry.IsReadyForStorage()
		if !ok {
			messager.setWarning(msg)
			sf.data.SetFocus(0)
			app.SetFocus(sf.data)
			return
		}
		ds, err := srv.StoreEntry(context.Background(), sf.record.typ, sf.record.title, sf.record.entry)
		if err != nil {
			messager.setError(err.Error())
			return
		}
		messager.setMessage("store: " + string(sf.record.typ) + " : " + ds.ID.String())
		sf.choice.SetCurrentOption(-1)
		pages.SwitchToPage(pageMenu)
	})
}

func (sf *storeForm) setChoice() {
	choice := tview.NewDropDown().SetLabel("Type").SetFieldWidth(5)
	choice.SetOptions(types[1:], func(typ string, index int) {
		if index >= 0 {
			sf.record = &record{typ: models.StringToMType(typ)}
			sf.updateForm()
			sf.updateTitleForm()
			switch sf.record.typ {
			case models.TEXT:
				pages.SwitchToPage(pageInputText)
				return
			case models.FILE:
				pages.SwitchToPage(pagePickFile)
				return
			}
			app.SetFocus(sf.data)
		}
	})
	sf.choice = choice
}

func (sf *storeForm) updateForm() {
	sf.updateFlex()
	sf.data.Clear(true)
	switch sf.record.typ {
	case models.PAIR:
		sf.updatePair()
	case models.TEXT:
		sf.updateText()
	case models.FILE:
		sf.updateFile("")
	case models.CARD:
		sf.updateCard()
	}
}

func (sf *storeForm) updateFlex() {
	sf.flex.Clear().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(sf.choice, 0, 1, false).
			AddItem(sf.data, 0, 15, false).
			AddItem(sf.title, 0, 3, false).
			AddItem(sf.helpInfo, 0, 1, false), 0, 1, false).
		AddItem(messager.flex, 0, 1, false)
}

func (sf *storeForm) updatePair() {
	pair := &models.Pair{}
	sf.record.entry = pair
	sf.data.AddInputField("Login", "", 20, nil, func(login string) {
		pair.Login = login
	})
	sf.data.AddInputField("Password", "", 20, nil, func(password string) {
		pair.Password = password
	}).SetInputCapture(sf.changeFocusForLastItem())
}

func (sf *storeForm) updateFile(name string) {

	sf.data.AddInputField("File name", name, 20, nil, func(name string) {
		sf.getPickFileFunc()(name)
	})

	sf.data.AddButton("File Open", func() {
		pages.SwitchToPage(pagePickFile)
	})
}

func (sf *storeForm) updateText() {
	sf.data.AddButton("Enter text", func() {
		pages.SwitchToPage(pageInputText)
	})
}

func (sf *storeForm) updateCard() {
	checkNumbers := func(_ string, ch rune) bool {
		return ch >= 48 && ch < 58
	}

	getUint := func(s string) uint64 {
		i, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			messager.setWarning(err.Error())
			return 0
		}
		return i
	}

	card := &models.Card{}
	sf.record.entry = card
	sf.data.AddInputField("Card owner", "", 20, nil, func(name string) {
		card.Owner = name
	})

	sf.data.AddInputField("Card number", "", 20, checkNumbers, func(number string) {
		card.Number = getUint(number)
	})

	sf.data.AddInputField("Bank", "", 20, nil, func(bank string) {
		card.Bank = bank
	})

	setMonths := func() []string {
		months := make([]string, 12)
		for i := 0; i < 12; i++ {
			months[i] = fmt.Sprintf("%d", i+1)
		}
		return months
	}

	sf.data.AddDropDown("Expiration month", setMonths(), -1, func(month string, index int) {
		if index < 0 {
			return
		}
		card.ExpirationMonth = uint8(getUint(month))
	})

	setYears := func() []string {
		years := make([]string, 10)
		now := time.Now().Year()
		for i := 0; i < 10; i++ {
			years[i] = fmt.Sprintf("%d", now+i)
		}
		return years
	}

	sf.data.AddDropDown("Expiration year", setYears(), -1, func(year string, index int) {
		if index < 0 {
			return
		}
		card.ExpirationYear = uint16(getUint(year))
	})

	sf.data.AddInputField("CVV", "", 4, checkNumbers, func(cvv string) {
		card.CVV = uint16(getUint(cvv))
	}).SetInputCapture(sf.changeFocusForLastItem())
}

func (sf *storeForm) changeFocusForLastItem() func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter || event.Key() == tcell.KeyTAB {
			formItem, _ := sf.data.GetFocusedItemIndex()
			if sf.data.GetFormItemCount() == formItem+1 {
				app.SetFocus(sf.title)
				return nil
			}
		}
		return event
	}
}
