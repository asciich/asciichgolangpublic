package commandexecutorflux

import (
	"context"

	"github.com/asciich/asciichgolangpublic/commandexecutor"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"
	"github.com/asciich/asciichgolangpublic/tracederrors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type CommandExecutorDeployedFlux struct {
	commandExecutor commandexecutor.CommandExecutor
	cluster         kubernetesinterfaces.KubernetesCluster
}

func (c *CommandExecutorDeployedFlux) GetCommandExecutor() (commandexecutor.CommandExecutor, error) {
	if c.commandExecutor == nil {
		return nil, tracederrors.TracedError("commandExecutor not set")
	}

	return c.commandExecutor, nil

}

func (c *CommandExecutorDeployedFlux) GitRepositoryExists(ctx context.Context, name string, namespace string) (bool, error) {
	return false, tracederrors.TracedErrorNotImplemented()
}

func (c *CommandExecutorDeployedFlux) GetGitRepositoryStatusMessage(ctx context.Context, name string, namespace string) (string, error) {
	return "", tracederrors.TracedErrorNotImplemented()
}

func (c *CommandExecutorDeployedFlux) DeleteGitRepository(ctx context.Context, name string, namespace string) error {
	return tracederrors.TracedErrorNotImplemented()
}

func (c *CommandExecutorDeployedFlux) WatchGitRepository(ctx context.Context, name string, namespace string, create func(*unstructured.Unstructured), update func(*unstructured.Unstructured), delete func(*unstructured.Unstructured)) error {
	return tracederrors.TracedErrorNotImplemented()
}
