package contextutils

import (
	"context"

	"github.com/spf13/cobra"
)

// Returns a child context of ctx with verbosity enabled according to `verbose`.
// If ctx is nil a new context with verbosity set is returned.
func WithVerbosityContextByBool(ctx context.Context, verbose bool) (ctxWithVerbosity context.Context) {
	if ctx == nil {
		return GetVerbosityContextByBool(verbose)
	}

	return context.WithValue(ctx, "verbose", verbose)
}

func GetVerbosityContextByCobraCmd(cmd *cobra.Command) (ctx context.Context) {
	if cmd == nil {
		return ContextSilent()
	}

	ctx = cmd.Context()

	if cmd.Flags().Lookup("verbose") == nil {
		return WithVerbosityContextByBool(ctx, false)
	}

	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		return WithVerbosityContextByBool(ctx, false)
	}

	return WithVerbosityContextByBool(ctx, verbose)
}

func GetVerbosityContextByBool(verbose bool) (ctx context.Context) {
	if verbose {
		return ContextVerbose()
	}

	return ContextSilent()
}

// context.Background() with verbose output enabled.
func ContextVerbose() (ctx context.Context) {
	return context.WithValue(context.Background(), "verbose", true)
}

// context.Background() with verbose output explicitly disabled.
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

// Returns a child context with verbosity enabled.
func WithVerbose(ctx context.Context) context.Context {
	if ctx == nil {
		return nil
	}

	return WithVerbosityContextByBool(ctx, true)
}

// Returns a child context with verbosity disabled.
func WithSilent(ctx context.Context) context.Context {
	if ctx == nil {
		return nil
	}

	return WithVerbosityContextByBool(ctx, false)
}
