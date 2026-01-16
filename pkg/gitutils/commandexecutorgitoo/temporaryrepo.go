package commandexecutorgitoo

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutortempfile"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func CreateLocalTemporaryRepository(ctx context.Context, options *parameteroptions.CreateRepositoryOptions) (gitinterfaces.GitRepository, error) {
	return CreateTemporaryRepository(ctx, commandexecutorexecoo.Exec(), options)
}

func CreateTemporaryRepository(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, options *parameteroptions.CreateRepositoryOptions) (gitinterfaces.GitRepository, error) {
	path, err := commandexecutortempfile.CreateEmptyTemporaryDirectory(ctx, commandExecutor)
	if err != nil {
		return nil, err
	}

	if options == nil {
		options = new(parameteroptions.CreateRepositoryOptions)
	}

	repo, err := New(commandExecutor, path)
	if err != nil {
		return nil, err
	}

	if options.InitializeWithEmptyCommit {
		err := repo.Init(ctx, options)
		if err != nil {
			return nil, err
		}
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return nil, err
	}

	logging.LogChangedByCtxf(ctx, "Created temporary git repository '%s' on '%s'.", path, hostDescription)

	return repo, nil
}

func CloneToTemporaryRepository(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, urlOrPathToClone string) (gitinterfaces.GitRepository, error) {
	if commandExecutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExecutor")
	}

	if urlOrPathToClone == "" {
		return nil, tracederrors.TracedErrorEmptyString("urlOrPathToClone")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Clone git repository '%s' to temporary repository on '%s' started.", urlOrPathToClone, hostDescription)

	repo, err := CreateTemporaryRepository(ctx, commandExecutor, &parameteroptions.CreateRepositoryOptions{})
	if err != nil {
		return nil, err
	}

	err = repo.CloneRepositoryByPathOrUrl(ctx, urlOrPathToClone)
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Clone git repository '%s' to temporary repository on '%s' finished.", urlOrPathToClone, hostDescription)

	return repo, err
}

func (g *GitRepository) CloneToTemporaryRepository(ctx context.Context) (gitinterfaces.GitRepository, error) {
	path, hostDescription, err := g.GetPathAndHostDescription()
	if err != nil {
		return nil, err
	}

	commandExecutor, err := g.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Cloning repository '%s' to temporary repository on '%s' started.", path, hostDescription)

	repo, err := CloneToTemporaryRepository(ctx, commandExecutor, path)
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Cloning repository '%s' to temporary repository on '%s' finished.", path, hostDescription)

	return repo, nil
}
