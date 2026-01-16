package commandexecutorgitoo

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func (g *GitRepository) CloneRepository(ctx context.Context, repository gitinterfaces.GitRepository) (err error) {
	if repository == nil {
		return tracederrors.TracedErrorNil("repository")
	}

	repoHostDescription, err := repository.GetHostDescription()
	if err != nil {
		return err
	}

	hostDescription, err := g.GetHostDescription()
	if err != nil {
		return err
	}

	if hostDescription != repoHostDescription {
		return tracederrors.TracedErrorf(
			"Only implemented for two repositories on the same host. But repository from host '%s' should be cloned to host '%s'",
			repoHostDescription,
			hostDescription,
		)
	}

	pathToClone, err := repository.GetPath()
	if err != nil {
		return err
	}

	return g.CloneRepositoryByPathOrUrl(ctx, pathToClone)
}

func (g *GitRepository) CloneRepositoryByPathOrUrl(ctx context.Context, pathOrUrlToClone string) (err error) {
	if pathOrUrlToClone == "" {
		return tracederrors.TracedErrorEmptyString("pathToClone")
	}

	path, hostDescription, err := g.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Cloning git repository '%s' to '%s' on '%s' started.", pathOrUrlToClone, path, hostDescription)

	isInitialized, err := g.IsInitialized(ctx)
	if err != nil {
		return err
	}

	if isInitialized {
		logging.LogInfof(
			"'%s' is already an initialized git repository on host '%s'. Skip clone.",
			path,
			hostDescription,
		)
	} else {
		commandExecutor, err := g.GetCommandExecutor()
		if err != nil {
			return err
		}

		_, err = commandExecutor.RunCommand(
			commandexecutorgeneric.WithLiveOutputOnStdoutIfVerbose(ctx),
			&parameteroptions.RunCommandOptions{
				Command: []string{"git", "clone", pathOrUrlToClone, path},
			},
		)
		if err != nil {
			return err
		}
	}

	logging.LogInfoByCtxf(ctx, "Cloning git repository '%s' to '%s' on host '%s' finished.", pathOrUrlToClone, path, hostDescription)

	return nil
}
