package gorad

import (
	"io"
	"log/slog"
	"os"
)

// write to stdout if len(path) == 0
func NewTextFileSlogHandler(path string, level slog.Leveler) *slog.TextHandler {
	var w io.Writer
	var err error
	if len(path) > 0 {
		w, err = os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
		if err != nil {
			panic(err)
		}
	} else {
		w = os.Stdout
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}
	return slog.NewTextHandler(w, opts)
}
