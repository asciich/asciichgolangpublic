package asciichgolangpublic

import (
	"fmt"
	"strings"
	"time"
)

// This is the GitRepository implementation based on a CommandExecutor (e.g. Bash, SSH...).
// This means it's a wrapper around the "git" binary which needs to be available.
// While very inefficient this solution can manage git repository on remote hosts, inside containers...
// which makes it very flexible.
//
// When dealing with locally available repositories it's recommended to use the LocalGitRepository
// implementation which uses go build in git functionality instead.
type CommandExecutorGitRepository struct {
	CommandExecutorDirectory
	GitRepositoryBase
}

func GetCommandExecutorGitRepositoryFromDirectory(directory Directory) (c *CommandExecutorGitRepository, err error) {
	if directory == nil {
		return nil, TracedErrorNil("directory")
	}

	commandExecutoryDirectory, ok := directory.(*CommandExecutorDirectory)
	if ok {
		commandExecutor, path, err := commandExecutoryDirectory.GetCommandExecutorAndDirPath()
		if err != nil {
			return nil, err
		}

		return GetCommandExecutorGitRepositoryByPath(commandExecutor, path)
	}

	localDirectory, ok := directory.(*LocalDirectory)
	if ok {
		path, err := localDirectory.GetPath()
		if err != nil {
			return nil, err
		}

		return GetLocalCommandExecutorGitRepositoryByPath(path)
	}

	unknownTypeName, err := Types().GetTypeName(directory)
	if err != nil {
		return nil, err
	}

	return nil, TracedErrorf(
		"Not implemented for directory type = '%s'",
		unknownTypeName,
	)
}

func MustGetCommandExecutorGitRepositoryFromDirectory(directory Directory) (c *CommandExecutorGitRepository) {
	c, err := GetCommandExecutorGitRepositoryFromDirectory(directory)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return c
}

func MustNewCommandExecutorGitRepository(commandExecutor CommandExecutor) (c *CommandExecutorGitRepository) {
	c, err := NewCommandExecutorGitRepository(commandExecutor)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return c
}

func NewCommandExecutorGitRepository(commandExecutor CommandExecutor) (c *CommandExecutorGitRepository, err error) {
	if commandExecutor == nil {
		return nil, TracedErrorNil("commandExecutor")
	}

	c = new(CommandExecutorGitRepository)

	err = c.CommandExecutorDirectory.SetCommandExecutor(commandExecutor)
	if err != nil {
		return nil, err
	}

	err = c.SetParentDirectoryForBaseClass(c)
	if err != nil {
		return nil, err
	}

	err = c.SetParentRepositoryForBaseClass(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// This function was only added to fulfil the current interface.
// On the long run this method has to be removed.
func (c *CommandExecutorGitRepository) GetAsLocalDirectory() (l *LocalDirectory, err error) {
	return nil, TracedErrorNotImplemented()
}

// This function was only added to fulfil the current interface.
// On the long run this method has to be removed.
func (c *CommandExecutorGitRepository) GetAsLocalGitRepository() (l *LocalGitRepository, err error) {
	return nil, TracedErrorNotImplemented()
}

func (c *CommandExecutorGitRepository) GetTagByHash(hash string) (tag *GitTag, err error) {
	if hash == "" {
		return nil, TracedErrorEmptyString("hash")
	}

	return &GitTag{
		gitRepository: c,
		hash:          hash,
	}, nil
}

func (c *CommandExecutorGitRepository) AddFileByPath(pathToAdd string, verbose bool) (err error) {
	if pathToAdd == "" {
		return TracedErrorEmptyString("pathToAdd")
	}

	_, err = c.RunGitCommand(
		[]string{"add", pathToAdd},
		verbose,
	)
	if err != nil {
		return err
	}

	if verbose {
		path, hostDescription, err := c.GetPathAndHostDescription()
		if err != nil {
			return err
		}

		LogChangedf(
			"Added '%s' to git repository '%s' on host '%s'.",
			pathToAdd,
			path,
			hostDescription,
		)
	}

	return nil
}

func (c *CommandExecutorGitRepository) CloneRepository(repository GitRepository, verbose bool) (err error) {
	if repository == nil {
		return TracedErrorNil("repository")
	}

	repoHostDescription, err := repository.GetHostDescription()
	if err != nil {
		return err
	}

	hostDescription, err := c.GetHostDescription()
	if err != nil {
		return err
	}

	if hostDescription != repoHostDescription {
		return TracedErrorf(
			"Only implemented for two repositories on the same host. But repository from host '%s' should be cloned to host '%s'",
			repoHostDescription,
			hostDescription,
		)
	}

	pathToClone, err := repository.GetPath()
	if err != nil {
		return err
	}

	return c.CloneRepositoryByPathOrUrl(pathToClone, verbose)
}

func (c *CommandExecutorGitRepository) CloneRepositoryByPathOrUrl(pathOrUrlToClone string, verbose bool) (err error) {
	if pathOrUrlToClone == "" {
		return TracedErrorEmptyString("pathToClone")
	}

	path, hostDescription, err := c.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	if verbose {
		LogInfof(
			"Cloning git repository '%s' to '%s' on '%s' started.",
			pathOrUrlToClone,
			path,
			hostDescription,
		)
	}

	isInitialized, err := c.IsInitialized(verbose)
	if err != nil {
		return err
	}

	if isInitialized {
		return TracedErrorf(
			"'%s' on host '%s' is already an initialized git repository. Clone of '%s' aborted.",
			path,
			hostDescription,
			pathOrUrlToClone,
		)
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return err
	}

	_, err = commandExecutor.RunCommand(
		&RunCommandOptions{
			Command:            []string{"git", "clone", pathOrUrlToClone, path},
			Verbose:            verbose,
			LiveOutputOnStdout: verbose,
		},
	)
	if err != nil {
		return err
	}

	if verbose {
		LogInfof(
			"Cloning git repository '%s' to '%s' on host '%s' finished.",
			pathOrUrlToClone,
			path,
			hostDescription,
		)
	}

	return nil
}

func (c *CommandExecutorGitRepository) Commit(commitOptions *GitCommitOptions) (createdCommit *GitCommit, err error) {
	if commitOptions == nil {
		return nil, TracedErrorNil("commitOptions")
	}

	commitCommand := []string{"commit"}

	if commitOptions.AllowEmpty {
		commitCommand = append(commitCommand, "--allow-empty")
	}

	message, err := commitOptions.GetMessage()
	if err != nil {
		return nil, err
	}

	commitCommand = append(commitCommand, "-m", message)

	_, err = c.RunGitCommand(
		commitCommand,
		commitOptions.Verbose,
	)
	if err != nil {
		return nil, err
	}

	createdCommit, err = c.GetCurrentCommit()
	if err != nil {
		return nil, err
	}

	if commitOptions.Verbose {
		path, hostDescription, err := c.GetPathAndHostDescription()
		if err != nil {
			return nil, err
		}

		createdHash, err := createdCommit.GetHash()
		if err != nil {
			return nil, err
		}

		LogChangedf(
			"Created commit '%s' in git repository '%s' on host '%s'.",
			createdHash,
			path,
			hostDescription,
		)
	}

	return createdCommit, nil
}

func (c *CommandExecutorGitRepository) CommitHasParentCommitByCommitHash(hash string) (hasParentCommit bool, err error) {
	return false, TracedErrorNotImplemented()
}

func (c *CommandExecutorGitRepository) GetAuthorEmailByCommitHash(hash string) (authorEmail string, err error) {
	return "", TracedErrorNotImplemented()
}

func (c *CommandExecutorGitRepository) GetAuthorStringByCommitHash(hash string) (authorEmail string, err error) {
	return "", TracedErrorNotImplemented()
}

func (c *CommandExecutorGitRepository) GetCommitAgeDurationByCommitHash(hash string) (ageDuration *time.Duration, err error) {
	return nil, TracedErrorNotImplemented()
}

func (c *CommandExecutorGitRepository) GetCommitAgeSecondsByCommitHash(hash string) (ageSeconds float64, err error) {
	return -1, TracedErrorNotImplemented()
}

func (c *CommandExecutorGitRepository) GetCommitByHash(hash string) (gitCommit *GitCommit, err error) {
	if hash == "" {
		return nil, TracedErrorEmptyString("hash")
	}

	gitCommit = NewGitCommit()

	err = gitCommit.SetGitRepo(c)
	if err != nil {
		return nil, err
	}

	err = gitCommit.SetHash(hash)
	if err != nil {
		return nil, err
	}

	return gitCommit, nil
}

func (c *CommandExecutorGitRepository) GetCommitMessageByCommitHash(hash string) (commitMessage string, err error) {
	return "", TracedErrorNotImplemented()
}

func (c *CommandExecutorGitRepository) GetCommitParentsByCommitHash(hash string, options *GitCommitGetParentsOptions) (commitParents []*GitCommit, err error) {
	return nil, TracedErrorNotImplemented()
}

func (c *CommandExecutorGitRepository) GetCommitTimeByCommitHash(hash string) (commitTime *time.Time, err error) {
	return nil, TracedErrorNotImplemented()
}

func (c *CommandExecutorGitRepository) GetCurrentCommit() (currentCommit *GitCommit, err error) {
	currentCommitHash, err := c.GetCurrentCommitHash()
	if err != nil {
		return nil, err
	}

	return c.GetCommitByHash(currentCommitHash)
}

func (c *CommandExecutorGitRepository) GetCurrentCommitHash() (currentCommitHash string, err error) {
	currentCommitHash, err = c.RunGitCommandAndGetStdoutAsString(
		[]string{"rev-parse", "HEAD"},
		false,
	)
	if err != nil {
		return "", err
	}

	currentCommitHash = strings.TrimSpace(currentCommitHash)

	return currentCommitHash, nil
}

func (c *CommandExecutorGitRepository) GetGitStatusOutput(verbose bool) (output string, err error) {
	return "", TracedErrorNotImplemented()
}

func (c *CommandExecutorGitRepository) GetRootDirectoryPath(verbose bool) (rootDirectoryPath string, err error) {
	path, hostDescription, err := c.GetPathAndHostDescription()
	if err != nil {
		return "", err
	}

	isBareRepository, err := c.IsBareRepository(verbose)
	if err != nil {
		return "", err
	}

	if isBareRepository {
		var cwd Directory

		cwd, err = GetCommandExecutorDirectoryByPath(
			c.commandExecutor,
			c.dirPath,
		)
		if err != nil {
			return "", err
		}

		for {
			filePaths, err := cwd.ListFilePaths(
				&ListFileOptions{
					NonRecursive:        true,
					Verbose:             verbose,
					ReturnRelativePaths: true,
				},
			)
			if err != nil {
				return "", err
			}

			if Slices().ContainsAllStrings(filePaths, []string{"config", "HEAD"}) {
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
		stdout, err := c.RunGitCommandAndGetStdoutAsString(
			[]string{"rev-parse", "--show-toplevel"},
			verbose,
		)
		if err != nil {
			return "", err
		}

		rootDirectoryPath = strings.TrimSpace(stdout)
	}

	if rootDirectoryPath == "" {
		return "", TracedErrorf(
			"rootDirectoryPath is empty string after evaluating root directory of git repository '%s' on host '%s'",
			path,
			hostDescription,
		)
	}

	if verbose {
		LogInfof(
			"Git repo root directory is '%s' on host '%s'.",
			rootDirectoryPath,
			hostDescription,
		)
	}

	return rootDirectoryPath, nil
}

func (c *CommandExecutorGitRepository) HasInitialCommit(verbose bool) (hasInitialCommit bool, err error) {
	path, hostDescription, err := c.GetPathAndHostDescription()
	if err != nil {
		return false, err
	}

	isInitialized, err := c.IsInitialized(verbose)
	if err != nil {
		return false, err
	}

	if !isInitialized {
		if verbose {
			LogInfof(
				"'%s' does not initialized as git repository on host '%s' and can therefore not have an initial commit.",
				path,
				hostDescription,
			)
		}
		return false, nil
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return false, err
	}

	stdout, err := commandExecutor.RunCommandAndGetStdoutAsString(
		&RunCommandOptions{
			Command: []string{
				"bash",
				"-c",
				fmt.Sprintf(
					"git -C '%s' rev-list --max-parents=0 HEAD &> /dev/null && echo yes || echo no",
					path,
				),
			},
			Verbose:            verbose,
			LiveOutputOnStdout: verbose,
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
		return false, TracedErrorf("Unexpected stdout='%s'", stdout)
	}

	if verbose {
		if hasInitialCommit {
			LogInfof(
				"Git repository '%s' on host '%s' has an initial commit",
				path,
				hostDescription,
			)
		} else {
			LogInfof(
				"Git repository '%s' on host '%s' has no initial commit",
				path,
				hostDescription,
			)
		}
	}

	return hasInitialCommit, nil
}

func (c *CommandExecutorGitRepository) HasUncommittedChanges(verbose bool) (hasUncommitedChanges bool, err error) {
	path, hostDescription, err := c.GetPathAndHostDescription()
	if err != nil {
		return false, err
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return false, err
	}

	commandOutput, err := commandExecutor.RunCommand(
		&RunCommandOptions{
			Command: []string{
				"bash",
				"-c",
				fmt.Sprintf(
					"cd '%s' && git diff && git diff --cached && git status --porcelain",
					path,
				),
			},
			Verbose: verbose,
		},
	)
	if err != nil {
		return false, err
	}

	isEmpty, err := commandOutput.IsStdoutAndStderrEmpty()
	if err != nil {
		return false, err
	}

	if !isEmpty {
		hasUncommitedChanges = true
	}

	if verbose {
		if hasUncommitedChanges {
			LogInfof(
				"Git repository '%s' on '%s' has uncommited changes.",
				path,
				hostDescription,
			)
		} else {
			LogInfof(
				"Git repository '%s' on '%s' has no uncommited changes.",
				path,
				hostDescription,
			)
		}
	}

	return hasUncommitedChanges, nil
}

func (c *CommandExecutorGitRepository) Init(options *CreateRepositoryOptions) (err error) {
	if options == nil {
		return TracedErrorNil("options")
	}

	path, hostDescription, err := c.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	isInitialized, err := c.IsInitialized(options.Verbose)
	if err != nil {
		return err
	}

	if isInitialized {
		if options.Verbose {
			LogInfof(
				"Git repository '%s' on host '%s' is already initialized.",
				path,
				hostDescription,
			)
		}
	} else {
		err = c.Create(options.Verbose)
		if err != nil {
			return err
		}

		commandToUse := []string{"init"}

		if options.BareRepository {
			commandToUse = append(commandToUse, "--bare")
		}

		_, err = c.RunGitCommand(
			commandToUse,
			options.Verbose,
		)
		if err != nil {
			return err
		}

		if options.Verbose {
			if options.BareRepository {
				LogChangedf(
					"Git repository '%s' on host '%s' initialized as bare repository.",
					path,
					hostDescription,
				)
			} else {
				LogChangedf(
					"Git repository '%s' on host '%s' initialized as non bare repository.",
					path,
					hostDescription,
				)
			}
		}
	}

	if options.InitializeWithDefaultAuthor {
		err = c.SetDefaultAuthor(options.Verbose)
		if err != nil {
			return err
		}
	}

	if options.InitializeWithEmptyCommit {
		hasInitialCommit, err := c.HasInitialCommit(options.Verbose)
		if err != nil {
			return err
		}

		if hasInitialCommit {
			LogInfof(
				"Repository '%s' on host '%s' has already an initial commit.",
				path,
				hostDescription,
			)
		} else {
			if options.BareRepository {
				temporaryClone, err := GitRepositories().CloneGitRepositoryToTemporaryDirectory(c, options.Verbose)
				if err != nil {
					return err
				}
				defer temporaryClone.Delete(options.Verbose)

				if options.InitializeWithDefaultAuthor {
					temporaryClone.SetGitConfig(
						&GitConfigSetOptions{
							Name:    GitRepositryDefualtAuthorName(),
							Email:   GitRepositryDefualtAuthorEmail(),
							Verbose: options.Verbose,
						},
					)
				}

				_, err = temporaryClone.CommitAndPush(
					&GitCommitOptions{
						Message:    GitRepositoryDefaultCommitMessageForInitializeWithEmptyCommit(),
						AllowEmpty: true,
						Verbose:    true,
					},
				)
				if err != nil {
					return err
				}
			} else {
				_, err = c.Commit(
					&GitCommitOptions{
						Message:    GitRepositoryDefaultCommitMessageForInitializeWithEmptyCommit(),
						AllowEmpty: true,
						Verbose:    true,
					},
				)
				if err != nil {
					return err
				}
			}

			if options.Verbose {
				LogChangedf(
					"Initialized repository '%s' on host '%s' with an empty commit.",
					path,
					hostDescription,
				)
			}
		}
	}

	return nil
}

func (c *CommandExecutorGitRepository) IsBareRepository(verbose bool) (isBare bool, err error) {
	path, hostDescription, err := c.GetPathAndHostDescription()
	if err != nil {
		return false, err
	}

	stdout, err := c.RunGitCommandAndGetStdoutAsString(
		[]string{"rev-parse", "--is-bare-repository"},
		verbose,
	)
	if err != nil {
		return false, err
	}

	stdout = strings.TrimSpace(stdout)

	if stdout == "false" {
		isBare = false
	} else if stdout == "true" {
		isBare = true
	} else {
		return false, TracedErrorf(
			"Unknown is-bare-repositoy output '%s'",
			stdout,
		)
	}

	if verbose {
		if isBare {
			LogInfof(
				"Git repository '%s' on host '%s' is a bare repository.",
				path,
				hostDescription,
			)
		} else {
			LogInfof(
				"Git repository '%s' on host '%s' is not a bare repository.",
				path,
				hostDescription,
			)
		}
	}

	return isBare, nil
}

func (c *CommandExecutorGitRepository) IsGitRepository(verbose bool) (isRepository bool, err error) {
	isInitalized, err := c.IsInitialized(verbose)
	if err != nil {
		return false, err
	}

	isRepository = isInitalized

	if verbose {
		path, hostDescription, err := c.GetPathAndHostDescription()
		if err != nil {
			return false, err
		}

		if isRepository {
			LogInfof(
				"'%s' on host '%s' is a git repository",
				path,
				hostDescription,
			)
		} else {
			LogInfof(
				"'%s' on host '%s' is not a git repository",
				path,
				hostDescription,
			)
		}
	}

	return isRepository, nil
}

func (c *CommandExecutorGitRepository) IsInitialized(verbose bool) (isInitialited bool, err error) {
	path, hostDescription, err := c.GetPathAndHostDescription()
	if err != nil {
		return false, err
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return false, err
	}

	exists, err := c.Exists(verbose)
	if err != nil {
		return false, err
	}

	if !exists {
		if verbose {
			LogInfof(
				"Git repository '%s' does not exist on host '%s' and is therefore not initalized.",
				path,
				hostDescription,
			)
		}
		return false, nil
	}

	stdout, err := commandExecutor.RunCommandAndGetStdoutAsString(
		&RunCommandOptions{
			Command: []string{
				"bash",
				"-c",
				fmt.Sprintf(
					"git -C '%s' rev-parse --is-inside-work-tree &> /dev/null && echo yes || echo no",
					path,
				),
			},
			Verbose:            verbose,
			LiveOutputOnStdout: verbose,
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
		return false, TracedErrorNilf(
			"Unexpected output='%s' when checking if git repository '%s' is initialized on host '%s'",
			stdout,
			path,
			hostDescription,
		)
	}

	if verbose {
		if isInitialited {
			LogInfof(
				"Git repository '%s' on host '%s' is initialized.",
				path,
				hostDescription,
			)
		} else {
			LogInfof(
				"'%s' is not an initialized git repository on host '%s'.",
				path,
				hostDescription,
			)
		}
	}

	return isInitialited, nil
}

func (c *CommandExecutorGitRepository) MustAddFileByPath(pathToAdd string, verbose bool) {
	err := c.AddFileByPath(pathToAdd, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorGitRepository) MustCloneRepository(repository GitRepository, verbose bool) {
	err := c.CloneRepository(repository, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorGitRepository) MustCloneRepositoryByPathOrUrl(pathOrUrlToClone string, verbose bool) {
	err := c.CloneRepositoryByPathOrUrl(pathOrUrlToClone, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorGitRepository) MustCommit(commitOptions *GitCommitOptions) (createdCommit *GitCommit) {
	createdCommit, err := c.Commit(commitOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return createdCommit
}

func (c *CommandExecutorGitRepository) MustCommitHasParentCommitByCommitHash(hash string) (hasParentCommit bool) {
	hasParentCommit, err := c.CommitHasParentCommitByCommitHash(hash)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return hasParentCommit
}

func (c *CommandExecutorGitRepository) MustGetAsLocalDirectory() (l *LocalDirectory) {
	l, err := c.GetAsLocalDirectory()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return l
}

func (c *CommandExecutorGitRepository) MustGetAsLocalGitRepository() (l *LocalGitRepository) {
	l, err := c.GetAsLocalGitRepository()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return l
}

func (c *CommandExecutorGitRepository) MustGetAuthorEmailByCommitHash(hash string) (authorEmail string) {
	authorEmail, err := c.GetAuthorEmailByCommitHash(hash)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return authorEmail
}

func (c *CommandExecutorGitRepository) MustGetAuthorStringByCommitHash(hash string) (authorEmail string) {
	authorEmail, err := c.GetAuthorStringByCommitHash(hash)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return authorEmail
}

func (c *CommandExecutorGitRepository) MustGetCommitAgeDurationByCommitHash(hash string) (ageDuration *time.Duration) {
	ageDuration, err := c.GetCommitAgeDurationByCommitHash(hash)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return ageDuration
}

func (c *CommandExecutorGitRepository) MustGetCommitAgeSecondsByCommitHash(hash string) (ageSeconds float64) {
	ageSeconds, err := c.GetCommitAgeSecondsByCommitHash(hash)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return ageSeconds
}

func (c *CommandExecutorGitRepository) MustGetCommitByHash(hash string) (gitCommit *GitCommit) {
	gitCommit, err := c.GetCommitByHash(hash)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitCommit
}

func (c *CommandExecutorGitRepository) MustGetCommitMessageByCommitHash(hash string) (commitMessage string) {
	commitMessage, err := c.GetCommitMessageByCommitHash(hash)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return commitMessage
}

func (c *CommandExecutorGitRepository) MustGetCommitParentsByCommitHash(hash string, options *GitCommitGetParentsOptions) (commitParents []*GitCommit) {
	commitParents, err := c.GetCommitParentsByCommitHash(hash, options)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return commitParents
}

func (c *CommandExecutorGitRepository) MustGetCommitTimeByCommitHash(hash string) (commitTime *time.Time) {
	commitTime, err := c.GetCommitTimeByCommitHash(hash)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return commitTime
}

func (c *CommandExecutorGitRepository) MustGetCurrentCommit() (currentCommit *GitCommit) {
	currentCommit, err := c.GetCurrentCommit()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return currentCommit
}

func (c *CommandExecutorGitRepository) MustGetCurrentCommitHash() (currentCommitHash string) {
	currentCommitHash, err := c.GetCurrentCommitHash()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return currentCommitHash
}

func (c *CommandExecutorGitRepository) MustGetGitStatusOutput(verbose bool) (output string) {
	output, err := c.GetGitStatusOutput(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return output
}

func (c *CommandExecutorGitRepository) MustGetRootDirectoryPath(verbose bool) (rootDirectoryPath string) {
	rootDirectoryPath, err := c.GetRootDirectoryPath(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return rootDirectoryPath
}

func (c *CommandExecutorGitRepository) MustGetTagByHash(hash string) (tag *GitTag) {
	tag, err := c.GetTagByHash(hash)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return tag
}

func (c *CommandExecutorGitRepository) MustHasInitialCommit(verbose bool) (hasInitialCommit bool) {
	hasInitialCommit, err := c.HasInitialCommit(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return hasInitialCommit
}

func (c *CommandExecutorGitRepository) MustHasUncommittedChanges(verbose bool) (hasUncommitedChanges bool) {
	hasUncommitedChanges, err := c.HasUncommittedChanges(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return hasUncommitedChanges
}

func (c *CommandExecutorGitRepository) MustInit(options *CreateRepositoryOptions) {
	err := c.Init(options)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorGitRepository) MustIsBareRepository(verbose bool) (isBare bool) {
	isBare, err := c.IsBareRepository(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return isBare
}

func (c *CommandExecutorGitRepository) MustIsGitRepository(verbose bool) (isRepository bool) {
	isRepository, err := c.IsGitRepository(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return isRepository
}

func (c *CommandExecutorGitRepository) MustIsInitialized(verbose bool) (isInitialited bool) {
	isInitialited, err := c.IsInitialized(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return isInitialited
}

func (c *CommandExecutorGitRepository) MustPull(verbose bool) {
	err := c.Pull(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorGitRepository) MustPush(verbose bool) {
	err := c.Push(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorGitRepository) MustRunGitCommand(gitCommand []string, verbose bool) (commandOutput *CommandOutput) {
	commandOutput, err := c.RunGitCommand(gitCommand, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return commandOutput
}

func (c *CommandExecutorGitRepository) MustRunGitCommandAndGetStdoutAsString(command []string, verbose bool) (stdout string) {
	stdout, err := c.RunGitCommandAndGetStdoutAsString(command, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return stdout
}

func (c *CommandExecutorGitRepository) MustSetDefaultAuthor(verbose bool) {
	err := c.SetDefaultAuthor(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorGitRepository) MustSetGitConfig(options *GitConfigSetOptions) {
	err := c.SetGitConfig(options)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorGitRepository) MustSetUserEmail(email string, verbose bool) {
	err := c.SetUserEmail(email, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorGitRepository) MustSetUserName(name string, verbose bool) {
	err := c.SetUserName(name, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorGitRepository) Pull(verbose bool) (err error) {
	path, hostDescription, err := c.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	if verbose {
		LogInfof(
			"Pull git repository '%s' on '%s' started.",
			path,
			hostDescription,
		)
	}

	_, err = c.RunGitCommand(
		[]string{"pull"},
		verbose,
	)
	if err != nil {
		return err
	}

	if verbose {
		LogInfof(
			"Pull git repository '%s' on '%s' finished.",
			path,
			hostDescription,
		)
	}

	return
}

func (c *CommandExecutorGitRepository) Push(verbose bool) (err error) {
	path, hostDescription, err := c.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	if verbose {
		LogInfof(
			"Push git repository '%s' on '%s' started.",
			path,
			hostDescription,
		)
	}

	_, err = c.RunGitCommand(
		[]string{"push"},
		verbose,
	)
	if err != nil {
		return err
	}

	if verbose {
		LogInfof(
			"Push git repository '%s' on '%s' finished.",
			path,
			hostDescription,
		)
	}

	return
}

func (c *CommandExecutorGitRepository) RunGitCommand(gitCommand []string, verbose bool) (commandOutput *CommandOutput, err error) {
	if len(gitCommand) <= 0 {
		return nil, TracedError("gitCommand has no elements")
	}

	path, err := c.GetPath()
	if err != nil {
		return nil, err
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	commandToUse := append([]string{"git", "-C", path}, gitCommand...)

	return commandExecutor.RunCommand(
		&RunCommandOptions{
			Command:            commandToUse,
			Verbose:            verbose,
			LiveOutputOnStdout: verbose,
		},
	)
}

func (c *CommandExecutorGitRepository) RunGitCommandAndGetStdoutAsString(command []string, verbose bool) (stdout string, err error) {
	commandOutput, err := c.RunGitCommand(command, verbose)
	if err != nil {
		return "", err
	}

	return commandOutput.GetStdoutAsString()
}

func (c *CommandExecutorGitRepository) SetDefaultAuthor(verbose bool) (err error) {
	err = c.SetUserName(GitRepositryDefualtAuthorName(), verbose)
	if err != nil {
		return err
	}

	err = c.SetUserEmail(GitRepositryDefualtAuthorEmail(), verbose)
	if err != nil {
		return err
	}

	if verbose {
		path, hostDescription, err := c.GetPathAndHostDescription()
		if err != nil {
			return err
		}

		LogChangedf(
			"Initialized git repository '%s' on '%s' with default author and email.",
			path,
			hostDescription,
		)
	}

	return nil
}

func (c *CommandExecutorGitRepository) SetGitConfig(options *GitConfigSetOptions) (err error) {
	if options == nil {
		return TracedErrorNil("options")
	}

	if options.Email != "" {
		err = c.SetUserEmail(options.Email, options.Verbose)
		if err != nil {
			return err
		}
	}

	if options.Name != "" {
		err = c.SetUserName(options.Name, options.Verbose)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *CommandExecutorGitRepository) SetUserEmail(email string, verbose bool) (err error) {
	if email == "" {
		return TracedErrorEmptyString("email")
	}

	_, err = c.RunGitCommand(
		[]string{"config", "user.email", email},
		verbose,
	)
	if err != nil {
		return err
	}

	if verbose {
		path, hostDescription, err := c.GetPathAndHostDescription()
		if err != nil {
			return err
		}

		LogInfof(
			"Set git user email to '%s' for git repository '%s' on host '%s'.",
			email,
			path,
			hostDescription,
		)
	}

	return nil
}

func (c *CommandExecutorGitRepository) SetUserName(name string, verbose bool) (err error) {
	if name == "" {
		return TracedErrorEmptyString("name")
	}

	_, err = c.RunGitCommand(
		[]string{"config", "user.name", name},
		verbose,
	)
	if err != nil {
		return err
	}

	if verbose {
		path, hostDescription, err := c.GetPathAndHostDescription()
		if err != nil {
			return err
		}

		LogInfof(
			"Set git user name to '%s' for git repository '%s' on host '%s'.",
			name,
			path,
			hostDescription,
		)
	}

	return nil
}
