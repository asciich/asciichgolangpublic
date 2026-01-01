package containerinterfaces_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containerinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/nativedocker"
)

func Test_ContainerIsACommandExecutor(t *testing.T) {
	var container containerinterfaces.Container
	
	// nativedocker.Container is used as example here.
	container = &nativedocker.Container{}

	var commandExectuor commandexecutorinterfaces.CommandExecutor

	// This fails to compile if the containerinterfaces.Container does not implement the CommandExecutor interface:
	commandExectuor = container

	require.NotNil(t, commandExectuor)
}
