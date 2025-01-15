package docker

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func getDockerImplementationByName(implementationName string) (docker Docker) {
	if implementationName == "commandExecutorDocker" {
		return MustGetLocalCommandExecutorDocker()
	}

	logging.LogFatalWithTracef("Unknown implementation name '%s'", implementationName)
	return nil
}

func TestDocker_GetHostName(t *testing.T) {

	tests := []struct {
		implementationName string
	}{
		{"commandExecutorDocker"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				docker := getDockerImplementationByName(tt.implementationName)

				assert.EqualValues(
					"localhost",
					docker.MustGetHostDescription(),
				)
			},
		)
	}
}
