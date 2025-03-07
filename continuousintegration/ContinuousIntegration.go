package continuousintegration

import (
	"os"
	"strings"
)

func IsRunningInCircleCi() (isRunningInGitlab bool) {
	return strings.ToLower(os.Getenv("CIRCLECI")) != ""
}

func IsRunningInContinuousIntegration() (isRunningInContinousIntegration bool) {
	if IsRunningInGitlab() {
		return true
	}

	if IsRunningInGithub() {
		return true
	}

	if IsRunningInCircleCi() {
		return true
	}

	if IsRunningInTravis() {
		return true
	}

	return false
}

func IsRunningInGithub() (isRunningInGitlab bool) {
	return strings.ToLower(os.Getenv("GITHUB_ACTIONS")) != ""
}

func IsRunningInGitlab() (isRunningInGitlab bool) {
	return strings.ToLower(os.Getenv("GITLAB_CI")) == "true"
}

func IsRunningInTravis() (isRunningInGitlab bool) {
	return strings.ToLower(os.Getenv("TRAVIS")) != ""
}
