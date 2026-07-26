package continuousintegration

import (
	"os"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/randomgenerator"
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
	const name = "kind-ci-cluster"

	if IsRunningInContinuousIntegration() {
		// On Github multiple create and delete of the same cluster lead to errors (unable to create cluster again).
		// Therefore we generate a new name for every test.
		return name + "-" + strings.ToLower(mustutils.Must(randomgenerator.GetRandomString(5)))
	}

	return name
}

// GetKindClusterNameByPackageName returns a KinD cluster name based on the Go package name.
// This ensures that parallel running tests in different packages do not interfere with each other.
func GetKindClusterNameByPackageName(packageName string) string {
	if packageName == "" {
		return GetDefaultKindClusterName()
	}

	// Replace special characters with dashes to create a valid cluster name
	clusterName := strings.ReplaceAll(packageName, "/", "-")
	clusterName = strings.ReplaceAll(clusterName, ".", "-")
	clusterName = strings.ReplaceAll(clusterName, "_", "-")

	// In CI, add a random suffix to avoid conflicts from previous failed runs
	if IsRunningInContinuousIntegration() {
		return clusterName + "-" + strings.ToLower(mustutils.Must(randomgenerator.GetRandomString(5)))
	}

	return clusterName
}
