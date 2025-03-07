package asciichgolangpublic

import (
	"strings"

	"github.com/asciich/asciichgolangpublic/files"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
	"gopkg.in/yaml.v3"
)

type GitlabCiYamlFile struct {
	files.File
}

func GetGitlabCiYamlDefaultBaseName() (defaultBaseName string) {
	return ".gitlab-ci.yml"
}

func GetGitlabCiYamlFileByFile(file files.File) (gitlabCiYamlFile *GitlabCiYamlFile, err error) {
	if file == nil {
		return nil, tracederrors.TracedErrorNil("file")
	}

	gitlabCiYamlFile = NewGitlabCiYamlFile()
	gitlabCiYamlFile.File = file

	return gitlabCiYamlFile, nil
}

func GetGitlabCiYamlFileByPath(filePath string) (gitlabCiYamlFile *GitlabCiYamlFile, err error) {
	if filePath == "" {
		return nil, tracederrors.TracedError("filePath is empty string")
	}

	localFile, err := files.GetLocalFileByPath(filePath)
	if err != nil {
		return nil, err
	}

	return GetGitlabCiYamlFileByFile(localFile)
}

func GetGitlabCiYamlFileInGitRepository(gitRepository GitRepository) (gitlabCiYamlFile *GitlabCiYamlFile, err error) {
	if gitRepository == nil {
		return nil, tracederrors.TracedErrorNil("gitRepository")
	}

	fileToUse, err := gitRepository.GetFileByPath(GetGitlabCiYamlDefaultBaseName())
	if err != nil {
		return nil, err
	}

	return GetGitlabCiYamlFileByFile(fileToUse)
}

func MustGetGitlabCiYamlFileByFile(file files.File) (gitlabCiYamlFile *GitlabCiYamlFile) {
	gitlabCiYamlFile, err := GetGitlabCiYamlFileByFile(file)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabCiYamlFile
}

func MustGetGitlabCiYamlFileByPath(filePath string) (gitlabCiYamlFile *GitlabCiYamlFile) {
	gitlabCiYamlFile, err := GetGitlabCiYamlFileByPath(filePath)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabCiYamlFile
}

func MustGetGitlabCiYamlFileInGitRepository(gitRepository GitRepository) (gitlabCiYamlFile *GitlabCiYamlFile) {
	gitlabCiYamlFile, err := GetGitlabCiYamlFileInGitRepository(gitRepository)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabCiYamlFile
}

func NewGitlabCiYamlFile() (g *GitlabCiYamlFile) {
	return new(GitlabCiYamlFile)
}

func (g *GitlabCiYamlFile) AddInclude(include *GitlabCiYamlInclude, verbose bool) (err error) {
	if include == nil {
		return tracederrors.TracedError("include is nil")
	}

	err = g.Create(verbose)
	if err != nil {
		return err
	}

	const ignoreVersion bool = true
	containsInclude, err := g.ContainsInclude(include, ignoreVersion, verbose)
	if err != nil {
		return err
	}

	path, err := g.GetPath()
	if err != nil {
		return err
	}

	fileToInclude, err := include.GetFile()
	if err != nil {
		return err
	}

	if containsInclude {
		if verbose {
			logging.LogInfof("File '%s' is already included in '%s'.", fileToInclude, path)
		}
		return nil
	}

	includes, err := g.GetIncludes(verbose)
	if err != nil {
		return err
	}

	includes = append(includes, include)

	err = g.RewriteIncludes(includes, verbose)
	if err != nil {
		return err
	}

	printableString, err := include.GetLoggableString()
	if err != nil {
		return err
	}

	if verbose {
		logging.LogChangedf("Added include '%s' to '%s'.", printableString, path)
	}

	return nil
}

func (g *GitlabCiYamlFile) ContainsInclude(include *GitlabCiYamlInclude, ignoreVersion bool, verbose bool) (containsInclude bool, err error) {
	if include == nil {
		return false, tracederrors.TracedError("include is nil")
	}

	includes, err := g.GetIncludes(verbose)
	if err != nil {
		return false, err
	}

	if !ignoreVersion {
		return false, tracederrors.TracedError("Not implemented for !ignoreVersion")
	}

	for _, toCheck := range includes {
		isEqual, err := toCheck.EqualsIgnoreVersion(include)
		if err != nil {
			return false, err
		}

		if isEqual {
			return true, nil
		}
	}

	return false, nil
}

func (g *GitlabCiYamlFile) GetIncludes(verbose bool) (includes []*GitlabCiYamlInclude, err error) {
	localPath, err := g.GetLocalPath()
	if err != nil {
		return nil, err
	}

	includeBlock, err := g.getIncludeBlock(verbose)
	if err != nil {
		return nil, err
	}

	type IncludesYaml struct {
		Includes []*GitlabCiYamlInclude `yaml:"include"`
	}

	includesYaml := new(IncludesYaml)

	includeBlock = strings.TrimSpace(includeBlock)

	err = yaml.Unmarshal([]byte(includeBlock), includesYaml)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Unable to parse inclues in gitlab-ci.yaml '%s': %w", localPath, err)
	}

	includes = includesYaml.Includes

	if verbose {
		logging.LogInfof("Found '%d' includes in gitlab-ci.yml '%s'.", len(includes), localPath)
	}

	return includes, nil
}

func (g *GitlabCiYamlFile) GetTextBlocksWithoutIncludes(verbose bool) (textBlocks []string, err error) {
	blocks, err := g.GetTextBlocks(verbose)
	if err != nil {
		return nil, err
	}

	textBlocks = []string{}
	for _, block := range blocks {
		trimmedBlock := strings.TrimSpace(block)
		if strings.HasPrefix(trimmedBlock, "include:") {
			continue
		}

		textBlocks = append(textBlocks, trimmedBlock)
	}

	return textBlocks, nil
}

func (g *GitlabCiYamlFile) MustAddInclude(include *GitlabCiYamlInclude, verbose bool) {
	err := g.AddInclude(include, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabCiYamlFile) MustContainsInclude(include *GitlabCiYamlInclude, ignoreVersion bool, verbose bool) (containsInclude bool) {
	containsInclude, err := g.ContainsInclude(include, ignoreVersion, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return containsInclude
}

func (g *GitlabCiYamlFile) MustGetIncludes(verbose bool) (includes []*GitlabCiYamlInclude) {
	includes, err := g.GetIncludes(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return includes
}

func (g *GitlabCiYamlFile) MustGetTextBlocksWithoutIncludes(verbose bool) (textBlocks []string) {
	textBlocks, err := g.GetTextBlocksWithoutIncludes(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return textBlocks
}

func (g *GitlabCiYamlFile) MustRewriteIncludes(includes []*GitlabCiYamlInclude, verbose bool) {
	err := g.RewriteIncludes(includes, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabCiYamlFile) RewriteIncludes(includes []*GitlabCiYamlInclude, verbose bool) (err error) {
	blocks, err := g.GetTextBlocksWithoutIncludes(verbose)
	if err != nil {
		return err
	}

	includeBlock, err := g.getIncludesAsTextBlock(includes)
	if err != nil {
		return err
	}

	blocksToWrite := append([]string{includeBlock}, blocks...)

	err = g.WriteTextBlocks(blocksToWrite, verbose)
	if err != nil {
		return err
	}

	path, err := g.GetPath()
	if err != nil {
		return err
	}

	if verbose {
		logging.LogInfof(
			"Added %d includes to '%s'",
			len(includes),
			path,
		)
	}

	return nil
}

func (g *GitlabCiYamlFile) getIncludeBlock(verbose bool) (includeBlock string, err error) {
	blocks, err := g.GetTextBlocks(verbose)
	if err != nil {
		return "", err
	}

	for _, block := range blocks {
		block = strings.TrimSpace(block)
		if strings.HasPrefix(block, "include:\n") {
			return block, nil
		}
	}

	path, err := g.GetPath()
	if err != nil {
		return "", err
	}

	if verbose {
		logging.LogInfof("No include blocks found in '%s'.", path)
	}

	return includeBlock, nil
}

func (g *GitlabCiYamlFile) getIncludesAsTextBlock(includes []*GitlabCiYamlInclude) (textBlock string, err error) {
	if includes == nil {
		return "", tracederrors.TracedError("includes is nil")
	}

	if len(includes) == 0 {
		return "", nil
	}

	textBlockBytes, err := yaml.Marshal(includes)
	if err != nil {
		return "", tracederrors.TracedError(err.Error())
	}

	textBlock = "include:\n" + string(textBlockBytes)

	return textBlock, nil
}
