package asciichgolangpublic

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/asciich/asciichgolangpublic/pkg/changesummary"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/dependencyutils/dependencyinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/pathsutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type PreCommitConfigFile struct {
	files.LocalFile
}

func GetPreCommitConfigByFile(file filesinterfaces.File) (preCommitConfigFile *PreCommitConfigFile, err error) {
	if file == nil {
		return nil, tracederrors.TracedErrorNil("file")
	}

	path, err := file.GetLocalPath()
	if err != nil {
		return nil, err
	}

	preCommitConfigFile = NewPreCommitConfigFile()
	err = preCommitConfigFile.SetPath(path)
	if err != nil {
		return nil, err
	}

	return preCommitConfigFile, nil
}

func GetPreCommitConfigByLocalPath(localPath string) (preCommitConfigFile *PreCommitConfigFile, err error) {
	if localPath == "" {
		return nil, tracederrors.TracedErrorEmptyString("localPath")
	}

	file, err := files.GetLocalFileByPath(localPath)
	if err != nil {
		return nil, err
	}

	preCommitConfigFile, err = GetPreCommitConfigByFile(file)
	if err != nil {
		return nil, err
	}

	return preCommitConfigFile, nil
}

func GetPreCommitConfigFileInGitRepository(gitRepository gitinterfaces.GitRepository) (preCommitConfigFile *PreCommitConfigFile, err error) {
	if gitRepository == nil {
		return nil, tracederrors.TracedErrorNil("gitRepository")
	}

	fileInRepo, err := gitRepository.GetFileByPath(PreCommit().GetDefaultConfigFileName())
	if err != nil {
		return nil, err
	}

	return GetPreCommitConfigByFile(fileInRepo)
}

func MustGetPreCommitConfigByFile(file filesinterfaces.File) (preCommitConfigFile *PreCommitConfigFile) {
	preCommitConfigFile, err := GetPreCommitConfigByFile(file)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return preCommitConfigFile
}

func MustGetPreCommitConfigByLocalPath(localPath string) (preCommitConfigFile *PreCommitConfigFile) {
	preCommitConfigFile, err := GetPreCommitConfigByLocalPath(localPath)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return preCommitConfigFile
}

func MustGetPreCommitConfigFileInGitRepository(gitRepository gitinterfaces.GitRepository) (preCommitConfigFile *PreCommitConfigFile) {
	preCommitConfigFile, err := GetPreCommitConfigFileInGitRepository(gitRepository)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return preCommitConfigFile
}

func NewPreCommitConfigFile() (preCommitConfigFile *PreCommitConfigFile) {
	preCommitConfigFile = new(PreCommitConfigFile)

	err := preCommitConfigFile.SetParentFileForBaseClass(preCommitConfigFile)
	if err != nil {
		logging.LogFatalWithTracef("internal error: '%v'", err)
	}

	return preCommitConfigFile
}

func (p *PreCommitConfigFile) GetAbsolutePath() (absolutePath string, err error) {
	path, err := p.GetPath()
	if err != nil {
		return "", err
	}

	if pathsutils.IsRelativePath(path) {
		return "", tracederrors.TracedErrorf(
			"Unable to get absolute path, '%s' is relative",
			path,
		)
	}

	return path, nil
}

func (p *PreCommitConfigFile) GetDependencies(ctx context.Context) (dependencies []dependencyinterfaces.Dependency, err error) {
	preCommitConfigFileContent, err := p.GetPreCommitConfigFileContent(ctx)
	if err != nil {
		return nil, err
	}

	dependencies, err = preCommitConfigFileContent.GetDependencies(ctx)
	if err != nil {
		return nil, err
	}

	localPath, err := p.GetLocalPath()
	if err != nil {
		return nil, err
	}

	asciichgolangpublicFile, err := files.GetLocalFileByPath(localPath)
	if err != nil {
		return nil, err
	}

	err = DependenciesSlice().AddSourceFileForEveryEntry(dependencies, asciichgolangpublicFile)
	if err != nil {
		return nil, err
	}

	return dependencies, err
}

func (p *PreCommitConfigFile) GetLocalPath() (localPath string, err error) {
	return p.GetPath()
}

func (p *PreCommitConfigFile) GetPreCommitConfigFileContent(ctx context.Context) (content *PreCommitConfigFileContent, err error) {
	contentString, err := p.ReadAsString()
	if err != nil {
		return nil, err
	}

	content = NewPreCommitConfigFileContent()
	err = content.LoadFromString(contentString)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (p *PreCommitConfigFile) GetUriAsString() (uri string, err error) {
	absoutePath, err := p.GetAbsolutePath()
	if err != nil {
		return "", err
	}

	uri = "file://" + absoutePath

	return uri, nil
}

func (p *PreCommitConfigFile) IsValidPreCommitConfigFile(ctx context.Context) (isValidPreCommitConfigFile bool, err error) {
	_, err = p.GetPreCommitConfigFileContent(ctx)
	if err != nil {
		if errors.Is(err, ErrorPreCommitConfigFileContentLoad) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (p *PreCommitConfigFile) UpdateDependencies(ctx context.Context, options *parameteroptions.UpdateDependenciesOptions) (changeSummary *changesummary.ChangeSummary, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	dependencies, err := p.GetDependencies(ctx)
	if err != nil {
		return nil, err
	}

	changeSummary = changesummary.NewChangeSummary()

	for _, dependency := range dependencies {
		singleUpdateSummary, err := p.UpdateDependency(ctx, dependency, options)
		if err != nil {
			return nil, err
		}

		err = changeSummary.AddChildSummary(singleUpdateSummary)
		if err != nil {
			return nil, err
		}
	}

	path, err := p.GetPath()
	if err != nil {
		return nil, err
	}

	if changeSummary.IsChanged() {
		logging.LogChangedByCtxf(ctx, "Updated dependencies in pre-commit config file '%s'.", path)
	} else {
		logging.LogInfoByCtxf(ctx, "All dependencies in pre-commit config file '%s' were already up to date.", path)
	}

	return changeSummary, nil
}

func (p *PreCommitConfigFile) UpdateDependency(ctx context.Context, dependency dependencyinterfaces.Dependency, options *parameteroptions.UpdateDependenciesOptions) (changeSummary *changesummary.ChangeSummary, err error) {
	if dependency == nil {
		return nil, tracederrors.TracedErrorNil("dependency")
	}

	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	gitRepoDependency, ok := dependency.(*DependencyGitRepository)
	if !ok {
		return nil, tracederrors.TracedErrorf("Not implemented for dependency type '%v'", reflect.TypeOf(dependency))
	}

	url, err := gitRepoDependency.GetUrl()
	if err != nil {
		return nil, err
	}

	newestVersion, err := gitRepoDependency.GetNewestVersionAsString(ctx, options.AuthenticationOptions)
	if err != nil {
		return nil, err
	}

	repoLine := fmt.Sprintf("- repo: %s", url)

	dependencyName, err := dependency.GetName()
	if err != nil {
		return nil, err
	}

	path, err := p.GetPath()
	if err != nil {
		return nil, err
	}

	changeSummary, err = p.ReplaceLineAfterLine(
		repoLine,
		fmt.Sprintf("  rev: \"%s\"", newestVersion),
		contextutils.GetVerboseFromContext(ctx),
	)
	if err != nil {
		return nil, err
	}

	if changeSummary.IsChanged() {
		logging.LogChangedf(
			"Dependency '%s' updated in '%s'.",
			dependencyName,
			path,
		)
	} else {
		logging.LogInfof(
			"Dependency '%s' already up to date in '%s'.",
			dependencyName,
			path,
		)
	}

	return changeSummary, nil
}

func (p *PreCommitConfigFile) WritePreCommitConfigFileContent(content *PreCommitConfigFileContent, verbose bool) (err error) {
	toWrite, err := content.GetAsString()
	if err != nil {
		return err
	}

	err = p.WriteString(contextutils.GetVerbosityContextByBool(verbose), toWrite, &filesoptions.WriteOptions{})
	if err != nil {
		return err
	}

	path, err := p.GetPath()
	if err != nil {
		return err
	}

	if verbose {
		logging.LogChangedf("Wrote content of pre-commit config file '%s'.", path)
	}

	return nil
}
