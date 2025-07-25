package asciichgolangpublic

import (
	"slices"

	"github.com/asciich/asciichgolangpublic/pkg/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type GitignoreFile struct {
	filesinterfaces.File
}

func GetGitignoreDefaultBaseName() (defaultBaseName string) {
	return ".gitignore"
}

func GetGitignoreFileByFile(fileToUse filesinterfaces.File) (gitignoreFile *GitignoreFile, err error) {
	if fileToUse == nil {
		return nil, tracederrors.TracedErrorEmptyString("fileToUse")
	}

	gitignoreFile = NewGitignoreFile()

	gitignoreFile.File = fileToUse

	return gitignoreFile, nil
}

func GetGitignoreFileByPath(filePath string) (gitignoreFile *GitignoreFile, err error) {
	if filePath == "" {
		return nil, tracederrors.TracedErrorEmptyString("filePath")
	}

	fileToUse, err := files.GetLocalFileByPath(filePath)
	if err != nil {
		return nil, err
	}

	return GetGitignoreFileByFile(fileToUse)
}

func GetGitignoreFileInGitRepository(gitRepository GitRepository) (gitignoreFile *GitignoreFile, err error) {
	if gitRepository == nil {
		return nil, tracederrors.TracedErrorNil("gitRepository")
	}

	fileToUse, err := gitRepository.GetFileByPath(GetGitignoreDefaultBaseName())
	if err != nil {
		return nil, err
	}

	return GetGitignoreFileByFile(fileToUse)
}

func NewGitignoreFile() (g *GitignoreFile) {
	return new(GitignoreFile)
}

func (g *GitignoreFile) AddDirToIgnore(pathToIgnore string, comment string, verbose bool) (err error) {
	if pathToIgnore == "" {
		return tracederrors.TracedError("pathToIgnore is empty string")
	}

	if comment == "" {
		return tracederrors.TracedError("comment is empty string")
	}

	pathToIgnore = stringsutils.EnsureSuffix(pathToIgnore, "/")

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
			logging.LogInfof(
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
		logging.LogChangedf(
			"Added '%s' to gitignore file '%s'.",
			pathToIgnore,
			path,
		)
	}

	return nil
}

func (g *GitignoreFile) AddFileToIgnore(pathToIgnore string, comment string, verbose bool) (err error) {
	if pathToIgnore == "" {
		return tracederrors.TracedError("pathToIgnore is empty string")
	}

	if comment == "" {
		return tracederrors.TracedError("comment is empty string")
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
			logging.LogInfof(
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
		logging.LogChangedf(
			"Added '%s' to gitignore file '%s'.",
			pathToIgnore,
			path,
		)
	}

	return nil
}

func (g *GitignoreFile) ContainsIgnore(pathToCheck string) (containsIgnore bool, err error) {
	if pathToCheck == "" {
		return false, tracederrors.TracedError("pathToCheck is empty string")
	}

	ignoredPaths, err := g.GetIgnoredPaths()
	if err != nil {
		return false, err
	}

	containsIgnore = slices.Contains(ignoredPaths, pathToCheck)

	return containsIgnore, nil
}

func (g *GitignoreFile) GetIgnoredPaths() (ignoredPaths []string, err error) {
	ignoredPaths, err = g.ReadAsLinesWithoutComments()
	if err != nil {
		return nil, err
	}

	return ignoredPaths, nil
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
		logging.LogInfof("Reformat gitignore file '%s' finished.", path)
	}

	return nil
}
