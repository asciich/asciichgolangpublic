package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
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
		return nil, tracederrors.TracedErrorEmptyString("repoUrl")
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

	return nil, tracederrors.TracedErrorf(
		"No pre-commit repo '%s' found.",
		repoUrl,
	)
}

func (p *PreCommitConfigFileConfig) GetRepos() (repos []PreCommitConfigFileConfigRepo, err error) {
	if p.Repos == nil {
		return nil, tracederrors.TracedErrorf("Repos not set")
	}

	if len(p.Repos) <= 0 {
		return nil, tracederrors.TracedErrorf("Repos has no elements")
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
		return tracederrors.TracedErrorf("repos is nil")
	}

	if len(repos) <= 0 {
		return tracederrors.TracedErrorf("repos has no elements")
	}

	p.Repos = repos

	return nil
}

func (p *PreCommitConfigFileConfig) SetRepositoryVersion(repoUrl string, newVersion string) (err error) {
	if repoUrl == "" {
		return tracederrors.TracedErrorEmptyString("repoUrl")
	}

	if newVersion == "" {
		return tracederrors.TracedErrorEmptyString("newVersion")
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
		return nil, tracederrors.TracedErrorf("Hooks not set")
	}

	if len(p.Hooks) <= 0 {
		return nil, tracederrors.TracedErrorf("Hooks has no elements")
	}

	return p.Hooks, nil
}

func (p *PreCommitConfigFileConfigRepo) GetRepo() (repo string, err error) {
	if p.Repo == "" {
		return "", tracederrors.TracedErrorf("Repo not set")
	}

	return p.Repo, nil
}

func (p *PreCommitConfigFileConfigRepo) GetRev() (rev string, err error) {
	if p.Rev == "" {
		return "", tracederrors.TracedErrorf("Rev not set")
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
		return tracederrors.TracedErrorf("hooks is nil")
	}

	if len(hooks) <= 0 {
		return tracederrors.TracedErrorf("hooks has no elements")
	}

	p.Hooks = hooks

	return nil
}

func (p *PreCommitConfigFileConfigRepo) SetRepo(repo string) (err error) {
	if repo == "" {
		return tracederrors.TracedErrorf("repo is empty string")
	}

	p.Repo = repo

	return nil
}

func (p *PreCommitConfigFileConfigRepo) SetRev(rev string) (err error) {
	if rev == "" {
		return tracederrors.TracedErrorf("rev is empty string")
	}

	p.Rev = rev

	return nil
}

func (p *PreCommitConfigFileConfigRepoHook) GetID() (iD string, err error) {
	if p.ID == "" {
		return "", tracederrors.TracedErrorf("ID not set")
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
		return tracederrors.TracedErrorf("iD is empty string")
	}

	p.ID = iD

	return nil
}
