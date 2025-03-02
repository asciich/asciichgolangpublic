package contextutils

import "context"

// Returns a child context of ctx with verbosity enabled according to `verbose`.
// If ctx is nil a new context with verbosity set is returned.
func WithVerbosityContextByBool(ctx context.Context, verbose bool) (ctxWithVerbosity context.Context) {
	if ctx == nil {
		return GetVerbosityContextByBool(verbose)
	}

	return context.WithValue(ctx, "verbose", verbose)
}

func GetVerbosityContextByBool(verbose bool) (ctx context.Context) {
	if verbose {
		return ContextVerbose()
	}

	return ContextSilent()
}

func ContextVerbose() (ctx context.Context) {
	return context.WithValue(context.Background(), "verbose", true)
}

func ContextSilent() (ctx context.Context) {
	return context.WithValue(context.Background(), "verbose", false)
}

func GetVerboseFromContext(ctx context.Context) (verbose bool) {
	if ctx == nil {
		return false
	}

	verboseValue := ctx.Value("verbose")
	if verboseValue == nil {
		return false
	}

	verbose, ok := verboseValue.(bool)
	if !ok {
		return false
	}

	return verbose
}
