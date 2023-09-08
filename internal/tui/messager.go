package tui

import "github.com/rivo/tview"

type messageTextView struct {
	*tview.TextView
	messageRingBuf
}

func newMessageTextView(messagesLimit int) *messageTextView {
	return &messageTextView{tview.NewTextView(), newMessageRingBuf(messagesLimit)}
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