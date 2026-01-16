package commandexecutorgitoo

import (
	"context"
	"fmt"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func (g *GitRepository) IsInitialized(ctx context.Context) (isInitialited bool, err error) {
	path, hostDescription, err := g.GetPathAndHostDescription()
	if err != nil {
		return false, err
	}

	commandExecutor, err := g.GetCommandExecutor()
	if err != nil {
		return false, err
	}

	exists, err := g.Exists(ctx)
	if err != nil {
		return false, err
	}

	if !exists {
		logging.LogInfoByCtxf(ctx, "Git repository '%s' does not exist on host '%s' and is therefore not initalized.", path, hostDescription)
		return false, nil
	}

	stdout, err := commandExecutor.RunCommandAndGetStdoutAsString(
		commandexecutorgeneric.WithLiveOutputOnStdoutIfVerbose(ctx),
		&parameteroptions.RunCommandOptions{
			Command: []string{
				"bash",
				"-c",
				fmt.Sprintf(
					"git -C '%s' rev-parse --is-inside-work-tree &> /dev/null && echo yes || echo no",
					path,
				),
			},
		},
	)
	if err != nil {
		return false, err
	}

	stdout = strings.TrimSpace(stdout)

	if stdout == "yes" {
		isInitialited = true
	} else if stdout == "no" {
		isInitialited = false
	} else {
		return false, tracederrors.TracedErrorNilf(
			"Unexpected output='%s' when checking if git repository '%s' is initialized on host '%s'",
			stdout,
			path,
			hostDescription,
		)
	}

	if isInitialited {
		logging.LogInfoByCtxf(ctx, "Git repository '%s' on host '%s' is initialized.", path, hostDescription)
	} else {
		logging.LogInfoByCtxf(ctx, "'%s' is not an initialized git repository on host '%s'.", path, hostDescription)
	}

	return isInitialited, nil
}

func (g *GitRepository) HasInitialCommit(ctx context.Context) (hasInitialCommit bool, err error) {
	path, hostDescription, err := g.GetPathAndHostDescription()
	if err != nil {
		return false, err
	}

	isInitialized, err := g.IsInitialized(ctx)
	if err != nil {
		return false, err
	}

	if !isInitialized {
		logging.LogInfoByCtxf(ctx, "'%s' does not initialized as git repository on host '%s' and can therefore not have an initial commit.", path, hostDescription)
		return false, nil
	}

	commandExecutor, err := g.GetCommandExecutor()
	if err != nil {
		return false, err
	}

	stdout, err := commandExecutor.RunCommandAndGetStdoutAsString(
		commandexecutorgeneric.WithLiveOutputOnStdoutIfVerbose(ctx),
		&parameteroptions.RunCommandOptions{
			Command: []string{
				"bash",
				"-c",
				fmt.Sprintf(
					"git -C '%s' rev-list --max-parents=0 HEAD &> /dev/null && echo yes || echo no",
					path,
				),
			},
		},
	)
	if err != nil {
		return false, err
	}

	stdout = strings.TrimSpace(stdout)

	if stdout == "yes" {
		hasInitialCommit = true
	} else if stdout == "no" {
		hasInitialCommit = false
	} else {
		return false, tracederrors.TracedErrorf("Unexpected stdout='%s'", stdout)
	}

	if hasInitialCommit {
		logging.LogInfoByCtxf(ctx, "Git repository '%s' on host '%s' has an initial commit", path, hostDescription)
	} else {
		logging.LogInfoByCtxf(ctx, "Git repository '%s' on host '%s' has no initial commit", path, hostDescription)
	}

	return hasInitialCommit, nil
}

func (g *GitRepository) Init(ctx context.Context, options *parameteroptions.CreateRepositoryOptions) (err error) {
	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	path, hostDescription, err := g.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	isInitialized, err := g.IsInitialized(ctx)
	if err != nil {
		return err
	}

	if isInitialized {
		logging.LogInfoByCtxf(ctx, "Git repository '%s' on host '%s' is already initialized.", path, hostDescription)
	} else {
		err = g.Create(ctx, &filesoptions.CreateOptions{})
		if err != nil {
			return err
		}

		commandToUse := []string{"init"}

		if options.BareRepository {
			commandToUse = append(commandToUse, "--bare")
		}

		_, err = g.RunGitCommand(ctx, commandToUse)
		if err != nil {
			return err
		}

		if options.BareRepository {
			logging.LogChangedByCtxf(ctx, "Git repository '%s' on host '%s' initialized as bare repository.", path, hostDescription)
		} else {
			logging.LogChangedByCtxf(ctx, "Git repository '%s' on host '%s' initialized as non bare repository.", path, hostDescription)
		}

	}

	if options.InitializeWithDefaultAuthor {
		err = g.SetDefaultAuthor(ctx)
		if err != nil {
			return err
		}
	}

	if options.InitializeWithEmptyCommit {
		hasInitialCommit, err := g.HasInitialCommit(ctx)
		if err != nil {
			return err
		}

		if hasInitialCommit {
			logging.LogInfoByCtxf(ctx, "Repository '%s' on host '%s' has already an initial commit.", path, hostDescription)
		} else {
			if options.BareRepository {
				temporaryClone, err := g.CloneToTemporaryRepository(ctx)
				if err != nil {
					return err
				}
				defer temporaryClone.Delete(ctx, &filesoptions.DeleteOptions{})

				if options.InitializeWithDefaultAuthor {
					err = temporaryClone.SetGitConfig(
						ctx,
						&gitparameteroptions.GitConfigSetOptions{
							Name:  gitgeneric.GitRepositryDefaultAuthorName(),
							Email: gitgeneric.GitRepositryDefaultAuthorEmail(),
						},
					)
					if err != nil {
						return err
					}
				}

				_, err = temporaryClone.CommitAndPush(
					ctx,
					&gitparameteroptions.GitCommitOptions{
						Message:    gitgeneric.GitRepositoryDefaultCommitMessageForInitializeWithEmptyCommit(),
						AllowEmpty: true,
					},
				)
				if err != nil {
					return err
				}
			} else {
				_, err = g.Commit(
					ctx,
					&gitparameteroptions.GitCommitOptions{
						Message:    gitgeneric.GitRepositoryDefaultCommitMessageForInitializeWithEmptyCommit(),
						AllowEmpty: true,
					},
				)
				if err != nil {
					return err
				}
			}

			logging.LogChangedByCtxf(ctx, "Initialized repository '%s' on host '%s' with an empty commit.", path, hostDescription)
		}
	}

	return nil
}
