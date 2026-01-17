package commandexecutorgitoo

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func (c *GitRepository) AddFileByPath(ctx context.Context, pathToAdd string) (err error) {
	if pathToAdd == "" {
		return tracederrors.TracedErrorEmptyString("pathToAdd")
	}

	_, err = c.RunGitCommand(ctx, []string{"add", pathToAdd})
	if err != nil {
		return err
	}

	path, hostDescription, err := c.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(ctx, "Added '%s' to git repository '%s' on host '%s'.", pathToAdd, path, hostDescription)

	return nil
}

func (g *GitRepository) FileByPathExists(ctx context.Context, path string) (exists bool, err error) {
	if path == "" {
		return false, tracederrors.TracedErrorEmptyString(path)
	}

	return g.FileInDirectoryExists(ctx, path)
}
