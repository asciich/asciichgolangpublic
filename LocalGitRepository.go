package asciichgolangpublic

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type LocalGitRepository struct {
	LocalDirectory
	GitRepositoryBase
}

func GetLocalGitReposioryFromDirectory(directory Directory) (repo GitRepository, err error) {
	if directory == nil {
		return nil, TracedErrorNil("directory")
	}

	isLocalDirectory, err := directory.IsLocalDirectory()
	if err != nil {
		return nil, err
	}

	if !isLocalDirectory {
		return nil, TracedError("Only local directories are supported.")
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

func GetLocalGitReposioryFromLocalDirectory(localDirectory *LocalDirectory) (l *LocalGitRepository, err error) {
	if localDirectory == nil {
		return nil, TracedErrorNil("directory")
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
		return nil, TracedErrorEmptyString("path")
	}

	l = NewLocalGitRepository()

	err = l.SetLocalPath(path)
	if err != nil {
		return nil, err
	}

	return l, nil
}

func MustGetLocalGitReposioryFromDirectory(directory Directory) (repo GitRepository) {
	repo, err := GetLocalGitReposioryFromDirectory(directory)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return repo
}

func MustGetLocalGitReposioryFromLocalDirectory(localDirectory *LocalDirectory) (l *LocalGitRepository) {
	l, err := GetLocalGitReposioryFromLocalDirectory(localDirectory)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return l
}

func MustGetLocalGitRepositoryByPath(path string) (l *LocalGitRepository) {
	l, err := GetLocalGitRepositoryByPath(path)
	if err != nil {
		LogGoErrorFatal(err)
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
		return TracedErrorNil("path")
	}

	worktree, err := l.GetGoGitWorktree()
	if err != nil {
		return err
	}

	_, err = worktree.Add(pathToAdd)
	if err != nil {
		return TracedErrorf("%w", err)
	}

	if verbose {
		path, err := l.GetPath()
		if err != nil {
			return err
		}

		LogChangedf(
			"Added file '%s' to git repository '%s' on localhost",
			pathToAdd,
			path,
		)
	}

	return nil
}

func (l *LocalGitRepository) CloneRepository(repository GitRepository, verbose bool) (err error) {
	if repository == nil {
		return TracedErrorNil("repository")
	}

	repositoryHostDescription, err := repository.GetHostDescription()
	if err != nil {
		return err
	}

	if repositoryHostDescription != "localhost" {
		return TracedErrorf(
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
		return TracedErrorEmptyString("pathToClone")
	}

	path, hostDescription, err := l.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	if verbose {
		LogInfof(
			"Cloning git repository '%s' to '%s' on host '%s' started.",
			urlOrPathToClone,
			path,
			hostDescription,
		)
	}

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
				LogInfof(
					"Remote repository '%s' is empty. Going to add remote for empty repository.",
					urlOrPathToClone,
				)
			}

			err = l.Init(
				&CreateRepositoryOptions{
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
			return TracedErrorf(
				"Clone '%s' failed: '%w'",
				urlOrPathToClone,
				err,
			)
		}
	}

	if verbose {
		LogInfof(
			"Cloning git repository '%s' to '%s' on host '%s' finished.",
			urlOrPathToClone,
			path,
			hostDescription,
		)
	}

	return nil
}

func (l *LocalGitRepository) Commit(commitOptions *GitCommitOptions) (createdCommit *GitCommit, err error) {
	if commitOptions == nil {
		return nil, TracedErrorNil("commitOptions")
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
		LogChangedf(
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
		return false, TracedErrorEmptyString("hash")
	}

	goGitCommit, err := l.GetGoGitCommitByCommitHash(hash)
	if err != nil {
		return false, err
	}

	hasParentCommit = goGitCommit.NumParents() > 0

	return hasParentCommit, nil
}

func (l *LocalGitRepository) CreateTag(options *GitRepositoryCreateTagOptions) (createdTag GitTag, err error) {
	if options == nil {
		return nil, TracedErrorNil("options")
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
		LogInfof(
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
		return nil, TracedErrorf(
			"Creating tag failed: %w",
			err,
		)
	}

	if options.Verbose {

		LogChangedf(
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

func (l *LocalGitRepository) FileByPathExists(path string, verbose bool) (exists bool, err error) {
	if path == "" {
		return false, TracedErrorEmptyString(path)
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
			toReturn := TracedErrorf("%w: repoPath='%s'", ErrGitRepositoryDoesNotExist, repoPath)
			return nil, Errors().AddErrorToUnwrapToTracedError(toReturn, err)
		}
		return nil, TracedErrorf("%w: repoPath='%s'", err, repoPath)
	}

	if goGitRepository == nil {
		return nil, TracedError("goGitRepository is nil after evaluation.")
	}

	return goGitRepository, nil
}

func (l *LocalGitRepository) GetAsLocalDirectory() (localDirectory *LocalDirectory, err error) {
	localPath, err := l.GetLocalPath()
	if err != nil {
		return nil, err
	}

	localDirectory, err = GetLocalDirectoryByPath(localPath)
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
		return "", TracedErrorEmptyString("hash")
	}

	goGitCommit, err := l.GetGoGitCommitByCommitHash(hash)
	if err != nil {
		return "", err
	}

	authorEmail = goGitCommit.Author.Email
	if err != nil {
		return "", err
	}

	return authorEmail, nil
}

func (l *LocalGitRepository) GetAuthorStringByCommitHash(hash string) (authorString string, err error) {
	if hash == "" {
		return "", TracedErrorEmptyString(hash)
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
		return nil, TracedErrorEmptyString("hash")
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
		return -1, TracedErrorEmptyString("hash")
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
		return nil, TracedErrorNil("goGitCommit")
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
		return nil, TracedErrorNil("goGitHash")
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
		return nil, TracedErrorf("%w", err)
	}

	err = gitCommit.SetHash(commitObject.Hash.String())
	if err != nil {
		return nil, err
	}

	return gitCommit, nil
}

func (l *LocalGitRepository) GetCommitByGoGitReference(goGitReference *plumbing.Reference) (gitCommit *GitCommit, err error) {
	if goGitReference == nil {
		return nil, TracedErrorNil("goGitReference")
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
		return "", TracedErrorEmptyString("hash")
	}

	g, err := l.GetGoGitCommitByCommitHash(hash)
	if err != nil {
		return "", err
	}

	commitMessage = g.Message

	return commitMessage, nil
}

func (l *LocalGitRepository) GetCommitParentsByCommitHash(hash string, options *GitCommitGetParentsOptions) (commitParents []*GitCommit, err error) {
	if hash == "" {
		return nil, TracedErrorEmptyString("hash")
	}

	if options == nil {
		return nil, TracedErrorNil("options")
	}

	goGitCommit, err := l.GetGoGitCommitByCommitHash(hash)
	if err != nil {
		return nil, err
	}

	parents := goGitCommit.Parents()
	for {
		parentToAdd, err := parents.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, TracedErrorf("Unable to get next parent: %w", err)
		}

		toAdd, err := l.GetCommitByGoGitCommit(parentToAdd)
		if err != nil {
			return nil, err
		}

		commitParents = append(commitParents, toAdd)

		if options.IncludeParentsOfParents {
			additionalParents, err := toAdd.GetParentCommits(&GitCommitGetParentsOptions{
				IncludeParentsOfParents: true,
			})
			if err != nil {
				return nil, err
			}

			commitParents = append(commitParents, additionalParents...)
		}
	}

	if options.Verbose {
		LogInfof("Collected '%d' parent commits for commit '%s'.", len(commitParents), hash)
	}

	return commitParents, nil
}

func (l *LocalGitRepository) GetCommitTimeByCommitHash(hash string) (commitTime *time.Time, err error) {
	if hash == "" {
		return nil, TracedErrorEmptyString("hash")
	}

	goGitCommit, err := l.GetGoGitCommitByCommitHash(hash)
	if err != nil {
		return nil, err
	}

	commitTime = &(goGitCommit.Author.When)

	return commitTime, nil
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

		LogInfof(
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

	return Strings().HexStringToBytes(currentHash)
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
		return nil, TracedErrorEmptyString("hash")
	}

	goGitRepo, err := l.GetAsGoGitRepository()
	if err != nil {
		return nil, err
	}

	pHash := plumbing.NewHash(hash)

	goGitCommit, err = goGitRepo.CommitObject(pHash)
	if err != nil {
		return nil, TracedErrorf("%w", err)
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
		return nil, TracedErrorf("%w", err)
	}

	if config == nil {
		return nil, TracedError("config is nil after evaluation")
	}

	return config, nil
}

func (l *LocalGitRepository) GetGoGitHashFromHashString(hashString string) (hash *plumbing.Hash, err error) {
	if hashString == "" {
		return nil, TracedErrorNil("hashString")
	}

	hashBytes, err := Strings().HexStringToBytes(hashString)
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
			toReturn := TracedErrorf("%w", ErrGitRepositoryHeadNotFound)
			return nil, Errors().AddErrorToUnwrapToTracedError(toReturn, err)
		}
		return nil, TracedErrorf("%w", err)
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
		return nil, TracedErrorf("%w", err)
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
		return "", TracedErrorf(
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

			return "", TracedErrorf(
				"Unable to get next tag: '%w'",
				err,
			)
		}

		name := tag.Name
		name = strings.TrimPrefix(name, "refs/tags/")

		if tagName == name {
			hash = string(tag.Target.String())

			if hash == "" {
				return "", TracedError(
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

	return "", TracedErrorf(
		"Unable to get hash for tag '%s' in local git repository '%s'.",
		tagName,
		path,
	)
}

func (l *LocalGitRepository) GetRootDirectory(verbose bool) (rootDirectory Directory, err error) {
	rootDirectoryPath, err := l.GetRootDirectoryPath(verbose)
	if err != nil {
		return nil, err
	}

	rootDirectory, err = GetLocalDirectoryByPath(rootDirectoryPath)
	if err != nil {
		return nil, err
	}

	return rootDirectory, nil
}

func (l *LocalGitRepository) GetRootDirectoryPath(verbose bool) (rootDirectoryPath string, err error) {
	pathToCheck, err := l.GetLocalPath()
	if err != nil {
		return "", err
	}

	searchedFromPath := pathToCheck

	for {
		localDirToCheck, err := GetLocalDirectoryByPath(pathToCheck)
		if err != nil {
			return "", nil
		}

		localPathToCheck, err := localDirToCheck.GetLocalPath()
		if err != nil {
			return "", nil
		}

		if localPathToCheck == "" || localPathToCheck == "/" {
			return "", TracedErrorf("Not inside a git repository. Searched from '%s'", searchedFromPath)
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
		return nil, TracedErrorEmptyString("tagName")
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
			LogInfof(
				"Local git repository '%s' has uncommited changes.",
				path,
			)
		} else {
			LogInfof(
				"Local git repository '%s' has no uncommited changes.",
				path,
			)
		}
	}

	return hasUncommittedChanges, nil
}

func (l *LocalGitRepository) Init(options *CreateRepositoryOptions) (err error) {
	if options == nil {
		return TracedErrorNil("options")
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
			LogInfof("Local git repository '%s' is already initialized.", repoPath)
		}
	} else {
		_, err = git.PlainInit(repoPath, options.BareRepository)
		if err != nil {
			return TracedErrorf("%w", err)
		}
		if options.Verbose {
			LogChangedf("Local git repository '%s' is initialized.", repoPath)
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
					&GitConfigSetOptions{
						Name:    "asciichgolangpublic git repo initializer",
						Email:   "asciichgolangpublic@example.net",
						Verbose: options.Verbose,
					},
				)
				if err != nil {
					return err
				}

				_, err = temporaryRepository.CommitAndPush(
					&GitCommitOptions{
						Message:    "Initial empty commit during repo initialization",
						AllowEmpty: true,
						Verbose:    true,
					},
				)
				if err != nil {
					return err
				}

				if options.Verbose {
					LogChangedf("Initialized bare repository '%s' with an empty commit.", repoPath)
				}
			} else {
				if options.InitializeWithDefaultAuthor {
					err = l.SetGitConfig(
						&GitConfigSetOptions{
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
				&GitConfigSetOptions{
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

func (l *LocalGitRepository) IsBareRepository(verbose bool) (isBareRepository bool, err error) {

	config, err := l.GetGoGitConfig()
	if err != nil {
		return false, err
	}

	isBareRepository = config.Core.IsBare

	if verbose {
		repoRoot, err := l.GetPath()
		if err != nil {
			return false, err
		}

		if isBareRepository {
			LogInfof("Git repository '%s' is a bare repository.", repoRoot)
		} else {
			LogInfof("Git repository '%s' is not a bare repository.", repoRoot)
		}
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
			LogInfof(
				"Directory '%s' is an initialized git repository.",
				path,
			)
		} else {
			LogInfof(
				"Directory '%s' is not an initialized git repository.",
				path,
			)
		}
	}

	return isInitialized, nil
}

func (l *LocalGitRepository) ListTagNames(verbose bool) (tagNames []string, err error) {
	nativeRepo, err := l.GetAsGoGitRepository()
	if err != nil {
		return nil, err
	}

	tags, err := nativeRepo.Tags()
	if err != nil {
		return nil, TracedErrorf(
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

			return nil, TracedErrorf(
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
		return nil, TracedErrorEmptyString("hash")
	}

	nativeRepo, err := l.GetAsGoGitRepository()
	if err != nil {
		return nil, err
	}

	nativeTagObjects, err := nativeRepo.TagObjects()
	if err != nil {
		return nil, TracedErrorf(
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

			return nil, TracedErrorf(
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
		LogGoErrorFatal(err)
	}
}

func (l *LocalGitRepository) MustCloneRepository(repository GitRepository, verbose bool) {
	err := l.CloneRepository(repository, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (l *LocalGitRepository) MustCloneRepositoryByPathOrUrl(pathToClone string, verbose bool) {
	err := l.CloneRepositoryByPathOrUrl(pathToClone, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (l *LocalGitRepository) MustCommit(commitOptions *GitCommitOptions) (createdCommit *GitCommit) {
	createdCommit, err := l.Commit(commitOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return createdCommit
}

func (l *LocalGitRepository) MustCommitHasParentCommitByCommitHash(hash string) (hasParentCommit bool) {
	hasParentCommit, err := l.CommitHasParentCommitByCommitHash(hash)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return hasParentCommit
}

func (l *LocalGitRepository) MustCreateTag(options *GitRepositoryCreateTagOptions) (createdTag GitTag) {
	createdTag, err := l.CreateTag(options)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return createdTag
}

func (l *LocalGitRepository) MustFileByPathExists(path string, verbose bool) (exists bool) {
	exists, err := l.FileByPathExists(path, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return exists
}

func (l *LocalGitRepository) MustGetAsGoGitRepository() (goGitRepository *git.Repository) {
	goGitRepository, err := l.GetAsGoGitRepository()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return goGitRepository
}

func (l *LocalGitRepository) MustGetAsLocalDirectory() (localDirectory *LocalDirectory) {
	localDirectory, err := l.GetAsLocalDirectory()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return localDirectory
}

func (l *LocalGitRepository) MustGetAsLocalGitRepository() (localGitRepository *LocalGitRepository) {
	localGitRepository, err := l.GetAsLocalGitRepository()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return localGitRepository
}

func (l *LocalGitRepository) MustGetAuthorEmailByCommitHash(hash string) (authorEmail string) {
	authorEmail, err := l.GetAuthorEmailByCommitHash(hash)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return authorEmail
}

func (l *LocalGitRepository) MustGetAuthorStringByCommitHash(hash string) (authorString string) {
	authorString, err := l.GetAuthorStringByCommitHash(hash)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return authorString
}

func (l *LocalGitRepository) MustGetCommitAgeDurationByCommitHash(hash string) (ageDuration *time.Duration) {
	ageDuration, err := l.GetCommitAgeDurationByCommitHash(hash)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return ageDuration
}

func (l *LocalGitRepository) MustGetCommitAgeSecondsByCommitHash(hash string) (ageSeconds float64) {
	ageSeconds, err := l.GetCommitAgeSecondsByCommitHash(hash)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return ageSeconds
}

func (l *LocalGitRepository) MustGetCommitByGoGitCommit(goGitCommit *object.Commit) (gitCommit *GitCommit) {
	gitCommit, err := l.GetCommitByGoGitCommit(goGitCommit)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitCommit
}

func (l *LocalGitRepository) MustGetCommitByGoGitHash(goGitHash *plumbing.Hash) (gitCommit *GitCommit) {
	gitCommit, err := l.GetCommitByGoGitHash(goGitHash)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitCommit
}

func (l *LocalGitRepository) MustGetCommitByGoGitReference(goGitReference *plumbing.Reference) (gitCommit *GitCommit) {
	gitCommit, err := l.GetCommitByGoGitReference(goGitReference)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitCommit
}

func (l *LocalGitRepository) MustGetCommitMessageByCommitHash(hash string) (commitMessage string) {
	commitMessage, err := l.GetCommitMessageByCommitHash(hash)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return commitMessage
}

func (l *LocalGitRepository) MustGetCommitParentsByCommitHash(hash string, options *GitCommitGetParentsOptions) (commitParents []*GitCommit) {
	commitParents, err := l.GetCommitParentsByCommitHash(hash, options)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return commitParents
}

func (l *LocalGitRepository) MustGetCommitTimeByCommitHash(hash string) (commitTime *time.Time) {
	commitTime, err := l.GetCommitTimeByCommitHash(hash)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return commitTime
}

func (l *LocalGitRepository) MustGetCurrentCommit(verbose bool) (gitCommit *GitCommit) {
	gitCommit, err := l.GetCurrentCommit(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitCommit
}

func (l *LocalGitRepository) MustGetCurrentCommitGoGitHash(verbose bool) (hash *plumbing.Hash) {
	hash, err := l.GetCurrentCommitGoGitHash(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return hash
}

func (l *LocalGitRepository) MustGetCurrentCommitHash(verbose bool) (commitHash string) {
	commitHash, err := l.GetCurrentCommitHash(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return commitHash
}

func (l *LocalGitRepository) MustGetCurrentCommitHashAsBytes(verbose bool) (hash []byte) {
	hash, err := l.GetCurrentCommitHashAsBytes(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return hash
}

func (l *LocalGitRepository) MustGetGitStatusOutput(verbose bool) (output string) {
	output, err := l.GetGitStatusOutput(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return output
}

func (l *LocalGitRepository) MustGetGitlabCiYamlFile() (gitlabCiYamlFile *GitlabCiYamlFile) {
	gitlabCiYamlFile, err := l.GetGitlabCiYamlFile()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabCiYamlFile
}

func (l *LocalGitRepository) MustGetGoGitCommitByCommitHash(hash string) (goGitCommit *object.Commit) {
	goGitCommit, err := l.GetGoGitCommitByCommitHash(hash)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return goGitCommit
}

func (l *LocalGitRepository) MustGetGoGitConfig() (config *config.Config) {
	config, err := l.GetGoGitConfig()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return config
}

func (l *LocalGitRepository) MustGetGoGitHashFromHashString(hashString string) (hash *plumbing.Hash) {
	hash, err := l.GetGoGitHashFromHashString(hashString)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return hash
}

func (l *LocalGitRepository) MustGetGoGitHead() (head *plumbing.Reference) {
	head, err := l.GetGoGitHead()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return head
}

func (l *LocalGitRepository) MustGetGoGitWorktree() (worktree *git.Worktree) {
	worktree, err := l.GetGoGitWorktree()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return worktree
}

func (l *LocalGitRepository) MustGetHashByTagName(tagName string) (hash string) {
	hash, err := l.GetHashByTagName(tagName)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return hash
}

func (l *LocalGitRepository) MustGetRootDirectory(verbose bool) (rootDirectory Directory) {
	rootDirectory, err := l.GetRootDirectory(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return rootDirectory
}

func (l *LocalGitRepository) MustGetRootDirectoryPath(verbose bool) (rootDirectoryPath string) {
	rootDirectoryPath, err := l.GetRootDirectoryPath(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return rootDirectoryPath
}

func (l *LocalGitRepository) MustGetTagByName(tagName string) (tag GitTag) {
	tag, err := l.GetTagByName(tagName)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return tag
}

func (l *LocalGitRepository) MustGitlabCiYamlFileExists(verbose bool) (gitlabCiYamlFileExists bool) {
	gitlabCiYamlFileExists, err := l.GitlabCiYamlFileExists(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabCiYamlFileExists
}

func (l *LocalGitRepository) MustHasInitialCommit(verbose bool) (hasInitialCommit bool) {
	hasInitialCommit, err := l.HasInitialCommit(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return hasInitialCommit
}

func (l *LocalGitRepository) MustHasNoUncommittedChanges(verbose bool) (hasUncommittedChanges bool) {
	hasUncommittedChanges, err := l.HasNoUncommittedChanges(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return hasUncommittedChanges
}

func (l *LocalGitRepository) MustHasUncommittedChanges(verbose bool) (hasUncommittedChanges bool) {
	hasUncommittedChanges, err := l.HasUncommittedChanges(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return hasUncommittedChanges
}

func (l *LocalGitRepository) MustInit(options *CreateRepositoryOptions) {
	err := l.Init(options)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (l *LocalGitRepository) MustIsBareRepository(verbose bool) (isBareRepository bool) {
	isBareRepository, err := l.IsBareRepository(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return isBareRepository
}

func (l *LocalGitRepository) MustIsGitRepository(verbose bool) (isGitRepository bool) {
	isGitRepository, err := l.IsGitRepository(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return isGitRepository
}

func (l *LocalGitRepository) MustIsInitialized(verbose bool) (isInitialized bool) {
	isInitialized, err := l.IsInitialized(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return isInitialized
}

func (l *LocalGitRepository) MustListTagNames(verbose bool) (tagNames []string) {
	tagNames, err := l.ListTagNames(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return tagNames
}

func (l *LocalGitRepository) MustListTags(verbose bool) (tags []GitTag) {
	tags, err := l.ListTags(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return tags
}

func (l *LocalGitRepository) MustListTagsForCommitHash(hash string, verbose bool) (tags []GitTag) {
	tags, err := l.ListTagsForCommitHash(hash, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return tags
}

func (l *LocalGitRepository) MustPull(verbose bool) {
	err := l.Pull(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (l *LocalGitRepository) MustPullUsingGitCli(verbose bool) {
	err := l.PullUsingGitCli(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (l *LocalGitRepository) MustPush(verbose bool) {
	err := l.Push(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (l *LocalGitRepository) MustRunGitCommand(gitCommand []string, verbose bool) (commandOutput *CommandOutput) {
	commandOutput, err := l.RunGitCommand(gitCommand, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return commandOutput
}

func (l *LocalGitRepository) MustRunGitCommandAndGetStdout(gitCommand []string, verbose bool) (commandOutput string) {
	commandOutput, err := l.RunGitCommandAndGetStdout(gitCommand, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return commandOutput
}

func (l *LocalGitRepository) MustSetGitConfig(options *GitConfigSetOptions) {
	err := l.SetGitConfig(options)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (l *LocalGitRepository) MustSetGitConfigByGoGitConfig(config *config.Config, verbose bool) {
	err := l.SetGitConfigByGoGitConfig(config, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (l *LocalGitRepository) MustSetRemote(remoteName string, remotUrl string, verbose bool) (remote *LocalGitRemote) {
	remote, err := l.SetRemote(remoteName, remotUrl, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return remote
}

func (l *LocalGitRepository) Pull(verbose bool) (err error) {
	worktree, err := l.GetGoGitWorktree()
	if err != nil {
		return err
	}

	err = worktree.Pull(&git.PullOptions{})
	if err != nil {
		return TracedErrorf("%w", err)
	}

	return nil
}

func (l *LocalGitRepository) PullUsingGitCli(verbose bool) (err error) {
	_, err = l.RunGitCommand([]string{"pull"}, verbose)
	if err != nil {
		return err
	}

	return nil
}

func (l *LocalGitRepository) Push(verbose bool) (err error) {
	goGitRepo, err := l.GetAsGoGitRepository()
	if err != nil {
		return err
	}

	err = goGitRepo.Push(&git.PushOptions{})
	if err != nil {
		return TracedErrorf("%w", err)
	}

	return nil
}

func (l *LocalGitRepository) RunGitCommand(gitCommand []string, verbose bool) (commandOutput *CommandOutput, err error) {
	if gitCommand == nil {
		return nil, TracedErrorEmptyString("gitCommand")
	}

	repoRootPath, err := l.GetRootDirectoryPath(verbose)
	if err != nil {
		return nil, err
	}

	gitCommandString, err := ShellLineHandler().Join(gitCommand)
	if err != nil {
		return nil, err
	}

	command := fmt.Sprintf(
		"git -C '%s' %s",
		repoRootPath,
		gitCommandString,
	)

	commandOutput, err = Bash().RunOneLiner(command, verbose)
	if err != nil {
		return nil, err
	}

	return commandOutput, nil
}

func (l *LocalGitRepository) RunGitCommandAndGetStdout(gitCommand []string, verbose bool) (commandOutput string, err error) {
	if len(gitCommand) <= 0 {
		return "", TracedError("gitCommand is empty")
	}

	output, err := l.RunGitCommand(gitCommand, verbose)
	if err != nil {
		return "", err
	}

	commandOutput, err = output.GetStdoutAsString()
	if err != nil {
		return "", err
	}

	return commandOutput, nil
}

func (l *LocalGitRepository) SetGitConfig(options *GitConfigSetOptions) (err error) {
	if options == nil {
		return TracedErrorNil("options")
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
			LogInfof("Email in git config of local repository '%s' is already '%s'.", repoPath, email)
		} else {
			config.Author.Email = email
			rewriteNeeded = true
			LogChangedf("Set email in git config of local repository '%s' to '%s'.", repoPath, email)
		}
	}

	if options.IsNameSet() {
		name, err := options.GetName()
		if err != nil {
			return err
		}

		if config.Author.Name == name {
			LogInfof("Author name in git config of local repository '%s' is already '%s'.", repoPath, name)
		} else {
			config.Author.Name = name
			rewriteNeeded = true
			LogChangedf("Set author name in git config of local repository '%s' to '%s'.", repoPath, name)
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
		return TracedErrorNil("config")
	}

	outFile, err := l.GetFileInDirectory(".git", "config")
	if err != nil {
		return err
	}

	configData, err := config.Marshal()
	if err != nil {
		return TracedErrorf("%w", err)
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
		LogChangedf("Wrote git config of local git repository '%s'.", path)
	}

	return nil
}

func (l *LocalGitRepository) SetRemote(remoteName string, remotUrl string, verbose bool) (remote *LocalGitRemote, err error) {
	if remoteName == "" {
		return nil, TracedErrorEmptyString("remoteName")
	}

	if remotUrl == "" {
		return nil, TracedErrorEmptyString("remotUrl")
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
		return nil, TracedErrorf("Create remote failed: '%w'", err)
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

		LogInfof(
			"Set remote '%s' with remote URL '%s' to local Git repository '%s'.",
			remoteName,
			remotUrl,
			localPath,
		)
	}

	return remote, err
}
