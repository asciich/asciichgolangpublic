package asciichgolangpublic

import (
	"context"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"gopkg.in/yaml.v3"
)

type GitlabCiYamlFile struct {
	filesinterfaces.File
}

func GetGitlabCiYamlDefaultBaseName() (defaultBaseName string) {
	return ".gitlab-ci.yml"
}

func GetGitlabCiYamlFileByFile(file filesinterfaces.File) (gitlabCiYamlFile *GitlabCiYamlFile, err error) {
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

func GetGitlabCiYamlFileInGitRepository(gitRepository gitinterfaces.GitRepository) (gitlabCiYamlFile *GitlabCiYamlFile, err error) {
	if gitRepository == nil {
		return nil, tracederrors.TracedErrorNil("gitRepository")
	}

	fileToUse, err := gitRepository.GetFileByPath(GetGitlabCiYamlDefaultBaseName())
	if err != nil {
		return nil, err
	}

	return GetGitlabCiYamlFileByFile(fileToUse)
}

func NewGitlabCiYamlFile() (g *GitlabCiYamlFile) {
	return new(GitlabCiYamlFile)
}

func (g *GitlabCiYamlFile) AddInclude(ctx context.Context, include *GitlabCiYamlInclude) (err error) {
	if include == nil {
		return tracederrors.TracedError("include is nil")
	}

	err = g.Create(ctx, &filesoptions.CreateOptions{})
	if err != nil {
		return err
	}

	const ignoreVersion bool = true
	containsInclude, err := g.ContainsInclude(ctx, include, ignoreVersion)
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
		logging.LogInfoByCtxf(ctx, "File '%s' is already included in '%s'.", fileToInclude, path)
		return nil
	}

	includes, err := g.GetIncludes(ctx)
	if err != nil {
		return err
	}

	includes = append(includes, include)

	err = g.RewriteIncludes(ctx, includes)
	if err != nil {
		return err
	}

	printableString, err := include.GetLoggableString()
	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(ctx, "Added include '%s' to '%s'.", printableString, path)

	return nil
}

func (g *GitlabCiYamlFile) ContainsInclude(ctx context.Context, include *GitlabCiYamlInclude, ignoreVersion bool) (containsInclude bool, err error) {
	if include == nil {
		return false, tracederrors.TracedError("include is nil")
	}

	includes, err := g.GetIncludes(ctx)
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

func (g *GitlabCiYamlFile) GetIncludes(ctx context.Context) (includes []*GitlabCiYamlInclude, err error) {
	localPath, err := g.GetLocalPath()
	if err != nil {
		return nil, err
	}

	includeBlock, err := g.getIncludeBlock(ctx)
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

	logging.LogInfoByCtxf(ctx, "Found '%d' includes in gitlab-ci.yml '%s'.", len(includes), localPath)

	return includes, nil
}

func (g *GitlabCiYamlFile) GetTextBlocksWithoutIncludes(ctx context.Context) (textBlocks []string, err error) {
	blocks, err := g.GetTextBlocks(contextutils.GetVerboseFromContext(ctx))
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

func (g *GitlabCiYamlFile) RewriteIncludes(ctx context.Context, includes []*GitlabCiYamlInclude) (err error) {
	blocks, err := g.GetTextBlocksWithoutIncludes(ctx)
	if err != nil {
		return err
	}

	includeBlock, err := g.getIncludesAsTextBlock(includes)
	if err != nil {
		return err
	}

	blocksToWrite := append([]string{includeBlock}, blocks...)

	err = g.WriteTextBlocks(blocksToWrite, contextutils.GetVerboseFromContext(ctx))
	if err != nil {
		return err
	}

	path, err := g.GetPath()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Added %d includes to '%s'", len(includes), path)

	return nil
}

func (g *GitlabCiYamlFile) getIncludeBlock(ctx context.Context) (includeBlock string, err error) {
	blocks, err := g.GetTextBlocks(contextutils.GetVerboseFromContext(ctx))
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

	logging.LogInfoByCtxf(ctx, "No include blocks found in '%s'.", path)

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
