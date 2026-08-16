package commandexecutorkind

import (
	"context"
	"slices"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils/kindparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type CommandExecutorKind struct {
	commandExecutor commandexecutorinterfaces.CommandExecutor
}

func GetCommandExecutorKind(commandExecutor commandexecutorinterfaces.CommandExecutor) (kind kindutils.Kind, err error) {
	if commandExecutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExecutor")
	}

	ret := NewCommandExecutorKind()

	err = ret.SetCommandExecutor(commandExecutor)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func NewCommandExecutorKind() (c *CommandExecutorKind) {
	return new(CommandExecutorKind)
}

func (c *CommandExecutorKind) ClusterByNameExists(ctx context.Context, clusterName string) (exists bool, err error) {
	if clusterName == "" {
		return false, tracederrors.TracedErrorEmptyString("clusterName")
	}

	clusterNames, err := c.ListClusterNames(ctx)
	if err != nil {
		return false, err
	}

	exists = slices.Contains(clusterNames, clusterName)

	return exists, nil
}

func (c *CommandExecutorKind) CreateCluster(ctx context.Context, options *kindparameteroptions.CreateClusterOptions) (cluster kubernetesinterfaces.KubernetesCluster, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (c *CommandExecutorKind) DeleteClusterByName(ctx context.Context, clusterName string) (err error) {
	if clusterName == "" {
		return tracederrors.TracedErrorEmptyString("clusterName")
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return err
	}

	_, err = commandExecutor.RunCommand(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{"kind", "delete", "cluster", "--name", clusterName},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *CommandExecutorKind) GetClusterByName(clusterName string) (cluster kubernetesinterfaces.KubernetesCluster, err error) {
	if clusterName == "" {
		return nil, tracederrors.TracedErrorEmptyString("clusterName")
	}

	return nil, tracederrors.TracedErrorNotImplemented()
}

func (c *CommandExecutorKind) ListClusterNames(ctx context.Context) (clusterNames []string, err error) {
	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	output, err := commandExecutor.RunCommandAndGetStdoutAsString(ctx, &parameteroptions.RunCommandOptions{
		Command: []string{"kind", "get", "clusters"},
	})
	if err != nil {
		return nil, err
	}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			clusterNames = append(clusterNames, line)
		}
	}

	return clusterNames, nil
}

func (c *CommandExecutorKind) GetCommandExecutor() (commandExecutor commandexecutorinterfaces.CommandExecutor, err error) {
	if c.commandExecutor == nil {
		return nil, tracederrors.TracedError("commandExecutor not set")
	}
	return c.commandExecutor, nil
}

func (c *CommandExecutorKind) SetCommandExecutor(commandExecutor commandexecutorinterfaces.CommandExecutor) (err error) {
	c.commandExecutor = commandExecutor
	return nil
}
