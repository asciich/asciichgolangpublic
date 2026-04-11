package contextutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func CheckContextStillAlive(ctx context.Context) error {
	if ctx == nil {
		return tracederrors.TracedErrorNil("ctx")
	}

	err := ctx.Err()
	if err != nil {
		return tracederrors.TracedErrorf("Context is not alive anymore: %w", err)
	}

	return nil
}
