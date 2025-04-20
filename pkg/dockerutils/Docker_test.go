package dockerutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func getDockerImplementationByName(implementationName string) (docker dockerutils.Docker) {
	if implementationName == "commandExecutorDocker" {
		return dockerutils.MustGetLocalCommandExecutorDocker()
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
				docker := getDockerImplementationByName(tt.implementationName)

				hostDesciption, err := docker.GetHostDescription()
				require.NoError(t, err)

				require.EqualValues(t, "localhost", hostDesciption)
			},
		)
	}
}
