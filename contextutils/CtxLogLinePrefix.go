package contextutils

import "context"

func WithLogLinePrefix(ctx context.Context, logLinePrefix string) (ctxWithLogLinePrefix context.Context) {
	if ctx == nil {
		ctx = ContextSilent()
	}

	return context.WithValue(ctx, "logLinePrefix", logLinePrefix)
}

func GetLogLinePrefixFromCtx(ctx context.Context) (logLinePrefix string) {
	if ctx == nil {
		return ""
	}

	lineValue := ctx.Value("logLinePrefix")

	logLinePrefix, ok := lineValue.(string)
	if !ok {
		return logLinePrefix
	}

	return logLinePrefix
}
