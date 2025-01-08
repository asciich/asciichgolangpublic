package asciichgolangpublic

import (
	"strings"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type GitlabRunner struct {
	gitlab            *GitlabInstance
	id                int
	cachedName        string
	cachedDescription string
}

func NewGitlabRunner() (gitlabRunner *GitlabRunner) {
	return new(GitlabRunner)
}

func (g *GitlabRunner) MustGetCachedDescription() (description string) {
	description, err := g.GetCachedDescription()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return description
}

func (g *GitlabRunner) MustGetCachedName() (name string) {
	name, err := g.GetCachedName()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return name
}

func (g *GitlabRunner) MustGetCachedNameOrDescription() (name string) {
	name, err := g.GetCachedNameOrDescription()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return name
}

func (g *GitlabRunner) MustGetGitlab() (gitlab *GitlabInstance) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlab
}

func (g *GitlabRunner) MustGetGitlabRunners() (gitlabRunners *GitlabRunnersService) {
	gitlabRunners, err := g.GetGitlabRunners()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabRunners
}

func (g *GitlabRunner) MustGetId() (id int) {
	id, err := g.GetId()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return id
}

func (g *GitlabRunner) MustGetNativeRunnersService() (nativeRunnersService *gitlab.RunnersService) {
	nativeRunnersService, err := g.GetNativeRunnersService()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return nativeRunnersService
}

func (g *GitlabRunner) MustIsStatusOk() (isStatusOk bool) {
	isStatusOk, err := g.IsStatusOk()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return isStatusOk
}

func (g *GitlabRunner) MustRemove(verbose bool) {
	err := g.Remove(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabRunner) MustResetRunnerToken() (runnerToken string) {
	runnerToken, err := g.ResetRunnerToken()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return runnerToken
}

func (g *GitlabRunner) MustSetCachedDescription(description string) {
	err := g.SetCachedDescription(description)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabRunner) MustSetCachedName(name string) {
	err := g.SetCachedName(name)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabRunner) MustSetGitlab(gitlab *GitlabInstance) {
	err := g.SetGitlab(gitlab)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabRunner) MustSetId(id int) {
	err := g.SetId(id)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (r *GitlabRunner) GetCachedDescription() (description string, err error) {
	if len(r.cachedDescription) <= 0 {
		return "", TracedError("cachedDescription not set")
	}

	return r.cachedDescription, nil
}

func (r *GitlabRunner) GetCachedName() (name string, err error) {
	if len(r.cachedName) <= 0 {
		return "", TracedError("cachedName not set")
	}

	return r.cachedName, nil
}

func (r *GitlabRunner) GetCachedNameOrDescription() (name string, err error) {
	if len(r.cachedName) > 0 {
		return r.cachedName, nil
	}

	if len(r.cachedDescription) > 0 {
		return r.cachedDescription, nil
	}

	return "", TracedError("Both cachedName and cachedDescription not set")
}

func (r *GitlabRunner) GetGitlabRunners() (gitlabRunners *GitlabRunnersService, err error) {
	gitlab, err := r.GetGitlab()
	if err != nil {
		return nil, err
	}

	gitlabRunners, err = gitlab.GetGitlabRunners()
	if err != nil {
		return nil, err
	}

	return gitlabRunners, nil
}

func (r *GitlabRunner) GetNativeRunnersService() (nativeRunnersService *gitlab.RunnersService, err error) {
	runners, err := r.GetGitlabRunners()
	if err != nil {
		return nil, err
	}

	nativeRunnersService, err = runners.GetNativeRunnersService()
	if err != nil {
		return nil, err
	}

	return nativeRunnersService, nil
}

func (r *GitlabRunner) IsCachedDescriptionSet() (isSet bool) {
	return len(r.cachedDescription) > 0
}

func (r *GitlabRunner) IsCachedDescriptionUnset() (isUnset bool) {
	return !r.IsCachedDescriptionSet()
}

func (r *GitlabRunner) IsCachedNameSet() (isCachedNameSet bool) {
	return len(r.cachedName) > 0
}

func (r *GitlabRunner) IsCachedNameUnset() (isCachedNameUnset bool) {
	return !r.IsCachedNameSet()
}

func (r *GitlabRunner) IsStatusOk() (isStatusOk bool, err error) {
	id, err := r.GetId()
	if err != nil {
		return false, err
	}

	nativeRunnerService, err := r.GetNativeRunnersService()
	if err != nil {
		return false, err
	}

	nativeDetails, _, err := nativeRunnerService.GetRunnerDetails(id)
	if err != nil {
		return false, TracedError(err.Error())
	}

	if !nativeDetails.Online {
		return false, nil
	}

	if nativeDetails.Paused {
		return false, nil
	}

	if nativeDetails.Status != "online" {
		return false, nil
	}

	return true, nil
}

func (r *GitlabRunner) Remove(verbose bool) (err error) {
	nativeRunnersService, err := r.GetNativeRunnersService()
	if err != nil {
		return err
	}

	runnerId, err := r.GetId()
	if err != nil {
		return err
	}

	_, err = nativeRunnersService.RemoveRunner(runnerId)
	if err != nil {
		return err
	}

	if verbose {
		LogChangedf("Runner with id '%d' removed.", runnerId)
	}

	return nil
}

func (r *GitlabRunner) ResetRunnerToken() (runnerToken string, err error) {
	id, err := r.GetId()
	if err != nil {
		return "", err
	}

	nativeRunnersService, err := r.GetNativeRunnersService()
	if err != nil {
		return "", err
	}

	nativeToken, _, err := nativeRunnersService.ResetRunnerAuthenticationToken(id)
	if err != nil {
		return "", TracedError(err.Error())
	}

	if nativeToken == nil {
		return "", TracedError("nativeToken is nil")
	}

	runnerToken = *nativeToken.Token
	runnerToken = strings.TrimSpace(runnerToken)

	return runnerToken, nil
}

func (r *GitlabRunner) SetCachedDescription(description string) (err error) {
	if len(description) <= 0 {
		return TracedError("description is empty string")
	}

	r.cachedDescription = description

	return nil
}

func (s *GitlabRunner) GetGitlab() (gitlab *GitlabInstance, err error) {
	if s.gitlab == nil {
		return nil, TracedError("gitlab not set")
	}

	return s.gitlab, nil
}

func (s *GitlabRunner) GetId() (id int, err error) {
	if s.id <= 0 {
		return -1, TracedError("id not set")
	}

	return s.id, nil
}

func (s *GitlabRunner) SetCachedName(name string) (err error) {
	if len(name) <= 0 {
		return TracedError("name is empty string")
	}

	s.cachedName = name

	return nil
}

func (s *GitlabRunner) SetGitlab(gitlab *GitlabInstance) (err error) {
	if gitlab == nil {
		return TracedError("gitlab is nil")
	}

	s.gitlab = gitlab

	return nil
}

func (s *GitlabRunner) SetId(id int) (err error) {
	if id <= 0 {
		return TracedErrorf("Invalid id '%d'", id)
	}

	s.id = id

	return nil
}
