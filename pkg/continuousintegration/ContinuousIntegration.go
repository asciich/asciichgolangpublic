package continuousintegration

import (
	"os"
	"strings"

	"github.com/asciich/asciichgolangpublic/randomgenerator"
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

func GetDefaultKindClusterName() string {
	const defaultClusterName = "kind-ci-cluster"

	if IsRunningInContinuousIntegration() {
		// On Github multiple create and delete of the same cluster lead to errors (unable to create cluster again).
		// Therefore we generate a new name for every test.
		return defaultClusterName + "-" + strings.ToLower(randomgenerator.MustGetRandomString(5))
	}

	return defaultClusterName
}
