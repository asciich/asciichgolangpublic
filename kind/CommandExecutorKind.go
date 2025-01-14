package kind

import (
	"github.com/asciich/asciichgolangpublic"
	"github.com/asciich/asciichgolangpublic/errors"
	"github.com/asciich/asciichgolangpublic/kubernetes"
	"github.com/asciich/asciichgolangpublic/logging"

	aslices "github.com/asciich/asciichgolangpublic/datatypes/slices"
)

type CommandExecutorKind struct {
	commandExecutor asciichgolangpublic.CommandExecutor
}

func GetCommandExecutorKind(commandExecutor asciichgolangpublic.CommandExecutor) (kind Kind, err error) {
	if commandExecutor == nil {
		return nil, errors.TracedErrorNil("commandExectuor")
	}

	toReturn := NewCommandExecutorKind()

	err = toReturn.SetCommandExecutor(commandExecutor)
	if err != nil {
		return nil, err
	}

	return toReturn, nil
}

func GetLocalCommandExecutorKind() (kind Kind, err error) {
	return GetCommandExecutorKind(asciichgolangpublic.Bash())
}

func MustGetCommandExecutorKind(commandExecutor asciichgolangpublic.CommandExecutor) (kind Kind) {
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
		return false, errors.TracedErrorEmptyString("clusterName")
	}

	clusterNames, err := c.ListClusterNames(false)
	if err != nil {
		return false, err
	}

	exists = aslices.ContainsString(clusterNames, clusterName)

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

func (c *CommandExecutorKind) CreateClusterByName(clusterName string, verbose bool) (cluster kubernetes.KubernetesCluster, err error) {
	if clusterName == "" {
		return nil, errors.TracedErrorEmptyString("clusterName")
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

		_, err = commandExecutor.RunCommand(
			&asciichgolangpublic.RunCommandOptions{
				Command:            []string{"kind", "create", "cluster", "--name", clusterName},
				Verbose:            verbose,
				LiveOutputOnStdout: verbose,
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
		return errors.TracedErrorEmptyString("clusterName")
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

		_, err = commandExecutor.RunCommand(
			&asciichgolangpublic.RunCommandOptions{
				Command:            []string{"kind", "delete", "cluster", "--name", clusterName},
				Verbose:            verbose,
				LiveOutputOnStdout: verbose,
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

func (c *CommandExecutorKind) GetClusterByName(clusterName string) (cluster kubernetes.KubernetesCluster, err error) {
	if clusterName == "" {
		return nil, errors.TracedErrorEmptyString("clusterName")
	}

	toReturn := NewKindCluster()

	err = toReturn.SetName(clusterName)
	if err != nil {
		return nil, err
	}

	err = toReturn.SetKind(c)
	if err != nil {
		return nil, err
	}

	return toReturn, nil
}

func (c *CommandExecutorKind) GetCommandExecutor() (commandExecutor asciichgolangpublic.CommandExecutor, err error) {

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
		&asciichgolangpublic.RunCommandOptions{
			Command: []string{"kind", "get", "clusters"},
			Verbose: false,
		},
	)
}

func (c *CommandExecutorKind) MustClusterByNameExists(clusterName string, verbose bool) (exists bool) {
	exists, err := c.ClusterByNameExists(clusterName, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return exists
}

func (c *CommandExecutorKind) MustCreateClusterByName(clusterName string, verbose bool) (cluster kubernetes.KubernetesCluster) {
	cluster, err := c.CreateClusterByName(clusterName, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return cluster
}

func (c *CommandExecutorKind) MustDeleteClusterByName(clusterName string, verbose bool) {
	err := c.DeleteClusterByName(clusterName, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorKind) MustGetClusterByName(clusterName string) (cluster kubernetes.KubernetesCluster) {
	cluster, err := c.GetClusterByName(clusterName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return cluster
}

func (c *CommandExecutorKind) MustGetCommandExecutor() (commandExecutor asciichgolangpublic.CommandExecutor) {
	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return commandExecutor
}

func (c *CommandExecutorKind) MustGetHostDescription() (hostDescription string) {
	hostDescription, err := c.GetHostDescription()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return hostDescription
}

func (c *CommandExecutorKind) MustListClusterNames(verbose bool) (clusterNames []string) {
	clusterNames, err := c.ListClusterNames(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return clusterNames
}

func (c *CommandExecutorKind) MustRunCommand(runOptions *asciichgolangpublic.RunCommandOptions) (commandOutput *asciichgolangpublic.CommandOutput) {
	commandOutput, err := c.RunCommand(runOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return commandOutput
}

func (c *CommandExecutorKind) MustRunCommandAndGetStdoutAsLines(runOptions *asciichgolangpublic.RunCommandOptions) (lines []string) {
	lines, err := c.RunCommandAndGetStdoutAsLines(runOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return lines
}

func (c *CommandExecutorKind) MustSetCommandExecutor(commandExecutor asciichgolangpublic.CommandExecutor) {
	err := c.SetCommandExecutor(commandExecutor)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorKind) RunCommand(runOptions *asciichgolangpublic.RunCommandOptions) (commandOutput *asciichgolangpublic.CommandOutput, err error) {
	if runOptions == nil {
		return nil, errors.TracedErrorNil("runOptions")
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	return commandExecutor.RunCommand(runOptions)
}

func (c *CommandExecutorKind) RunCommandAndGetStdoutAsLines(runOptions *asciichgolangpublic.RunCommandOptions) (lines []string, err error) {
	if runOptions == nil {
		return nil, errors.TracedErrorNil("runOptions")
	}

	commandOutput, err := c.RunCommand(runOptions)
	if err != nil {
		return nil, err
	}

	return commandOutput.GetStdoutAsLines(false)
}

func (c *CommandExecutorKind) SetCommandExecutor(commandExecutor asciichgolangpublic.CommandExecutor) (err error) {
	c.commandExecutor = commandExecutor

	return nil
}
