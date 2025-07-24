package commandexecutorkubernetes

import (
	"context"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/fileformats/yamlutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type CommandExecutorObject struct {
	commandExecutor commandexecutorinterfaces.CommandExecutor
	name            string
	typeName        string
	namespace       kubernetesinterfaces.Namespace
}

func GetCommandExecutorObject(commandExectutor commandexecutorinterfaces.CommandExecutor, namespace kubernetesinterfaces.Namespace, objectName string, objectType string) (object kubernetesinterfaces.Object, err error) {
	if commandExectutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExectutor")
	}

	if namespace == nil {
		return nil, tracederrors.TracedErrorNil("namespace")
	}

	if objectName == "" {
		return nil, tracederrors.TracedErrorEmptyString("objectName")
	}

	if objectType == "" {
		return nil, tracederrors.TracedErrorEmptyString("objectType")
	}

	toReturn := NewCommandExecutorObject()

	err = toReturn.SetCommandExecutor(commandExectutor)
	if err != nil {
		return nil, err
	}

	err = toReturn.SetNamespace(namespace)
	if err != nil {
		return nil, err
	}

	err = toReturn.SetName(objectName)
	if err != nil {
		return nil, err
	}

	err = toReturn.SetTypeName(objectType)
	if err != nil {
		return nil, err
	}

	return toReturn, nil
}

func NewCommandExecutorObject() (c *CommandExecutorObject) {
	return new(CommandExecutorObject)
}

func (c *CommandExecutorObject) CreateByYamlString(ctx context.Context, options *kubernetesparameteroptions.CreateObjectOptions) (err error) {
	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	yamlString, err := options.GetYamlString()
	if err != nil {
		return err
	}

	objectName, objectType, namespaceName, clusterName, err := c.GetObjectAndTypeAndNamespaceAndClusterName()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(
		ctx,
		"Create kubernetes object by yaml '%s/%s' in namespace '%s' in cluster '%s' started.",
		objectType,
		objectName,
		namespaceName,
		clusterName,
	)

	yamlString, err = yamlutils.RunYqQueryAginstYamlStringAsString(yamlString, ".metadata.name=\""+objectName+"\"")
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
		logging.LogInfoByCtx(ctx, "Skip ensure namespace exists when creating object by yaml string.")
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
		"kubernetes object '%s/%s' by yam in namespace '%s' in cluster '%s' created and updated.",
		objectType,
		objectName,
		namespaceName,
		clusterName,
	)

	logging.LogInfoByCtxf(
		ctx,
		"Create kubernetes object '%s/%s' by yaml in namespace '%s' in cluster '%s' finished.",
		objectType,
		objectName,
		namespaceName,
		clusterName,
	)

	return nil
}

func (c *CommandExecutorObject) Delete(ctx context.Context) (err error) {
	objectName, objectTypeName, namespaceName, clusterName, err := c.GetObjectAndTypeAndNamespaceAndClusterName()
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
			objectName,
			objectTypeName,
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

		cmd = append(cmd, "--namespace", namespaceName, "delete", objectTypeName, objectName)

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
			objectName,
			objectTypeName,
			namespaceName,
			clusterName,
		)
	} else {
		logging.LogInfof(
			"Object '%s/%s' already absent in namespace '%s' on kubernetes cluster '%s'.",
			objectName,
			objectTypeName,
			namespaceName,
			clusterName,
		)
	}

	return nil
}

func (c *CommandExecutorObject) EnsureNamespaceExists(ctx context.Context) (err error) {
	namespace, err := c.GetNamespace()
	if err != nil {
		return err
	}

	return namespace.Create(ctx)
}

func (c *CommandExecutorObject) Exists(ctx context.Context) (exists bool, err error) {
	objectName, objectType, namespaceName, clusterName, err := c.GetObjectAndTypeAndNamespaceAndClusterName()
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
	cmd = append(cmd, "--namespace", namespaceName, objectType, objectName)

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
			"Kubernetes object '%s/%s' in namespace '%s' in cluster '%s' exists.",
			objectType,
			objectName,
			namespaceName,
			clusterName,
		)
	} else {
		logging.LogInfoByCtxf(
			ctx,
			"Kubernetes object '%s/%s' in namespace '%s' in cluster '%s' does not exist.",
			objectType,
			objectName,
			namespaceName,
			clusterName,
		)
	}

	return exists, nil
}

func (c *CommandExecutorObject) GetAsYamlString() (yamlString string, err error) {
	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return "", err
	}

	contextName, err := c.GetKubectlContext(contextutils.ContextSilent())
	if err != nil {
		return "", err
	}

	objectName, objectType, namspaceName, clusterName, err := c.GetObjectAndTypeAndNamespaceAndClusterName()
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
				objectType,
				objectName,
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
			"yamlString is empty string after evaluation. Tried to get object type '%s' named '%s' in namespace '%s' in cluster '%s'.",
			objectType,
			objectName,
			namspaceName,
			clusterName,
		)
	}

	return yamlString, nil
}

func (c *CommandExecutorObject) GetClusterName() (clusterName string, err error) {
	namespace, err := c.GetNamespace()
	if err != nil {
		return "", err
	}

	return namespace.GetClusterName()
}

func (c *CommandExecutorObject) GetCommandExecutor() (commandExecutor commandexecutorinterfaces.CommandExecutor, err error) {

	return c.commandExecutor, nil
}

func (c *CommandExecutorObject) GetKubectlContext(ctx context.Context) (contextName string, err error) {
	namespace, err := c.GetNamespace()
	if err != nil {
		return "", err
	}

	return namespace.GetKubectlContext(ctx)
}

func (c *CommandExecutorObject) GetName() (name string, err error) {
	if c.name == "" {
		return "", tracederrors.TracedErrorf("name not set")
	}

	return c.name, nil
}

func (c *CommandExecutorObject) GetNamespace() (namespace kubernetesinterfaces.Namespace, err error) {
	if c.namespace == nil {
		return nil, err
	}

	return c.namespace, nil
}

func (c *CommandExecutorObject) GetNamespaceName() (namsepaceName string, err error) {
	namespace, err := c.GetNamespace()
	if err != nil {
		return "", err
	}

	return namespace.GetName()
}

func (c *CommandExecutorObject) GetObjectAndTypeAndNamespaceAndClusterName() (objectName string, objectTypeName string, namespaceName string, clusterName string, err error) {
	objectName, err = c.GetName()
	if err != nil {
		return "", "", "", "", err
	}

	objectTypeName, err = c.GetTypeName()
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

	return objectName, objectTypeName, namespaceName, clusterName, err
}

func (c *CommandExecutorObject) GetTypeName() (typeName string, err error) {
	if c.typeName == "" {
		return "", tracederrors.TracedErrorf("typeName not set")
	}

	return c.typeName, nil
}

func (c *CommandExecutorObject) SetCommandExecutor(commandExecutor commandexecutorinterfaces.CommandExecutor) (err error) {
	c.commandExecutor = commandExecutor

	return nil
}

func (c *CommandExecutorObject) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorf("name is empty string")
	}

	c.name = name

	return nil
}

func (c *CommandExecutorObject) SetNamespace(namespace kubernetesinterfaces.Namespace) (err error) {
	c.namespace = namespace

	return nil
}

func (c *CommandExecutorObject) SetTypeName(typeName string) (err error) {
	if typeName == "" {
		return tracederrors.TracedErrorf("typeName is empty string")
	}

	c.typeName = typeName

	return nil
}

func (c *CommandExecutorObject) SetApiVersion(apiVersion string) error {
	return tracederrors.TracedErrorNotImplemented()
}
