package app_helpers

import (
	"context"
	"io"
	"log/slog"	
	"strings"
	"sync"
)

type CustomTextHandler struct {
	h   slog.Handler
	mu  *sync.Mutex
	out io.Writer
}

func NewCustomTextHandler(o io.Writer, opts *slog.HandlerOptions) *CustomTextHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &CustomTextHandler{
		out: o,
		h: slog.NewTextHandler(o, &slog.HandlerOptions{
			Level:       opts.Level,
			AddSource:   opts.AddSource,
			ReplaceAttr: nil,
		}),
		mu: &sync.Mutex{},
	}
}

func (h *CustomTextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.h.Enabled(ctx, level)
}

func (h *CustomTextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &CustomTextHandler{h: h.h.WithAttrs(attrs), out: h.out, mu: h.mu}
}

func (h *CustomTextHandler) WithGroup(name string) slog.Handler {
	return &CustomTextHandler{h: h.h.WithGroup(name), out: h.out, mu: h.mu}
}

func (h *CustomTextHandler) Handle(ctx context.Context, r slog.Record) error {
	//add message to values
	strs := []string{r.Message+"\n"}

	if r.NumAttrs() != 0 {
		r.Attrs(func(a slog.Attr) bool {
			strs = append(strs, a.Value.String())
			return true
		})
	}

	result := strings.Join(strs, " ")
	b := []byte(result)

	h.mu.Lock()
	defer h.mu.Unlock()

	_, err := h.out.Write(b)

	return err
}