package docker

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/asciich/asciichgolangpublic"
	"github.com/asciich/asciichgolangpublic/containers"
	"github.com/asciich/asciichgolangpublic/logging"
)

func getDockerContainerToTest(implementationName string, containerName string) (container containers.Container) {
	if implementationName == "commandExectuorDockerContainer" {
		return MustGetLocalCommandExecutorDocker().MustGetContainerByName(containerName)
	}

	logging.LogFatalWithTracef("Unkown implementaion name: '%s'", implementationName)

	return nil
}

func TestContainers_IsHostRunning(t *testing.T) {

	tests := []struct {
		implementationName string
	}{
		{"commandExectuorDockerContainer"},
	}

	for _, tt := range tests {
		t.Run(
			asciichgolangpublic.MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				container := getDockerContainerToTest(tt.implementationName, "thisContainerDoesNotRun")

				assert.False(container.MustIsRunning(verbose))
			},
		)
	}
}
