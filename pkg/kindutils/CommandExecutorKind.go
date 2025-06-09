package kindutils

import (
	"context"
	"os"
	"slices"

	"github.com/asciich/asciichgolangpublic/commandexecutor"
	"github.com/asciich/asciichgolangpublic/files"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/continuousintegration"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubeconfigutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"
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

func (c *CommandExecutorKind) ClusterByNameExists(ctx context.Context, clusterName string) (exists bool, err error) {
	if clusterName == "" {
		return false, tracederrors.TracedErrorEmptyString("clusterName")
	}

	clusterNames, err := c.ListClusterNames(contextutils.WithSilent(ctx))
	if err != nil {
		return false, err
	}

	exists = slices.Contains(clusterNames, clusterName)

	hostDescription, err := c.GetHostDescription()
	if err != nil {
		return false, err
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Kind cluster '%s' on host '%s' exists.", clusterName, hostDescription)
	} else {
		logging.LogInfoByCtxf(ctx, "Kind cluster '%s' on host '%s' does not exist.", clusterName, hostDescription)
	}

	return exists, nil
}

func (c *CommandExecutorKind) EnsureKubectlConfigPresent(ctx context.Context, clusterName string) error {
	if clusterName == "" {
		return tracederrors.TracedErrorEmptyString("clusterName")
	}

	logging.LogInfoByCtxf(ctx, "Enusre kubectl config present for kind cluster '%s' started.", clusterName)

	if os.Getenv("KUBECONFIG") != "" {
		return tracederrors.TracedErrorf("Not implemented when 'KUBECONFIG' env var is set. But KUBECONFIG is set to '%s'.", os.Getenv("KUBECONFIG"))
	}

	path, err := kubeconfigutils.GetDefaultKubeConfigPath(ctx)
	if err != nil {
		return err
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return err
	}

	kubeConfigFile, err := files.GetCommandExecutorFileByPath(commandExecutor, path)
	if err != nil {
		return err
	}

	exists, err := kubeConfigFile.Exists(contextutils.GetVerboseFromContext(ctx))
	if err != nil {
		return err
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Kube config file '%s' is already present.", path)
	} else {
		config, err := commandExecutor.RunCommandAndGetStdoutAsString(ctx, &parameteroptions.RunCommandOptions{
			Command: []string{"kind", "get", "kubeconfig", "--name", clusterName},
		})
		if err != nil {
			return err
		}

		err = kubeConfigFile.WriteString(config, contextutils.GetVerboseFromContext(ctx))
		if err != nil {
			return err
		}

		logging.LogChangedByCtxf(ctx, "Created kube config file '%s' with config to access kind cluster '%s'.", kubeConfigFile, clusterName)
	}

	logging.LogInfoByCtxf(ctx, "Enusre kubectl config present for kind cluster '%s' finished.", clusterName)

	return nil
}

func (c *CommandExecutorKind) CreateClusterByName(ctx context.Context, clusterName string) (cluster kubernetesinterfaces.KubernetesCluster, err error) {
	if clusterName == "" {
		return nil, tracederrors.TracedErrorEmptyString("clusterName")
	}

	exists, err := c.ClusterByNameExists(ctx, clusterName)
	if err != nil {
		return nil, err
	}

	hostDescription, err := c.GetHostDescription()
	if err != nil {
		return nil, err
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Kind cluster '%s' on host '%s' already exists.", clusterName, hostDescription)
	} else {
		commandExecutor, err := c.GetCommandExecutor()
		if err != nil {
			return nil, err
		}

		logging.LogInfoByCtxf(ctx, "Going to create kind cluster '%s'. This may take a while...", clusterName)

		_, err = commandExecutor.RunCommand(
			commandexecutor.WithLiveOutputOnStdout(ctx),
			&parameteroptions.RunCommandOptions{
				Command: []string{"kind", "create", "cluster", "--name", clusterName, "2>&1"},
			},
		)
		if err != nil {
			if continuousintegration.IsRunningInContinuousIntegration() {
				logging.LogInfoByCtxf(ctx, "Retry kind cluster '%s' creation in CI.", clusterName)

				// Show known clusters for debugging
				knownClusterNames, err := c.ListClusterNames(ctx)
				if err != nil {
					return nil, err
				}

				msg := "Known cluster names before retry:\n\n"
				for _, knowCluster := range knownClusterNames {
					msg += " - '" +knowCluster+ "'\n"
				}
				logging.LogInfoByCtx(ctx, msg)

				err = c.DeleteClusterByName(ctx, clusterName)
				if err != nil {
					return nil, err
				}

				_, err = commandExecutor.RunCommand(
					commandexecutor.WithLiveOutputOnStdout(ctx),
					&parameteroptions.RunCommandOptions{
						Command: []string{"kind", "create", "cluster", "--name", clusterName, "2>&1"},
					},
				)
				if err != nil {
					return nil, err
				}
			}

			return nil, err
		}

		logging.LogChangedByCtxf(ctx, "Kind cluster '%s' on host '%s' created.", clusterName, hostDescription)
	}

	err = c.EnsureKubectlConfigPresent(ctx, clusterName)
	if err != nil {
		return nil, err
	}

	cluster, err = c.GetClusterByName(clusterName)
	if err != nil {
		return nil, err
	}

	err = cluster.CheckAccessible(ctx)
	if err != nil {
		return nil, err
	}

	return cluster, nil
}

func (c *CommandExecutorKind) DeleteClusterByName(ctx context.Context, clusterName string) (err error) {
	if clusterName == "" {
		return tracederrors.TracedErrorEmptyString("clusterName")
	}

	exists, err := c.ClusterByNameExists(ctx, clusterName)
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
			commandexecutor.WithLiveOutputOnStdout(ctx),
			&parameteroptions.RunCommandOptions{
				Command: []string{"kind", "delete", "cluster", "--name", clusterName},
			},
		)
		if err != nil {
			return err
		}

		logging.LogChangedByCtxf(ctx, "Kind cluster '%s' on host '%s' deleted.", clusterName, hostDescription)
	} else {
		logging.LogInfoByCtxf(ctx, "Kind cluster '%s' on host '%s' already absent.", clusterName, hostDescription)
	}

	return nil
}

func (c *CommandExecutorKind) GetClusterByName(clusterName string) (cluster kubernetesinterfaces.KubernetesCluster, err error) {
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

func (c *CommandExecutorKind) ListClusterNames(ctx context.Context) (clusterNames []string, err error) {
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
