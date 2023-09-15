package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type inputText struct {
	grid *tview.Grid
	*tview.TextArea
}

func newInputText(pickText func(path string)) *inputText {

	textArea := tview.NewTextArea().
		SetPlaceholder("Enter text here...")
	textArea.SetTitle("Text Area").SetBorder(true).SetFocusFunc(func() { textArea.SetText("", false) })
	helpInfo := tview.NewTextView().
		SetText(" Press Ctrl-S or F2 for saving text, press F1 for help, press Ctrl-C to exit")
	grid := tview.NewGrid().
		SetRows(0, 1).
		AddItem(textArea, 0, 0, 1, 2, 0, 0, true).
		AddItem(helpInfo, 1, 0, 1, 1, 0, 0, false)
	addHelp()
	textArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyF1 {
			pages.ShowPage(pageInputText)
			return nil
		}
		if event.Key() == tcell.KeyF2 || event.Key() == tcell.KeyCtrlS {
			text := textArea.GetText()
			pickText(text)
			return nil
		}
		return event
	})
	return &inputText{grid: grid, TextArea: textArea}
}

func addHelp() {
	help1 := tview.NewTextView().
		SetDynamicColors(true).
		SetText(`[green]Navigation

[yellow]Left arrow[white]: Move left.
[yellow]Right arrow[white]: Move right.
[yellow]Down arrow[white]: Move down.
[yellow]Up arrow[white]: Move up.
[yellow]Ctrl-A, Home[white]: Move to the beginning of the current line.
[yellow]Ctrl-E, End[white]: Move to the end of the current line.
[yellow]Ctrl-F, page down[white]: Move down by one page.
[yellow]Ctrl-B, page up[white]: Move up by one page.
[yellow]Alt-Up arrow[white]: Scroll the page up.
[yellow]Alt-Down arrow[white]: Scroll the page down.
[yellow]Alt-Left arrow[white]: Scroll the page to the left.
[yellow]Alt-Right arrow[white]:  Scroll the page to the right.
[yellow]Alt-B, Ctrl-Left arrow[white]: Move back by one word.
[yellow]Alt-F, Ctrl-Right arrow[white]: Move forward by one word.

[blue]Press Enter for more help, press Escape to return.`)
	help2 := tview.NewTextView().
		SetDynamicColors(true).
		SetText(`[green]Editing[white]

Type to enter text.
[yellow]Ctrl-H, Backspace[white]: Delete the left character.
[yellow]Ctrl-D, Delete[white]: Delete the right character.
[yellow]Ctrl-K[white]: Delete until the end of the line.
[yellow]Ctrl-W[white]: Delete the rest of the word.
[yellow]Ctrl-U[white]: Delete the current line.

[blue]Press Enter for more help, press Escape to return.`)
	help3 := tview.NewTextView().
		SetDynamicColors(true).
		SetText(`[green]Selecting Text[white]

Move while holding Shift or drag the mouse.
Double-click to select a word.
[yellow]Ctrl-L[white] to select entire text.

[green]Clipboard

[yellow]Ctrl-Q[white]: Copy.
[yellow]Ctrl-X[white]: Cut.
[yellow]Ctrl-V[white]: Paste.
		
[green]Undo

[yellow]Ctrl-Z[white]: Undo.
[yellow]Ctrl-Y[white]: Redo.

[blue]Press Enter for more help, press Escape to return.`)
	help := tview.NewFrame(help1).
		SetBorders(1, 1, 0, 0, 2, 2)
	help.SetBorder(true).
		SetTitle("Help").
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyEscape {
				pages.SwitchToPage(pageInputText)
				return nil
			} else if event.Key() == tcell.KeyEnter {
				switch {
				case help.GetPrimitive() == help1:
					help.SetPrimitive(help2)
				case help.GetPrimitive() == help2:
					help.SetPrimitive(help3)
				case help.GetPrimitive() == help3:
					help.SetPrimitive(help1)
				}
				return nil
			}
			return event
		})
	pages.AddPage("textHelp", tview.NewGrid().
		SetColumns(0, 64, 0).
		SetRows(0, 22, 0).
		AddItem(help, 1, 1, 1, 1, 0, 0, true), true, false)
}
