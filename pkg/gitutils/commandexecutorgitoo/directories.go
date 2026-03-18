package commandexecutorgitoo

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func (g *GitRepository) GetDirectoryByPath(ctx context.Context, pathToSubDir ...string) (subDir filesinterfaces.Directory, err error) {
	if len(pathToSubDir) <= 0 {
		return nil, tracederrors.TracedError("pathToSubdir has no elements")
	}

	return g.GetSubDirectory(ctx, pathToSubDir...)
}
