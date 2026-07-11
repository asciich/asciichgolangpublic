package commandexecutorkubernetes

import (
	"context"
	"fmt"
	"reflect"
	"slices"
	"sort"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kuberneteserrors"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type CommandExecutorNamespace struct {
	name              string
	kubernetesCluster kubernetesinterfaces.KubernetesCluster
}

func (c *CommandExecutorNamespace) GetKubernetesCluster() (kubernetesinterfaces.KubernetesCluster, error) {
	if c.kubernetesCluster == nil {
		return nil, tracederrors.TracedError("kubernetesCluster not set")
	}

	return c.kubernetesCluster, nil
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

	roleNames, err := c.ListRoleNames(contextutils.WithSilent(ctx))
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

func (c *CommandExecutorNamespace) RunCommand(ctx context.Context, runCommandOptions *parameteroptions.RunCommandOptions) (commandOutput *commandoutput.CommandOutput, err error) {
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

func (c *CommandExecutorNamespace) DeleteSecretByName(ctx context.Context, secretName string) error {
	if secretName == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	contextName, err := c.GetCachedKubectlContext(ctx)
	if err != nil {
		return err
	}

	namespaceName, err := c.GetName()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Delete secret '%s' in namespace '%s' of kubernetes '%s' started.", secretName, namespaceName, contextName)

	_, err = c.RunCommandAndGetStdoutAsLines(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{
				"kubectl",
				"--context",
				contextName,
				"--namespace",
				namespaceName,
				"delete",
				"secret",
				secretName,
			},
		},
	)
	if err == nil {
		logging.LogChangedByCtxf(ctx, "Secret '%s' in namespace '%s' of kubernetes '%s' deleted.", secretName, namespaceName, contextName)
	} else {
		notFoundMessage := fmt.Sprintf("Error from server (NotFound): secrets \"%s\" not found", secretName)
		if strings.Contains(err.Error(), notFoundMessage) {
			logging.LogInfoByCtxf(ctx, "Secret '%s' in namespace '%s' of kubernetes '%s' already absent. Skip delete.", secretName, namespaceName, contextName)
		} else {
			return err
		}
	}

	logging.LogInfoByCtxf(ctx, "Delete secret '%s' in namespace '%s' of kubernetes '%s' finished.", secretName, namespaceName, contextName)

	return nil
}

func (c *CommandExecutorNamespace) SecretByNameExists(ctx context.Context, secretName string) (bool, error) {
	if secretName == "" {
		return false, tracederrors.TracedErrorEmptyString("name")
	}

	contextName, err := c.GetCachedKubectlContext(ctx)
	if err != nil {
		return false, err
	}

	namespaceName, err := c.GetName()
	if err != nil {
		return false, err
	}

	var exists bool
	_, err = c.RunCommandAndGetStdoutAsLines(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{
				"kubectl",
				"--context",
				contextName,
				"--namespace",
				namespaceName,
				"get",
				"secret",
				secretName,
			},
		},
	)
	if err == nil {
		exists = true
	} else {
		expectedNotFoundMessage := fmt.Sprintf("Error from server (NotFound): secrets \"%s\" not found", secretName)
		if !strings.Contains(err.Error(), expectedNotFoundMessage) {
			return false, err
		}
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Secret '%s' in namespace '%s' of kubernetes '%s' exists.", secretName, namespaceName, contextName)
	} else {
		logging.LogInfoByCtxf(ctx, "Secret '%s' in namespace '%s' of kubernetes '%s' does not exist.", secretName, namespaceName, contextName)
	}

	return exists, nil
}

func (c *CommandExecutorNamespace) ListSecretNames(ctx context.Context) ([]string, error) {
	contextName, err := c.GetCachedKubectlContext(ctx)
	if err != nil {
		return nil, err
	}

	namespaceName, err := c.GetName()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "List secret names in namespace '%s' of kubernetes '%s' started.", namespaceName, contextName)

	lines, err := c.RunCommandAndGetStdoutAsLines(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{
				"kubectl",
				"--context",
				contextName,
				"--namespace",
				namespaceName,
				"get",
				"secrets",
				"-o",
				"name",
			},
		},
	)
	if err != nil {
		return nil, err
	}

	names := []string{}
	for _, l := range lines {
		if !strings.HasPrefix(l, "secret/") {
			continue
		}

		names = append(names, strings.TrimPrefix(l, "secret/"))
	}


	sort.Strings(names)

	logging.LogInfoByCtxf(ctx, "List secret names in namespace '%s' of kubernetes '%s' finished.", namespaceName, contextName)

	return names, err
}

func (c *CommandExecutorNamespace) GetSecretByName(name string) (secret kubernetesinterfaces.Secret, err error) {
	if name == "" {
		return nil, tracederrors.TracedErrorEmptyString("name")
	}

	return &CommandExecutorSecret{
		name:      name,
		namespace: c,
	}, nil
}

func (c *CommandExecutorNamespace) CreateSecret(ctx context.Context, secretName string, options *kubernetesparameteroptions.CreateSecretOptions) (createdSecret kubernetesinterfaces.Secret, err error) {
	if secretName == "" {
		return nil, tracederrors.TracedErrorEmptyString("secretName")
	}

	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	contextName, err := c.GetCachedKubectlContext(ctx)
	if err != nil {
		return nil, err
	}

	namespaceName, err := c.GetName()
	if err != nil {
		return nil, err
	}

	secretData, err := options.GetSecretData()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Create secret '%s' in namespace '%s' of kubernetes '%s' started.", secretName, namespaceName, contextName)

	err = c.Create(ctx)
	if err != nil {
		return nil, err
	}

	secret, err := c.GetSecretByName(secretName)
	if err != nil {
		return nil, err
	}

	var create bool
	currentData, err := secret.Read(ctx)
	if err != nil {
		if kuberneteserrors.IsSecretNotFoundError(err) {
			create = true
		} else {
			return nil, err
		}
	}

	if currentData != nil {
		if reflect.DeepEqual(currentData, secretData) {
			logging.LogInfoByCtxf(ctx, "Secret '%s' in namespace '%s' of kubernetes '%s' created.", secretName, namespaceName, contextName)
		} else {
			err = secret.Delete(ctx)
			if err != nil {
				return nil, err
			}
			create = true
		}
	}

	if create {
		command := []string{
			"kubectl",
			"--context",
			contextName,
			"--namespace",
			namespaceName,
			"create",
			"secret",
			"generic",
			secretName,
		}

		for key, value := range secretData {
			command = append(command, fmt.Sprintf("--from-literal=%s=%s", key, string(value)))
		}

		_, err = c.RunCommand(
			contextutils.WithSilent(ctx), // do not expose secrets
			&parameteroptions.RunCommandOptions{
				Command: command,
			},
		)
		if err != nil {
			return nil, tracederrors.TracedErrorf("Failed to create secret '%s' in namespace '%s' of kubernetes '%s': %w", secretName, namespaceName, contextName, err)
		}

		logging.LogChangedByCtxf(ctx, "Secret '%s' in namespace '%s' of kubernetes '%s' created.", secretName, namespaceName, contextName)
	}

	logging.LogInfoByCtxf(ctx, "Create secret '%s' in namespace '%s' of kubernetes '%s' finished.", secretName, namespaceName, contextName)

	return secret, nil
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
