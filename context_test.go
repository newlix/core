package core_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/newlix/core"
	"github.com/tj/assert"
)

func TestRequestContext(t *testing.T) {
	t.Run("round trip", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/test", nil)
		ctx := core.NewRequestContext(context.Background(), r)
		got, ok := core.RequestFromContext(ctx)
		assert.True(t, ok)
		assert.Equal(t, r, got)
	})

	t.Run("missing request returns false", func(t *testing.T) {
		_, ok := core.RequestFromContext(context.Background())
		assert.Equal(t, false, ok)
	})
}
