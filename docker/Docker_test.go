package docker

import (
	"testing"

	"github.com/stretchr/testify/require"
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
				require := require.New(t)

				docker := getDockerImplementationByName(tt.implementationName)

				require.EqualValues(
					"localhost",
					docker.MustGetHostDescription(),
				)
			},
		)
	}
}
