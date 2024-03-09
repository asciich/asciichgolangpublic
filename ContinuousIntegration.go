package asciichgolangpublic

import (
	"os"
	"strings"
)

type ContinuousIntegrationService struct {
}

func ContinuousIntegration() (continuousIntegration *ContinuousIntegrationService) {
	return NewContinuousIntegrationService()
}

func NewContinuousIntegrationService() (continuousIntegration *ContinuousIntegrationService) {
	return new(ContinuousIntegrationService)
}

func (c *ContinuousIntegrationService) IsRunningInCircleCi() (isRunningInGitlab bool) {
	return strings.ToLower(os.Getenv("CIRCLECI")) != ""
}

func (c *ContinuousIntegrationService) IsRunningInContinuousIntegration() (isRunningInContinousIntegration bool) {
	if c.IsRunningInGitlab() {
		return true
	}

	if c.IsRunningInGithub() {
		return true
	}

	if c.IsRunningInCircleCi() {
		return true
	}

	if c.IsRunningInTravis() {
		return true
	}

	return false
}

func (c *ContinuousIntegrationService) IsRunningInGithub() (isRunningInGitlab bool) {
	return strings.ToLower(os.Getenv("GITHUB_ACTIONS")) != ""
}

func (c *ContinuousIntegrationService) IsRunningInGitlab() (isRunningInGitlab bool) {
	return strings.ToLower(os.Getenv("GITLAB_CI")) == "true"
}

func (c *ContinuousIntegrationService) IsRunningInTravis() (isRunningInGitlab bool) {
	return strings.ToLower(os.Getenv("TRAVIS")) != ""
}
