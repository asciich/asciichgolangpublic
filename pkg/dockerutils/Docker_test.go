package dockerutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/commandexecutordocker"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/dockerinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/nativedocker"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func getDockerImplementationByName(implementationName string) (docker dockerinterfaces.Docker) {
	if implementationName == "commandExecutorDocker" {
		return mustutils.Must(commandexecutordocker.GetLocalCommandExecutorDocker())
	}

	if implementationName == "nativeDocker" {
		return nativedocker.NewDocker()
	}

	logging.LogFatalWithTracef("Unknown implementation name '%s'", implementationName)
	return nil
}

func TestDocker_GetHostDescription(t *testing.T) {

	tests := []struct {
		implementationName string
	}{
		{"commandExecutorDocker"},
		{"nativeDocker"},
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

func Test_ListContainerNames(t *testing.T) {
	t.Run("", func(t *testing.T) {
		list, err := dockerutils.ListContainerNames(getCtx())
		require.NoError(t, err)
		require.NotNil(t, list)
	})
}
