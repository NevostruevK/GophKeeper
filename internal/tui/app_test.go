package tui

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunApp(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	setSimulationScreen()

	go func() {
		err := Run(ctx, nil, "", "")
		require.NoError(t, err)
	}()
	defer app.Stop()

}
