package asciichgolangpublic

import (
	"fmt"
	"slices"
	"sort"

	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorbashoo"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type GitlabRunnersService struct {
	gitlab *GitlabInstance
}

func NewGitlabRunners() (gitlabRunners *GitlabRunnersService) {
	return new(GitlabRunnersService)
}

func NewGitlabRunnersService() (g *GitlabRunnersService) {
	return new(GitlabRunnersService)
}

// According to the documentation this only works when logged in as admin:
// https://github.com/xanzy/go-gitlab/blob/master/runners.go#L126
func (s *GitlabRunnersService) GetRunnerList() (runners []*GitlabRunner, err error) {
	g, err := s.GetGitlab()
	if err != nil {
		return nil, err
	}

	nativeRunnersService, err := s.GetNativeRunnersService()
	if err != nil {
		return nil, err
	}

	nativeRunners, _, err := nativeRunnersService.ListAllRunners(&gitlab.ListRunnersOptions{})
	if err != nil {
		return nil, err
	}

	runners = []*GitlabRunner{}
	for _, nativeRunner := range nativeRunners {
		nameToAdd := nativeRunner.Name
		descriptionToAdd := nativeRunner.Description
		idToAdd := nativeRunner.ID

		runnerToAdd := NewGitlabRunner()
		err = runnerToAdd.SetGitlab(g)
		if err != nil {
			return nil, err
		}

		if len(nameToAdd) > 0 {
			err = runnerToAdd.SetCachedName(nameToAdd)
			if err != nil {
				return nil, err
			}
		}

		err = runnerToAdd.SetId(idToAdd)
		if err != nil {
			return nil, err
		}

		if len(descriptionToAdd) > 0 {
			err = runnerToAdd.SetCachedDescription(descriptionToAdd)
			if err != nil {
				return nil, err
			}
		}

		runners = append(runners, runnerToAdd)
	}

	return runners, nil
}

func (g *GitlabRunnersService) MustAddRunner(newRunnerOptions *GitlabAddRunnerOptions) (createdRunner *GitlabRunner) {
	createdRunner, err := g.AddRunner(newRunnerOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return createdRunner
}

func (g *GitlabRunnersService) MustCheckRunnerStatusOk(runnerName string, verbose bool) (isRunnerOk bool) {
	isRunnerOk, err := g.CheckRunnerStatusOk(runnerName, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isRunnerOk
}

func (g *GitlabRunnersService) MustGetApiV4Url() (apiV4Url string) {
	apiV4Url, err := g.GetApiV4Url()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return apiV4Url
}

func (g *GitlabRunnersService) MustGetCurrentlyUsedAccessToken() (gitlabAccessToken string) {
	gitlabAccessToken, err := g.GetCurrentlyUsedAccessToken()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlabAccessToken
}

func (g *GitlabRunnersService) MustGetFqdn() (fqdn string) {
	fqdn, err := g.GetFqdn()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return fqdn
}

func (g *GitlabRunnersService) MustGetGitlab() (gitlab *GitlabInstance) {
	gitlab, err := g.GetGitlab()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitlab
}

func (g *GitlabRunnersService) MustGetNativeClient() (nativeClient *gitlab.Client) {
	nativeClient, err := g.GetNativeClient()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nativeClient
}

func (g *GitlabRunnersService) MustGetNativeRunnersService() (nativeRunnersService *gitlab.RunnersService) {
	nativeRunnersService, err := g.GetNativeRunnersService()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nativeRunnersService
}

func (g *GitlabRunnersService) MustGetRunnerByName(runnerName string) (runner *GitlabRunner) {
	runner, err := g.GetRunnerByName(runnerName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return runner
}

func (g *GitlabRunnersService) MustGetRunnerList() (runners []*GitlabRunner) {
	runners, err := g.GetRunnerList()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return runners
}

func (g *GitlabRunnersService) MustGetRunnerNamesList() (runnerNames []string) {
	runnerNames, err := g.GetRunnerNamesList()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return runnerNames
}

func (g *GitlabRunnersService) MustIsRunnerStatusOk(runnerName string, verbose bool) (isStatusOk bool) {
	isStatusOk, err := g.IsRunnerStatusOk(runnerName, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isStatusOk
}

func (g *GitlabRunnersService) MustRemoveAllRunners(verbose bool) {
	err := g.RemoveAllRunners(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabRunnersService) MustRunnerByNameExists(runnerName string) (exists bool) {
	exists, err := g.RunnerByNameExists(runnerName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return exists
}

func (g *GitlabRunnersService) MustSetGitlab(gitlab *GitlabInstance) {
	err := g.SetGitlab(gitlab)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (r *GitlabRunnersService) GetFqdn() (fqdn string, err error) {
	gitlab, err := r.GetGitlab()
	if err != nil {
		return "", err
	}

	fqdn, err = gitlab.GetFqdn()
	if err != nil {
		return "", err
	}

	return fqdn, nil
}

func (r *GitlabRunnersService) RemoveAllRunners(verbose bool) (err error) {
	fqdn, err := r.GetFqdn()
	if err != nil {
		return err
	}

	if verbose {
		logging.LogInfof("Delete all gitlab '%s' runners started.", fqdn)
	}

	runners, err := r.GetRunnerList()
	if err != nil {
		return err
	}
	for _, runner := range runners {
		err = runner.Remove(verbose)
		if err != nil {
			return err
		}
	}

	if verbose {
		logging.LogInfof("Delete all gitlab '%s' runners finished.", fqdn)
	}

	return nil
}

func (s *GitlabRunnersService) AddRunner(newRunnerOptions *GitlabAddRunnerOptions) (createdRunner *GitlabRunner, err error) {
	if newRunnerOptions == nil {
		return nil, tracederrors.TracedError("newRunnerOptions is nil")
	}

	runnerName, err := newRunnerOptions.GetRunnerName()
	if err != nil {
		return nil, err
	}

	apiV4Url, err := s.GetApiV4Url()
	if err != nil {
		return nil, err
	}

	tagsCommaSeperated, err := newRunnerOptions.GetTagsCommaSeparated()
	if err != nil {
		return nil, err
	}

	privateToken, err := s.GetCurrentlyUsedAccessToken()
	if err != nil {
		return nil, err
	}

	runnerExists, err := s.RunnerByNameExists(runnerName)
	if err != nil {
		return nil, err
	}

	if runnerExists {
		if newRunnerOptions.Verbose {
			logging.LogInfof("Gitlab runner '%s' already exists.", runnerName)
		}
	} else {
		addRunnerCommand := []string{
			"curl",
			"-sX", "POST",
			fmt.Sprintf("%s/user/runners", apiV4Url),
			"--data", "runner_type=instance_type",
			"--data", "description=" + runnerName,
			"--data", "tag_list=" + tagsCommaSeperated,
			"--data", "run_untagged=false",
			"--header", "PRIVATE-TOKEN: " + privateToken,
		}

		_, err = commandexecutorbashoo.Bash().RunCommand(
			contextutils.GetVerbosityContextByBool(newRunnerOptions.Verbose),
			&parameteroptions.RunCommandOptions{
				Command: addRunnerCommand,
			},
		)
		if err != nil {
			return nil, err
		}
		if newRunnerOptions.Verbose {
			logging.LogChangedf("Registered/ created new gitlab runner '%s'", runnerName)
		}
	}

	createdRunner, err = s.GetRunnerByName(runnerName)
	if err != nil {
		return nil, err
	}

	return createdRunner, nil
}

func (s *GitlabRunnersService) CheckRunnerStatusOk(runnerName string, verbose bool) (isRunnerOk bool, err error) {
	if len(runnerName) <= 0 {
		return false, tracederrors.TracedError("runnerName is empty string")
	}

	isRunnerOk, err = s.IsRunnerStatusOk(runnerName, verbose)
	if err != nil {
		return false, err
	}

	if !isRunnerOk {
		return false, tracederrors.TracedErrorf("Runner '%s' status is NOT ok", runnerName)
	}

	return isRunnerOk, nil
}

func (s *GitlabRunnersService) GetApiV4Url() (apiV4Url string, err error) {
	gitlab, err := s.GetGitlab()
	if err != nil {
		return "", err
	}

	apiV4Url, err = gitlab.GetApiV4Url()
	if err != nil {
		return "", err
	}

	return apiV4Url, nil
}

func (s *GitlabRunnersService) GetCurrentlyUsedAccessToken() (gitlabAccessToken string, err error) {
	gitlab, err := s.GetGitlab()
	if err != nil {
		return "", err
	}

	gitlabAccessToken, err = gitlab.GetCurrentlyUsedAccessToken()
	if err != nil {
		return "", err
	}

	return gitlabAccessToken, nil
}

func (s *GitlabRunnersService) GetGitlab() (gitlab *GitlabInstance, err error) {
	if s.gitlab == nil {
		return nil, tracederrors.TracedError("gitlab not set")
	}

	return s.gitlab, nil
}

func (s *GitlabRunnersService) GetNativeClient() (nativeClient *gitlab.Client, err error) {
	gitlab, err := s.GetGitlab()
	if err != nil {
		return nil, err
	}

	nativeClient, err = gitlab.GetNativeClient()
	if err != nil {
		return nil, err
	}

	return nativeClient, nil
}

func (s *GitlabRunnersService) GetNativeRunnersService() (nativeRunnersService *gitlab.RunnersService, err error) {
	nativeClient, err := s.GetNativeClient()
	if err != nil {
		return nil, err
	}

	nativeRunnersService = nativeClient.Runners
	return nativeRunnersService, nil
}

func (s *GitlabRunnersService) GetRunnerByName(runnerName string) (runner *GitlabRunner, err error) {
	if len(runnerName) <= 0 {
		return nil, tracederrors.TracedError("runnerName is empty string")
	}

	runners, err := s.GetRunnerList()
	if err != nil {
		return nil, err
	}

	for _, runner := range runners {
		if runner.IsCachedNameUnset() {
			if runner.IsCachedDescriptionUnset() {
				continue
			}
		}

		nameToCheck, err := runner.GetCachedNameOrDescription()
		if err != nil {
			return nil, err
		}

		if nameToCheck == runnerName {
			return runner, nil
		}
	}

	return nil, tracederrors.TracedErrorf("runner '%s' not found.", runnerName)
}

func (s *GitlabRunnersService) GetRunnerNamesList() (runnerNames []string, err error) {
	runners, err := s.GetRunnerList()
	if err != nil {
		return nil, err
	}

	runnerNames = []string{}
	for _, runner := range runners {
		nameToAdd, err := runner.GetCachedNameOrDescription()
		if err != nil {
			return nil, err
		}

		runnerNames = append(runnerNames, nameToAdd)
	}

	sort.Strings(runnerNames)

	return runnerNames, err
}

func (s *GitlabRunnersService) IsRunnerStatusOk(runnerName string, verbose bool) (isStatusOk bool, err error) {
	if len(runnerName) <= 0 {
		return false, tracederrors.TracedError("runnerName is empty string")
	}

	runnerExists, err := s.RunnerByNameExists(runnerName)
	if err != nil {
		return false, err
	}

	if !runnerExists {
		if verbose {
			logging.LogInfof("Runner '%s' does not exists and therefore status is not ok", runnerName)
		}
		return false, nil
	}

	runner, err := s.GetRunnerByName(runnerName)
	if err != nil {
		return false, err
	}

	isStatusOk, err = runner.IsStatusOk()
	if err != nil {
		return false, err
	}

	return isStatusOk, nil
}

func (s *GitlabRunnersService) RunnerByNameExists(runnerName string) (exists bool, err error) {
	if len(runnerName) <= 0 {
		return false, tracederrors.TracedError("runnerName is emtpy string")
	}

	runnerNames, err := s.GetRunnerNamesList()
	if err != nil {
		return false, err
	}

	return slices.Contains(runnerNames, runnerName), nil
}

func (s *GitlabRunnersService) SetGitlab(gitlab *GitlabInstance) (err error) {
	if gitlab == nil {
		return tracederrors.TracedError("gitlab is nil")
	}

	s.gitlab = gitlab

	return nil
}
