package core_test

import (
	"errors"
	"fmt"
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

	t.Run("with a wrapped core error", func(t *testing.T) {
		w := httptest.NewRecorder()
		core.WriteError(w, fmt.Errorf("wrap: %w", core.Error(422, "bad input")))
		if w.Code != 422 {
			t.Errorf("status code = %v, want 422", w.Code)
		}
	})

	t.Run("with zero-status ServerError defaults to 500", func(t *testing.T) {
		w := httptest.NewRecorder()
		core.WriteError(w, core.ServerError{})
		if w.Code != 500 {
			t.Errorf("status code = %v, want 500", w.Code)
		}
	})
}

func TestBadRequest(t *testing.T) {
	err := core.BadRequest("invalid input")
	var se core.ServerError
	if !errors.As(err, &se) {
		t.Fatal("expected ServerError")
	}
	if se.Status != 400 {
		t.Errorf("status = %d, want 400", se.Status)
	}
	if se.Message != "invalid input" {
		t.Errorf("message = %q, want %q", se.Message, "invalid input")
	}
}
