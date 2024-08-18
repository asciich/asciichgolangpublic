package asciichgolangpublic

import (
	"errors"
	"fmt"
	"reflect"
)

type PreCommitConfigFile struct {
	LocalFile
}

func GetPreCommitConfigByFile(file File) (preCommitConfigFile *PreCommitConfigFile, err error) {
	if file == nil {
		return nil, TracedErrorNil("file")
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
		return nil, TracedErrorEmptyString("localPath")
	}

	file, err := NewLocalFileByPath(localPath)
	if err != nil {
		return nil, err
	}

	preCommitConfigFile, err = GetPreCommitConfigByFile(file)
	if err != nil {
		return nil, err
	}

	return preCommitConfigFile, nil
}

func MustGetPreCommitConfigByFile(file File) (preCommitConfigFile *PreCommitConfigFile) {
	preCommitConfigFile, err := GetPreCommitConfigByFile(file)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return preCommitConfigFile
}

func MustGetPreCommitConfigByLocalPath(localPath string) (preCommitConfigFile *PreCommitConfigFile) {
	preCommitConfigFile, err := GetPreCommitConfigByLocalPath(localPath)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return preCommitConfigFile
}

func NewPreCommitConfigFile() (preCommitConfigFile *PreCommitConfigFile) {
	preCommitConfigFile = new(PreCommitConfigFile)

	err := preCommitConfigFile.SetParentFileForBaseClass(preCommitConfigFile)
	if err != nil {
		LogFatalWithTracef("internal error: '%v'", err)
	}

	return preCommitConfigFile
}

func (p *PreCommitConfigFile) GetAbsolutePath() (absolutePath string, err error) {
	path, err := p.GetPath()
	if err != nil {
		return "", err
	}

	if Paths().IsRelativePath(path) {
		return "", TracedErrorf(
			"Unable to get absolute path, '%s' is relative",
			path,
		)
	}

	return path, nil
}

func (p *PreCommitConfigFile) GetDependencies(verbose bool) (dependencies []Dependency, err error) {
	preCommitConfigFileContent, err := p.GetPreCommitConfigFileContent(verbose)
	if err != nil {
		return nil, err
	}

	dependencies, err = preCommitConfigFileContent.GetDependencies(verbose)
	if err != nil {
		return nil, err
	}

	localPath, err := p.GetLocalPath()
	if err != nil {
		return nil, err
	}

	asciichgolangpublicFile, err := GetLocalFileByPath(localPath)
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

func (p *PreCommitConfigFile) GetPreCommitConfigFileContent(verbose bool) (content *PreCommitConfigFileContent, err error) {
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

func (p *PreCommitConfigFile) IsValidPreCommitConfigFile(verbose bool) (isValidPreCommitConfigFile bool, err error) {
	_, err = p.GetPreCommitConfigFileContent(verbose)
	if err != nil {
		if errors.Is(err, ErrorPreCommitConfigFileContentLoad) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (p *PreCommitConfigFile) MustGetAbsolutePath() (absolutePath string) {
	absolutePath, err := p.GetAbsolutePath()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return absolutePath
}

func (p *PreCommitConfigFile) MustGetDependencies(verbose bool) (dependencies []Dependency) {
	dependencies, err := p.GetDependencies(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return dependencies
}

func (p *PreCommitConfigFile) MustGetLocalPath() (localPath string) {
	localPath, err := p.GetLocalPath()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return localPath
}

func (p *PreCommitConfigFile) MustGetPreCommitConfigFileContent(verbose bool) (content *PreCommitConfigFileContent) {
	content, err := p.GetPreCommitConfigFileContent(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return content
}

func (p *PreCommitConfigFile) MustGetUriAsString() (uri string) {
	uri, err := p.GetUriAsString()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return uri
}

func (p *PreCommitConfigFile) MustIsValidPreCommitConfigFile(verbose bool) (isValidPreCommitConfigFile bool) {
	isValidPreCommitConfigFile, err := p.IsValidPreCommitConfigFile(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return isValidPreCommitConfigFile
}

func (p *PreCommitConfigFile) MustUpdateDependencies(options *UpdateDependenciesOptions) (changeSummary *ChangeSummary) {
	changeSummary, err := p.UpdateDependencies(options)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return changeSummary
}

func (p *PreCommitConfigFile) MustUpdateDependency(dependency Dependency, options *UpdateDependenciesOptions) (changeSummary *ChangeSummary) {
	changeSummary, err := p.UpdateDependency(dependency, options)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return changeSummary
}

func (p *PreCommitConfigFile) MustWritePreCommitConfigFileContent(content *PreCommitConfigFileContent, verbose bool) {
	err := p.WritePreCommitConfigFileContent(content, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (p *PreCommitConfigFile) UpdateDependencies(options *UpdateDependenciesOptions) (changeSummary *ChangeSummary, err error) {
	if options == nil {
		return nil, TracedErrorNil("options")
	}

	dependencies, err := p.GetDependencies(options.Verbose)
	if err != nil {
		return nil, err
	}

	changeSummary = NewChangeSummary()

	for _, dependency := range dependencies {
		singleUpdateSummary, err := p.UpdateDependency(dependency, options)
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

	if options.Verbose {
		if changeSummary.IsChanged() {
			LogChangedf("Updated dependencies in pre-commit config file '%s'.", path)
		} else {
			LogInfof("All dependencies in pre-commit config file '%s' were already up to date.", path)
		}
	}

	return changeSummary, nil
}

func (p *PreCommitConfigFile) UpdateDependency(dependency Dependency, options *UpdateDependenciesOptions) (changeSummary *ChangeSummary, err error) {
	if dependency == nil {
		return nil, TracedErrorNil("dependency")
	}

	if options == nil {
		return nil, TracedErrorNil("options")
	}

	gitRepoDependency, ok := dependency.(*DependencyGitRepository)
	if !ok {
		return nil, TracedErrorf("Not implemented for dependency type '%v'", reflect.TypeOf(dependency))
	}

	url, err := gitRepoDependency.GetUrl()
	if err != nil {
		return nil, err
	}

	newestVersion, err := gitRepoDependency.GetNewestVersionAsString(options.AuthenticationOptions, options.Verbose)
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
		options.Verbose,
	)
	if err != nil {
		return nil, err
	}

	if changeSummary.IsChanged() {
		LogChangedf(
			"Dependency '%s' updated in '%s'.",
			dependencyName,
			path,
		)
	} else {
		LogInfof(
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

	err = p.WriteString(toWrite, verbose)
	if err != nil {
		return err
	}

	path, err := p.GetPath()
	if err != nil {
		return err
	}

	if verbose {
		LogChangedf("Wrote content of pre-commit config file '%s'.", path)
	}

	return nil
}
