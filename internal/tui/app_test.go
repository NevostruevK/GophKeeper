package tui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunApp(t *testing.T) {
	tui := NewTui(nil, "", "")
	assert.NotNil(t, tui)
}
