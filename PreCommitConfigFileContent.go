package asciichgolangpublic

import (
	"context"
	"errors"
	"reflect"

	"github.com/asciich/asciichgolangpublic/pkg/changesummary"
	"github.com/asciich/asciichgolangpublic/pkg/dependencyutils/dependencyinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions/authenticationoptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
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
		return "", tracederrors.TracedError(err)
	}

	contentString = string(contentBytes)

	return contentString, nil
}

func (p *PreCommitConfigFileContent) GetConfig() (config *PreCommitConfigFileConfig, err error) {
	if p.config == nil {
		return nil, tracederrors.TracedErrorf("config not set")
	}

	return p.config, nil
}

func (p *PreCommitConfigFileContent) GetDependencies(ctx context.Context) (dependencies []dependencyinterfaces.Dependency, err error) {
	config, err := p.GetConfig()
	if err != nil {
		return nil, err
	}

	repos, err := config.GetRepos()
	if err != nil {
		return nil, err
	}

	dependencies = []dependencyinterfaces.Dependency{}
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
		return tracederrors.TracedErrorEmptyString("toLoad")
	}

	content := NewPreCommitConfigFileConfig()
	err = yaml.Unmarshal([]byte(toLoad), &content)
	if err != nil {
		return tracederrors.TracedErrorf(
			"%w: %w",
			ErrorPreCommitConfigFileContentLoad,
			err,
		)
	}

	err = p.SetConfig(content)
	if err != nil {
		return tracederrors.TracedErrorf(
			"%w: %w",
			ErrorPreCommitConfigFileContentLoad,
			err,
		)
	}

	return nil
}

func (p *PreCommitConfigFileContent) SetConfig(config *PreCommitConfigFileConfig) (err error) {
	if config == nil {
		return tracederrors.TracedErrorf("config is nil")
	}

	p.config = config

	return nil
}

func (p *PreCommitConfigFileContent) UpdateDependency(ctx context.Context, dependency dependencyinterfaces.Dependency, authOptions []authenticationoptions.AuthenticationOption) (changeSummary *changesummary.ChangeSummary, err error) {
	if dependency == nil {
		return nil, tracederrors.TracedErrorNil("dependency")
	}

	gitRepoDependency, ok := dependency.(*DependencyGitRepository)
	if !ok {
		return nil, tracederrors.TracedErrorf("Not implemented for dependency type '%v'", reflect.TypeOf(dependency))
	}

	repoUrl, err := gitRepoDependency.GetUrl()
	if err != nil {
		return nil, err
	}

	isUpdateAvailable, err := gitRepoDependency.IsUpdateAvailable(ctx, authOptions)
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
		newVersion, err := dependency.GetNewestVersionAsString(ctx, authOptions)
		if err != nil {
			return nil, err
		}

		err = config.SetRepositoryVersion(repoUrl, newVersion)
		if err != nil {
			return nil, err
		}

		changeSummary.SetIsChanged(true)

		logging.LogChangedByCtxf(ctx, "Dependency '%s' updated in pre-commit config file content to '%s'.", dependencyName, newVersion)
	} else {
		logging.LogInfoByCtxf(ctx, "Dependency '%s' is already up to date in pre-commit config file content.", dependencyName)
	}

	return changeSummary, nil
}
