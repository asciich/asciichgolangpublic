package asciichgolangpublic

import (
	astrings "github.com/asciich/asciichgolangpublic/datatypes/strings"
)

type GitignoreFile struct {
	File
}

func GetGitignoreDefaultBaseName() (defaultBaseName string) {
	return ".gitignore"
}

func GetGitignoreFileByFile(fileToUse File) (gitignoreFile *GitignoreFile, err error) {
	if fileToUse == nil {
		return nil, TracedErrorEmptyString("fileToUse")
	}

	gitignoreFile = NewGitignoreFile()

	gitignoreFile.File = fileToUse

	return gitignoreFile, nil
}

func GetGitignoreFileByPath(filePath string) (gitignoreFile *GitignoreFile, err error) {
	if filePath == "" {
		return nil, TracedErrorEmptyString("filePath")
	}

	fileToUse, err := GetLocalFileByPath(filePath)
	if err != nil {
		return nil, err
	}

	return GetGitignoreFileByFile(fileToUse)
}

func GetGitignoreFileInGitRepository(gitRepository GitRepository) (gitignoreFile *GitignoreFile, err error) {
	if gitRepository == nil {
		return nil, TracedErrorNil("gitRepository")
	}

	fileToUse, err := gitRepository.GetFileByPath(GetGitignoreDefaultBaseName())
	if err != nil {
		return nil, err
	}

	return GetGitignoreFileByFile(fileToUse)
}

func MustGetGitignoreFileByFile(fileToUse File) (gitignoreFile *GitignoreFile) {
	gitignoreFile, err := GetGitignoreFileByFile(fileToUse)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitignoreFile
}

func MustGetGitignoreFileByPath(filePath string) (gitignoreFile *GitignoreFile) {
	gitignoreFile, err := GetGitignoreFileByPath(filePath)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitignoreFile
}

func MustGetGitignoreFileInGitRepository(gitRepository GitRepository) (gitignoreFile *GitignoreFile) {
	gitignoreFile, err := GetGitignoreFileInGitRepository(gitRepository)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitignoreFile
}

func NewGitignoreFile() (g *GitignoreFile) {
	return new(GitignoreFile)
}

func (g *GitignoreFile) AddDirToIgnore(pathToIgnore string, comment string, verbose bool) (err error) {
	if pathToIgnore == "" {
		return TracedError("pathToIgnore is empty string")
	}

	if comment == "" {
		return TracedError("comment is empty string")
	}

	pathToIgnore = astrings.EnsureSuffix(pathToIgnore, "/")

	err = g.Create(verbose)
	if err != nil {
		return err
	}

	containsIgnore, err := g.ContainsIgnore(pathToIgnore)
	if err != nil {
		return err
	}

	path, err := g.GetPath()
	if err != nil {
		return err
	}

	if containsIgnore {
		if verbose {
			LogInfof(
				"Gitignore file '%s' already contains ignore entry for '%s'.",
				path,
				pathToIgnore,
			)
		}
		return nil
	}

	err = g.AppendLine("# "+comment, verbose)
	if err != nil {
		return err
	}

	err = g.AppendLine(pathToIgnore, verbose)
	if err != nil {
		return err
	}

	if verbose {
		LogChangedf(
			"Added '%s' to gitignore file '%s'.",
			pathToIgnore,
			path,
		)
	}

	return nil
}

func (g *GitignoreFile) AddFileToIgnore(pathToIgnore string, comment string, verbose bool) (err error) {
	if pathToIgnore == "" {
		return TracedError("pathToIgnore is empty string")
	}

	if comment == "" {
		return TracedError("comment is empty string")
	}

	err = g.Create(verbose)
	if err != nil {
		return err
	}

	containsIgnore, err := g.ContainsIgnore(pathToIgnore)
	if err != nil {
		return err
	}

	path, err := g.GetPath()
	if err != nil {
		return err
	}

	if containsIgnore {
		if verbose {
			LogInfof(
				"Gitignore file '%s' already contains ignore entry for '%s'.",
				path,
				pathToIgnore,
			)
		}
		return nil
	}

	err = g.AppendLine("# "+comment, verbose)
	if err != nil {
		return err
	}

	err = g.AppendLine(pathToIgnore, verbose)
	if err != nil {
		return err
	}

	if verbose {
		LogChangedf(
			"Added '%s' to gitignore file '%s'.",
			pathToIgnore,
			path,
		)
	}

	return nil
}

func (g *GitignoreFile) ContainsIgnore(pathToCheck string) (containsIgnore bool, err error) {
	if pathToCheck == "" {
		return false, TracedError("pathToCheck is empty string")
	}

	ignoredPaths, err := g.GetIgnoredPaths()
	if err != nil {
		return false, err
	}

	containsIgnore = Slices().ContainsString(ignoredPaths, pathToCheck)

	return containsIgnore, nil
}

func (g *GitignoreFile) GetIgnoredPaths() (ignoredPaths []string, err error) {
	ignoredPaths, err = g.ReadAsLinesWithoutComments()
	if err != nil {
		return nil, err
	}

	return ignoredPaths, nil
}

func (g *GitignoreFile) MustAddDirToIgnore(pathToIgnore string, comment string, verbose bool) {
	err := g.AddDirToIgnore(pathToIgnore, comment, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitignoreFile) MustAddFileToIgnore(pathToIgnore string, comment string, verbose bool) {
	err := g.AddFileToIgnore(pathToIgnore, comment, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitignoreFile) MustContainsIgnore(pathToCheck string) (containsIgnore bool) {
	containsIgnore, err := g.ContainsIgnore(pathToCheck)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return containsIgnore
}

func (g *GitignoreFile) MustGetIgnoredPaths() (ignoredPaths []string) {
	ignoredPaths, err := g.GetIgnoredPaths()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return ignoredPaths
}

func (g *GitignoreFile) MustReformat(verbose bool) {
	err := g.Reformat(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitignoreFile) Reformat(verbose bool) (err error) {
	path, err := g.GetPath()
	if err != nil {
		return err
	}

	err = g.TrimSpacesAtBeginningOfFile(verbose)
	if err != nil {
		return err
	}

	if verbose {
		LogInfof("Reformat gitignore file '%s' finished.", path)
	}

	return nil
}
