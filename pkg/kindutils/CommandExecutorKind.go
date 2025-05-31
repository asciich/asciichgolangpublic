package kindutils

import (
	"context"
	"slices"

	"github.com/asciich/asciichgolangpublic/commandexecutor"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type CommandExecutorKind struct {
	commandExecutor commandexecutor.CommandExecutor
}

func GetCommandExecutorKind(commandExecutor commandexecutor.CommandExecutor) (kind Kind, err error) {
	if commandExecutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExectuor")
	}

	toReturn := NewCommandExecutorKind()

	err = toReturn.SetCommandExecutor(commandExecutor)
	if err != nil {
		return nil, err
	}

	return toReturn, nil
}

func GetLocalCommandExecutorKind() (kind Kind, err error) {
	return GetCommandExecutorKind(commandexecutor.Bash())
}

func MustGetCommandExecutorKind(commandExecutor commandexecutor.CommandExecutor) (kind Kind) {
	kind, err := GetCommandExecutorKind(commandExecutor)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return kind
}

func MustGetLocalCommandExecutorKind() (kind Kind) {
	kind, err := GetLocalCommandExecutorKind()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return kind
}

func NewCommandExecutorKind() (c *CommandExecutorKind) {
	return new(CommandExecutorKind)
}

func (c *CommandExecutorKind) ClusterByNameExists(clusterName string, verbose bool) (exists bool, err error) {
	if clusterName == "" {
		return false, tracederrors.TracedErrorEmptyString("clusterName")
	}

	clusterNames, err := c.ListClusterNames(false)
	if err != nil {
		return false, err
	}

	exists = slices.Contains(clusterNames, clusterName)

	if verbose {
		hostDescription, err := c.GetHostDescription()
		if err != nil {
			return false, err
		}

		if exists {
			logging.LogInfof(
				"Kind cluster '%s' on host '%s' exists.",
				clusterName,
				hostDescription,
			)
		} else {
			logging.LogInfof(
				"Kind cluster '%s' on host '%s' does not exist.",
				clusterName,
				hostDescription,
			)
		}
	}

	return exists, nil
}

func (c *CommandExecutorKind) CreateClusterByName(clusterName string, verbose bool) (cluster kubernetesutils.KubernetesCluster, err error) {
	if clusterName == "" {
		return nil, tracederrors.TracedErrorEmptyString("clusterName")
	}

	exists, err := c.ClusterByNameExists(clusterName, false)
	if err != nil {
		return nil, err
	}

	hostDescription, err := c.GetHostDescription()
	if err != nil {
		return nil, err
	}

	if exists {
		if verbose {
			logging.LogInfof(
				"Kind cluster '%s' on host '%s' already exists.",
				clusterName,
				hostDescription,
			)

		}
	} else {
		commandExecutor, err := c.GetCommandExecutor()
		if err != nil {
			return nil, err
		}

		if verbose {
			logging.LogInfof(
				"Going to create kind cluster '%s'. This may take a while...",
				clusterName,
			)
		}

		ctx := contextutils.GetVerbosityContextByBool(verbose)
		_, err = commandExecutor.RunCommand(
			commandexecutor.WithLiveOutputOnStdout(ctx),
			&parameteroptions.RunCommandOptions{
				Command: []string{"kind", "create", "cluster", "--name", clusterName},
			},
		)
		if err != nil {
			return nil, err
		}

		if verbose {
			logging.LogChangedf(
				"Kind cluster '%s' on host '%s' created.",
				clusterName,
				hostDescription,
			)
		}
	}

	return c.GetClusterByName(clusterName)
}

func (c *CommandExecutorKind) DeleteClusterByName(clusterName string, verbose bool) (err error) {
	if clusterName == "" {
		return tracederrors.TracedErrorEmptyString("clusterName")
	}

	exists, err := c.ClusterByNameExists(clusterName, false)
	if err != nil {
		return err
	}

	hostDescription, err := c.GetHostDescription()
	if err != nil {
		return err
	}

	if exists {
		commandExecutor, err := c.GetCommandExecutor()
		if err != nil {
			return err
		}

		ctx := contextutils.GetVerbosityContextByBool(verbose)
		_, err = commandExecutor.RunCommand(
			commandexecutor.WithLiveOutputOnStdout(ctx),
			&parameteroptions.RunCommandOptions{
				Command: []string{"kind", "delete", "cluster", "--name", clusterName},
			},
		)
		if err != nil {
			return err
		}

		if verbose {
			logging.LogChangedf(
				"Kind cluster '%s' on host '%s' deleted.",
				clusterName,
				hostDescription,
			)
		}
	} else {
		if verbose {
			logging.LogInfof(
				"Kind cluster '%s' on host '%s' already absent.",
				clusterName,
				hostDescription,
			)

		}
	}

	return nil
}

func (c *CommandExecutorKind) GetClusterByName(clusterName string) (cluster kubernetesutils.KubernetesCluster, err error) {
	if clusterName == "" {
		return nil, tracederrors.TracedErrorEmptyString("clusterName")
	}

	toReturn := NewCommandExecutorKindCluster()

	err = toReturn.SetName("kind-" + clusterName)
	if err != nil {
		return nil, err
	}

	err = toReturn.SetKind(c)
	if err != nil {
		return nil, err
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	err = toReturn.SetCommandExecutor(commandExecutor)
	if err != nil {
		return nil, err
	}

	return toReturn, nil
}

func (c *CommandExecutorKind) GetCommandExecutor() (commandExecutor commandexecutor.CommandExecutor, err error) {
	if c.commandExecutor == nil {
		return nil, tracederrors.TracedError("commandExecutor not set")
	}
	return c.commandExecutor, nil
}

func (c *CommandExecutorKind) GetHostDescription() (hostDescription string, err error) {
	commandExector, err := c.GetCommandExecutor()
	if err != nil {
		return "", err
	}

	return commandExector.GetHostDescription()
}

func (c *CommandExecutorKind) ListClusterNames(verbose bool) (clusterNames []string, err error) {
	return c.RunCommandAndGetStdoutAsLines(
		contextutils.ContextSilent(),
		&parameteroptions.RunCommandOptions{
			Command: []string{"kind", "get", "clusters"},
		},
	)
}

func (c *CommandExecutorKind) RunCommand(ctx context.Context, runOptions *parameteroptions.RunCommandOptions) (commandOutput *commandexecutor.CommandOutput, err error) {
	if runOptions == nil {
		return nil, tracederrors.TracedErrorNil("runOptions")
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	return commandExecutor.RunCommand(ctx, runOptions)
}

func (c *CommandExecutorKind) RunCommandAndGetStdoutAsLines(ctx context.Context, runOptions *parameteroptions.RunCommandOptions) (lines []string, err error) {
	if runOptions == nil {
		return nil, tracederrors.TracedErrorNil("runOptions")
	}

	commandOutput, err := c.RunCommand(ctx, runOptions)
	if err != nil {
		return nil, err
	}

	return commandOutput.GetStdoutAsLines(false)
}

func (c *CommandExecutorKind) SetCommandExecutor(commandExecutor commandexecutor.CommandExecutor) (err error) {
	c.commandExecutor = commandExecutor

	return nil
}
