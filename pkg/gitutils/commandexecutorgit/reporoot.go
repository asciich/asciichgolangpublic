package commandexecutorgit

import (
	"context"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorbashoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func GetRepositoryRootPathByPath(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, path string) (repoRootPath string, err error) {
	if commandExecutor == nil {
		return "", tracederrors.TracedErrorNil("commandExecutor")
	}

	if path == "" {
		return "", tracederrors.TracedErrorEmptyString("path")
	}

	repoRootPath, err = commandexecutorbashoo.Bash().RunCommandAndGetStdoutAsString(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{"git", "-C", path, "rev-parse", "--show-toplevel"},
		},
	)
	if err != nil {
		return "", err
	}

	repoRootPath = strings.TrimSpace(repoRootPath)

	repoRootDir, err := files.GetLocalDirectoryByPath(ctx, repoRootPath)
	if err != nil {
		return "", err
	}

	exists, err := repoRootDir.Exists(ctx)
	if err != nil {
		return "", err
	}

	if !exists {
		return "", tracederrors.TracedErrorf(
			"internal error: repoRootDir '%s' points to an non existent path after evaluation",
			repoRootPath,
		)
	}

	logging.LogInfoByCtxf(ctx, "Found git repository root directory '%s' for local path '%s'.", repoRootPath, path)

	return repoRootPath, nil
}
