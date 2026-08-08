package logging

import (
	"bytes"
	"context"
	"sync"
)

// LogRecorder captures log messages in memory while still emitting them normally.
// This is useful for testing scenarios where you need to assert on log output.
type LogRecorder struct {
	buffer bytes.Buffer
	mu     sync.Mutex
}

// String returns all captured log messages as a string.
func (lr *LogRecorder) String() string {
	lr.mu.Lock()
	defer lr.mu.Unlock()
	return lr.buffer.String()
}

// Write implements io.Writer to capture log messages.
func (lr *LogRecorder) Write(p []byte) (n int, err error) {
	lr.mu.Lock()
	defer lr.mu.Unlock()
	return lr.buffer.Write(p)
}

type logRecorderContextKey struct{}

// WithLogRecorder returns a child context with a LogRecorder and the recorder itself.
// Logs emitted using the returned context will be captured in memory while still
// being emitted normally.
func WithLogRecorder(ctx context.Context) (context.Context, *LogRecorder) {
	if ctx == nil {
		ctx = context.Background()
	}

	logRecorder := &LogRecorder{}
	ctxWithRecorder := context.WithValue(ctx, logRecorderContextKey{}, logRecorder)

	return ctxWithRecorder, logRecorder
}

// getLogRecorderFromCtx retrieves the LogRecorder from context if present.
func getLogRecorderFromCtx(ctx context.Context) *LogRecorder {
	if ctx == nil {
		return nil
	}

	recorder, ok := ctx.Value(logRecorderContextKey{}).(*LogRecorder)
	if !ok {
		return nil
	}

	return recorder
}
