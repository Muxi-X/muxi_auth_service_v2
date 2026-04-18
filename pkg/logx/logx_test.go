package logx

import (
	"log/slog"
	"testing"
)

func TestParseLevel(t *testing.T) {
	if got := parseLevel("DEBUG"); got != slog.LevelDebug {
		t.Fatalf("expected debug, got %v", got)
	}
	if got := parseLevel("warn"); got != slog.LevelWarn {
		t.Fatalf("expected warn, got %v", got)
	}
	if got := parseLevel("ERROR"); got != slog.LevelError {
		t.Fatalf("expected error, got %v", got)
	}
	if got := parseLevel("unexpected"); got != slog.LevelInfo {
		t.Fatalf("expected info fallback, got %v", got)
	}
}
