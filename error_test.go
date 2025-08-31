package core_test

import (
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/newlix/core"
)

// Test error reporting.
func TestWriteError(t *testing.T) {
	t.Run("with a regular error", func(t *testing.T) {
		want := "boom"
		w := httptest.NewRecorder()
		core.WriteError(w, errors.New(want))
		if w.Code != 500 {
			t.Errorf("status code = %v, want 500", w.Code)
		}
		body := strings.TrimSpace(w.Body.String())
		if body != want {
			t.Errorf("body = %q, want %q", body, want)
		}
	})

	t.Run("with a core error", func(t *testing.T) {
		want := "Invalid team slug"
		w := httptest.NewRecorder()
		core.WriteError(w, core.Error(400, want))
		if w.Code != 400 {
			t.Errorf("status code = %v, want 400", w.Code)
		}
		body := strings.TrimSpace(w.Body.String())
		if body != want {
			t.Errorf("body = %q, want %q", body, want)
		}
	})
}
