package tui

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRingBuf(t *testing.T) {
	messages := []string{"message1", "message2", "message3", "message4"}
	limit := len(messages) - 1
	t.Run("write less than limit", func(t *testing.T) {
		rb := newMessageRingBuf(limit)
		for i := 0; i < limit-1; i++ {
			rb.add(messages[i])
		}
		expected := strings.Join(messages[:limit-1], "\n") + "\n"
		actual := rb.read()
		assert.Equal(t, expected, actual)
	})
	t.Run("write limit messages", func(t *testing.T) {
		rb := newMessageRingBuf(limit)
		for i := 0; i < limit; i++ {
			rb.add(messages[i])
		}
		expected := strings.Join(messages[:limit], "\n") + "\n"
		actual := rb.read()
		assert.Equal(t, expected, actual)
	})
	t.Run("write more than limit", func(t *testing.T) {
		rb := newMessageRingBuf(limit)
		for i := 0; i < limit+1; i++ {
			rb.add(messages[i])
		}
		expected := strings.Join(messages[1:limit+1], "\n") + "\n"
		actual := rb.read()
		assert.Equal(t, expected, actual)
	})
}
func TestMessager(t *testing.T) {
	setSimulationScreen()
	const (
		testMessage = "test message"
	)
	messager = newMessageTextView(1)
	t.Run("set message", func(t *testing.T) {
		messager.setMessage(testMessage)
		actual := messager.GetText(false)
		expected := testMessage + "\n"
		assert.Equal(t, expected, actual)
		messager.SetText("")
	})
	t.Run("set warning", func(t *testing.T) {
		messager.setWarning(testMessage)
		actual := messager.GetText(false)
		expected := "[Warning] " + testMessage + "\n"
		assert.Equal(t, expected, actual)
	})
	t.Run("set error", func(t *testing.T) {
		messager.setError(testMessage)
		actual := messager.GetText(false)
		expected := "[Error] " + testMessage + "\n"
		assert.Equal(t, expected, actual)
	})
}
