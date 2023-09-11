package tui

import (
	"fmt"

	"github.com/rivo/tview"
)

type messageTextView struct {
	flex  *tview.Flex
	about *tview.TextView
	*tview.TextView
	messageRingBuf
}

func newMessageTextView(messagesLimit int) *messageTextView {
	about := tview.NewTextView().SetDynamicColors(true) //.
	messager := tview.NewTextView()
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(about, 0, 1, false).
		AddItem(messager, 0, 10, false)
	return &messageTextView{flex, about, messager, newMessageRingBuf(messagesLimit)}
}
func (mtv *messageTextView) setAbout(version, buildTime string) {
	mtv.about.SetText(fmt.Sprintf("[green]Version:[white] %s [green]BuildTime:[white] %s", version, buildTime))
}

func (mtv *messageTextView) setMessage(msg string) {
	mtv.add(msg)
	mtv.SetText(mtv.read())
}

func (mtv *messageTextView) setWarning(msg string) {
	mtv.add("[Warning] " + msg)
	mtv.SetText(mtv.read())
}

func (mtv *messageTextView) setError(msg string) {
	mtv.add("[Error] " + msg)
	mtv.SetText(mtv.read())
}

type messageRingBuf struct {
	messagesLimit int
	buf           []string
	size          int
	head          int
}

func newMessageRingBuf(messagesLimit int) messageRingBuf {
	return messageRingBuf{
		messagesLimit: messagesLimit,
		buf:           make([]string, messagesLimit),
	}
}

func (mrb *messageRingBuf) add(msg string) {
	if mrb.size == mrb.messagesLimit {
		mrb.buf[mrb.head] = msg
		mrb.head = (mrb.head + 1) % mrb.messagesLimit
	} else {
		mrb.buf[mrb.size] = msg
		mrb.size++
	}
}

func (mrb *messageRingBuf) read() string {
	var s string
	head := mrb.head
	for i := 0; i < mrb.size; i++ {
		s = s + mrb.buf[head] + "\n"
		head = (head + 1) % mrb.messagesLimit
	}
	return s
}
