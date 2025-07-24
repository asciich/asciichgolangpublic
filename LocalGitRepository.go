package asciichgolangpublic

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"

	"github.com/asciich/asciichgolangpublic/files"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorbashoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/shellutils/shelllinehandler"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type LocalGitRepository struct {
	files.LocalDirectory
	GitRepositoryBase
}

func GetLocalGitReposioryFromDirectory(directory files.Directory) (repo GitRepository, err error) {
	if directory == nil {
		return nil, tracederrors.TracedErrorNil("directory")
	}

	isLocalDirectory, err := directory.IsLocalDirectory()
	if err != nil {
		return nil, err
	}

	if !isLocalDirectory {
		return nil, tracederrors.TracedError("Only local directories are supported.")
	}

	localPath, err := directory.GetLocalPath()
	if err != nil {
		return nil, err
	}

	repo, err = GetLocalGitRepositoryByPath(localPath)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func GetLocalGitReposioryFromLocalDirectory(localDirectory *files.LocalDirectory) (l *LocalGitRepository, err error) {
	if localDirectory == nil {
		return nil, tracederrors.TracedErrorNil("directory")
	}

	localPath, err := localDirectory.GetLocalPath()
	if err != nil {
		return nil, err
	}

	l, err = GetLocalGitRepositoryByPath(localPath)
	if err != nil {
		return nil, err
	}

	return l, nil
}

func GetLocalGitRepositoryByPath(path string) (l *LocalGitRepository, err error) {
	if path == "" {
		return nil, tracederrors.TracedErrorEmptyString("path")
	}

	l = NewLocalGitRepository()

	err = l.SetLocalPath(path)
	if err != nil {
		return nil, err
	}

	return l, nil
}

func MustGetLocalGitReposioryFromDirectory(directory files.Directory) (repo GitRepository) {
	repo, err := GetLocalGitReposioryFromDirectory(directory)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return repo
}

func MustGetLocalGitReposioryFromLocalDirectory(localDirectory *files.LocalDirectory) (l *LocalGitRepository) {
	l, err := GetLocalGitReposioryFromLocalDirectory(localDirectory)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return l
}

func MustGetLocalGitRepositoryByPath(path string) (l *LocalGitRepository) {
	l, err := GetLocalGitRepositoryByPath(path)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return l
}

func NewLocalGitRepository() (l *LocalGitRepository) {
	l = new(LocalGitRepository)

	err := l.SetParentDirectoryForBaseClass(l)
	if err != nil {
		panic(err)
	}

	err = l.SetParentRepositoryForBaseClass(l)
	if err != nil {
		panic(err)
	}

	return l
}

// TODO remove: LocalGitRepository should purely base on goGit, not by calling the git binary.
func (l *LocalGitRepository) RunGitCommand(ctx context.Context, gitCommand []string) (commandOutput *commandoutput.CommandOutput, err error) {
	if gitCommand == nil {
		return nil, tracederrors.TracedErrorEmptyString("gitCommand")
	}

	repoRootPath, err := l.GetRootDirectoryPath(ctx)
	if err != nil {
		return nil, err
	}

	gitCommandString, err := shelllinehandler.Join(gitCommand)
	if err != nil {
		return nil, err
	}

	command := fmt.Sprintf(
		"git -C '%s' %s",
		repoRootPath,
		gitCommandString,
	)

	commandOutput, err = commandexecutorbashoo.Bash().RunOneLiner(ctx, command)
	if err != nil {
		return nil, err
	}

	return commandOutput, nil
}

func (c *LocalGitRepository) HasInitialCommit(verbose bool) (hasInitialCommit bool, err error) {
	_, err = c.GetCurrentCommit(verbose)
	if err != nil {
		if errors.Is(err, ErrGitRepositoryDoesNotExist) { // The repository does not even exist.
			return false, nil
		}

		if errors.Is(err, ErrGitRepositoryHeadNotFound) { // The repository exists but has no initial commit and therefore no head is found.
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (l *LocalGitRepository) AddFileByPath(pathToAdd string, verbose bool) (err error) {
	if pathToAdd == "" {
		return tracederrors.TracedErrorNil("path")
	}

	worktree, err := l.GetGoGitWorktree()
	if err != nil {
		return err
	}

	_, err = worktree.Add(pathToAdd)
	if err != nil {
		return tracederrors.TracedErrorf("%w", err)
	}

	if verbose {
		path, err := l.GetPath()
		if err != nil {
			return err
		}

		logging.LogChangedf(
			"Added file '%s' to git repository '%s' on localhost",
			pathToAdd,
			path,
		)
	}

	return nil
}

func (l *LocalGitRepository) AddRemote(remoteOptions *gitparameteroptions.GitRemoteAddOptions) (err error) {
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

	repoPath, err := l.GetPath()
	if err != nil {
		return err
	}

	remoteExists, err := l.RemoteConfigurationExists(
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
		err = l.RemoveRemoteByName(remoteName, remoteOptions.Verbose)
		if err != nil {
			return err
		}

		// TODO reimplement without calling the git command.
		_, err = l.RunGitCommand(
			contextutils.GetVerbosityContextByBool(remoteOptions.Verbose),
			[]string{"remote", "add", remoteName, remoteUrl},
		)
		if err != nil {
			return err
		}

		if remoteOptions.Verbose {
			logging.LogChangedf("Added remote '%s' as '%s' to repository '%s'.", remoteUrl, remoteName, repoPath)
		}
	}

	return nil
}

func (l *LocalGitRepository) CheckoutBranchByName(name string, verbose bool) (err error) {
	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	currentBranchName, err := l.GetCurrentBranchName(verbose)
	if err != nil {
		return err
	}

	path, hostDescription, err := l.GetPathAndHostDescription()
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
		_, err := l.RunGitCommand(
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

func (l *LocalGitRepository) CloneRepository(repository GitRepository, verbose bool) (err error) {
	if repository == nil {
		return tracederrors.TracedErrorNil("repository")
	}

	repositoryHostDescription, err := repository.GetHostDescription()
	if err != nil {
		return err
	}

	if repositoryHostDescription != "localhost" {
		return tracederrors.TracedErrorf(
			"Only cloning from local repositories is implemented at the moment but got '%s'",
			repositoryHostDescription,
		)
	}

	pathToClone, err := repository.GetPath()
	if err != nil {
		return err
	}

	return l.CloneRepositoryByPathOrUrl(pathToClone, verbose)
}

func (l *LocalGitRepository) CloneRepositoryByPathOrUrl(urlOrPathToClone string, verbose bool) (err error) {
	if urlOrPathToClone == "" {
		return tracederrors.TracedErrorEmptyString("pathToClone")
	}

	path, hostDescription, err := l.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	if verbose {
		logging.LogInfof(
			"Cloning git repository '%s' to '%s' on host '%s' started.",
			urlOrPathToClone,
			path,
			hostDescription,
		)
	}

	isInitialzed, err := l.IsInitialized(verbose)
	if err != nil {
		return err
	}

	if isInitialzed {
		logging.LogInfof(
			"'%s' is already an initialized git repository on host '%s'. Skip clone.",
			path,
			hostDescription,
		)
	} else {
		const isBare = false
		_, err = git.PlainClone(
			path,
			isBare,
			&git.CloneOptions{
				URL: urlOrPathToClone,
			},
		)
		if err != nil {
			if err.Error() == "remote repository is empty" {
				if verbose {
					logging.LogInfof(
						"Remote repository '%s' is empty. Going to add remote for empty repository.",
						urlOrPathToClone,
					)
				}

				err = l.Init(
					&parameteroptions.CreateRepositoryOptions{
						Verbose:                   verbose,
						BareRepository:            isBare,
						InitializeWithEmptyCommit: false,
					})
				if err != nil {
					return err
				}

				_, err = l.SetRemote("origin", urlOrPathToClone, verbose)
				if err != nil {
					return err
				}
			} else {
				return tracederrors.TracedErrorf(
					"Clone '%s' failed: '%w'",
					urlOrPathToClone,
					err,
				)
			}
		}
	}

	if verbose {
		logging.LogInfof(
			"Cloning git repository '%s' to '%s' on host '%s' finished.",
			urlOrPathToClone,
			path,
			hostDescription,
		)
	}

	return nil
}

func (l *LocalGitRepository) Commit(commitOptions *gitparameteroptions.GitCommitOptions) (createdCommit *GitCommit, err error) {
	if commitOptions == nil {
		return nil, tracederrors.TracedErrorNil("commitOptions")
	}

	worktree, err := l.GetGoGitWorktree()
	if err != nil {
		return nil, err
	}

	commitMessage, err := commitOptions.GetMessage()
	if err != nil {
		return nil, err
	}

	hash, err := worktree.Commit(
		commitMessage,
		&git.CommitOptions{
			AllowEmptyCommits: commitOptions.GetAllowEmpty(),
		},
	)
	if err != nil {
		return nil, err
	}

	createdCommit, err = l.GetCommitByGoGitHash(&hash)
	if err != nil {
		return nil, err
	}

	path, err := l.GetLocalPath()
	if err != nil {
		return nil, err
	}

	if commitOptions.Verbose {
		logging.LogChangedf(
			"Created commit '%s' with hash '%s' in git repository '%s'.",
			commitMessage,
			hash.String(),
			path,
		)
	}

	return createdCommit, nil
}

func (l *LocalGitRepository) CommitHasParentCommitByCommitHash(hash string) (hasParentCommit bool, err error) {
	if hash == "" {
		return false, tracederrors.TracedErrorEmptyString("hash")
	}

	goGitCommit, err := l.GetGoGitCommitByCommitHash(hash)
	if err != nil {
		return false, err
	}

	hasParentCommit = goGitCommit.NumParents() > 0

	return hasParentCommit, nil
}

func (l *LocalGitRepository) CreateBranch(createOptions *parameteroptions.CreateBranchOptions) (err error) {
	if createOptions == nil {
		return tracederrors.TracedErrorNil("createOptions")
	}

	name, err := createOptions.GetName()
	if err != nil {
		return err
	}

	branchExists, err := l.BranchByNameExists(name, createOptions.Verbose)
	if err != nil {
		return err
	}

	path, hostDescription, err := l.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	if branchExists {
		if createOptions.Verbose {
			logging.LogInfof(
				"Branch '%s' already exists in git repository '%s' on host '%s'.",
				name,
				path,
				hostDescription,
			)
		}
	} else {
		/* TODO implement using native GoGit
		worktree, err := l.GetGoGitWorktree()
		if err != nil {
			return err
		}

		err = worktree.Checkout(&git.CheckoutOptions{
			Create: true,
			Branch: plumbing.ReferenceName("ref/branch/" + name),
		})

		if err != nil {
			return tracederrors.TracedErrorf(
				"Unable to create branch '%s' in git repository '%s' on host '%s': '%w'",
				name,
				path,
				hostDescription,
				err,
			)
		}
		*/
		l.RunGitCommand(
			contextutils.GetVerbosityContextByBool(createOptions.Verbose),
			[]string{"checkout", "-b", name},
		)

		if createOptions.Verbose {
			logging.LogChangedf(
				"Branch '%s' in git repository '%s' on host '%s' created.",
				name,
				path,
				hostDescription,
			)
		}
	}

	return nil
}

func (l *LocalGitRepository) CreateTag(options *gitparameteroptions.GitRepositoryCreateTagOptions) (createdTag GitTag, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	tagName, err := options.GetTagName()
	if err != nil {
		return nil, err
	}

	goRepo, err := l.GetAsGoGitRepository()
	if err != nil {
		return nil, err
	}

	hashToTag := ""
	if options.IsCommitHashSet() {
		hashToTag, err = options.GetCommitHash()
		if err != nil {
			return nil, err
		}
	} else {
		hashToTag, err = l.GetCurrentCommitHash(options.Verbose)
		if err != nil {
			return nil, err
		}
	}

	path, err := l.GetPath()
	if err != nil {
		return nil, err
	}

	createTagOptions := &git.CreateTagOptions{}

	if options.IsTagCommentSet() {
		tagComment, err := options.GetTagComment()
		if err != nil {
			return nil, err
		}
		createTagOptions.Message = tagComment
	} else {
		createTagOptions.Message = tagName
	}

	if options.Verbose {
		logging.LogInfof(
			"Going to create tag '%s' on commit '%s' in local git repository '%s'.",
			tagName,
			hashToTag,
			path,
		)
	}

	goHash, err := l.GetGoGitHashFromHashString(hashToTag)
	if err != nil {
		return nil, err
	}

	_, err = goRepo.CreateTag(
		tagName,
		*goHash,
		createTagOptions,
	)
	if err != nil {
		return nil, tracederrors.TracedErrorf(
			"Creating tag failed: %w",
			err,
		)
	}

	if options.Verbose {

		logging.LogChangedf(
			"Created tag '%s' in local git repository '%s'.",
			tagName,
			path,
		)
	}

	createdTag, err = l.GetTagByName(tagName)
	if err != nil {
		return nil, err
	}

	return createdTag, nil
}

func (l *LocalGitRepository) DeleteBranchByName(name string, verbose bool) (err error) {
	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	branchExists, err := l.BranchByNameExists(name, verbose)
	if err != nil {
		return err
	}

	path, hostDescription, err := l.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	if branchExists {
		_, err := l.RunGitCommand(
			contextutils.GetVerbosityContextByBool(verbose),
			[]string{"branch", "-D", name},
		)
		if err != nil {
			return err
		}

		/* TODO implement using native go git
		goGitRepo, err := l.GetAsGoGitRepository()
		if err != nil {
			return err
		}

		err = goGitRepo.DeleteBranch("refs/heads/" + name)
		if err != nil {
			return tracederrors.TracedErrorf(
				"Delete branch '%s' in git repository '%s' on host '%s' failed: '%w'",
				name,
				path,
				hostDescription,
				err,
			)
		}
		*/

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

func (l *LocalGitRepository) Fetch(verbose bool) (err error) {
	_, err = l.RunGitCommand(
		contextutils.GetVerbosityContextByBool(verbose),
		[]string{"fetch"},
	)
	if err != nil {
		return err
	}

	if verbose {
		path, hostDescription, err := l.GetPathAndHostDescription()
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

func (l *LocalGitRepository) FileByPathExists(path string, verbose bool) (exists bool, err error) {
	if path == "" {
		return false, tracederrors.TracedErrorEmptyString(path)
	}

	return l.FileInDirectoryExists(verbose, path)
}

func (l *LocalGitRepository) GetAsGoGitRepository() (goGitRepository *git.Repository, err error) {
	repoPath, err := l.GetLocalPath()
	if err != nil {
		return nil, err
	}

	goGitRepository, err = git.PlainOpen(repoPath)
	if err != nil {
		if errors.Is(err, git.ErrRepositoryNotExists) {
			toReturn := tracederrors.TracedErrorf("%w: repoPath='%s'", ErrGitRepositoryDoesNotExist, repoPath)
			return nil, tracederrors.AddErrorToUnwrapToTracedError(toReturn, err)
		}
		return nil, tracederrors.TracedErrorf("%w: repoPath='%s'", err, repoPath)
	}

	if goGitRepository == nil {
		return nil, tracederrors.TracedError("goGitRepository is nil after evaluation.")
	}

	return goGitRepository, nil
}

func (l *LocalGitRepository) GetAsLocalDirectory() (localDirectory *files.LocalDirectory, err error) {
	localPath, err := l.GetLocalPath()
	if err != nil {
		return nil, err
	}

	localDirectory, err = files.GetLocalDirectoryByPath(localPath)
	if err != nil {
		return nil, err
	}

	return localDirectory, nil
}

func (l *LocalGitRepository) GetAsLocalGitRepository() (localGitRepository *LocalGitRepository, err error) {
	localPath, err := l.GetLocalPath()
	if err != nil {
		return nil, err
	}

	localGitRepository, err = GetLocalGitRepositoryByPath(localPath)
	if err != nil {
		return nil, err
	}

	return localGitRepository, nil
}

func (l *LocalGitRepository) GetAuthorEmailByCommitHash(hash string) (authorEmail string, err error) {
	if hash == "" {
		return "", tracederrors.TracedErrorEmptyString("hash")
	}

	goGitCommit, err := l.GetGoGitCommitByCommitHash(hash)
	if err != nil {
		return "", err
	}

	authorEmail = goGitCommit.Author.Email

	return authorEmail, nil
}

func (l *LocalGitRepository) GetAuthorStringByCommitHash(hash string) (authorString string, err error) {
	if hash == "" {
		return "", tracederrors.TracedErrorEmptyString(hash)
	}

	goGitCommit, err := l.GetGoGitCommitByCommitHash(hash)
	if err != nil {
		return "", err
	}

	authorString = goGitCommit.Author.String()

	return authorString, nil
}

func (l *LocalGitRepository) GetCommitAgeDurationByCommitHash(hash string) (ageDuration *time.Duration, err error) {
	if hash == "" {
		return nil, tracederrors.TracedErrorEmptyString("hash")
	}

	commitTime, err := l.GetCommitTimeByCommitHash(hash)
	if err != nil {
		return nil, err
	}

	ageDurationNonPtr := time.Since(*commitTime)
	ageDuration = &ageDurationNonPtr

	return ageDuration, nil
}

func (l *LocalGitRepository) GetCommitAgeSecondsByCommitHash(hash string) (ageSeconds float64, err error) {
	if hash == "" {
		return -1, tracederrors.TracedErrorEmptyString("hash")
	}

	ageDuration, err := l.GetCommitAgeDurationByCommitHash(hash)
	if err != nil {
		return -1, err
	}

	ageSeconds = ageDuration.Seconds()

	return ageSeconds, nil
}

func (l *LocalGitRepository) GetCommitByGoGitCommit(goGitCommit *object.Commit) (gitCommit *GitCommit, err error) {
	if goGitCommit == nil {
		return nil, tracederrors.TracedErrorNil("goGitCommit")
	}

	hash := goGitCommit.Hash

	gitCommit, err = l.GetCommitByGoGitHash(&hash)
	if err != nil {
		return nil, err
	}

	return gitCommit, nil
}

func (l *LocalGitRepository) GetCommitByGoGitHash(goGitHash *plumbing.Hash) (gitCommit *GitCommit, err error) {
	if goGitHash == nil {
		return nil, tracederrors.TracedErrorNil("goGitHash")
	}

	gitCommit = NewGitCommit()

	err = gitCommit.SetGitRepo(l)
	if err != nil {
		return nil, err
	}

	goGitRepo, err := l.GetAsGoGitRepository()
	if err != nil {
		return nil, err
	}

	commitObject, err := goGitRepo.CommitObject(*goGitHash)
	if err != nil {
		return nil, tracederrors.TracedErrorf("%w", err)
	}

	err = gitCommit.SetHash(commitObject.Hash.String())
	if err != nil {
		return nil, err
	}

	return gitCommit, nil
}

func (l *LocalGitRepository) GetCommitByGoGitReference(goGitReference *plumbing.Reference) (gitCommit *GitCommit, err error) {
	if goGitReference == nil {
		return nil, tracederrors.TracedErrorNil("goGitReference")
	}

	hash := goGitReference.Hash()

	gitCommit, err = l.GetCommitByGoGitHash(&hash)
	if err != nil {
		return nil, err
	}

	return gitCommit, nil
}

func (l *LocalGitRepository) GetCommitMessageByCommitHash(hash string) (commitMessage string, err error) {
	if hash == "" {
		return "", tracederrors.TracedErrorEmptyString("hash")
	}

	g, err := l.GetGoGitCommitByCommitHash(hash)
	if err != nil {
		return "", err
	}

	commitMessage = g.Message

	return commitMessage, nil
}

func (l *LocalGitRepository) GetCommitParentsByCommitHash(hash string, options *parameteroptions.GitCommitGetParentsOptions) (commitParents []*GitCommit, err error) {
	if hash == "" {
		return nil, tracederrors.TracedErrorEmptyString("hash")
	}

	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	goGitCommit, err := l.GetGoGitCommitByCommitHash(hash)
	if err != nil {
		return nil, err
	}

	parents := goGitCommit.Parents()
	defer parents.Close()

	for {
		parentToAdd, err := parents.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, tracederrors.TracedErrorf("Unable to get next parent: %w", err)
		}

		toAdd, err := l.GetCommitByGoGitCommit(parentToAdd)
		if err != nil {
			return nil, err
		}

		commitParents = append(commitParents, toAdd)

		if options.IncludeParentsOfParents {
			additionalParents, err := toAdd.GetParentCommits(&parameteroptions.GitCommitGetParentsOptions{
				IncludeParentsOfParents: true,
			})
			if err != nil {
				return nil, err
			}

			commitParents = append(commitParents, additionalParents...)
		}
	}

	if options.Verbose {
		logging.LogInfof("Collected '%d' parent commits for commit '%s'.", len(commitParents), hash)
	}

	return commitParents, nil
}

func (l *LocalGitRepository) GetCommitTimeByCommitHash(hash string) (commitTime *time.Time, err error) {
	if hash == "" {
		return nil, tracederrors.TracedErrorEmptyString("hash")
	}

	goGitCommit, err := l.GetGoGitCommitByCommitHash(hash)
	if err != nil {
		return nil, err
	}

	commitTime = &(goGitCommit.Author.When)

	return commitTime, nil
}

func (l *LocalGitRepository) GetCurrentBranchName(verbose bool) (branchName string, err error) {
	goGitRepo, err := l.GetAsGoGitRepository()
	if err != nil {
		return "", err
	}

	head, err := goGitRepo.Head()
	if err != nil {
		return "", tracederrors.TracedErrorf("Unable to get head: '%w'", err)
	}
	branchName = head.String()
	branchName = stringsutils.SplitAndGetLastElement(branchName, " ")
	branchName = stringsutils.SplitAndGetLastElement(branchName, "/")
	branchName = strings.TrimSpace(branchName)

	path, hostDescription, err := l.GetPathAndHostDescription()
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

func (l *LocalGitRepository) GetCurrentCommit(verbose bool) (gitCommit *GitCommit, err error) {
	head, err := l.GetGoGitHead()
	if err != nil {
		return nil, err
	}

	gitCommit, err = l.GetCommitByGoGitReference(head)
	if err != nil {
		return nil, err
	}

	if verbose {
		hash, err := gitCommit.GetHash()
		if err != nil {
			return nil, err
		}

		path, err := l.GetPath()
		if err != nil {
			return nil, err
		}

		logging.LogInfof(
			"Current commit in local git repository '%s' has hash '%s'.",
			path,
			hash,
		)
	}

	return gitCommit, nil
}

func (l *LocalGitRepository) GetCurrentCommitGoGitHash(verbose bool) (hash *plumbing.Hash, err error) {
	currentHashBytes, err := l.GetCurrentCommitHashAsBytes(verbose)
	if err != nil {
		return nil, err
	}

	hashValue := plumbing.Hash(currentHashBytes)

	return &hashValue, nil
}

func (l *LocalGitRepository) GetCurrentCommitHash(verbose bool) (commitHash string, err error) {
	commit, err := l.GetCurrentCommit(verbose)
	if err != nil {
		return "", err
	}

	commitHash, err = commit.GetHash()
	if err != nil {
		return "", err
	}

	return commitHash, nil
}

func (l *LocalGitRepository) GetCurrentCommitHashAsBytes(verbose bool) (hash []byte, err error) {
	currentHash, err := l.GetCurrentCommitHash(verbose)
	if err != nil {
		return nil, err
	}

	return stringsutils.HexStringToBytes(currentHash)
}

func (l *LocalGitRepository) GetDirectoryByPath(pathToSubDir ...string) (subDir files.Directory, err error) {
	if len(pathToSubDir) <= 0 {
		return nil, tracederrors.TracedError("pathToSubdir has no elements")
	}

	return l.GetSubDirectory(pathToSubDir...)
}

func (l *LocalGitRepository) GetGitStatusOutput(verbose bool) (output string, err error) {
	output, err = l.RunGitCommandAndGetStdout([]string{"status"}, verbose)
	if err != nil {
		return "", err
	}

	return output, nil
}

func (l *LocalGitRepository) GetGitlabCiYamlFile() (gitlabCiYamlFile *GitlabCiYamlFile, err error) {
	ciYamlFile, err := l.GetFileInDirectory(Gitlab().GetDefaultGitlabCiYamlFileName())
	if err != nil {
		return nil, err
	}

	gitlabCiYamlFile, err = GetGitlabCiYamlFileByFile(ciYamlFile)
	if err != nil {
		return nil, err
	}

	return gitlabCiYamlFile, nil
}

func (l *LocalGitRepository) GetGoGitCommitByCommitHash(hash string) (goGitCommit *object.Commit, err error) {
	if hash == "" {
		return nil, tracederrors.TracedErrorEmptyString("hash")
	}

	goGitRepo, err := l.GetAsGoGitRepository()
	if err != nil {
		return nil, err
	}

	pHash := plumbing.NewHash(hash)

	goGitCommit, err = goGitRepo.CommitObject(pHash)
	if err != nil {
		return nil, tracederrors.TracedErrorf("%w", err)
	}

	return goGitCommit, err
}

func (l *LocalGitRepository) GetGoGitConfig() (config *config.Config, err error) {
	goGitRepo, err := l.GetAsGoGitRepository()
	if err != nil {
		return nil, err
	}

	config, err = goGitRepo.Config()
	if err != nil {
		return nil, tracederrors.TracedErrorf("%w", err)
	}

	if config == nil {
		return nil, tracederrors.TracedError("config is nil after evaluation")
	}

	return config, nil
}

func (l *LocalGitRepository) GetGoGitHashFromHashString(hashString string) (hash *plumbing.Hash, err error) {
	if hashString == "" {
		return nil, tracederrors.TracedErrorNil("hashString")
	}

	hashBytes, err := stringsutils.HexStringToBytes(hashString)
	if err != nil {
		return nil, err
	}

	hashValue := plumbing.Hash(hashBytes)

	return &hashValue, err
}

func (l *LocalGitRepository) GetGoGitHead() (head *plumbing.Reference, err error) {
	goGitRepo, err := l.GetAsGoGitRepository()
	if err != nil {
		return nil, err
	}

	head, err = goGitRepo.Head()
	if err != nil {
		if errors.Is(err, plumbing.ErrReferenceNotFound) {
			toReturn := tracederrors.TracedErrorf("%w", ErrGitRepositoryHeadNotFound)
			return nil, tracederrors.AddErrorToUnwrapToTracedError(toReturn, err)
		}
		return nil, tracederrors.TracedErrorf("%w", err)
	}

	return head, nil
}

func (l *LocalGitRepository) GetGoGitWorktree() (worktree *git.Worktree, err error) {
	goGitRepo, err := l.GetAsGoGitRepository()
	if err != nil {
		return nil, err
	}

	worktree, err = goGitRepo.Worktree()
	if err != nil {
		return nil, tracederrors.TracedErrorf("%w", err)
	}

	return worktree, nil
}

func (l *LocalGitRepository) GetHashByTagName(tagName string) (hash string, err error) {
	if tagName == "" {
		return "", err
	}

	nativeRepo, err := l.GetAsGoGitRepository()
	if err != nil {
		return "", err
	}

	nativeTagObjects, err := nativeRepo.TagObjects()
	if err != nil {
		return "", tracederrors.TracedErrorf(
			"Unable to get native tags: %w",
			err,
		)
	}
	defer nativeTagObjects.Close()

	for {
		tag, err := nativeTagObjects.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return "", tracederrors.TracedErrorf(
				"Unable to get next tag: '%w'",
				err,
			)
		}

		name := tag.Name
		name = strings.TrimPrefix(name, "refs/tags/")

		if tagName == name {
			hash = string(tag.Target.String())

			if hash == "" {
				return "", tracederrors.TracedError(
					"hash is empty string after evaluation.",
				)
			}

			return hash, nil
		}
	}

	path, err := l.GetPath()
	if err != nil {
		return "", err
	}

	return "", tracederrors.TracedErrorf(
		"Unable to get hash for tag '%s' in local git repository '%s'.",
		tagName,
		path,
	)
}

func (l *LocalGitRepository) GetRemoteConfigs(verbose bool) (remoteConfigs []*GitRemoteConfig, err error) {
	// TODO reimplement without calling the git binary.
	output, err := l.RunGitCommand(
		contextutils.GetVerbosityContextByBool(verbose),
		[]string{"remote", "-v"},
	)
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

func (l *LocalGitRepository) GetRootDirectory(ctx context.Context) (rootDirectory files.Directory, err error) {
	rootDirectoryPath, err := l.GetRootDirectoryPath(ctx)
	if err != nil {
		return nil, err
	}

	rootDirectory, err = files.GetLocalDirectoryByPath(rootDirectoryPath)
	if err != nil {
		return nil, err
	}

	return rootDirectory, nil
}

func (l *LocalGitRepository) GetRootDirectoryPath(ctx context.Context) (rootDirectoryPath string, err error) {
	pathToCheck, err := l.GetLocalPath()
	if err != nil {
		return "", err
	}

	searchedFromPath := pathToCheck

	for {
		localDirToCheck, err := files.GetLocalDirectoryByPath(pathToCheck)
		if err != nil {
			return "", nil
		}

		localPathToCheck, err := localDirToCheck.GetLocalPath()
		if err != nil {
			return "", nil
		}

		if localPathToCheck == "" || localPathToCheck == "/" {
			return "", tracederrors.TracedErrorf("Not inside a git repository. Searched from '%s'", searchedFromPath)
		}

		// local git repository
		dotGitExists, err := localDirToCheck.SubDirectoryExists(".git", false)
		if err != nil {
			return "", nil
		}

		if dotGitExists {
			return pathToCheck, nil
		}

		// bare git repository
		if filepath.Base(localPathToCheck) != ".git" {
			refsExists, err := localDirToCheck.SubDirectoryExists("refs", false)
			if err != nil {
				return "", nil
			}

			objectsExists, err := localDirToCheck.SubDirectoryExists("objects", false)
			if err != nil {
				return "", nil
			}

			if refsExists && objectsExists {
				return pathToCheck, nil
			}
		}

		pathToCheck = filepath.Dir(pathToCheck)
	}
}

func (l *LocalGitRepository) GetTagByName(tagName string) (tag GitTag, err error) {
	if tagName == "" {
		return nil, tracederrors.TracedErrorEmptyString("tagName")
	}

	tagToReturn := NewGitRepositoryTag()

	err = tagToReturn.SetName(tagName)
	if err != nil {
		return nil, err
	}

	err = tagToReturn.SetGitRepository(l)
	if err != nil {
		return nil, err
	}

	return tagToReturn, nil
}

func (l *LocalGitRepository) GitlabCiYamlFileExists(verbose bool) (gitlabCiYamlFileExists bool, err error) {
	gitlabCiYamlFile, err := l.GetGitlabCiYamlFile()
	if err != nil {
		return false, err
	}

	gitlabCiYamlFileExists, err = gitlabCiYamlFile.Exists(verbose)
	if err != nil {
		return false, err
	}

	return gitlabCiYamlFileExists, nil
}

func (l *LocalGitRepository) HasNoUncommittedChanges(verbose bool) (hasUncommittedChanges bool, err error) {
	hasUncommittedChanges, err = l.HasUncommittedChanges(verbose)
	if err != nil {
		return false, err
	}

	return !hasUncommittedChanges, nil
}

func (l *LocalGitRepository) HasUncommittedChanges(verbose bool) (hasUncommittedChanges bool, err error) {
	worktree, err := l.GetGoGitWorktree()
	if err != nil {
		return false, err
	}

	gitStatus, err := worktree.Status()
	if err != nil {
		return false, err
	}

	if !gitStatus.IsClean() {
		hasUncommittedChanges = true
	}

	if verbose {
		path, err := l.GetPath()
		if err != nil {
			return false, err
		}

		if hasUncommittedChanges {
			logging.LogInfof(
				"Local git repository '%s' has uncommited changes.",
				path,
			)
		} else {
			logging.LogInfof(
				"Local git repository '%s' has no uncommited changes.",
				path,
			)
		}
	}

	return hasUncommittedChanges, nil
}

func (l *LocalGitRepository) Init(options *parameteroptions.CreateRepositoryOptions) (err error) {
	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	repoPath, err := l.GetLocalPath()
	if err != nil {
		return err
	}

	isInitialized, err := l.IsInitialized(options.Verbose)
	if err != nil {
		return err
	}

	if isInitialized {
		if options.Verbose {
			logging.LogInfof("Local git repository '%s' is already initialized.", repoPath)
		}
	} else {
		_, err = git.PlainInit(repoPath, options.BareRepository)
		if err != nil {
			return tracederrors.TracedErrorf("%w", err)
		}
		if options.Verbose {
			logging.LogChangedf("Local git repository '%s' is initialized.", repoPath)
		}
	}

	if options.InitializeWithEmptyCommit {
		hasInitialCommit, err := l.HasInitialCommit(options.Verbose)
		if err != nil {
			return err
		}

		if !hasInitialCommit {
			if options.BareRepository {
				temporaryRepository, err := GitRepositories().CloneGitRepositoryToTemporaryDirectory(
					l,
					options.Verbose,
				)
				if err != nil {
					return err
				}
				defer temporaryRepository.Delete(options.Verbose)

				err = temporaryRepository.SetGitConfig(
					&gitparameteroptions.GitConfigSetOptions{
						Name:    "asciichgolangpublic git repo initializer",
						Email:   "asciichgolangpublic@example.net",
						Verbose: options.Verbose,
					},
				)
				if err != nil {
					return err
				}

				_, err = temporaryRepository.CommitAndPush(
					&gitparameteroptions.GitCommitOptions{
						Message:    "Initial empty commit during repo initialization",
						AllowEmpty: true,
						Verbose:    true,
					},
				)
				if err != nil {
					return err
				}

				if options.Verbose {
					logging.LogChangedf("Initialized bare repository '%s' with an empty commit.", repoPath)
				}
			} else {
				if options.InitializeWithDefaultAuthor {
					err = l.SetGitConfig(
						&gitparameteroptions.GitConfigSetOptions{
							Name:    GitRepositryDefaultAuthorName(),
							Email:   GitRepositryDefaultAuthorEmail(),
							Verbose: options.Verbose,
						},
					)
					if err != nil {
						return err
					}
				}

				if options.InitializeWithEmptyCommit {
					_, err = l.Commit(
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
						"Initialized local repository '%s' with an empty commit.",
						repoPath,
					)
				}
			}
		}
	}

	if !options.BareRepository {
		if options.InitializeWithDefaultAuthor {
			err = l.SetGitConfig(
				&gitparameteroptions.GitConfigSetOptions{
					Name:    "asciichgolangpublic git repo initializer",
					Email:   "asciichgolangpublic@example.net",
					Verbose: options.Verbose,
				},
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (l *LocalGitRepository) IsBareRepository(ctx context.Context) (isBareRepository bool, err error) {

	config, err := l.GetGoGitConfig()
	if err != nil {
		return false, err
	}

	isBareRepository = config.Core.IsBare

	repoRoot, err := l.GetPath()
	if err != nil {
		return false, err
	}

	if isBareRepository {
		logging.LogInfoByCtxf(ctx, "Git repository '%s' is a bare repository.", repoRoot)
	} else {
		logging.LogInfoByCtxf(ctx, "Git repository '%s' is not a bare repository.", repoRoot)
	}

	return isBareRepository, nil
}

func (l *LocalGitRepository) IsGitRepository(verbose bool) (isGitRepository bool, err error) {
	isInitialited, err := l.IsInitialized(verbose)
	if err != nil {
		return false, err
	}

	return isInitialited, nil
}

func (l *LocalGitRepository) IsInitialized(verbose bool) (isInitialized bool, err error) {
	isInitialized = true

	_, err = l.GetAsGoGitRepository()
	if err != nil {
		if errors.Is(err, ErrGitRepositoryDoesNotExist) {
			isInitialized = false
		} else {
			return false, err
		}
	}

	if verbose {
		path, err := l.GetPath()
		if err != nil {
			return false, err
		}

		if isInitialized {
			logging.LogInfof(
				"Directory '%s' is an initialized git repository.",
				path,
			)
		} else {
			logging.LogInfof(
				"Directory '%s' is not an initialized git repository.",
				path,
			)
		}
	}

	return isInitialized, nil
}

func (l *LocalGitRepository) ListBranchNames(verbose bool) (branchNames []string, err error) {
	goRepo, err := l.GetAsGoGitRepository()
	if err != nil {
		return nil, err
	}

	branches, err := goRepo.Branches()
	if err != nil {
		return nil, err
	}
	defer branches.Close()

	branchNames = []string{}
	for {
		branch, err := branches.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, tracederrors.TracedErrorf("Unable to get next parent: %w", err)
		}

		nameToAdd := branch.Name().String()
		nameToAdd = strings.TrimPrefix(nameToAdd, "refs/heads/")

		branchNames = append(branchNames, nameToAdd)
	}

	sort.Strings(branchNames)

	if verbose {
		path, hostDescripton, err := l.GetPathAndHostDescription()
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

func (l *LocalGitRepository) ListTagNames(verbose bool) (tagNames []string, err error) {
	nativeRepo, err := l.GetAsGoGitRepository()
	if err != nil {
		return nil, err
	}

	tags, err := nativeRepo.Tags()
	if err != nil {
		return nil, tracederrors.TracedErrorf(
			"Unable to get native tags: %w",
			err,
		)
	}
	defer tags.Close()

	tagNames = []string{}
	for {
		tag, err := tags.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, tracederrors.TracedErrorf(
				"Unable to get next tag: '%w'",
				err,
			)
		}

		toAdd := tag.Name().String()

		toAdd = strings.TrimPrefix(toAdd, "refs/tags/")

		tagNames = append(tagNames, toAdd)
	}

	return tagNames, nil
}

func (l *LocalGitRepository) ListTags(verbose bool) (tags []GitTag, err error) {
	tagNames, err := l.ListTagNames(verbose)
	if err != nil {
		return nil, err
	}

	tags = []GitTag{}
	for _, name := range tagNames {
		toAdd, err := l.GetTagByName(name)
		if err != nil {
			return nil, err
		}

		tags = append(tags, toAdd)
	}

	return tags, nil
}

func (l *LocalGitRepository) ListTagsForCommitHash(hash string, verbose bool) (tags []GitTag, err error) {
	if hash == "" {
		return nil, tracederrors.TracedErrorEmptyString("hash")
	}

	nativeRepo, err := l.GetAsGoGitRepository()
	if err != nil {
		return nil, err
	}

	nativeTagObjects, err := nativeRepo.TagObjects()
	if err != nil {
		return nil, tracederrors.TracedErrorf(
			"Unable to get native tags: %w",
			err,
		)
	}
	defer nativeTagObjects.Close()

	tags = []GitTag{}
	for {
		tag, err := nativeTagObjects.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, tracederrors.TracedErrorf(
				"Unable to get next tag: '%w'",
				err,
			)
		}

		if tag.Target.String() == hash {
			nameToAdd := strings.TrimPrefix(tag.Name, "refs/tags/")

			toAdd, err := l.GetTagByName(nameToAdd)
			if err != nil {
				return nil, err
			}

			tags = append(tags, toAdd)
		}
	}

	return tags, nil
}

func (l *LocalGitRepository) MustAddFileByPath(pathToAdd string, verbose bool) {
	err := l.AddFileByPath(pathToAdd, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (l *LocalGitRepository) MustAddRemote(remoteOptions *gitparameteroptions.GitRemoteAddOptions) {
	err := l.AddRemote(remoteOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (l *LocalGitRepository) MustCheckoutBranchByName(name string, verbose bool) {
	err := l.CheckoutBranchByName(name, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (l *LocalGitRepository) MustCloneRepository(repository GitRepository, verbose bool) {
	err := l.CloneRepository(repository, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (l *LocalGitRepository) MustCloneRepositoryByPathOrUrl(pathToClone string, verbose bool) {
	err := l.CloneRepositoryByPathOrUrl(pathToClone, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (l *LocalGitRepository) MustCommit(commitOptions *gitparameteroptions.GitCommitOptions) (createdCommit *GitCommit) {
	createdCommit, err := l.Commit(commitOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return createdCommit
}

func (l *LocalGitRepository) MustCommitHasParentCommitByCommitHash(hash string) (hasParentCommit bool) {
	hasParentCommit, err := l.CommitHasParentCommitByCommitHash(hash)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return hasParentCommit
}

func (l *LocalGitRepository) MustCreateBranch(createOptions *parameteroptions.CreateBranchOptions) {
	err := l.CreateBranch(createOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (l *LocalGitRepository) MustCreateTag(options *gitparameteroptions.GitRepositoryCreateTagOptions) (createdTag GitTag) {
	createdTag, err := l.CreateTag(options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return createdTag
}

func (l *LocalGitRepository) MustDeleteBranchByName(name string, verbose bool) {
	err := l.DeleteBranchByName(name, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (l *LocalGitRepository) MustFetch(verbose bool) {
	err := l.Fetch(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (l *LocalGitRepository) MustFileByPathExists(path string, verbose bool) (exists bool) {
	exists, err := l.FileByPathExists(path, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return exists
}

func (l *LocalGitRepository) MustGetAsGoGitRepository() (goGitRepository *git.Repository) {
	goGitRepository, err := l.GetAsGoGitRepository()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return goGitRepository
}

func (l *LocalGitRepository) MustGetAsLocalDirectory() (localDirectory *files.LocalDirectory) {
	localDirectory, err := l.GetAsLocalDirectory()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return localDirectory
}

func (l *LocalGitRepository) MustGetAsLocalGitRepository() (localGitRepository *LocalGitRepository) {
	localGitRepository, err := l.GetAsLocalGitRepository()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return localGitRepository
}

func (l *LocalGitRepository) MustGetAuthorEmailByCommitHash(hash string) (authorEmail string) {
	authorEmail, err := l.GetAuthorEmailByCommitHash(hash)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return authorEmail
}

func (l *LocalGitRepository) MustGetAuthorStringByCommitHash(hash string) (authorString string) {
	authorString, err := l.GetAuthorStringByCommitHash(hash)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return authorString
}

func (l *LocalGitRepository) MustGetCommitAgeDurationByCommitHash(hash string) (ageDuration *time.Duration) {
	ageDuration, err := l.GetCommitAgeDurationByCommitHash(hash)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return ageDuration
}

func (l *LocalGitRepository) MustGetCommitAgeSecondsByCommitHash(hash string) (ageSeconds float64) {
	ageSeconds, err := l.GetCommitAgeSecondsByCommitHash(hash)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return ageSeconds
}

func (l *LocalGitRepository) MustGetCommitByGoGitCommit(goGitCommit *object.Commit) (gitCommit *GitCommit) {
	gitCommit, err := l.GetCommitByGoGitCommit(goGitCommit)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitCommit
}

func (l *LocalGitRepository) MustGetCommitByGoGitHash(goGitHash *plumbing.Hash) (gitCommit *GitCommit) {
	gitCommit, err := l.GetCommitByGoGitHash(goGitHash)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitCommit
}

func (l *LocalGitRepository) MustGetCommitByGoGitReference(goGitReference *plumbing.Reference) (gitCommit *GitCommit) {
	gitCommit, err := l.GetCommitByGoGitReference(goGitReference)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitCommit
}

func (l *LocalGitRepository) MustGetCommitMessageByCommitHash(hash string) (commitMessage string) {
	commitMessage, err := l.GetCommitMessageByCommitHash(hash)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return commitMessage
}

func (l *LocalGitRepository) MustGetCommitParentsByCommitHash(hash string, options *parameteroptions.GitCommitGetParentsOptions) (commitParents []*GitCommit) {
	commitParents, err := l.GetCommitParentsByCommitHash(hash, options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return commitParents
}

func (l *LocalGitRepository) MustGetCommitTimeByCommitHash(hash string) (commitTime *time.Time) {
	commitTime, err := l.GetCommitTimeByCommitHash(hash)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return commitTime
}

func (l *LocalGitRepository) MustGetCurrentBranchName(verbose bool) (branchName string) {
	branchName, err := l.GetCurrentBranchName(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return branchName
}

func (l *LocalGitRepository) MustGetCurrentCommit(verbose bool) (gitCommit *GitCommit) {
	gitCommit, err := l.GetCurrentCommit(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitCommit
}

func (l *LocalGitRepository) MustGetCurrentCommitGoGitHash(verbose bool) (hash *plumbing.Hash) {
	hash, err := l.GetCurrentCommitGoGitHash(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return hash
}

func (l *LocalGitRepository) MustGetCurrentCommitHash(verbose bool) (commitHash string) {
	commitHash, err := l.GetCurrentCommitHash(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return commitHash
}

func (l *LocalGitRepository) MustGetCurrentCommitHashAsBytes(verbose bool) (hash []byte) {
	hash, err := l.GetCurrentCommitHashAsBytes(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return hash
}

func (l *LocalGitRepository) MustGetDirectoryByPath(pathToSubDir ...string) (subDir files.Directory) {
	subDir, err := l.GetDirectoryByPath(pathToSubDir...)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return subDir
}

func (l *LocalGitRepository) MustGetGitStatusOutput(verbose bool) (output string) {
	output, err := l.GetGitStatusOutput(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return output
}

func (l *LocalGitRepository) MustGetGitlabCiYamlFile() (gitlabCiYamlFile *GitlabCiYamlFile) {
	gitlabCiYamlFile, err := l.GetGitlabCiYamlFile()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabCiYamlFile
}

func (l *LocalGitRepository) MustGetGoGitCommitByCommitHash(hash string) (goGitCommit *object.Commit) {
	goGitCommit, err := l.GetGoGitCommitByCommitHash(hash)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return goGitCommit
}

func (l *LocalGitRepository) MustGetGoGitConfig() (config *config.Config) {
	config, err := l.GetGoGitConfig()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return config
}

func (l *LocalGitRepository) MustGetGoGitHashFromHashString(hashString string) (hash *plumbing.Hash) {
	hash, err := l.GetGoGitHashFromHashString(hashString)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return hash
}

func (l *LocalGitRepository) MustGetGoGitHead() (head *plumbing.Reference) {
	head, err := l.GetGoGitHead()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return head
}

func (l *LocalGitRepository) MustGetGoGitWorktree() (worktree *git.Worktree) {
	worktree, err := l.GetGoGitWorktree()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return worktree
}

func (l *LocalGitRepository) MustGetHashByTagName(tagName string) (hash string) {
	hash, err := l.GetHashByTagName(tagName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return hash
}

func (l *LocalGitRepository) MustGetRemoteConfigs(verbose bool) (remoteConfigs []*GitRemoteConfig) {
	remoteConfigs, err := l.GetRemoteConfigs(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return remoteConfigs
}

func (l *LocalGitRepository) MustGetTagByName(tagName string) (tag GitTag) {
	tag, err := l.GetTagByName(tagName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return tag
}

func (l *LocalGitRepository) MustGitlabCiYamlFileExists(verbose bool) (gitlabCiYamlFileExists bool) {
	gitlabCiYamlFileExists, err := l.GitlabCiYamlFileExists(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabCiYamlFileExists
}

func (l *LocalGitRepository) MustHasInitialCommit(verbose bool) (hasInitialCommit bool) {
	hasInitialCommit, err := l.HasInitialCommit(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return hasInitialCommit
}

func (l *LocalGitRepository) MustHasNoUncommittedChanges(verbose bool) (hasUncommittedChanges bool) {
	hasUncommittedChanges, err := l.HasNoUncommittedChanges(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return hasUncommittedChanges
}

func (l *LocalGitRepository) MustHasUncommittedChanges(verbose bool) (hasUncommittedChanges bool) {
	hasUncommittedChanges, err := l.HasUncommittedChanges(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return hasUncommittedChanges
}

func (l *LocalGitRepository) MustInit(options *parameteroptions.CreateRepositoryOptions) {
	err := l.Init(options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (l *LocalGitRepository) MustIsGitRepository(verbose bool) (isGitRepository bool) {
	isGitRepository, err := l.IsGitRepository(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isGitRepository
}

func (l *LocalGitRepository) MustIsInitialized(verbose bool) (isInitialized bool) {
	isInitialized, err := l.IsInitialized(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isInitialized
}

func (l *LocalGitRepository) MustListBranchNames(verbose bool) (branchNames []string) {
	branchNames, err := l.ListBranchNames(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return branchNames
}

func (l *LocalGitRepository) MustListTagNames(verbose bool) (tagNames []string) {
	tagNames, err := l.ListTagNames(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return tagNames
}

func (l *LocalGitRepository) MustListTags(verbose bool) (tags []GitTag) {
	tags, err := l.ListTags(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return tags
}

func (l *LocalGitRepository) MustListTagsForCommitHash(hash string, verbose bool) (tags []GitTag) {
	tags, err := l.ListTagsForCommitHash(hash, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return tags
}

func (l *LocalGitRepository) MustPullFromRemote(pullOptions *GitPullFromRemoteOptions) {
	err := l.PullFromRemote(pullOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (l *LocalGitRepository) MustPullUsingGitCli(verbose bool) {
	err := l.PullUsingGitCli(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (l *LocalGitRepository) MustPushTagsToRemote(remoteName string, verbose bool) {
	err := l.PushTagsToRemote(remoteName, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (l *LocalGitRepository) MustPushToRemote(remoteName string, verbose bool) {
	err := l.PushToRemote(remoteName, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (l *LocalGitRepository) MustRemoteByNameExists(remoteName string, verbose bool) (remoteExists bool) {
	remoteExists, err := l.RemoteByNameExists(remoteName, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return remoteExists
}

func (l *LocalGitRepository) MustRemoteConfigurationExists(config *GitRemoteConfig, verbose bool) (exists bool) {
	exists, err := l.RemoteConfigurationExists(config, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return exists
}

func (l *LocalGitRepository) MustRemoveRemoteByName(remoteNameToRemove string, verbose bool) {
	err := l.RemoveRemoteByName(remoteNameToRemove, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (l *LocalGitRepository) MustRunGitCommandAndGetStdout(gitCommand []string, verbose bool) (commandOutput string) {
	commandOutput, err := l.RunGitCommandAndGetStdout(gitCommand, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return commandOutput
}

func (l *LocalGitRepository) MustSetGitConfig(options *gitparameteroptions.GitConfigSetOptions) {
	err := l.SetGitConfig(options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (l *LocalGitRepository) MustSetGitConfigByGoGitConfig(config *config.Config, verbose bool) {
	err := l.SetGitConfigByGoGitConfig(config, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (l *LocalGitRepository) MustSetRemote(remoteName string, remotUrl string, verbose bool) (remote *LocalGitRemote) {
	remote, err := l.SetRemote(remoteName, remotUrl, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return remote
}

func (l *LocalGitRepository) MustSetRemoteUrl(remoteUrl string, verbose bool) {
	err := l.SetRemoteUrl(remoteUrl, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (l *LocalGitRepository) Pull(ctx context.Context) (err error) {
	worktree, err := l.GetGoGitWorktree()
	if err != nil {
		return err
	}

	alreadyUpToDate := false
	err = worktree.Pull(&git.PullOptions{})
	if err != nil {
		if errors.Is(err, git.NoErrAlreadyUpToDate) {
			alreadyUpToDate = true
		} else {
			return tracederrors.TracedErrorf("%w", err)
		}
	}

	repoPath, err := l.GetPath()
	if err != nil {
		return err
	}

	if alreadyUpToDate {
		logging.LogInfoByCtxf(ctx, "Pull git repository '%s'. Already up to date.", repoPath)
	} else {
		logging.LogInfoByCtxf(ctx, "Pulled git repository '%s'.", repoPath)
	}

	return nil
}

func (l *LocalGitRepository) PullFromRemote(pullOptions *GitPullFromRemoteOptions) (err error) {
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

	path, hostDescription, err := l.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	// TODO implement without calling the git binary.
	_, err = l.RunGitCommand(
		contextutils.GetVerbosityContextByBool(pullOptions.Verbose),
		[]string{"pull", remoteName, branchName},
	)
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

func (l *LocalGitRepository) PullUsingGitCli(verbose bool) (err error) {
	_, err = l.RunGitCommand(
		contextutils.GetVerbosityContextByBool(verbose),
		[]string{"pull"},
	)
	if err != nil {
		return err
	}

	return nil
}

func (l *LocalGitRepository) Push(ctx context.Context) (err error) {
	goGitRepo, err := l.GetAsGoGitRepository()
	if err != nil {
		return err
	}

	err = goGitRepo.Push(&git.PushOptions{})
	alreadyUpToDate := false
	if err != nil {
		if errors.Is(err, git.NoErrAlreadyUpToDate) {
			alreadyUpToDate = true
		} else {
			return tracederrors.TracedErrorf("local git repository push failed: %w", err)
		}
	}

	repoPath, err := l.GetPath()
	if err != nil {
		return err
	}

	if alreadyUpToDate {
		logging.LogInfoByCtxf(ctx, "Push git repository '%s'. Already up to date.", repoPath)
	} else {
		logging.LogInfoByCtxf(ctx, "Pushed git repository '%s'.", repoPath)
	}

	return nil
}

func (l *LocalGitRepository) PushTagsToRemote(remoteName string, verbose bool) (err error) {
	if len(remoteName) <= 0 {
		return tracederrors.TracedError("remoteName is empty string")
	}

	path, hostDescription, err := l.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	// TODO: Implemnet without calling git binary
	_, err = l.RunGitCommand(
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

func (l *LocalGitRepository) PushToRemote(remoteName string, verbose bool) (err error) {
	if len(remoteName) <= 0 {
		return tracederrors.TracedError("remoteName is empty string")
	}

	// TODO: Implement without calling git binary
	_, err = l.RunGitCommand(
		contextutils.GetVerbosityContextByBool(verbose),
		[]string{"push", remoteName},
	)
	if err != nil {
		return err
	}

	if verbose {
		path, hostDescription, err := l.GetPathAndHostDescription()
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

func (l *LocalGitRepository) RemoteByNameExists(remoteName string, verbose bool) (remoteExists bool, err error) {
	if len(remoteName) <= 0 {
		return false, fmt.Errorf("remoteName is empty string")
	}

	remoteConfigs, err := l.GetRemoteConfigs(verbose)
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

func (l *LocalGitRepository) RemoteConfigurationExists(config *GitRemoteConfig, verbose bool) (exists bool, err error) {
	if config == nil {
		return false, tracederrors.TracedError("config is nil")
	}

	remoteConfigs, err := l.GetRemoteConfigs(verbose)
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

func (l *LocalGitRepository) RemoveRemoteByName(remoteNameToRemove string, verbose bool) (err error) {
	if len(remoteNameToRemove) <= 0 {
		return tracederrors.TracedError("remoteNameToRemove is empty string")
	}

	remoteExists, err := l.RemoteByNameExists(remoteNameToRemove, verbose)
	if err != nil {
		return err
	}

	repoDirPath, err := l.GetPath()
	if err != nil {
		return err
	}

	if remoteExists {
		// TODO: reimplement without calling the git binary.
		_, err := l.RunGitCommand(
			contextutils.GetVerbosityContextByBool(verbose),
			[]string{"remote", "remove", remoteNameToRemove},
		)
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

func (l *LocalGitRepository) RunGitCommandAndGetStdout(gitCommand []string, verbose bool) (commandOutput string, err error) {
	if len(gitCommand) <= 0 {
		return "", tracederrors.TracedError("gitCommand is empty")
	}

	output, err := l.RunGitCommand(
		contextutils.GetVerbosityContextByBool(verbose),
		gitCommand,
	)
	if err != nil {
		return "", err
	}

	commandOutput, err = output.GetStdoutAsString()
	if err != nil {
		return "", err
	}

	return commandOutput, nil
}

func (l *LocalGitRepository) SetGitConfig(options *gitparameteroptions.GitConfigSetOptions) (err error) {
	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	repoPath, err := l.GetLocalPath()
	if err != nil {
		return err
	}

	config, err := l.GetGoGitConfig()
	if err != nil {
		return err
	}

	rewriteNeeded := false
	if options.IsEmailSet() {
		email, err := options.GetEmail()
		if err != nil {
			return err
		}

		if config.Author.Email == email {
			logging.LogInfof("Email in git config of local repository '%s' is already '%s'.", repoPath, email)
		} else {
			config.Author.Email = email
			rewriteNeeded = true
			logging.LogChangedf("Set email in git config of local repository '%s' to '%s'.", repoPath, email)
		}
	}

	if options.IsNameSet() {
		name, err := options.GetName()
		if err != nil {
			return err
		}

		if config.Author.Name == name {
			logging.LogInfof("Author name in git config of local repository '%s' is already '%s'.", repoPath, name)
		} else {
			config.Author.Name = name
			rewriteNeeded = true
			logging.LogChangedf("Set author name in git config of local repository '%s' to '%s'.", repoPath, name)
		}
	}

	if rewriteNeeded {
		err = l.SetGitConfigByGoGitConfig(config, options.Verbose)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *LocalGitRepository) SetGitConfigByGoGitConfig(config *config.Config, verbose bool) (err error) {
	if config == nil {
		return tracederrors.TracedErrorNil("config")
	}

	outFile, err := l.GetFileInDirectory(".git", "config")
	if err != nil {
		return err
	}

	configData, err := config.Marshal()
	if err != nil {
		return tracederrors.TracedErrorf("%w", err)
	}

	const verboseWrite bool = false
	err = outFile.WriteBytes(configData, verboseWrite)
	if err != nil {
		return err
	}

	path, err := l.GetLocalPath()
	if err != nil {
		return err
	}

	if verbose {
		logging.LogChangedf("Wrote git config of local git repository '%s'.", path)
	}

	return nil
}

func (l *LocalGitRepository) SetRemote(remoteName string, remotUrl string, verbose bool) (remote *LocalGitRemote, err error) {
	if remoteName == "" {
		return nil, tracederrors.TracedErrorEmptyString("remoteName")
	}

	if remotUrl == "" {
		return nil, tracederrors.TracedErrorEmptyString("remotUrl")
	}

	goGitRepo, err := l.GetAsGoGitRepository()
	if err != nil {
		return nil, err
	}

	nativeRemote, err := goGitRepo.CreateRemote(&config.RemoteConfig{
		Name: remoteName,
		URLs: []string{remotUrl},
	})
	if err != nil {
		return nil, tracederrors.TracedErrorf("Create remote failed: '%w'", err)
	}

	remote, err = NewLocalGitRemoteByNativeGoGitRemote(nativeRemote)
	if err != nil {
		return nil, err
	}

	if verbose {
		localPath, err := l.GetLocalPath()
		if err != nil {
			return nil, err
		}

		logging.LogInfof(
			"Set remote '%s' with remote URL '%s' to local Git repository '%s'.",
			remoteName,
			remotUrl,
			localPath,
		)
	}

	return remote, err
}

func (l *LocalGitRepository) SetRemoteUrl(remoteUrl string, verbose bool) (err error) {
	remoteUrl = strings.TrimSpace(remoteUrl)
	if len(remoteUrl) <= 0 {
		return tracederrors.TracedError("remoteUrl is empty string")
	}

	name := "origin"

	// TODO: Implement without calling the git binary
	_, err = l.RunGitCommand(
		contextutils.GetVerbosityContextByBool(verbose),
		[]string{"remote", "set-url", name, remoteUrl},
	)
	if err != nil {
		return err
	}

	if verbose {
		path, hostDescription, err := l.GetPathAndHostDescription()
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
