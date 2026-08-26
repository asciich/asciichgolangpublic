package nativegitoo

import (
	"context"

	"github.com/go-git/go-git/v5"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func (n *NativeGitRepository) Init(ctx context.Context, options *parameteroptions.CreateRepositoryOptions) (err error) {
	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	path, hostDescription, err := n.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	isInitialized, err := n.IsInitialized(ctx)
	if err != nil {
		return err
	}

	if isInitialized {
		logging.LogInfoByCtxf(ctx, "Git repository '%s' on host '%s' is already initialized.", path, hostDescription)
	} else {
		err = n.Create(ctx, &filesoptions.CreateOptions{})
		if err != nil {
			return err
		}

		_, err = git.PlainInit(path, options.BareRepository)
		if err != nil {
			return tracederrors.TracedErrorf("Init git repository '%s' on '%s' failed: %w", path, hostDescription, err)
		}

		if options.BareRepository {
			logging.LogChangedByCtxf(ctx, "Git repository '%s' on host '%s' initialized as bare repository.", path, hostDescription)
		} else {
			logging.LogChangedByCtxf(ctx, "Git repository '%s' on host '%s' initialized as non bare repository.", path, hostDescription)
		}
	}

	if options.InitializeWithDefaultAuthor {
		err = n.SetGitConfig(
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

	if options.InitializeWithEmptyCommit {
		hasInitialCommit, err := n.HasInitialCommit(ctx)
		if err != nil {
			return err
		}

		if hasInitialCommit {
			logging.LogInfoByCtxf(ctx, "Repository '%s' on host '%s' has already an initial commit.", path, hostDescription)
		} else {
			if options.BareRepository {
				temporaryClone, err := n.CloneToTemporaryRepository(ctx)
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
				_, err = n.Commit(
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
