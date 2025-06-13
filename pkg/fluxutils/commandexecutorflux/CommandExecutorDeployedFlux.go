package commandexecutorflux

import (
	"github.com/asciich/asciichgolangpublic/commandexecutor"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type CommandExecutorDeployedFlux struct {
	commandExecutor   commandexecutor.CommandExecutor
	kubernetesCluster kubernetesinterfaces.KubernetesCluster
}

func (c *CommandExecutorDeployedFlux) GetCommandExecutor() (commandexecutor.CommandExecutor, error) {
	if c.commandExecutor == nil {
		return nil, tracederrors.TracedError("commandExecutor not set")
	}

	return c.commandExecutor, nil
}
