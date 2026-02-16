package app_helpers

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"io"
	"log/slog"
	"testing"
	"time"
)

func TestNewCustomTextHandler_NilOptions(t *testing.T) {
	var buf bytes.Buffer
	handler := NewCustomTextHandler(&buf, nil)
	if handler == nil {
		t.Fatal("Expected handler to be non-nil")
	}
	if handler.out != &buf {
		t.Errorf("Expected out to be set to provided writer")
	}
	if handler.mu == nil {
		t.Errorf("Expected mu to be initialized")
	}
}

func TestNewCustomTextHandler_CleanLogLine(t *testing.T) {
	assert := assert.New(t)
	var ctx context.Context
	var buf bytes.Buffer

	ctx = context.Background()
	debugLine := "Oh! attaboy! library should not change this line at all"
	record := slog.NewRecord(time.Now(), slog.LevelDebug, debugLine, 0)
	handler := NewCustomTextHandler(&buf, nil)
	err := handler.Handle(ctx, record)
	if err != nil {
		t.Fatalf("Failed to handle log record: %v", err)
	}

	all, err := io.ReadAll(&buf)
	if err != nil {
		panic(err)
	}

	s := string(all)
	assert.Equal(s, debugLine+"\n")
}
