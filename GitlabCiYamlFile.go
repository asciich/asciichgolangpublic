package asciichgolangpublic

import (
	"strings"

	"gopkg.in/yaml.v3"
)

type GitlabCiYamlFile struct {
	LocalFile
}

func GetGitlabCiYamlFileByFile(file File) (gitlabCiYamlFile *GitlabCiYamlFile, err error) {
	if file == nil {
		return nil, TracedErrorNil("file")
	}

	path, err := file.GetLocalPath()
	if err != nil {
		return nil, err
	}

	gitlabCiYamlFile, err = GetGitlabCiYamlFileByPath(path)
	if err != nil {
		return nil, err
	}

	return gitlabCiYamlFile, nil
}

func GetGitlabCiYamlFileByLocalFile(localFile *LocalFile) (gitlabCiYamlFile *GitlabCiYamlFile, err error) {
	if localFile == nil {
		return nil, TracedErrorNil("localFile")
	}

	path, err := localFile.GetLocalPath()
	if err != nil {
		return nil, err
	}

	gitlabCiYamlFile, err = GetGitlabCiYamlFileByPath(path)
	if err != nil {
		return nil, err
	}

	return gitlabCiYamlFile, nil
}

func GetGitlabCiYamlFileByPath(filePath string) (gitlabCiYamlFile *GitlabCiYamlFile, err error) {
	if filePath == "" {
		return nil, TracedError("filePath is empty string")
	}

	gitlabCiYamlFile = NewGitlabCiYamlFile()
	err = gitlabCiYamlFile.SetPath(filePath)
	if err != nil {
		return nil, err
	}

	err = gitlabCiYamlFile.SetParentFileForBaseClass(gitlabCiYamlFile)
	if err != nil {
		return nil, err
	}

	return gitlabCiYamlFile, nil
}

func MustGetGitlabCiYamlFileByFile(file File) (gitlabCiYamlFile *GitlabCiYamlFile) {
	gitlabCiYamlFile, err := GetGitlabCiYamlFileByFile(file)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabCiYamlFile
}

func MustGetGitlabCiYamlFileByLocalFile(localFile *LocalFile) (gitlabCiYamlFile *GitlabCiYamlFile) {
	gitlabCiYamlFile, err := GetGitlabCiYamlFileByLocalFile(localFile)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabCiYamlFile
}

func MustGetGitlabCiYamlFileByPath(filePath string) (gitlabCiYamlFile *GitlabCiYamlFile) {
	gitlabCiYamlFile, err := GetGitlabCiYamlFileByPath(filePath)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabCiYamlFile
}

func NewGitlabCiYamlFile() (g *GitlabCiYamlFile) {
	return new(GitlabCiYamlFile)
}

func (g *GitlabCiYamlFile) AddInclude(include *GitlabCiYamlInclude, verbose bool) (err error) {
	if include == nil {
		return TracedError("include is nil")
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
			LogInfof("File '%s' is already included in '%s'.", fileToInclude, path)
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
		LogChangedf("Added include '%s' to '%s'.", printableString, path)
	}

	return nil
}

func (g *GitlabCiYamlFile) ContainsInclude(include *GitlabCiYamlInclude, ignoreVersion bool, verbose bool) (containsInclude bool, err error) {
	if include == nil {
		return false, TracedError("include is nil")
	}

	includes, err := g.GetIncludes(verbose)
	if err != nil {
		return false, err
	}

	if !ignoreVersion {
		return false, TracedError("Not implemented for !ignoreVersion")
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
	includeBlock, err := g.getIncludeBlock(verbose)
	if err != nil {
		return nil, err
	}

	includeBlock = strings.TrimSpace(includeBlock)

	includes = []*GitlabCiYamlInclude{}
	blockToAdd := NewGitlabCiYamlInclude()
	for _, line := range Strings().SplitLines(includeBlock) {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			continue
		}

		if strings.HasPrefix(trimmedLine, "include:") {
			continue
		}

		if strings.HasPrefix(trimmedLine, "-") {
			if blockToAdd.IsNonEmpty() {
				includes = append(includes, blockToAdd)
			}

			blockToAdd = NewGitlabCiYamlInclude()
		}

		keyValueLine := strings.TrimSpace(strings.TrimPrefix(trimmedLine, "-"))
		splitted := strings.Split(keyValueLine, ":")
		if len(splitted) != 2 {
			return nil, TracedErrorf(
				"Unexpected splitted '%v' for line '%s'",
				splitted,
				keyValueLine,
			)
		}

		key := strings.TrimSpace(splitted[0])
		value := strings.TrimSpace(splitted[1])

		if key == "project" {
			err = blockToAdd.SetProject(value)
			if err != nil {
				return nil, err
			}
		} else if key == "ref" {
			err = blockToAdd.SetRef(value)
			if err != nil {
				return nil, err
			}
		} else if key == "file" {
			err = blockToAdd.SetFile(value)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, TracedErrorf("Unknown key: '%s' in line '%s'", splitted[0], keyValueLine)
		}
	}

	if blockToAdd.IsNonEmpty() {
		includes = append(includes, blockToAdd)
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
		LogGoErrorFatal(err)
	}
}

func (g *GitlabCiYamlFile) MustContainsInclude(include *GitlabCiYamlInclude, ignoreVersion bool, verbose bool) (containsInclude bool) {
	containsInclude, err := g.ContainsInclude(include, ignoreVersion, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return containsInclude
}

func (g *GitlabCiYamlFile) MustGetIncludes(verbose bool) (includes []*GitlabCiYamlInclude) {
	includes, err := g.GetIncludes(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return includes
}

func (g *GitlabCiYamlFile) MustGetTextBlocksWithoutIncludes(verbose bool) (textBlocks []string) {
	textBlocks, err := g.GetTextBlocksWithoutIncludes(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return textBlocks
}

func (g *GitlabCiYamlFile) MustRewriteIncludes(includes []*GitlabCiYamlInclude, verbose bool) {
	err := g.RewriteIncludes(includes, verbose)
	if err != nil {
		LogGoErrorFatal(err)
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
		LogInfof(
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
		LogInfof("No include blocks found in '%s'.", path)
	}

	return includeBlock, nil
}

func (g *GitlabCiYamlFile) getIncludesAsTextBlock(includes []*GitlabCiYamlInclude) (textBlock string, err error) {
	if includes == nil {
		return "", TracedError("includes is nil")
	}

	if len(includes) == 0 {
		return "", nil
	}

	textBlockBytes, err := yaml.Marshal(includes)
	if err != nil {
		return "", TracedError(err.Error())
	}

	textBlock = "include:\n" + string(textBlockBytes)

	return textBlock, nil
}
