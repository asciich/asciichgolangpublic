package contextutils

import "context"

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
