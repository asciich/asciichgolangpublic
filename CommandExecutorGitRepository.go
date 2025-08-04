package asciichgolangpublic

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/slicesutils"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// This is the GitRepository implementation based on a CommandExecutor (e.g. Bash, SSH...).
// This means it's a wrapper around the "git" binary which needs to be available.
// While very inefficient this solution can manage git repository on remote hosts, inside containers...
// which makes it very flexible.
//
// When dealing with locally available repositories it's recommended to use the LocalGitRepository
// implementation which uses go build in git functionality instead.
type CommandExecutorGitRepository struct {
	files.CommandExecutorDirectory
	GitRepositoryBase
}

func GetCommandExecutorGitRepositoryFromDirectory(directory filesinterfaces.Directory) (c *CommandExecutorGitRepository, err error) {
	if directory == nil {
		return nil, tracederrors.TracedErrorNil("directory")
	}

	commandExecutoryDirectory, ok := directory.(*files.CommandExecutorDirectory)
	if ok {
		commandExecutor, path, err := commandExecutoryDirectory.GetCommandExecutorAndDirPath()
		if err != nil {
			return nil, err
		}

		return GetCommandExecutorGitRepositoryByPath(commandExecutor, path)
	}

	localDirectory, ok := directory.(*files.LocalDirectory)
	if ok {
		path, err := localDirectory.GetPath()
		if err != nil {
			return nil, err
		}

		return GetLocalCommandExecutorGitRepositoryByPath(path)
	}

	unknownTypeName, err := datatypes.GetTypeName(directory)
	if err != nil {
		return nil, err
	}

	return nil, tracederrors.TracedErrorf(
		"Not implemented for directory type = '%s'",
		unknownTypeName,
	)
}

func NewCommandExecutorGitRepository(commandExecutor commandexecutorinterfaces.CommandExecutor) (c *CommandExecutorGitRepository, err error) {
	if commandExecutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExecutor")
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
func (c *CommandExecutorGitRepository) GetAsLocalDirectory() (l *files.LocalDirectory, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

// This function was only added to fulfil the current interface.
// On the long run this method has to be removed.
func (c *CommandExecutorGitRepository) GetAsLocalGitRepository() (l *LocalGitRepository, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (c *CommandExecutorGitRepository) AddFileByPath(pathToAdd string, verbose bool) (err error) {
	if pathToAdd == "" {
		return tracederrors.TracedErrorEmptyString("pathToAdd")
	}

	_, err = c.RunGitCommand(contextutils.GetVerbosityContextByBool(verbose), []string{"add", pathToAdd})
	if err != nil {
		return err
	}

	if verbose {
		path, hostDescription, err := c.GetPathAndHostDescription()
		if err != nil {
			return err
		}

		logging.LogChangedf(
			"Added '%s' to git repository '%s' on host '%s'.",
			pathToAdd,
			path,
			hostDescription,
		)
	}

	return nil
}

func (c *CommandExecutorGitRepository) AddRemote(remoteOptions *gitparameteroptions.GitRemoteAddOptions) (err error) {
	if remoteOptions == nil {
		return tracederrors.TracedError("remoteOptions is nil")
	}

	remoteName, err := remoteOptions.GetRemoteName()
	if err != nil {
		return err
	}

	remoteUrl, err := remoteOptions.GetRemoteUrl()
	if err != nil {
		return err
	}

	repoPath, err := c.GetPath()
	if err != nil {
		return err
	}

	remoteExists, err := c.RemoteConfigurationExists(
		&GitRemoteConfig{
			RemoteName: remoteName,
			UrlFetch:   remoteUrl,
			UrlPush:    remoteUrl,
		},
		remoteOptions.Verbose,
	)
	if err != nil {
		return err
	}

	if remoteExists {
		if remoteOptions.Verbose {
			logging.LogInfof("Remote '%s' as '%s' to repository '%s' already exists.", remoteUrl, remoteName, repoPath)
		}
	} else {
		err = c.RemoveRemoteByName(remoteName, remoteOptions.Verbose)
		if err != nil {
			return err
		}

		_, err = c.RunGitCommand(contextutils.GetVerbosityContextByBool(remoteOptions.Verbose), []string{"remote", "add", remoteName, remoteUrl})
		if err != nil {
			return err
		}

		if remoteOptions.Verbose {
			logging.LogChangedf("Added remote '%s' as '%s' to repository '%s'.", remoteUrl, remoteName, repoPath)
		}
	}

	return nil
}

func (c *CommandExecutorGitRepository) CheckoutBranchByName(name string, verbose bool) (err error) {
	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	currentBranchName, err := c.GetCurrentBranchName(verbose)
	if err != nil {
		return err
	}

	path, hostDescription, err := c.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	if currentBranchName == name {
		if verbose {
			logging.LogInfof(
				"Git repository '%s' on host '%s' is already checked out on branch '%s'.",
				path,
				hostDescription,
				name,
			)
		}
	} else {
		_, err := c.RunGitCommand(
			contextutils.GetVerbosityContextByBool(verbose),
			[]string{"checkout", name},
		)
		if err != nil {
			return err
		}

		if verbose {
			logging.LogChangedf(
				"Git repository '%s' on host '%s' checked out on branch '%s'.",
				path,
				hostDescription,
				name,
			)
		}
	}

	return nil
}

func (c *CommandExecutorGitRepository) CloneRepository(repository GitRepository, verbose bool) (err error) {
	if repository == nil {
		return tracederrors.TracedErrorNil("repository")
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

	return c.CloneRepositoryByPathOrUrl(pathToClone, verbose)
}

func (c *CommandExecutorGitRepository) CloneRepositoryByPathOrUrl(pathOrUrlToClone string, verbose bool) (err error) {
	if pathOrUrlToClone == "" {
		return tracederrors.TracedErrorEmptyString("pathToClone")
	}

	path, hostDescription, err := c.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	if verbose {
		logging.LogInfof(
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
		logging.LogInfof(
			"'%s' is already an initialized git repository on host '%s'. Skip clone.",
			path,
			hostDescription,
		)
	} else {
		commandExecutor, err := c.GetCommandExecutor()
		if err != nil {
			return err
		}

		ctx := contextutils.GetVerbosityContextByBool(verbose)
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

	if verbose {
		logging.LogInfof(
			"Cloning git repository '%s' to '%s' on host '%s' finished.",
			pathOrUrlToClone,
			path,
			hostDescription,
		)
	}

	return nil
}

func (c *CommandExecutorGitRepository) Commit(commitOptions *gitparameteroptions.GitCommitOptions) (createdCommit *GitCommit, err error) {
	if commitOptions == nil {
		return nil, tracederrors.TracedErrorNil("commitOptions")
	}

	commitCommand := []string{"commit"}

	if commitOptions.AllowEmpty {
		commitCommand = append(commitCommand, "--allow-empty")
	}

	if commitOptions.CommitAllChanges {
		commitCommand = append(commitCommand, "--all")
	}

	message, err := commitOptions.GetMessage()
	if err != nil {
		return nil, err
	}

	commitCommand = append(commitCommand, "-m", message)

	_, err = c.RunGitCommand(
		contextutils.GetVerbosityContextByBool(commitOptions.Verbose),
		commitCommand,
	)
	if err != nil {
		return nil, err
	}

	createdCommit, err = c.GetCurrentCommit(commitOptions.Verbose)
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

		logging.LogChangedf(
			"Created commit '%s' in git repository '%s' on host '%s'.",
			createdHash,
			path,
			hostDescription,
		)
	}

	return createdCommit, nil
}

func (c *CommandExecutorGitRepository) CommitHasParentCommitByCommitHash(hash string) (hasParentCommit bool, err error) {
	return false, tracederrors.TracedErrorNotImplemented()
}

func (c *CommandExecutorGitRepository) CreateBranch(createOptions *parameteroptions.CreateBranchOptions) (err error) {
	if createOptions == nil {
		return tracederrors.TracedErrorNil("createOptions")
	}

	name, err := createOptions.GetName()
	if err != nil {
		return err
	}

	branchExists, err := c.BranchByNameExists(name, createOptions.Verbose)
	if err != nil {
		return err
	}

	path, hostDescription, err := c.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	if branchExists {
		logging.LogInfof(
			"Branch '%s' already exists in git repository '%s' on host '%s'.",
			name,
			path,
			hostDescription,
		)
	} else {
		_, err = c.RunGitCommand(
			contextutils.GetVerbosityContextByBool(createOptions.Verbose),
			[]string{"checkout", "-b", name},
		)
		if err != nil {
			return err
		}

		logging.LogChangedf(
			"Branch '%s' in git repository '%s' on host '%s' created.",
			name,
			path,
			hostDescription,
		)
	}

	return nil
}

func (c *CommandExecutorGitRepository) CreateTag(options *gitparameteroptions.GitRepositoryCreateTagOptions) (createdTag GitTag, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	tagName, err := options.GetTagName()
	if err != nil {
		return nil, err
	}

	tagMessage := tagName
	if options.IsTagCommentSet() {
		tagMessage, err = options.GetTagComment()
		if err != nil {
			return nil, err
		}
	}

	hashToTag := ""
	if options.IsCommitHashSet() {
		hashToTag, err = options.GetCommitHash()
		if err != nil {
			return nil, err
		}
	} else {
		hashToTag, err = c.GetCurrentCommitHash(options.Verbose)
		if err != nil {
			return nil, err
		}
	}

	_, err = c.RunGitCommand(
		contextutils.GetVerbosityContextByBool(options.Verbose),
		[]string{"tag", "-a", tagName, hashToTag, "-m", tagMessage},
	)
	if err != nil {
		return nil, err
	}

	path, hostDescription, err := c.GetPathAndHostDescription()
	if err != nil {
		return nil, err
	}

	createdTag, err = c.GetTagByName(tagName)
	if err != nil {
		return nil, err
	}

	if options.Verbose {
		logging.LogChangedf(
			"Created tag '%s' for commit '%s' in git repository '%s' on host '%s'.",
			tagName,
			hashToTag,
			path,
			hostDescription,
		)
	}

	return createdTag, nil
}

func (c *CommandExecutorGitRepository) DeleteBranchByName(name string, verbose bool) (err error) {
	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	branchExists, err := c.BranchByNameExists(name, verbose)
	if err != nil {
		return err
	}

	path, hostDescription, err := c.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	if branchExists {
		_, err := c.RunGitCommand(
			contextutils.GetVerbosityContextByBool(verbose),
			[]string{"branch", "-D", name},
		)
		if err != nil {
			return err
		}

		if verbose {
			logging.LogChangedf(
				"Branch '%s' in git repository '%s' on host '%s' deleted.",
				name,
				path,
				hostDescription,
			)
		}

	} else {
		if verbose {
			logging.LogInfof(
				"Branch '%s' in git repository '%s' on host '%s' is already absent. Skip delete.",
				name,
				path,
				hostDescription,
			)
		}
	}

	return nil
}

func (c *CommandExecutorGitRepository) Fetch(verbose bool) (err error) {
	_, err = c.RunGitCommand(
		contextutils.GetVerbosityContextByBool(verbose),
		[]string{"fetch"},
	)
	if err != nil {
		return err
	}

	if verbose {
		path, hostDescription, err := c.GetPathAndHostDescription()
		if err != nil {
			return err
		}

		logging.LogChangedf(
			"Fetched git repository '%s' on host '%s'",
			path,
			hostDescription,
		)
	}

	return nil
}

func (c *CommandExecutorGitRepository) FileByPathExists(path string, verbose bool) (exists bool, err error) {
	if path == "" {
		return false, tracederrors.TracedErrorEmptyString(path)
	}

	return c.FileInDirectoryExists(verbose, path)
}

func (c *CommandExecutorGitRepository) GetAuthorEmailByCommitHash(hash string) (authorEmail string, err error) {
	return "", tracederrors.TracedErrorNotImplemented()
}

func (c *CommandExecutorGitRepository) GetAuthorStringByCommitHash(hash string) (authorEmail string, err error) {
	return "", tracederrors.TracedErrorNotImplemented()
}

func (c *CommandExecutorGitRepository) GetCommitAgeDurationByCommitHash(hash string) (ageDuration *time.Duration, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (c *CommandExecutorGitRepository) GetCommitAgeSecondsByCommitHash(hash string) (ageSeconds float64, err error) {
	return -1, tracederrors.TracedErrorNotImplemented()
}

func (c *CommandExecutorGitRepository) GetCommitByHash(hash string) (gitCommit *GitCommit, err error) {
	if hash == "" {
		return nil, tracederrors.TracedErrorEmptyString("hash")
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
	if hash == "" {
		return "", tracederrors.TracedErrorEmptyString("hash")
	}

	stdout, err := c.RunGitCommandAndGetStdoutAsString(
		contextutils.ContextSilent(),
		[]string{"log", "-n", "1", "--pretty=format:%s", hash},
	)
	if err != nil {
		return "", err
	}

	commitMessage = strings.TrimSpace(stdout)

	if commitMessage == "" {
		path, hostDescription, err := c.GetPathAndHostDescription()
		if err != nil {
			return "", err
		}

		return "", tracederrors.TracedErrorf(
			"Unable to get commit message for hash '%s' in git repository '%s' on host '%s'. commitMessage is empty string after evaluation.",
			hash,
			path,
			hostDescription,
		)
	}

	return commitMessage, nil
}

func (c *CommandExecutorGitRepository) GetCommitParentsByCommitHash(hash string, options *parameteroptions.GitCommitGetParentsOptions) (commitParents []*GitCommit, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (c *CommandExecutorGitRepository) GetCommitTimeByCommitHash(hash string) (commitTime *time.Time, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (c *CommandExecutorGitRepository) GetCurrentBranchName(verbose bool) (branchName string, err error) {
	stdout, err := c.RunGitCommandAndGetStdoutAsString(
		contextutils.GetVerbosityContextByBool(verbose),
		[]string{"rev-parse", "--abbrev-ref", "HEAD"},
	)
	if err != nil {
		return "", err
	}

	branchName = strings.TrimSpace(stdout)

	path, hostDescription, err := c.GetPathAndHostDescription()
	if err != nil {
		return "", err
	}

	if branchName == "" {
		return "", tracederrors.TracedErrorf(
			"Unable to get branch name for git repository '%s' on host '%s'. branchName is empty string after evaluation.",
			path,
			hostDescription,
		)
	}

	if verbose {
		logging.LogInfof(
			"Branch '%s' is currently checked out in git repository '%s' on host '%s'.",
			branchName,
			path,
			hostDescription,
		)
	}

	return branchName, nil
}

func (c *CommandExecutorGitRepository) GetCurrentCommit(verbose bool) (currentCommit *GitCommit, err error) {
	currentCommitHash, err := c.GetCurrentCommitHash(verbose)
	if err != nil {
		return nil, err
	}

	currentCommit, err = c.GetCommitByHash(currentCommitHash)
	if err != nil {
		return nil, err
	}

	if verbose {
		path, hostDescription, err := c.GetPathAndHostDescription()
		if err != nil {
			return nil, err
		}

		logging.LogInfof(
			"Current commit of git repository '%s' on host '%s' has hash '%s'.",
			path,
			hostDescription,
			currentCommitHash,
		)
	}

	return currentCommit, nil
}

func (c *CommandExecutorGitRepository) GetCurrentCommitHash(verbose bool) (currentCommitHash string, err error) {
	currentCommitHash, err = c.RunGitCommandAndGetStdoutAsString(
		contextutils.GetVerbosityContextByBool(verbose),
		[]string{"rev-parse", "HEAD"},
	)
	if err != nil {
		return "", err
	}

	currentCommitHash = strings.TrimSpace(currentCommitHash)

	return currentCommitHash, nil
}

func (c *CommandExecutorGitRepository) GetDirectoryByPath(pathToSubDir ...string) (subDir filesinterfaces.Directory, err error) {
	if len(pathToSubDir) <= 0 {
		return nil, tracederrors.TracedError("pathToSubdir has no elements")
	}

	return c.GetSubDirectory(pathToSubDir...)
}

func (c *CommandExecutorGitRepository) GetGitStatusOutput(verbose bool) (output string, err error) {
	return "", tracederrors.TracedErrorNotImplemented()
}

func (c *CommandExecutorGitRepository) GetHashByTagName(tagName string) (hash string, err error) {
	if tagName == "" {
		return "", tracederrors.TracedErrorEmptyString("tagName")
	}

	stdoutLines, err := c.RunGitCommandAndGetStdoutAsLines(
		contextutils.ContextSilent(),
		[]string{"show-ref", "--dereference", tagName},
	)
	if err != nil {
		return "", err
	}

	for _, line := range stdoutLines {
		if strings.HasSuffix(line, "{}") {
			hash = strings.Split(line, " ")[0]
			break
		}
	}

	hash = strings.TrimSpace(hash)

	if hash == "" {
		return "", tracederrors.TracedError("hash is empty string after evaluation")
	}

	return hash, nil
}

func (c *CommandExecutorGitRepository) GetRemoteConfigs(verbose bool) (remoteConfigs []*GitRemoteConfig, err error) {
	output, err := c.RunGitCommand(contextutils.GetVerbosityContextByBool(verbose), []string{"remote", "-v"})
	if err != nil {
		return nil, err
	}

	outputLines, err := output.GetStdoutAsLines(false)
	if err != nil {
		return nil, err
	}

	remoteConfigs = []*GitRemoteConfig{}
	for _, line := range outputLines {
		line = strings.TrimSpace(line)
		if len(line) <= 0 {
			continue
		}

		lineCleaned := strings.ReplaceAll(line, "\t", " ")

		splitted := stringsutils.SplitAtSpacesAndRemoveEmptyStrings(lineCleaned)
		if len(splitted) != 3 {
			return nil, tracederrors.TracedErrorf("Unable to parse '%s' as remote. splitted is '%v'", line, splitted)
		}

		remoteName := splitted[0]
		remoteUrl := splitted[1]
		remoteDirection := splitted[2]

		var remoteToModify *GitRemoteConfig = nil
		for _, existingRemote := range remoteConfigs {
			if existingRemote.RemoteName == remoteName {
				remoteToModify = existingRemote
			}
		}

		if remoteToModify == nil {
			remoteToAdd := NewGitRemoteConfig()
			remoteToAdd.RemoteName = remoteName
			remoteConfigs = append(remoteConfigs, remoteToAdd)
			remoteToModify = remoteToAdd
		}

		if remoteDirection == "(fetch)" {
			remoteToModify.UrlFetch = remoteUrl
		} else if remoteDirection == "(push)" {
			remoteToModify.UrlPush = remoteUrl
		} else {
			return nil, tracederrors.TracedErrorf("Unknown remoteDirection='%s'", remoteDirection)
		}
	}

	return remoteConfigs, nil
}

func (c *CommandExecutorGitRepository) GetRootDirectory(ctx context.Context) (rootDirectory filesinterfaces.Directory, err error) {
	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	rootDirPath, err := c.GetRootDirectoryPath(ctx)
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

func (c *CommandExecutorGitRepository) GetRootDirectoryPath(ctx context.Context) (rootDirectoryPath string, err error) {
	path, hostDescription, err := c.GetPathAndHostDescription()
	if err != nil {
		return "", err
	}

	isBareRepository, err := c.IsBareRepository(ctx)
	if err != nil {
		return "", err
	}

	if isBareRepository {
		var cwd filesinterfaces.Directory

		commandExecutor, err := c.GetCommandExecutor()
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
		stdout, err := c.RunGitCommandAndGetStdoutAsString(
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

func (c *CommandExecutorGitRepository) GetTagByName(name string) (tag GitTag, err error) {
	if name == "" {
		return nil, tracederrors.TracedErrorEmptyString("name")
	}

	toReturn := NewGitRepositoryTag()

	err = toReturn.SetName(name)
	if err != nil {
		return nil, err
	}

	err = toReturn.SetName(name)
	if err != nil {
		return nil, err
	}

	err = toReturn.SetGitRepository(c)
	if err != nil {
		return nil, err
	}

	return toReturn, nil
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
			logging.LogInfof(
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

	ctx := contextutils.GetVerbosityContextByBool(verbose)
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

	if verbose {
		if hasInitialCommit {
			logging.LogInfof(
				"Git repository '%s' on host '%s' has an initial commit",
				path,
				hostDescription,
			)
		} else {
			logging.LogInfof(
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
		contextutils.GetVerbosityContextByBool(verbose),
		&parameteroptions.RunCommandOptions{
			Command: []string{
				"bash",
				"-c",
				fmt.Sprintf(
					"cd '%s' && git diff && git diff --cached && git status --porcelain",
					path,
				),
			},
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
			logging.LogInfof(
				"Git repository '%s' on '%s' has uncommited changes.",
				path,
				hostDescription,
			)
		} else {
			logging.LogInfof(
				"Git repository '%s' on '%s' has no uncommited changes.",
				path,
				hostDescription,
			)
		}
	}

	return hasUncommitedChanges, nil
}

func (c *CommandExecutorGitRepository) Init(options *parameteroptions.CreateRepositoryOptions) (err error) {
	if options == nil {
		return tracederrors.TracedErrorNil("options")
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
			logging.LogInfof(
				"Git repository '%s' on host '%s' is already initialized.",
				path,
				hostDescription,
			)
		}
	} else {
		err = c.Create(contextutils.GetVerbosityContextByBool(options.Verbose))
		if err != nil {
			return err
		}

		commandToUse := []string{"init"}

		if options.BareRepository {
			commandToUse = append(commandToUse, "--bare")
		}

		_, err = c.RunGitCommand(
			contextutils.GetVerbosityContextByBool(options.Verbose),
			commandToUse,
		)
		if err != nil {
			return err
		}

		if options.Verbose {
			if options.BareRepository {
				logging.LogChangedf(
					"Git repository '%s' on host '%s' initialized as bare repository.",
					path,
					hostDescription,
				)
			} else {
				logging.LogChangedf(
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
			logging.LogInfof(
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
						&gitparameteroptions.GitConfigSetOptions{
							Name:    GitRepositryDefaultAuthorName(),
							Email:   GitRepositryDefaultAuthorEmail(),
							Verbose: options.Verbose,
						},
					)
				}

				_, err = temporaryClone.CommitAndPush(
					&gitparameteroptions.GitCommitOptions{
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
					&gitparameteroptions.GitCommitOptions{
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
				logging.LogChangedf(
					"Initialized repository '%s' on host '%s' with an empty commit.",
					path,
					hostDescription,
				)
			}
		}
	}

	return nil
}

func (c *CommandExecutorGitRepository) IsBareRepository(ctx context.Context) (isBare bool, err error) {
	path, hostDescription, err := c.GetPathAndHostDescription()
	if err != nil {
		return false, err
	}

	stdout, err := c.RunGitCommandAndGetStdoutAsString(ctx, []string{"rev-parse", "--is-bare-repository"})
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
			logging.LogInfof(
				"'%s' on host '%s' is a git repository",
				path,
				hostDescription,
			)
		} else {
			logging.LogInfof(
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
			logging.LogInfof(
				"Git repository '%s' does not exist on host '%s' and is therefore not initalized.",
				path,
				hostDescription,
			)
		}
		return false, nil
	}

	ctx := contextutils.GetVerbosityContextByBool(verbose)
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

	if verbose {
		if isInitialited {
			logging.LogInfof(
				"Git repository '%s' on host '%s' is initialized.",
				path,
				hostDescription,
			)
		} else {
			logging.LogInfof(
				"'%s' is not an initialized git repository on host '%s'.",
				path,
				hostDescription,
			)
		}
	}

	return isInitialited, nil
}

func (c *CommandExecutorGitRepository) ListBranchNames(verbose bool) (branchNames []string, err error) {
	lines, err := c.RunGitCommandAndGetStdoutAsLines(contextutils.GetVerbosityContextByBool(verbose), []string{"branch"})
	if err != nil {
		return nil, err
	}

	branchNames = []string{}

	for _, line := range lines {
		line = strings.TrimPrefix(line, "* ")
		line = strings.TrimSpace(line)

		branchNames = append(branchNames, line)
	}

	sort.Strings(branchNames)

	if verbose {
		path, hostDescripton, err := c.GetPathAndHostDescription()
		if err != nil {
			return nil, err
		}

		logging.LogInfof(
			"Found '%d' branches in git repository '%s' on host '%s'.",
			len(branchNames),
			path,
			hostDescripton,
		)
	}

	return branchNames, nil
}

func (c *CommandExecutorGitRepository) ListTagNames(verbose bool) (tagNames []string, err error) {
	return c.RunGitCommandAndGetStdoutAsLines(
		contextutils.ContextSilent(), // Do not clutter output by pritning all tags.
		[]string{"tag"},
	)
}

func (c *CommandExecutorGitRepository) ListTags(verbose bool) (tags []GitTag, err error) {
	tagNames, err := c.ListTagNames(verbose)
	if err != nil {
		return nil, err
	}

	tags = []GitTag{}
	for _, name := range tagNames {
		toAdd, err := c.GetTagByName(name)
		if err != nil {
			return nil, err
		}

		tags = append(tags, toAdd)
	}

	return tags, nil
}

func (c *CommandExecutorGitRepository) ListTagsForCommitHash(hash string, verbose bool) (tags []GitTag, err error) {
	if hash == "" {
		return nil, tracederrors.TracedErrorEmptyString("hash")
	}

	tagNames, err := c.RunGitCommandAndGetStdoutAsLines(
		contextutils.ContextSilent(),
		[]string{"tag", "--points-at", "HEAD"},
	)
	if err != nil {
		return nil, err
	}

	tags = []GitTag{}
	for _, name := range tagNames {
		toAdd, err := c.GetTagByName(name)
		if err != nil {
			return nil, err
		}

		tags = append(tags, toAdd)
	}

	return tags, nil
}

func (c *CommandExecutorGitRepository) Pull(ctx context.Context) (err error) {
	path, hostDescription, err := c.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Pull git repository '%s' on '%s' started.", path, hostDescription)

	_, err = c.RunGitCommand(
		ctx,
		[]string{"pull"},
	)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Pull git repository '%s' on '%s' finished.", path, hostDescription)

	return
}

func (c *CommandExecutorGitRepository) PullFromRemote(pullOptions *GitPullFromRemoteOptions) (err error) {
	if pullOptions == nil {
		return tracederrors.TracedError("pullOptions not set")
	}

	remoteName, err := pullOptions.GetRemoteName()
	if err != nil {
		return err
	}

	branchName, err := pullOptions.GetBranchName()
	if err != nil {
		return err
	}

	if len(remoteName) <= 0 {
		return tracederrors.TracedError("remoteName is empty string")
	}

	path, hostDescription, err := c.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	_, err = c.RunGitCommand(contextutils.GetVerbosityContextByBool(pullOptions.Verbose), []string{"pull", remoteName, branchName})
	if err != nil {
		return err
	}

	if pullOptions.Verbose {
		logging.LogInfof(
			"Pulled git repository '%s' on host '%s' from remote '%s'.",
			path,
			hostDescription,
			remoteName,
		)
	}

	return nil
}

func (c *CommandExecutorGitRepository) Push(ctx context.Context) (err error) {
	path, hostDescription, err := c.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Push git repository '%s' on '%s' started.", path, hostDescription)

	_, err = c.RunGitCommand(
		ctx,
		[]string{"push"},
	)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Push git repository '%s' on '%s' finished.", path, hostDescription)

	return
}

func (c *CommandExecutorGitRepository) PushTagsToRemote(remoteName string, verbose bool) (err error) {
	if len(remoteName) <= 0 {
		return tracederrors.TracedError("remoteName is empty string")
	}

	path, hostDescription, err := c.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	_, err = c.RunGitCommand(
		contextutils.GetVerbosityContextByBool(verbose),
		[]string{"push", remoteName, "--tags"},
	)
	if err != nil {
		return err
	}

	if verbose {
		logging.LogInfof(
			"Pushed tags of git repository '%s' on host '%s' to remote '%s'.",
			path,
			hostDescription,
			remoteName,
		)
	}

	return nil
}

func (c *CommandExecutorGitRepository) PushToRemote(remoteName string, verbose bool) (err error) {
	if len(remoteName) <= 0 {
		return tracederrors.TracedError("remoteName is empty string")
	}

	_, err = c.RunGitCommand(contextutils.GetVerbosityContextByBool(verbose), []string{"push", remoteName})
	if err != nil {
		return err
	}

	if verbose {
		path, hostDescription, err := c.GetPathAndHostDescription()
		if err != nil {
			return err
		}

		logging.LogInfof(
			"Pushed git repository '%s' on host '%s' to remote '%s'.",
			path,
			hostDescription,
			remoteName,
		)
	}

	return nil
}

func (c *CommandExecutorGitRepository) RemoteByNameExists(remoteName string, verbose bool) (remoteExists bool, err error) {
	if len(remoteName) <= 0 {
		return false, fmt.Errorf("remoteName is empty string")
	}

	remoteConfigs, err := c.GetRemoteConfigs(verbose)
	if err != nil {
		return false, err
	}

	for _, toCheck := range remoteConfigs {
		if toCheck.RemoteName == remoteName {
			return true, nil
		}
	}

	return false, nil
}

func (c *CommandExecutorGitRepository) RemoteConfigurationExists(config *GitRemoteConfig, verbose bool) (exists bool, err error) {
	if config == nil {
		return false, tracederrors.TracedError("config is nil")
	}

	remoteConfigs, err := c.GetRemoteConfigs(verbose)
	if err != nil {
		return false, err
	}

	for _, toCheck := range remoteConfigs {
		if config.Equals(toCheck) {
			return true, nil
		}
	}

	return false, nil
}

func (c *CommandExecutorGitRepository) RemoveRemoteByName(remoteNameToRemove string, verbose bool) (err error) {
	if len(remoteNameToRemove) <= 0 {
		return tracederrors.TracedError("remoteNameToRemove is empty string")
	}

	remoteExists, err := c.RemoteByNameExists(remoteNameToRemove, verbose)
	if err != nil {
		return err
	}

	repoDirPath, err := c.GetPath()
	if err != nil {
		return err
	}

	if remoteExists {
		_, err := c.RunGitCommand(contextutils.GetVerbosityContextByBool(verbose), []string{"remote", "remove", remoteNameToRemove})
		if err != nil {
			return err
		}

		if verbose {
			logging.LogChangedf("Remote '%s' for repository '%s' removed.", remoteNameToRemove, repoDirPath)
		}
	} else {
		if verbose {
			logging.LogInfof("Remote '%s' for repository '%s' already deleted.", remoteNameToRemove, repoDirPath)
		}
	}

	return nil
}

func (c *CommandExecutorGitRepository) RunGitCommand(ctx context.Context, gitCommand []string) (commandOutput *commandoutput.CommandOutput, err error) {
	if len(gitCommand) <= 0 {
		return nil, tracederrors.TracedError("gitCommand has no elements")
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
		commandexecutorgeneric.WithLiveOutputOnStdoutIfVerbose(ctx),
		&parameteroptions.RunCommandOptions{
			Command: commandToUse,
		},
	)
}

func (c *CommandExecutorGitRepository) RunGitCommandAndGetStdoutAsLines(ctx context.Context, command []string) (lines []string, err error) {
	if command == nil {
		return nil, tracederrors.TracedErrorNil("command")
	}

	output, err := c.RunGitCommand(ctx, command)
	if err != nil {
		return nil, err
	}

	lines, err = output.GetStdoutAsLines(true)
	if err != nil {
		return nil, err
	}

	return lines, nil
}

func (c *CommandExecutorGitRepository) RunGitCommandAndGetStdoutAsString(ctx context.Context, command []string) (stdout string, err error) {
	commandOutput, err := c.RunGitCommand(ctx, command)
	if err != nil {
		return "", err
	}

	return commandOutput.GetStdoutAsString()
}

func (c *CommandExecutorGitRepository) SetDefaultAuthor(verbose bool) (err error) {
	err = c.SetUserName(GitRepositryDefaultAuthorName(), verbose)
	if err != nil {
		return err
	}

	err = c.SetUserEmail(GitRepositryDefaultAuthorEmail(), verbose)
	if err != nil {
		return err
	}

	if verbose {
		path, hostDescription, err := c.GetPathAndHostDescription()
		if err != nil {
			return err
		}

		logging.LogChangedf(
			"Initialized git repository '%s' on '%s' with default author and email.",
			path,
			hostDescription,
		)
	}

	return nil
}

func (c *CommandExecutorGitRepository) SetGitConfig(options *gitparameteroptions.GitConfigSetOptions) (err error) {
	if options == nil {
		return tracederrors.TracedErrorNil("options")
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

func (c *CommandExecutorGitRepository) SetRemoteUrl(remoteUrl string, verbose bool) (err error) {
	remoteUrl = strings.TrimSpace(remoteUrl)
	if len(remoteUrl) <= 0 {
		return tracederrors.TracedError("remoteUrl is empty string")
	}

	name := "origin"

	_, err = c.RunGitCommand(contextutils.GetVerbosityContextByBool(verbose), []string{"remote", "set-url", name, remoteUrl})
	if err != nil {
		return err
	}

	if verbose {
		path, hostDescription, err := c.GetPathAndHostDescription()
		if err != nil {
			return err
		}

		logging.LogChangedf(
			"Set remote Url for '%v' in git repository '%v' on host '%s' to '%v'.",
			name,
			path,
			hostDescription,
			remoteUrl,
		)
	}

	return nil
}

func (c *CommandExecutorGitRepository) SetUserEmail(email string, verbose bool) (err error) {
	if email == "" {
		return tracederrors.TracedErrorEmptyString("email")
	}

	_, err = c.RunGitCommand(contextutils.GetVerbosityContextByBool(verbose), []string{"config", "user.email", email})
	if err != nil {
		return err
	}

	if verbose {
		path, hostDescription, err := c.GetPathAndHostDescription()
		if err != nil {
			return err
		}

		logging.LogInfof(
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
		return tracederrors.TracedErrorEmptyString("name")
	}

	_, err = c.RunGitCommand(contextutils.GetVerbosityContextByBool(verbose), []string{"config", "user.name", name})
	if err != nil {
		return err
	}

	if verbose {
		path, hostDescription, err := c.GetPathAndHostDescription()
		if err != nil {
			return err
		}

		logging.LogInfof(
			"Set git user name to '%s' for git repository '%s' on host '%s'.",
			name,
			path,
			hostDescription,
		)
	}

	return nil
}
