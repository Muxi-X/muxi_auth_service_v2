package logx

import (
	"testing"

	"go.uber.org/zap/zapcore"
)

func TestParseLevel(t *testing.T) {
	if got := parseLevel("DEBUG"); got != zapcore.DebugLevel {
		t.Fatalf("expected debug, got %v", got)
	}
	if got := parseLevel("warn"); got != zapcore.WarnLevel {
		t.Fatalf("expected warn, got %v", got)
	}
	if got := parseLevel("ERROR"); got != zapcore.ErrorLevel {
		t.Fatalf("expected error, got %v", got)
	}
	if got := parseLevel("unexpected"); got != zapcore.InfoLevel {
		t.Fatalf("expected info fallback, got %v", got)
	}
}

func TestSanitizeKeyValues(t *testing.T) {
	tests := []struct {
		name string
		in   []any
		want int
	}{
		{
			name: "even key values keep length",
			in:   []any{"key", "value"},
			want: 2,
		},
		{
			name: "odd key values append placeholder",
			in:   []any{"key"},
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizeKeyValues(tt.in)
			if len(got) != tt.want {
				t.Fatalf("expected length %d, got %d", tt.want, len(got))
			}
		})
	}
}
