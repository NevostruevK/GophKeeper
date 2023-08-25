package auth_test

import (
	"testing"
	"time"

	"github.com/NevostruevK/GophKeeper/internal/api/grpc/auth"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTManager_Verify(t *testing.T) {
	testID := uuid.New()
	t.Run("Verify ok", func(t *testing.T) {
		m := auth.NewJWTManager("secret_key", time.Hour)
		accseeToken, err := m.Generate(testID)
		require.NoError(t, err)
		ID, err := m.Verify(accseeToken)
		require.NoError(t, err)
		assert.Equal(t, testID, ID)
	})
	/*
		t.Run("invalide token err", func(t *testing.T) {
			m := auth.NewJWTManager("secret_key", time.Nanosecond)
			accseeToken, err := m.Generate(testID)
			require.NoError(t, err)
			time.Sleep(time.Second)
			ID, err := m.Verify(accseeToken)
			require.NoError(t, err)
			assert.Equal(t, testID, ID)
		})
	*/
}
