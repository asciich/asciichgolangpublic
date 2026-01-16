package commandexecutorgitoo

import (
	"context"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/slicesutils"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfileoo"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type GitRepository struct {
	commandexecutorfileoo.Directory
	gitgeneric.GitRepositoryBase
}

func NewFromDirectory(dir filesinterfaces.Directory) (*GitRepository, error) {
	if dir == nil {
		return nil, tracederrors.TracedErrorNil("dir")
	}

	path, err := dir.GetPath()
	if err != nil {
		return nil, err
	}

	cedir, ok := dir.(*commandexecutorfileoo.Directory)
	if ok {
		commandExecutor, err := cedir.GetCommandExecutor()
		if err != nil {
			return nil, err
		}

		return New(commandExecutor, path)
	}

	hostDescription, err := dir.GetHostDescription()
	if err != nil {
		return nil, err
	}

	if hostDescription == "localhost" {
		return NewLocalGitRepository(path)
	}

	typeName, err := datatypes.GetTypeName(dir)
	if err != nil {
		return nil, err
	}

	return nil, tracederrors.TracedErrorf("Unable to get new directory from: dir of type '%s' on '%s'", typeName, hostDescription)
}

func NewLocalGitRepository(path string) (*GitRepository, error) {
	if path == "" {
		return nil, tracederrors.TracedErrorEmptyString("path")
	}

	return New(commandexecutorexecoo.Exec(), path)
}

func New(commandExecutor commandexecutorinterfaces.CommandExecutor, path string) (*GitRepository, error) {
	repo := new(GitRepository)

	err := repo.SetParentDirectoryForBaseClass(repo)
	if err != nil {
		return nil, err
	}

	err = repo.SetParentRepositoryForBaseClass(repo)
	if err != nil {
		return nil, err
	}

	err = repo.SetPath(path)
	if err != nil {
		return nil, err
	}

	err = repo.SetCommandExecutor(commandExecutor)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (g *GitRepository) GetRootDirectory(ctx context.Context) (rootDirectory filesinterfaces.Directory, err error) {
	commandExecutor, err := g.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	rootDirPath, err := g.GetRootDirectoryPath(ctx)
	if err != nil {
		return nil, err
	}

	rootDirectory, err = files.GetCommandExecutorDirectoryByPath(
		commandExecutor,
		rootDirPath,
	)
	if err != nil {
		return nil, err
	}

	return rootDirectory, nil
}

func (g *GitRepository) GetRootDirectoryPath(ctx context.Context) (rootDirectoryPath string, err error) {
	path, hostDescription, err := g.GetPathAndHostDescription()
	if err != nil {
		return "", err
	}

	isBareRepository, err := g.IsBareRepository(ctx)
	if err != nil {
		return "", err
	}

	if isBareRepository {
		var cwd filesinterfaces.Directory

		commandExecutor, err := g.GetCommandExecutor()
		if err != nil {
			return "", err
		}

		cwd, err = files.GetCommandExecutorDirectoryByPath(
			commandExecutor,
			path,
		)
		if err != nil {
			return "", err
		}

		for {
			filePaths, err := cwd.ListFilePaths(
				ctx,
				&parameteroptions.ListFileOptions{
					NonRecursive:        true,
					ReturnRelativePaths: true,
				},
			)
			if err != nil {
				return "", err
			}

			if slicesutils.ContainsAllStrings(filePaths, []string{"config", "HEAD"}) {
				rootDirectoryPath, err = cwd.GetPath()
				if err != nil {
					return "", err
				}
			}

			if rootDirectoryPath != "" {
				break
			}

			cwd, err = cwd.GetParentDirectory()
			if err != nil {
				return "", err
			}
		}
	} else {
		stdout, err := g.RunGitCommandAndGetStdoutAsString(
			ctx,
			[]string{"rev-parse", "--show-toplevel"},
		)
		if err != nil {
			return "", err
		}

		rootDirectoryPath = strings.TrimSpace(stdout)
	}

	if rootDirectoryPath == "" {
		return "", tracederrors.TracedErrorf(
			"rootDirectoryPath is empty string after evaluating root directory of git repository '%s' on host '%s'",
			path,
			hostDescription,
		)
	}

	logging.LogInfoByCtxf(ctx, "Git repo root directory is '%s' on host '%s'.", rootDirectoryPath, hostDescription)

	return rootDirectoryPath, nil
}

func (g *GitRepository) IsBareRepository(ctx context.Context) (isBare bool, err error) {
	path, hostDescription, err := g.GetPathAndHostDescription()
	if err != nil {
		return false, err
	}

	stdout, err := g.RunGitCommandAndGetStdoutAsString(ctx, []string{"rev-parse", "--is-bare-repository"})
	if err != nil {
		return false, err
	}

	stdout = strings.TrimSpace(stdout)

	if stdout == "false" {
		isBare = false
	} else if stdout == "true" {
		isBare = true
	} else {
		return false, tracederrors.TracedErrorf(
			"Unknown is-bare-repositoy output '%s'",
			stdout,
		)
	}

	if isBare {
		logging.LogInfoByCtxf(ctx, "Git repository '%s' on host '%s' is a bare repository.", path, hostDescription)
	} else {
		logging.LogInfoByCtxf(ctx, "Git repository '%s' on host '%s' is not a bare repository.", path, hostDescription)
	}

	return isBare, nil
}

func (g *GitRepository) SetDefaultAuthor(ctx context.Context) (err error) {
	err = g.SetUserName(ctx, gitgeneric.GitRepositryDefaultAuthorName())
	if err != nil {
		return err
	}

	err = g.SetUserEmail(ctx, gitgeneric.GitRepositryDefaultAuthorEmail())
	if err != nil {
		return err
	}

	path, hostDescription, err := g.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(ctx, "Initialized git repository '%s' on '%s' with default author and email.", path, hostDescription)

	return nil
}

func (g *GitRepository) SetUserEmail(ctx context.Context, email string) (err error) {
	if email == "" {
		return tracederrors.TracedErrorEmptyString("email")
	}

	_, err = g.RunGitCommand(ctx, []string{"config", "user.email", email})
	if err != nil {
		return err
	}

	path, hostDescription, err := g.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Set git user email to '%s' for git repository '%s' on host '%s'.", email, path, hostDescription)

	return nil
}

func (g *GitRepository) SetUserName(ctx context.Context, name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	_, err = g.RunGitCommand(ctx, []string{"config", "user.name", name})
	if err != nil {
		return err
	}

	path, hostDescription, err := g.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Set git user name to '%s' for git repository '%s' on host '%s'.", name, path,
		hostDescription,
	)

	return nil
}

func (g *GitRepository) IsGitRepository(ctx context.Context) (isRepository bool, err error) {
	isInitalized, err := g.IsInitialized(ctx)
	if err != nil {
		return false, err
	}

	isRepository = isInitalized

	path, hostDescription, err := g.GetPathAndHostDescription()
	if err != nil {
		return false, err
	}

	if isRepository {
		logging.LogInfoByCtxf(ctx, "'%s' on host '%s' is a git repository", path, hostDescription)
	} else {
		logging.LogInfoByCtxf(ctx, "'%s' on host '%s' is not a git repository", path, hostDescription)
	}

	return isRepository, nil
}

func (g *GitRepository) Pull(ctx context.Context) (err error) {
	path, hostDescription, err := g.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Pull git repository '%s' on '%s' started.", path, hostDescription)

	_, err = g.RunGitCommand(
		ctx,
		[]string{"pull"},
	)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Pull git repository '%s' on '%s' finished.", path, hostDescription)

	return
}

func (g *GitRepository) SetGitConfig(ctx context.Context, options *gitparameteroptions.GitConfigSetOptions) (err error) {
	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	if options.Email != "" {
		err = g.SetUserEmail(ctx, options.Email)
		if err != nil {
			return err
		}
	}

	if options.Name != "" {
		err = g.SetUserName(ctx, options.Name)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *GitRepository) SetRemoteUrl(ctx context.Context, remoteUrl string) (err error) {
	remoteUrl = strings.TrimSpace(remoteUrl)
	if len(remoteUrl) <= 0 {
		return tracederrors.TracedError("remoteUrl is empty string")
	}

	name := "origin"

	_, err = g.RunGitCommand(ctx, []string{"remote", "set-url", name, remoteUrl})
	if err != nil {
		return err
	}

	path, hostDescription, err := g.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(ctx, "Set remote Url for '%v' in git repository '%v' on host '%s' to '%v'.", name, path, hostDescription, remoteUrl)

	return nil
}
