package server

import "testing"

func TestHandler(t *testing.T) {
	h := NewHandler("test")
	result := h.Handle("request")
	if result == "" {
		t.Error("expected non-empty result")
	}
}

func TestHandlerName(t *testing.T) {
	h := NewHandler("myhandler")
	got := h.Handle("ping")
	expected := "myhandler: handled ping"
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}

// testHelperUnused is dead code (test scope)
func testHelperUnused() string {
	return "unused test helper"
}
