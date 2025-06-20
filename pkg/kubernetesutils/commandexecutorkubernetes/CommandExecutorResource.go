package commandexecutorkubernetes

import (
	"context"
	"strings"

	"github.com/asciich/asciichgolangpublic/commandexecutor"
	"github.com/asciich/asciichgolangpublic/fileformats/yamlutils"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type CommandExecutorResource struct {
	commandExecutor commandexecutor.CommandExecutor
	name            string
	typeName        string
	namespace       kubernetesinterfaces.Namespace
}

func GetCommandExecutorResource(commandExectutor commandexecutor.CommandExecutor, namespace kubernetesinterfaces.Namespace, resourceName string, resourceType string) (resource kubernetesinterfaces.Resource, err error) {
	if commandExectutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExectutor")
	}

	if namespace == nil {
		return nil, tracederrors.TracedErrorNil("namespace")
	}

	if resourceName == "" {
		return nil, tracederrors.TracedErrorEmptyString("resourceName")
	}

	if resourceType == "" {
		return nil, tracederrors.TracedErrorEmptyString("resourceType")
	}

	toReturn := NewCommandExecutorResource()

	err = toReturn.SetCommandExecutor(commandExectutor)
	if err != nil {
		return nil, err
	}

	err = toReturn.SetNamespace(namespace)
	if err != nil {
		return nil, err
	}

	err = toReturn.SetName(resourceName)
	if err != nil {
		return nil, err
	}

	err = toReturn.SetTypeName(resourceType)
	if err != nil {
		return nil, err
	}

	return toReturn, nil
}

func NewCommandExecutorResource() (c *CommandExecutorResource) {
	return new(CommandExecutorResource)
}

func (c *CommandExecutorResource) CreateByYamlString(ctx context.Context, options *kubernetesparameteroptions.CreateResourceOptions) (err error) {
	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	yamlString, err := options.GetYamlString()
	if err != nil {
		return err
	}

	resourceName, resourceType, namespaceName, clusterName, err := c.GetResourceAndTypeAndNamespaceAndClusterName()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(
		ctx,
		"Create kubernetes resource by yaml '%s/%s' in namespace '%s' in cluster '%s' started.",
		resourceType,
		resourceName,
		namespaceName,
		clusterName,
	)

	yamlString, err = yamlutils.RunYqQueryAginstYamlStringAsString(yamlString, ".metadata.name=\""+resourceName+"\"")
	if err != nil {
		return err
	}

	yamlString, err = yamlutils.RunYqQueryAginstYamlStringAsString(yamlString, ".metadata.namespace=\""+namespaceName+"\"")
	if err != nil {
		return err
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return err
	}

	if options.SkipNamespaceCreation {
		logging.LogInfoByCtx(ctx, "Skip ensure namespace exists when creating resource by yaml string.")
	} else {
		err = c.EnsureNamespaceExists(ctx)
		if err != nil {
			return err
		}
	}

	cmd := []string{"kubectl"}
	if kubernetesutils.IsInClusterAuthenticationAvailable(ctx) {
		logging.LogInfoByCtxf(ctx, "Kubernetes in cluster authentication is used. Skip validation of kubectlContext.")
	} else {
		kubectlContext, err := c.GetKubectlContext(ctx)
		if err != nil {
			return err
		}

		cmd = append(cmd, "--context", kubectlContext)
	}

	cmd = append(cmd, "apply", "-f", "-")

	_, err = commandExecutor.RunCommand(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command:     cmd,
			StdinString: yamlString,
		},
	)
	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(
		ctx,
		"kubernetes resource '%s/%s' by yam in namespace '%s' in cluster '%s' created and updated.",
		resourceType,
		resourceName,
		namespaceName,
		clusterName,
	)

	logging.LogInfoByCtxf(
		ctx,
		"Create kubernetes resource '%s/%s' by yaml in namespace '%s' in cluster '%s' finished.",
		resourceType,
		resourceName,
		namespaceName,
		clusterName,
	)

	return nil
}

func (c *CommandExecutorResource) Delete(ctx context.Context) (err error) {
	resourceName, resourceTypeName, namespaceName, clusterName, err := c.GetResourceAndTypeAndNamespaceAndClusterName()
	if err != nil {
		return err
	}

	exists, err := c.Exists(ctx)
	if err != nil {
		return err
	}

	if exists {
		logging.LogInfoByCtxf(
			ctx,
			"Going to delete '%s/%s' in namespace '%s' on kubernetes cluster '%s'.",
			resourceName,
			resourceTypeName,
			namespaceName,
			clusterName,
		)

		commandExecutor, err := c.GetCommandExecutor()
		if err != nil {
			return err
		}

		kubectlContext, err := c.GetKubectlContext(ctx)
		if err != nil {
			return err
		}

		cmd := []string{"kubectl"}
		if kubernetesutils.IsInClusterAuthenticationAvailable(ctx) {
			logging.LogInfoByCtxf(ctx, "Kubernetes in cluster authentication is used. Skip validation of kubectlContext '%s'", kubectlContext)
		} else {
			cmd = append(cmd, "--context", kubectlContext)
		}

		cmd = append(cmd, "--namespace", namespaceName, "delete", resourceTypeName, resourceName)

		_, err = commandExecutor.RunCommand(
			ctx,
			&parameteroptions.RunCommandOptions{
				Command: cmd,
			},
		)
		if err != nil {
			return err
		}

		logging.LogChangedByCtxf(
			ctx,
			"Delete '%s/%s' in namespace '%s' on kubernetes cluster '%s'.",
			resourceName,
			resourceTypeName,
			namespaceName,
			clusterName,
		)
	} else {
		logging.LogInfof(
			"Resource '%s/%s' already absent in namespace '%s' on kubernetes cluster '%s'.",
			resourceName,
			resourceTypeName,
			namespaceName,
			clusterName,
		)
	}

	return nil
}

func (c *CommandExecutorResource) EnsureNamespaceExists(ctx context.Context) (err error) {
	namespace, err := c.GetNamespace()
	if err != nil {
		return err
	}

	return namespace.Create(ctx)
}

func (c *CommandExecutorResource) Exists(ctx context.Context) (exists bool, err error) {
	resourceName, resourceType, namespaceName, clusterName, err := c.GetResourceAndTypeAndNamespaceAndClusterName()
	if err != nil {
		return false, err
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return false, err
	}

	kubectlContext, err := c.GetKubectlContext(ctx)
	if err != nil {
		return false, err
	}

	cmd := []string{"kubectl", "get"}
	if kubernetesutils.IsInClusterAuthenticationAvailable(ctx) {
		logging.LogInfoByCtxf(ctx, "Kubernetes in cluster authentication is used. Skip validation of kubectlContext '%s'", kubectlContext)
	} else {
		cmd = append(cmd, "--context", kubectlContext)
	}
	cmd = append(cmd, "--namespace", namespaceName, resourceType, resourceName)

	_, err = commandExecutor.RunCommand(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: cmd,
		},
	)
	if err != nil {
		if strings.Contains(err.Error(), "(NotFound)") {
			exists = false
		} else {
			return false, err
		}
	} else {
		exists = true
	}

	if exists {
		logging.LogInfoByCtxf(
			ctx,
			"Kubernetes resource '%s/%s' in namespace '%s' in cluster '%s' exists.",
			resourceType,
			resourceName,
			namespaceName,
			clusterName,
		)
	} else {
		logging.LogInfoByCtxf(
			ctx,
			"Kubernetes resource '%s/%s' in namespace '%s' in cluster '%s' does not exist.",
			resourceType,
			resourceName,
			namespaceName,
			clusterName,
		)
	}

	return exists, nil
}

func (c *CommandExecutorResource) GetAsYamlString() (yamlString string, err error) {
	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return "", err
	}

	contextName, err := c.GetKubectlContext(contextutils.ContextSilent())
	if err != nil {
		return "", err
	}

	resourceName, resourceType, namspaceName, clusterName, err := c.GetResourceAndTypeAndNamespaceAndClusterName()
	if err != nil {
		return "", err
	}

	yamlString, err = commandExecutor.RunCommandAndGetStdoutAsString(
		contextutils.ContextSilent(),
		&parameteroptions.RunCommandOptions{
			Command: []string{
				"kubectl",
				"get",
				"--context",
				contextName,
				"--namespace",
				namspaceName,
				resourceType,
				resourceName,
				"-o",
				"yaml",
			},
		},
	)
	if err != nil {
		return "", err
	}

	if yamlString == "" {
		return "", tracederrors.TracedErrorf(
			"yamlString is empty string after evaluation. Tried to get resource type '%s' named '%s' in namespace '%s' in cluster '%s'.",
			resourceType,
			resourceName,
			namspaceName,
			clusterName,
		)
	}

	return yamlString, nil
}

func (c *CommandExecutorResource) GetClusterName() (clusterName string, err error) {
	namespace, err := c.GetNamespace()
	if err != nil {
		return "", err
	}

	return namespace.GetClusterName()
}

func (c *CommandExecutorResource) GetCommandExecutor() (commandExecutor commandexecutor.CommandExecutor, err error) {

	return c.commandExecutor, nil
}

func (c *CommandExecutorResource) GetKubectlContext(ctx context.Context) (contextName string, err error) {
	namespace, err := c.GetNamespace()
	if err != nil {
		return "", err
	}

	return namespace.GetKubectlContext(ctx)
}

func (c *CommandExecutorResource) GetName() (name string, err error) {
	if c.name == "" {
		return "", tracederrors.TracedErrorf("name not set")
	}

	return c.name, nil
}

func (c *CommandExecutorResource) GetNamespace() (namespace kubernetesinterfaces.Namespace, err error) {
	if c.namespace == nil {
		return nil, err
	}

	return c.namespace, nil
}

func (c *CommandExecutorResource) GetNamespaceName() (namsepaceName string, err error) {
	namespace, err := c.GetNamespace()
	if err != nil {
		return "", err
	}

	return namespace.GetName()
}

func (c *CommandExecutorResource) GetResourceAndTypeAndNamespaceAndClusterName() (resourceName string, resourceTypeName string, namespaceName string, clusterName string, err error) {
	resourceName, err = c.GetName()
	if err != nil {
		return "", "", "", "", err
	}

	resourceTypeName, err = c.GetTypeName()
	if err != nil {
		return "", "", "", "", err
	}

	namespaceName, err = c.GetNamespaceName()
	if err != nil {
		return "", "", "", "", err
	}

	clusterName, err = c.GetClusterName()
	if err != nil {
		return "", "", "", "", err
	}

	return resourceName, resourceTypeName, namespaceName, clusterName, err
}

func (c *CommandExecutorResource) GetTypeName() (typeName string, err error) {
	if c.typeName == "" {
		return "", tracederrors.TracedErrorf("typeName not set")
	}

	return c.typeName, nil
}

func (c *CommandExecutorResource) SetCommandExecutor(commandExecutor commandexecutor.CommandExecutor) (err error) {
	c.commandExecutor = commandExecutor

	return nil
}

func (c *CommandExecutorResource) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorf("name is empty string")
	}

	c.name = name

	return nil
}

func (c *CommandExecutorResource) SetNamespace(namespace kubernetesinterfaces.Namespace) (err error) {
	c.namespace = namespace

	return nil
}

func (c *CommandExecutorResource) SetTypeName(typeName string) (err error) {
	if typeName == "" {
		return tracederrors.TracedErrorf("typeName is empty string")
	}

	c.typeName = typeName

	return nil
}

func (c *CommandExecutorResource) SetApiVersion(apiVersion string) (error) {
	return tracederrors.TracedErrorNotImplemented()
}