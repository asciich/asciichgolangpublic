package installutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/installutils/installoptions"
	"github.com/asciich/asciichgolangpublic/pkg/installutils/nativeinstall"
)

// Install installs binaries from various sources (URL or local path).
// This is a convenience function that delegates to the nativeinstall implementation.
// For remote installation via command executor, use commandexecutorinstall.Install.
func Install(ctx context.Context, options *installoptions.InstallOptions) error {
	return nativeinstall.Install(ctx, options)
}
