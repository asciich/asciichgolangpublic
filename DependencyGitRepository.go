package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/changesummary"
	"github.com/asciich/asciichgolangpublic/files"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/parameteroptions/authenticationoptions"
	"github.com/asciich/asciichgolangpublic/pkg/versionutils"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

// Represents a dependency to (another) git repository.
type DependencyGitRepository struct {
	url           string
	versionString string
	sourceFiles   []files.File

	// If defined the url will not be used to get the newest version automatically.
	// Instead this targetVersionString will become the newest available version and will be set in the sourceFiles.
	targetVersionString string
}

func NewDependencyGitRepository() (d *DependencyGitRepository) {
	return new(DependencyGitRepository)
}

func (d *DependencyGitRepository) AddSourceFile(sourceFile files.File) (err error) {
	if sourceFile == nil {
		return tracederrors.TracedErrorNil("sourceFile")
	}

	d.sourceFiles = append(d.sourceFiles, sourceFile)

	return nil
}

func (d *DependencyGitRepository) GetName() (name string, err error) {
	return d.GetUrl()
}

func (d *DependencyGitRepository) GetNewestVersion(authOptions []authenticationoptions.AuthenticationOption, verbose bool) (newestVersion versionutils.Version, err error) {
	url, err := d.GetUrl()
	if err != nil {
		return nil, err
	}

	if d.IsTargetVersionSet() {
		newestVersion, err = d.GetTargetVersion()
		if err != nil {
			return nil, err
		}

		targetVersionString, err := newestVersion.GetAsString()
		if err != nil {
			return nil, err
		}

		if verbose {
			logging.LogInfof(
				"Newest version for '%s' is set by already defined target version '%s'",
				url,
				targetVersionString,
			)
		}

		return newestVersion, nil
	}

	gitlabProject, err := GetGitlabProjectByUrlFromString(url, authOptions, verbose)
	if err != nil {
		return nil, err
	}

	newestVersion, err = gitlabProject.GetNewestVersion(verbose)
	if err != nil {
		return nil, err
	}

	name, err := d.GetName()
	if err != nil {
		return nil, err
	}

	newestVersionString, err := newestVersion.GetAsString()
	if err != nil {
		return nil, err
	}

	if verbose {
		logging.LogInfof(
			"Newest version of git repository dependency '%s' is '%s'.",
			name,
			newestVersionString,
		)
	}

	return newestVersion, err
}

func (d *DependencyGitRepository) GetNewestVersionAsString(authOptions []authenticationoptions.AuthenticationOption, verbose bool) (newestVersionString string, err error) {
	newestVersion, err := d.GetNewestVersion(authOptions, verbose)
	if err != nil {
		return "", err
	}

	newestVersionString, err = newestVersion.GetAsString()
	if err != nil {
		return "", err
	}

	if newestVersionString == "" {
		return "", tracederrors.TracedError(
			"Unable to get newest version string, newestVersionString is empty string after evaluation",
		)
	}

	return newestVersionString, nil
}

func (d *DependencyGitRepository) GetSourceFiles() (sourceFiles []files.File, err error) {
	if d.sourceFiles == nil {
		return nil, tracederrors.TracedErrorf("sourceFiles not set")
	}

	if len(d.sourceFiles) <= 0 {
		return nil, tracederrors.TracedErrorf("sourceFiles has no elements")
	}

	return d.sourceFiles, nil
}

func (d *DependencyGitRepository) GetTargetVersion() (targetVersion versionutils.Version, err error) {
	targetVersionString, err := d.GetTargetVersionString()
	if err != nil {
		return nil, err
	}

	targetVersion, err = versionutils.ReadFromString(targetVersionString)
	if err != nil {
		return nil, err
	}

	return targetVersion, nil
}

func (d *DependencyGitRepository) GetTargetVersionString() (targetVersionString string, err error) {
	if d.targetVersionString == "" {
		return "", tracederrors.TracedErrorf("targetVersionString not set")
	}

	return d.targetVersionString, nil
}

func (d *DependencyGitRepository) GetUrl() (url string, err error) {
	if d.url == "" {
		return "", tracederrors.TracedErrorf("url not set")
	}

	return d.url, nil
}

func (d *DependencyGitRepository) GetVersion() (version versionutils.Version, err error) {
	versionString, err := d.GetVersionString()
	if err != nil {
		return nil, err
	}

	version, err = versionutils.ReadFromString(versionString)
	if err != nil {
		return nil, err
	}

	return version, nil
}

func (d *DependencyGitRepository) GetVersionString() (versionString string, err error) {
	if d.versionString == "" {
		return "", tracederrors.TracedErrorf("versionString not set")
	}

	return d.versionString, nil
}

func (d *DependencyGitRepository) IsAtLeastOneSourceFileSet() (isSet bool) {
	return len(d.sourceFiles) > 0
}

func (d *DependencyGitRepository) IsTargetVersionSet() (isSet bool) {
	return d.targetVersionString != ""
}

func (d *DependencyGitRepository) IsUpdateAvailable(authOptions []authenticationoptions.AuthenticationOption, verbose bool) (isUpdateAvailable bool, err error) {
	if d.IsVersionStringUnset() {
		return true, nil
	}

	currentVersion, err := d.GetVersion()
	if err != nil {
		return false, err
	}

	newestVersionString, err := d.GetNewestVersion(authOptions, verbose)
	if err != nil {
		return false, err
	}

	isUpdateAvailable = !currentVersion.Equals(newestVersionString)

	dependencyName, err := d.GetName()
	if err != nil {
		return false, err
	}

	if verbose {
		if isUpdateAvailable {
			logging.LogChangedf(
				"Update available for dependency '%s'. Current version is '%s' but newest version is '%s'.",
				dependencyName,
				currentVersion,
				newestVersionString,
			)
		} else {
			logging.LogInfof(
				"No Update available for dependency '%s'. Current version is '%s' and newest version is '%s'.",
				dependencyName,
				currentVersion,
				newestVersionString,
			)
		}
	}

	return isUpdateAvailable, nil
}

func (d *DependencyGitRepository) IsVersionStringUnset() (isUnset bool) {
	return d.versionString == ""
}

func (d *DependencyGitRepository) MustAddSourceFile(sourceFile files.File) {
	err := d.AddSourceFile(sourceFile)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (d *DependencyGitRepository) MustGetName() (name string) {
	name, err := d.GetName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return name
}

func (d *DependencyGitRepository) MustGetNewestVersion(authOptions []authenticationoptions.AuthenticationOption, verbose bool) (newestVersion versionutils.Version) {
	newestVersion, err := d.GetNewestVersion(authOptions, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return newestVersion
}

func (d *DependencyGitRepository) MustGetNewestVersionAsString(authOptions []authenticationoptions.AuthenticationOption, verbose bool) (newestVersion string) {
	newestVersion, err := d.GetNewestVersionAsString(authOptions, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return newestVersion
}

func (d *DependencyGitRepository) MustGetSourceFiles() (sourceFiles []files.File) {
	sourceFiles, err := d.GetSourceFiles()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return sourceFiles
}

func (d *DependencyGitRepository) MustGetTargetVersion() (targeVersion versionutils.Version) {
	targeVersion, err := d.GetTargetVersion()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return targeVersion
}

func (d *DependencyGitRepository) MustGetTargetVersionString() (targetVersionString string) {
	targetVersionString, err := d.GetTargetVersionString()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return targetVersionString
}

func (d *DependencyGitRepository) MustGetUrl() (url string) {
	url, err := d.GetUrl()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return url
}

func (d *DependencyGitRepository) MustGetVersion() (version versionutils.Version) {
	version, err := d.GetVersion()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return version
}

func (d *DependencyGitRepository) MustGetVersionString() (versionString string) {
	versionString, err := d.GetVersionString()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return versionString
}

func (d *DependencyGitRepository) MustIsUpdateAvailable(authOptions []authenticationoptions.AuthenticationOption, verbose bool) (isUpdateAvailable bool) {
	isUpdateAvailable, err := d.IsUpdateAvailable(authOptions, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isUpdateAvailable
}

func (d *DependencyGitRepository) MustSetSourceFiles(sourceFiles []files.File) {
	err := d.SetSourceFiles(sourceFiles)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (d *DependencyGitRepository) MustSetTargetVersionString(targetVersionString string) {
	err := d.SetTargetVersionString(targetVersionString)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (d *DependencyGitRepository) MustSetUrl(url string) {
	err := d.SetUrl(url)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (d *DependencyGitRepository) MustSetVersionString(versionString string) {
	err := d.SetVersionString(versionString)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (d *DependencyGitRepository) MustUpdate(options *parameteroptions.UpdateDependenciesOptions) (changeSummary *changesummary.ChangeSummary) {
	changeSummary, err := d.Update(options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return changeSummary
}

func (d *DependencyGitRepository) MustUpdateVersionByStringInSourceFile(version string, sourceFile files.File, options *parameteroptions.UpdateDependenciesOptions) (changeSummary *changesummary.ChangeSummary) {
	changeSummary, err := d.UpdateVersionByStringInSourceFile(version, sourceFile, options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return changeSummary
}

func (d *DependencyGitRepository) SetSourceFiles(sourceFiles []files.File) (err error) {
	if sourceFiles == nil {
		return tracederrors.TracedErrorf("sourceFiles is nil")
	}

	if len(sourceFiles) <= 0 {
		return tracederrors.TracedErrorf("sourceFiles has no elements")
	}

	d.sourceFiles = sourceFiles

	return nil
}

func (d *DependencyGitRepository) SetTargetVersionString(targetVersionString string) (err error) {
	if targetVersionString == "" {
		return tracederrors.TracedErrorf("targetVersionString is empty string")
	}

	d.targetVersionString = targetVersionString

	return nil
}

func (d *DependencyGitRepository) SetUrl(url string) (err error) {
	if url == "" {
		return tracederrors.TracedErrorf("url is empty string")
	}

	d.url = url

	return nil
}

func (d *DependencyGitRepository) SetVersionString(versionString string) (err error) {
	if versionString == "" {
		return tracederrors.TracedErrorf("versionString is empty string")
	}

	d.versionString = versionString

	return nil
}

func (d *DependencyGitRepository) Update(options *parameteroptions.UpdateDependenciesOptions) (changeSummary *changesummary.ChangeSummary, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	name, err := d.GetName()
	if err != nil {
		return nil, err
	}

	changeSummary = changesummary.NewChangeSummary()

	if options.Verbose {
		logging.LogInfof("Update git repository dependency '%s' started.", name)
	}

	if !d.IsAtLeastOneSourceFileSet() {
		return nil, tracederrors.TracedErrorf("No source files set for git repository dependency '%s'", name)
	}

	latestVersion, err := d.GetNewestVersionAsString(
		options.AuthenticationOptions,
		options.Verbose,
	)
	if err != nil {
		return nil, err
	}

	sourceFiles, err := d.GetSourceFiles()
	if err != nil {
		return nil, err
	}

	for _, sourceFile := range sourceFiles {
		sourceFileSummary, err := d.UpdateVersionByStringInSourceFile(latestVersion, sourceFile, options)
		if err != nil {
			return nil, err
		}

		err = changeSummary.AddChildSummary(sourceFileSummary)
		if err != nil {
			return nil, err
		}
	}

	if options.Verbose {
		logging.LogByChangeSummaryf(changeSummary, "Update git repository dependency '%s' finished.", name)
	}

	return changeSummary, nil
}

func (d *DependencyGitRepository) UpdateVersionByStringInSourceFile(version string, sourceFile files.File, options *parameteroptions.UpdateDependenciesOptions) (changeSummary *changesummary.ChangeSummary, err error) {
	if version == "" {
		return nil, tracederrors.TracedErrorEmptyString("version")
	}

	if sourceFile == nil {
		return nil, tracederrors.TracedErrorNil("sourceFile")
	}

	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	changeSummary = changesummary.NewChangeSummary()

	name, err := d.GetName()
	if err != nil {
		return nil, err
	}

	sourceFileUri, err := sourceFile.GetUriAsString()
	if err != nil {
		return nil, err
	}

	if options.Verbose {
		logging.LogInfof(
			"Update of git repository dependency '%s' in '%s' started.",
			name,
			sourceFileUri,
		)
	}

	preCommitConfigFile, err := PreCommit().GetAsPreCommitConfigFileOrNilIfContentIsInvalid(sourceFile, options.Verbose)
	if err != nil {
		return nil, err
	}

	if preCommitConfigFile != nil {
		fileChangeSummary, err := preCommitConfigFile.UpdateDependency(d, options)
		if err != nil {
			return nil, err
		}

		err = changeSummary.AddChildSummary(fileChangeSummary)
		if err != nil {
			return nil, err
		}

		return fileChangeSummary, nil
	}

	return nil, tracederrors.TracedErrorf("Not implemneted for '%s'", sourceFileUri)
}
