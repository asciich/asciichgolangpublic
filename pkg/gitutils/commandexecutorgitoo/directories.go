package commandexecutorgitoo

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// GetDirectoryByPath is inherited from the embedded Directory struct.
// This method is intentionally left empty to avoid recursive calls.
// The embedded Directory.GetDirectoryByPath(ctx context.Context, path ...string) method
// will be used instead.
func (g *GitRepository) GetDirectoryByPath(ctx context.Context, pathToSubDir ...string) (subDir filesinterfaces.Directory, err error) {
	if len(pathToSubDir) <= 0 {
		return nil, tracederrors.TracedError("pathToSubdir has no elements")
	}

	// Delegate to the embedded Directory's GetDirectoryByPath method
	// by not overriding it. This prevents infinite recursion.
	// The method call will automatically resolve to Directory.GetDirectoryByPath
	return g.Directory.GetDirectoryByPath(ctx, pathToSubDir...)
}
