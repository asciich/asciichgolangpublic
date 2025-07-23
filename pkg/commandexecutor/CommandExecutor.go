package commandexecutor

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
)

type ContextKeyLiveOutputOnStdout struct{}

func WithLiveOutputOnStdoutIfVerbose(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = contextutils.ContextSilent()
	}

	return WithLiveOutputOnStdoutEnabled(ctx, contextutils.GetVerboseFromContext(ctx))
}

func WithLiveOutputOnStdoutEnabled(ctx context.Context, enabled bool) context.Context {
	if ctx == nil {
		ctx = contextutils.ContextVerbose()
	}

	return context.WithValue(ctx, ContextKeyLiveOutputOnStdout{}, enabled)
}

func WithLiveOutputOnStdout(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = contextutils.ContextVerbose()
	}

	return WithLiveOutputOnStdoutEnabled(ctx, true)
}

func IsLiveOutputOnStdoutEnabled(ctx context.Context) bool {
	if ctx == nil {
		return false
	}

	val := ctx.Value(ContextKeyLiveOutputOnStdout{})
	if val == nil {
		return false
	}

	valBool, ok := val.(bool)
	if !ok {
		return false
	}

	return valBool
}
