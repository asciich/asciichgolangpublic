package commandexecutorkubernetes

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/asciich/asciichgolangpublic/datatypes"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type CommandExecutorNamespace struct {
	name              string
	kubernetesCluster kubernetesinterfaces.KubernetesCluster
}

func NewCommandExecutorNamespace() (c *CommandExecutorNamespace) {
	return new(CommandExecutorNamespace)
}

func (c *CommandExecutorNamespace) Create(ctx context.Context) (err error) {
	name, err := c.GetName()
	if err != nil {
		return err
	}

	cluster, err := c.GetKubernetesCluster()
	if err != nil {
		return err
	}

	_, err = cluster.CreateNamespaceByName(ctx, name)
	if err != nil {
		return err
	}

	return nil
}

func (c *CommandExecutorNamespace) CreateRole(ctx context.Context, createOptions *kubernetesparameteroptions.CreateRoleOptions) (createdRole kubernetesinterfaces.Role, err error) {
	if createOptions == nil {
		return nil, tracederrors.TracedErrorNil("createOptions")
	}

	roleName, err := createOptions.GetName()
	if err != nil {
		return nil, err
	}

	namespaceName, err := c.GetName()
	if err != nil {
		return nil, err
	}

	cluserName, err := c.GetClusterName()
	if err != nil {
		return nil, err
	}

	context, err := c.GetCachedKubectlContext(ctx)
	if err != nil {
		return nil, err
	}

	exists, err := c.RoleByNameExists(ctx, roleName)
	if err != nil {
		return nil, err
	}

	if exists {
		logging.LogInfof(
			"Role '%s' in namespace '%s' in kubernetes cluster '%s' already exists.",
			roleName,
			namespaceName,
			cluserName,
		)
	} else {
		command := []string{
			"kubectl",
			"--context",
			context,
			"--namespace",
			namespaceName,
			"create",
			"role",
			roleName,
		}

		verbs := createOptions.Verbs
		if len(verbs) > 0 {
			command = append(
				command,
				fmt.Sprintf(
					"--verb=%s",
					strings.Join(verbs, ","),
				),
			)
		}

		resources := createOptions.Resorces
		if len(resources) > 0 {
			command = append(
				command,
				fmt.Sprintf(
					"--resource=%s",
					strings.Join(resources, ","),
				),
			)
		}

		_, err = c.RunCommand(
			ctx,
			&parameteroptions.RunCommandOptions{
				Command: command,
			},
		)
		if err != nil {
			return nil, err
		}
	}

	return c.GetRoleByName(roleName)
}

func (c *CommandExecutorNamespace) DeleteRoleByName(ctx context.Context, name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	namespaceName, err := c.GetName()
	if err != nil {
		return err
	}

	exists, err := c.RoleByNameExists(ctx, name)
	if err != nil {
		return err
	}

	clusterName, err := c.GetClusterName()
	if err != nil {
		return err
	}

	if exists {
		context, err := c.GetCachedKubectlContext(ctx)
		if err != nil {
			return err
		}

		_, err = c.RunCommand(
			ctx,
			&parameteroptions.RunCommandOptions{
				Command: []string{
					"kubectl",
					"--context",
					context,
					"--namespace",
					namespaceName,
					"delete",
					"role",
					name,
				},
			},
		)
		if err != nil {
			return err
		}

		logging.LogChangedByCtxf(ctx, "Role '%s' in namespace '%s' in kubernetes cluster '%s' deleted.", name, namespaceName, clusterName)
	} else {
		logging.LogChangedByCtxf(ctx, "Role '%s' in namespace '%s' in kubernetes cluster '%s' already absent.", name, namespaceName, clusterName)
	}

	return nil
}

func (c *CommandExecutorNamespace) GetCachedKubectlContext(ctx context.Context) (context string, err error) {
	kubernetes, err := c.GetKubernetesCluster()
	if err != nil {
		return "", err
	}

	commandExecutorKubernetes, ok := kubernetes.(*CommandExecutorKubernetes)
	if !ok {
		typeName, err := datatypes.GetTypeName(kubernetes)
		if err != nil {
			return "", err
		}

		return "", tracederrors.TracedErrorNilf(
			"Unable to get kubectl context. unexpected kubernetes type '%s'",
			typeName,
		)
	}

	return commandExecutorKubernetes.GetCachedKubectlContext(ctx)
}

func (c *CommandExecutorNamespace) GetClusterName() (clusterName string, err error) {
	kubernetesCluster, err := c.GetKubernetesCluster()
	if err != nil {
		return "", err
	}

	return kubernetesCluster.GetName()
}

func (c *CommandExecutorNamespace) GetCommandExecutor() (commandExecutor commandexecutorinterfaces.CommandExecutor, err error) {
	kubernetes, err := c.GetKubernetesCluster()
	if err != nil {
		return nil, err
	}

	commandExecutorKubernetes, ok := kubernetes.(*CommandExecutorKubernetes)
	if !ok {
		typeName, err := datatypes.GetTypeName(kubernetes)
		if err != nil {
			return nil, err
		}

		return nil, tracederrors.TracedErrorNilf(
			"Unable to get command executor. unexpected kubernetes type '%s'",
			typeName,
		)
	}

	return commandExecutorKubernetes.GetCommandExecutor()
}

func (c *CommandExecutorNamespace) GetKubectlContext(ctx context.Context) (contextName string, err error) {
	cluster, err := c.GetKubernetesCluster()
	if err != nil {
		return "", err
	}

	return cluster.GetKubectlContext(ctx)
}

func (c *CommandExecutorNamespace) GetKubernetesCluster() (kubernetesCluster kubernetesinterfaces.KubernetesCluster, err error) {

	return c.kubernetesCluster, nil
}

func (c *CommandExecutorNamespace) GetName() (name string, err error) {
	if c.name == "" {
		return "", tracederrors.TracedErrorf("name not set")
	}

	return c.name, nil
}

func (c *CommandExecutorNamespace) GetObjectByNames(objectName string, objectType string) (object kubernetesinterfaces.Object, err error) {
	if objectName == "" {
		return nil, tracederrors.TracedErrorEmptyString("objectName")
	}

	if objectType == "" {
		return nil, tracederrors.TracedErrorEmptyString("objectType")
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	return GetCommandExecutorObject(commandExecutor, c, objectName, objectType)
}

func (c *CommandExecutorNamespace) GetRoleByName(name string) (role kubernetesinterfaces.Role, err error) {
	if name == "" {
		return nil, tracederrors.TracedErrorEmptyString("name")
	}

	toReturn := NewCommandExecutorRole()

	err = toReturn.SetName(name)
	if err != nil {
		return nil, err
	}

	err = toReturn.SetNamespace(c)
	if err != nil {
		return nil, err
	}

	return toReturn, nil
}

func (c *CommandExecutorNamespace) ListRoleNames(ctx context.Context) (roleNames []string, err error) {
	context, err := c.GetCachedKubectlContext(ctx)
	if err != nil {
		return nil, err
	}

	name, err := c.GetName()
	if err != nil {
		return nil, err
	}

	lines, err := c.RunCommandAndGetStdoutAsLines(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{
				"kubectl",
				"--context",
				context,
				"--namespace",
				name,
				"get",
				"roles",
				"-o",
				"name",
			},
		},
	)
	if err != nil {
		return nil, err
	}

	roleNames = []string{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		splitted := strings.Split(line, "/")
		if len(splitted) != 2 {
			return nil, tracederrors.TracedErrorf(
				"Unable to get role name out of line='%s'.",
				line,
			)
		}

		roleName := splitted[1]
		if roleName == "" {
			return nil, tracederrors.TracedErrorf(
				"roleName is empty stiring after evaluation of line '%s'",
				line,
			)
		}

		roleNames = append(roleNames, roleName)
	}

	return roleNames, nil
}

func (c *CommandExecutorNamespace) RoleByNameExists(ctx context.Context, name string) (exists bool, err error) {
	if name == "" {
		return false, tracederrors.TracedErrorEmptyString("name")
	}

	roleNames, err := c.ListRoleNames(contextutils.WithVerbosityContextByBool(ctx, false))
	if err != nil {
		return false, err
	}

	exists = slices.Contains(roleNames, name)

	clusterName, err := c.GetClusterName()
	if err != nil {
		return false, err
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Role '%s' in kubernetes cluster '%s' exists.", name, clusterName)
	} else {
		logging.LogInfoByCtxf(ctx, "Role '%s' in kubernetes cluster '%s' does not exist.", name, clusterName)
	}

	return exists, nil
}

func (c *CommandExecutorNamespace) RunCommand(ctx context.Context, runCommandOptions *parameteroptions.RunCommandOptions) (commandOutput *commandexecutorgeneric.CommandOutput, err error) {
	if runCommandOptions == nil {
		return nil, tracederrors.TracedErrorNil("runCommandOptions")
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	return commandExecutor.RunCommand(ctx, runCommandOptions)
}

func (c *CommandExecutorNamespace) RunCommandAndGetStdoutAsLines(ctx context.Context, runCommandOptions *parameteroptions.RunCommandOptions) (lines []string, err error) {
	if runCommandOptions == nil {
		return nil, tracederrors.TracedErrorNil("runCommandOptions")
	}

	commandOutput, err := c.RunCommand(ctx, runCommandOptions)
	if err != nil {
		return nil, err
	}

	return commandOutput.GetStdoutAsLines(false)
}

func (c *CommandExecutorNamespace) SetKubernetesCluster(kubernetesCluster kubernetesinterfaces.KubernetesCluster) (err error) {
	c.kubernetesCluster = kubernetesCluster

	return nil
}

func (c *CommandExecutorNamespace) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorf("name is empty string")
	}

	c.name = name

	return nil
}

func (c *CommandExecutorNamespace) DeleteSecretByName(ctx context.Context, name string) (err error) {
	return tracederrors.TracedErrorNotImplemented()
}

func (c *CommandExecutorNamespace) SecretByNameExists(ctx context.Context, name string) (bool, error) {
	return false, tracederrors.TracedErrorNotImplemented()
}

func (c *CommandExecutorNamespace) GetSecretByName(name string) (secret kubernetesinterfaces.Secret, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (c *CommandExecutorNamespace) CreateSecret(ctx context.Context, name string, options *kubernetesparameteroptions.CreateSecretOptions) (createdSecret kubernetesinterfaces.Secret, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (c *CommandExecutorNamespace) DeleteConfigMapByName(ctx context.Context, name string) (err error) {
	return tracederrors.TracedErrorNotImplemented()
}

func (c *CommandExecutorNamespace) ConfigMapByNameExists(ctx context.Context, name string) (bool, error) {
	return false, tracederrors.TracedErrorNotImplemented()
}

func (c *CommandExecutorNamespace) GetConfigMapByName(name string) (secret kubernetesinterfaces.ConfigMap, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (c *CommandExecutorNamespace) CreateConfigMap(ctx context.Context, name string, options *kubernetesparameteroptions.CreateConfigMapOptions) (createdConfigMap kubernetesinterfaces.ConfigMap, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (c *CommandExecutorNamespace) WatchConfigMap(ctx context.Context, name string, onCreate func(kubernetesinterfaces.ConfigMap), onUpdate func(kubernetesinterfaces.ConfigMap), onDelete func(kubernetesinterfaces.ConfigMap)) error {
	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	return tracederrors.TracedErrorNotImplemented()
}

func (c *CommandExecutorNamespace) WaitUntilAllPodsInNamespaceAreRunning(ctx context.Context, options *kubernetesparameteroptions.WaitForPodsOptions) error {
	return tracederrors.TracedErrorNotImplemented()
}

func (c *CommandExecutorNamespace) CreateObject(ctx context.Context, options *kubernetesparameteroptions.CreateObjectOptions) (kubernetesinterfaces.Object, error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (c *CommandExecutorNamespace) Exists(ctx context.Context) (bool, error) {
	namespaceName, err := c.GetName()
	if err != nil {
		return false, err
	}

	cluster, err := c.GetKubernetesCluster()
	if err != nil {
		return false, err
	}

	return cluster.NamespaceByNameExists(ctx, namespaceName)
}
