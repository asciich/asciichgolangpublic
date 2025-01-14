package asciichgolangpublic

import (
	"errors"
	"reflect"

	"github.com/asciich/asciichgolangpublic/changesummary"
	aerrors "github.com/asciich/asciichgolangpublic/errors"
	"github.com/asciich/asciichgolangpublic/logging"
	"gopkg.in/yaml.v3"
)

var ErrorPreCommitConfigFileContentLoad = errors.New("failed to load preCommitConfigFileContent")

type PreCommitConfigFileContent struct {
	config *PreCommitConfigFileConfig
}

func NewPreCommitConfigFileContent() (p *PreCommitConfigFileContent) {
	return new(PreCommitConfigFileContent)
}

func (p *PreCommitConfigFileContent) GetAsString() (contentString string, err error) {
	config, err := p.GetConfig()
	if err != nil {
		return "", err
	}

	contentBytes, err := yaml.Marshal(config)
	if err != nil {
		return "", aerrors.TracedError(err)
	}

	contentString = string(contentBytes)

	return contentString, nil
}

func (p *PreCommitConfigFileContent) GetConfig() (config *PreCommitConfigFileConfig, err error) {
	if p.config == nil {
		return nil, aerrors.TracedErrorf("config not set")
	}

	return p.config, nil
}

func (p *PreCommitConfigFileContent) GetDependencies(verbose bool) (dependencies []Dependency, err error) {
	config, err := p.GetConfig()
	if err != nil {
		return nil, err
	}

	repos, err := config.GetRepos()
	if err != nil {
		return nil, err
	}

	dependencies = []Dependency{}
	for _, repo := range repos {
		repoUrl, err := repo.GetRepo()
		if err != nil {
			return nil, err
		}

		versionString, err := repo.GetRev()
		if err != nil {
			return nil, err
		}

		toAdd := NewDependencyGitRepository()
		err = toAdd.SetUrl(repoUrl)
		if err != nil {
			return nil, err
		}

		err = toAdd.SetVersionString(versionString)
		if err != nil {
			return nil, err
		}

		dependencies = append(dependencies, toAdd)
	}

	return dependencies, nil
}

func (p *PreCommitConfigFileContent) LoadFromString(toLoad string) (err error) {
	if toLoad == "" {
		return aerrors.TracedErrorEmptyString("toLoad")
	}

	content := NewPreCommitConfigFileConfig()
	err = yaml.Unmarshal([]byte(toLoad), &content)
	if err != nil {
		return aerrors.TracedErrorf(
			"%w: %w",
			ErrorPreCommitConfigFileContentLoad,
			err,
		)
	}

	err = p.SetConfig(content)
	if err != nil {
		return aerrors.TracedErrorf(
			"%w: %w",
			ErrorPreCommitConfigFileContentLoad,
			err,
		)
	}

	return nil
}

func (p *PreCommitConfigFileContent) MustGetAsString() (contentString string) {
	contentString, err := p.GetAsString()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return contentString
}

func (p *PreCommitConfigFileContent) MustGetConfig() (config *PreCommitConfigFileConfig) {
	config, err := p.GetConfig()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return config
}

func (p *PreCommitConfigFileContent) MustGetDependencies(verbose bool) (dependencies []Dependency) {
	dependencies, err := p.GetDependencies(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return dependencies
}

func (p *PreCommitConfigFileContent) MustLoadFromString(toLoad string) {
	err := p.LoadFromString(toLoad)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (p *PreCommitConfigFileContent) MustSetConfig(config *PreCommitConfigFileConfig) {
	err := p.SetConfig(config)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (p *PreCommitConfigFileContent) MustUpdateDependency(dependency Dependency, authOptions []AuthenticationOption, verbose bool) (changeSummary *changesummary.ChangeSummary) {
	changeSummary, err := p.UpdateDependency(dependency, authOptions, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return changeSummary
}

func (p *PreCommitConfigFileContent) SetConfig(config *PreCommitConfigFileConfig) (err error) {
	if config == nil {
		return aerrors.TracedErrorf("config is nil")
	}

	p.config = config

	return nil
}

func (p *PreCommitConfigFileContent) UpdateDependency(dependency Dependency, authOptions []AuthenticationOption, verbose bool) (changeSummary *changesummary.ChangeSummary, err error) {
	if dependency == nil {
		return nil, aerrors.TracedErrorNil("dependency")
	}

	gitRepoDependency, ok := dependency.(*DependencyGitRepository)
	if !ok {
		return nil, aerrors.TracedErrorf("Not implemented for dependency type '%v'", reflect.TypeOf(dependency))
	}

	repoUrl, err := gitRepoDependency.GetUrl()
	if err != nil {
		return nil, err
	}

	isUpdateAvailable, err := gitRepoDependency.IsUpdateAvailable(authOptions, verbose)
	if err != nil {
		return nil, err
	}

	config, err := p.GetConfig()
	if err != nil {
		return nil, err
	}

	changeSummary = changesummary.NewChangeSummary()

	dependencyName, err := gitRepoDependency.GetName()
	if err != nil {
		return nil, err
	}

	if isUpdateAvailable {
		newVersion, err := dependency.GetNewestVersionAsString(authOptions, verbose)
		if err != nil {
			return nil, err
		}

		err = config.SetRepositoryVersion(repoUrl, newVersion)
		if err != nil {
			return nil, err
		}

		changeSummary.SetIsChanged(true)

		if verbose {
			logging.LogChangedf(
				"Dependency '%s' updated in pre-commit config file content to '%s'.",
				dependencyName,
				newVersion,
			)
		}

	} else {
		if verbose {
			logging.LogInfof(
				"Dependency '%s' is already up to date in pre-commit config file content.",
				dependencyName,
			)
		}
	}

	return changeSummary, nil
}
