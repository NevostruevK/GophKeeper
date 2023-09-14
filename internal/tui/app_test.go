package tui

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunApp(t *testing.T) {
	setSimulationScreen()

	go func() {
		err := Run(nil, "", "")
		require.NoError(t, err)
	}()
	defer app.Stop()

}
