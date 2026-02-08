package headscaleutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}