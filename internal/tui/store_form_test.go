package tui

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreForm_setChoice(t *testing.T) {
	service, server, client, err := startService()
	require.NoError(t, err)
	defer func() {
		err = stopService(server, client)
		require.NoError(t, err)
	}()
	ch := make(chan any)
	go func() {
		err := run(service, "", "", ch)
		require.NoError(t, err)
	}()
	defer app.Stop()
	<-ch
	storeFormTest := newStoreForm()

	t.Run("set pair", func(t *testing.T) {
		storeFormTest.choice.SetCurrentOption(0)
		actual := storeFormTest.data.GetFormItem(0).GetLabel()
		expected := "Login"
		assert.Equal(t, expected, actual)
	})
	t.Run("set text", func(t *testing.T) {
		storeFormTest.choice.SetCurrentOption(1)
		actual, _ := pages.GetFrontPage()
		expected := pageInputText
		assert.Equal(t, expected, actual)
	})
	t.Run("set file", func(t *testing.T) {
		storeFormTest.choice.SetCurrentOption(2)
		actual, _ := pages.GetFrontPage()
		expected := pagePickFile
		assert.Equal(t, expected, actual)
	})
	t.Run("set card", func(t *testing.T) {
		storeFormTest.choice.SetCurrentOption(3)
		actual := storeFormTest.data.GetFormItem(0).GetLabel()
		expected := "Card owner"
		assert.Equal(t, expected, actual)
	})
}
