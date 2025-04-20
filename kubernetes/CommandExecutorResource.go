package kubernetes

import (
	"strings"

	"github.com/asciich/asciichgolangpublic/commandexecutor"
	"github.com/asciich/asciichgolangpublic/fileformats/yamlutils"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type CommandExecutorResource struct {
	commandExecutor commandexecutor.CommandExecutor
	name            string
	typeName        string
	namespace       Namespace
}

func GetCommandExecutorResource(commandExectutor commandexecutor.CommandExecutor, namespace Namespace, resourceName string, resourceType string) (resource Resource, err error) {
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

func MustGetCommandExecutorResource(commandExectutor commandexecutor.CommandExecutor, namespace Namespace, resourceName string, resourceType string) (resource Resource) {
	resource, err := GetCommandExecutorResource(commandExectutor, namespace, resourceName, resourceType)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return resource
}

func NewCommandExecutorResource() (c *CommandExecutorResource) {
	return new(CommandExecutorResource)
}

func (c *CommandExecutorResource) CreateByYamlString(yamlString string, verbose bool) (err error) {
	if yamlString == "" {
		return tracederrors.TracedErrorEmptyString("yamlString")
	}

	resourceName, resourceType, namespaceName, clusterName, err := c.GetResourceAndTypeAndNamespaceAndClusterName()
	if err != nil {
		return err
	}

	if verbose {
		logging.LogInfof(
			"Create kubernetes resource by yaml '%s/%s' in namespace '%s' in cluster '%s' started.",
			resourceType,
			resourceName,
			namespaceName,
			clusterName,
		)
	}

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

	context, err := c.GetKubectlContext(verbose)
	if err != nil {
		return err
	}

	err = c.EnsureNamespaceExists(verbose)
	if err != nil {
		return err
	}

	_, err = commandExecutor.RunCommand(
		contextutils.GetVerbosityContextByBool(verbose),
		&parameteroptions.RunCommandOptions{
			Command:     []string{"kubectl", "--context", context, "apply", "-f", "-"},
			StdinString: yamlString,
		},
	)
	if err != nil {
		return err
	}

	if verbose {
		logging.LogChangedf(
			"kubernetes resource '%s/%s' by yam in namespace '%s' in cluster '%s' created and updated.",
			resourceType,
			resourceName,
			namespaceName,
			clusterName,
		)
	}

	if verbose {
		logging.LogInfof(
			"Create kubernetes resource '%s/%s' by yaml in namespace '%s' in cluster '%s' finished.",
			resourceType,
			resourceName,
			namespaceName,
			clusterName,
		)
	}

	return nil
}

func (c *CommandExecutorResource) Delete(verbose bool) (err error) {
	resourceName, resourceTypeName, namespaceName, clusterName, err := c.GetResourceAndTypeAndNamespaceAndClusterName()
	if err != nil {
		return err
	}

	exists, err := c.Exists(verbose)
	if err != nil {
		return err
	}

	if exists {
		if verbose {
			logging.LogInfof(
				"Going to delete '%s/%s' in namespace '%s' on kubernetes cluster '%s'.",
				resourceName,
				resourceTypeName,
				namespaceName,
				clusterName,
			)
		}

		commandExecutor, err := c.GetCommandExecutor()
		if err != nil {
			return err
		}

		contextName, err := c.GetKubectlContext(verbose)
		if err != nil {
			return err
		}

		_, err = commandExecutor.RunCommand(
			contextutils.GetVerbosityContextByBool(verbose),
			&parameteroptions.RunCommandOptions{
				Command: []string{"kubectl", "--context", contextName, "--namespace", namespaceName, "delete", resourceTypeName, resourceName},
			},
		)
		if err != nil {
			return err
		}

		if verbose {
			logging.LogChangedf(
				"Delete '%s/%s' in namespace '%s' on kubernetes cluster '%s'.",
				resourceName,
				resourceTypeName,
				namespaceName,
				clusterName,
			)
		}
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

func (c *CommandExecutorResource) EnsureNamespaceExists(verbose bool) (err error) {
	namespace, err := c.GetNamespace()
	if err != nil {
		return err
	}

	return namespace.Create(verbose)
}

func (c *CommandExecutorResource) Exists(verbose bool) (exists bool, err error) {
	resourceName, resourceType, namespaceName, clusterName, err := c.GetResourceAndTypeAndNamespaceAndClusterName()
	if err != nil {
		return false, err
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return false, err
	}

	context, err := c.GetKubectlContext(verbose)
	if err != nil {
		return false, err
	}

	_, err = commandExecutor.RunCommand(
		contextutils.GetVerbosityContextByBool(verbose),
		&parameteroptions.RunCommandOptions{
			Command: []string{"kubectl", "get", "--context", context, "--namespace", namespaceName, resourceType, resourceName},
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

	if verbose {
		if exists {
			logging.LogInfof(
				"Kubernetes resource '%s/%s' in namespace '%s' in cluster '%s' exists.",
				resourceType,
				resourceName,
				namespaceName,
				clusterName,
			)
		} else {
			logging.LogInfof(
				"Kubernetes resource '%s/%s' in namespace '%s' in cluster '%s' does not exist.",
				resourceType,
				resourceName,
				namespaceName,
				clusterName,
			)
		}
	}

	return exists, nil
}

func (c *CommandExecutorResource) GetAsYamlString() (yamlString string, err error) {
	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return "", err
	}

	contextName, err := c.GetKubectlContext(false)
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

func (c *CommandExecutorResource) GetKubectlContext(verbose bool) (contextName string, err error) {
	namespace, err := c.GetNamespace()
	if err != nil {
		return "", err
	}

	return namespace.GetKubectlContext(verbose)
}

func (c *CommandExecutorResource) GetName() (name string, err error) {
	if c.name == "" {
		return "", tracederrors.TracedErrorf("name not set")
	}

	return c.name, nil
}

func (c *CommandExecutorResource) GetNamespace() (namespace Namespace, err error) {
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

func (c *CommandExecutorResource) MustCreateByYamlString(yamlString string, verbose bool) {
	err := c.CreateByYamlString(yamlString, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorResource) MustDelete(verbose bool) {
	err := c.Delete(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorResource) MustEnsureNamespaceExists(verbose bool) {
	err := c.EnsureNamespaceExists(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorResource) MustExists(verbose bool) (exists bool) {
	exists, err := c.Exists(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return exists
}

func (c *CommandExecutorResource) MustGetAsYamlString() (yamlString string) {
	yamlString, err := c.GetAsYamlString()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return yamlString
}

func (c *CommandExecutorResource) MustGetClusterName() (clusterName string) {
	clusterName, err := c.GetClusterName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return clusterName
}

func (c *CommandExecutorResource) MustGetCommandExecutor() (commandExecutor commandexecutor.CommandExecutor) {
	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return commandExecutor
}

func (c *CommandExecutorResource) MustGetKubectlContext(verbose bool) (contextName string) {
	contextName, err := c.GetKubectlContext(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return contextName
}

func (c *CommandExecutorResource) MustGetName() (name string) {
	name, err := c.GetName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return name
}

func (c *CommandExecutorResource) MustGetNamespace() (namespace Namespace) {
	namespace, err := c.GetNamespace()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return namespace
}

func (c *CommandExecutorResource) MustGetNamespaceName() (namsepaceName string) {
	namsepaceName, err := c.GetNamespaceName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return namsepaceName
}

func (c *CommandExecutorResource) MustGetResourceAndTypeAndNamespaceAndClusterName() (resourceName string, resourceTypeName string, namespaceName string, clusterName string) {
	resourceName, resourceTypeName, namespaceName, clusterName, err := c.GetResourceAndTypeAndNamespaceAndClusterName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return resourceName, resourceTypeName, namespaceName, clusterName
}

func (c *CommandExecutorResource) MustGetTypeName() (typeName string) {
	typeName, err := c.GetTypeName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return typeName
}

func (c *CommandExecutorResource) MustSetCommandExecutor(commandExecutor commandexecutor.CommandExecutor) {
	err := c.SetCommandExecutor(commandExecutor)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorResource) MustSetName(name string) {
	err := c.SetName(name)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorResource) MustSetNamespace(namespace Namespace) {
	err := c.SetNamespace(namespace)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorResource) MustSetTypeName(typeName string) {
	err := c.SetTypeName(typeName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
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

func (c *CommandExecutorResource) SetNamespace(namespace Namespace) (err error) {
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
