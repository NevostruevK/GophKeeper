package cut

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCut(t *testing.T) {
	t.Run("smaller than limit", func(t *testing.T) {
		out := "smaller than limit"
		in := Cut(out, 1000)
		assert.Equal(t, out, in)
	})
	t.Run("equal limit", func(t *testing.T) {
		out := "smaller than limit"
		in := Cut(out, len(out))
		assert.Equal(t, out, in)
	})
	t.Run("bigger than limit", func(t *testing.T) {
		out := "bigger than limit"
		limit := 6
		in := Cut(out, limit)
		assert.Equal(t, limit, len(in))
		assert.True(t, strings.HasPrefix(out, in))
	})
}
