package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/errors"
	"github.com/asciich/asciichgolangpublic/logging"
)

type PreCommitConfigFileConfig struct {
	Repos []PreCommitConfigFileConfigRepo `yaml:"repos"`
}

type PreCommitConfigFileConfigRepo struct {
	Repo  string                              `yaml:"repo"`
	Rev   string                              `yaml:"rev"`
	Hooks []PreCommitConfigFileConfigRepoHook `yaml:"hooks"`
}

type PreCommitConfigFileConfigRepoHook struct {
	ID string `yaml:"id"`
}

func NewPreCommitConfigFileConfig() (p *PreCommitConfigFileConfig) {
	return new(PreCommitConfigFileConfig)
}

func NewPreCommitConfigFileConfigRepo() (p *PreCommitConfigFileConfigRepo) {
	return new(PreCommitConfigFileConfigRepo)
}

func NewPreCommitConfigFileConfigRepoHook() (p *PreCommitConfigFileConfigRepoHook) {
	return new(PreCommitConfigFileConfigRepoHook)
}

func (p *PreCommitConfigFileConfig) GetRepoByUrl(repoUrl string) (repo *PreCommitConfigFileConfigRepo, err error) {
	if repoUrl == "" {
		return nil, errors.TracedErrorEmptyString("repoUrl")
	}

	repos, err := p.GetRepos()
	if err != nil {
		return nil, err
	}

	for _, repoToCheck := range repos {
		url, err := repoToCheck.GetRepo()
		if err != nil {
			return nil, err
		}

		if url == repoUrl {
			return &repoToCheck, nil
		}
	}

	return nil, errors.TracedErrorf(
		"No pre-commit repo '%s' found.",
		repoUrl,
	)
}

func (p *PreCommitConfigFileConfig) GetRepos() (repos []PreCommitConfigFileConfigRepo, err error) {
	if p.Repos == nil {
		return nil, errors.TracedErrorf("Repos not set")
	}

	if len(p.Repos) <= 0 {
		return nil, errors.TracedErrorf("Repos has no elements")
	}

	return p.Repos, nil
}

func (p *PreCommitConfigFileConfig) MustGetRepoByUrl(repoUrl string) (repo *PreCommitConfigFileConfigRepo) {
	repo, err := p.GetRepoByUrl(repoUrl)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return repo
}

func (p *PreCommitConfigFileConfig) MustGetRepos() (repos []PreCommitConfigFileConfigRepo) {
	repos, err := p.GetRepos()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return repos
}

func (p *PreCommitConfigFileConfig) MustSetRepos(repos []PreCommitConfigFileConfigRepo) {
	err := p.SetRepos(repos)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (p *PreCommitConfigFileConfig) MustSetRepositoryVersion(repoUrl string, newVersion string) {
	err := p.SetRepositoryVersion(repoUrl, newVersion)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (p *PreCommitConfigFileConfig) SetRepos(repos []PreCommitConfigFileConfigRepo) (err error) {
	if repos == nil {
		return errors.TracedErrorf("repos is nil")
	}

	if len(repos) <= 0 {
		return errors.TracedErrorf("repos has no elements")
	}

	p.Repos = repos

	return nil
}

func (p *PreCommitConfigFileConfig) SetRepositoryVersion(repoUrl string, newVersion string) (err error) {
	if repoUrl == "" {
		return errors.TracedErrorEmptyString("repoUrl")
	}

	if newVersion == "" {
		return errors.TracedErrorEmptyString("newVersion")
	}

	repo, err := p.GetRepoByUrl(repoUrl)
	if err != nil {
		return err
	}

	err = repo.SetRev(newVersion)
	if err != nil {
		return err
	}

	return nil
}

func (p *PreCommitConfigFileConfigRepo) GetHooks() (hooks []PreCommitConfigFileConfigRepoHook, err error) {
	if p.Hooks == nil {
		return nil, errors.TracedErrorf("Hooks not set")
	}

	if len(p.Hooks) <= 0 {
		return nil, errors.TracedErrorf("Hooks has no elements")
	}

	return p.Hooks, nil
}

func (p *PreCommitConfigFileConfigRepo) GetRepo() (repo string, err error) {
	if p.Repo == "" {
		return "", errors.TracedErrorf("Repo not set")
	}

	return p.Repo, nil
}

func (p *PreCommitConfigFileConfigRepo) GetRev() (rev string, err error) {
	if p.Rev == "" {
		return "", errors.TracedErrorf("Rev not set")
	}

	return p.Rev, nil
}

func (p *PreCommitConfigFileConfigRepo) MustGetHooks() (hooks []PreCommitConfigFileConfigRepoHook) {
	hooks, err := p.GetHooks()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return hooks
}

func (p *PreCommitConfigFileConfigRepo) MustGetRepo() (repo string) {
	repo, err := p.GetRepo()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return repo
}

func (p *PreCommitConfigFileConfigRepo) MustGetRev() (rev string) {
	rev, err := p.GetRev()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return rev
}

func (p *PreCommitConfigFileConfigRepo) MustSetHooks(hooks []PreCommitConfigFileConfigRepoHook) {
	err := p.SetHooks(hooks)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (p *PreCommitConfigFileConfigRepo) MustSetRepo(repo string) {
	err := p.SetRepo(repo)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (p *PreCommitConfigFileConfigRepo) MustSetRev(rev string) {
	err := p.SetRev(rev)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (p *PreCommitConfigFileConfigRepo) SetHooks(hooks []PreCommitConfigFileConfigRepoHook) (err error) {
	if hooks == nil {
		return errors.TracedErrorf("hooks is nil")
	}

	if len(hooks) <= 0 {
		return errors.TracedErrorf("hooks has no elements")
	}

	p.Hooks = hooks

	return nil
}

func (p *PreCommitConfigFileConfigRepo) SetRepo(repo string) (err error) {
	if repo == "" {
		return errors.TracedErrorf("repo is empty string")
	}

	p.Repo = repo

	return nil
}

func (p *PreCommitConfigFileConfigRepo) SetRev(rev string) (err error) {
	if rev == "" {
		return errors.TracedErrorf("rev is empty string")
	}

	p.Rev = rev

	return nil
}

func (p *PreCommitConfigFileConfigRepoHook) GetID() (iD string, err error) {
	if p.ID == "" {
		return "", errors.TracedErrorf("ID not set")
	}

	return p.ID, nil
}

func (p *PreCommitConfigFileConfigRepoHook) MustGetID() (iD string) {
	iD, err := p.GetID()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return iD
}

func (p *PreCommitConfigFileConfigRepoHook) MustSetID(iD string) {
	err := p.SetID(iD)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (p *PreCommitConfigFileConfigRepoHook) SetID(iD string) (err error) {
	if iD == "" {
		return errors.TracedErrorf("iD is empty string")
	}

	p.ID = iD

	return nil
}
