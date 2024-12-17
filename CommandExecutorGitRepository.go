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

	stdout, err := c.RunGitCommandAndGetStdoutAsString(
		[]string{"log", "--all", "-n1"},
		verbose,
	)
	if err != nil {
		if verbose {
			commandExecutor, err := c.GetCommandExecutor()
			if err != nil {
				return false, err
			}

			commandExecutor.RunCommand(
				&RunCommandOptions{
					Command:            []string{"ls", path},
					Verbose:            verbose,
					LiveOutputOnStdout: verbose,
				},
			)
		}

		return false, err
	}

	stdout = strings.TrimSpace(stdout)

	hasInitialCommit = stdout != ""

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

	if options.InitializeWithEmptyCommit {
		hasInitialCommit, err := c.HasInitialCommit(options.Verbose)
		if err != nil {
			return err
		}

		if !hasInitialCommit {
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
				"sh",
				"-c",
				fmt.Sprintf(
					"git -C '%s' rev-parse --is-bare-repository &> /dev/null && echo yes || echo no",
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

func (c *CommandExecutorGitRepository) MustHasInitialCommit(verbose bool) (hasInitialCommit bool) {
	hasInitialCommit, err := c.HasInitialCommit(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return hasInitialCommit
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

func (c *CommandExecutorGitRepository) MustIsInitialized(verbose bool) (isInitialited bool) {
	isInitialited, err := c.IsInitialized(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return isInitialited
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
