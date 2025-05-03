package kubernetesutils

import (
	"context"
	"slices"
	"sort"
	"strings"

	"github.com/asciich/asciichgolangpublic/commandexecutor"
	"github.com/asciich/asciichgolangpublic/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type CommandExecutorKubernetes struct {
	name              string
	commandExecutor   commandexecutor.CommandExecutor
	cachedContextName string
}

func GetCommandExecutorKubernetsByName(commandExecutor commandexecutor.CommandExecutor, clusterName string) (kubernetes KubernetesCluster, err error) {
	if commandExecutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExecutor")
	}

	if clusterName == "" {
		return nil, tracederrors.TracedErrorEmptyString("clusterName")
	}

	toReturn := NewCommandExecutorKubernetes()

	err = toReturn.SetCommandExecutor(commandExecutor)
	if err != nil {
		return nil, err
	}

	err = toReturn.SetName(clusterName)
	if err != nil {
		return nil, err
	}

	return toReturn, nil
}

func GetLocalCommandExecutorKubernetesByName(clusterName string) (kubernetes KubernetesCluster, err error) {
	if clusterName == "" {
		return nil, tracederrors.TracedErrorEmptyString("clusterName")
	}

	return GetCommandExecutorKubernetsByName(commandexecutor.Bash(), clusterName)
}

func MustGetCommandExecutorKubernetsByName(commandExecutor commandexecutor.CommandExecutor, clusterName string) (kubernetes KubernetesCluster) {
	kubernetes, err := GetCommandExecutorKubernetsByName(commandExecutor, clusterName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return kubernetes
}

func MustGetLocalCommandExecutorKubernetesByName(clusterName string) (kubernetes KubernetesCluster) {
	kubernetes, err := GetLocalCommandExecutorKubernetesByName(clusterName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return kubernetes
}

func NewCommandExecutorKubernetes() (c *CommandExecutorKubernetes) {
	return new(CommandExecutorKubernetes)
}

// Returns the kubernetes cluster name
func (c *CommandExecutorKubernetes) GetName() (name string, err error) {
	if c.name == "" {
		return "", tracederrors.TracedErrorf("name not set")
	}

	return c.name, nil
}

func (c *CommandExecutorKubernetes) CreateNamespaceByName(ctx context.Context, name string) (createdNamespace Namespace, err error) {
	if name == "" {
		return nil, tracederrors.TracedErrorEmptyString("name")
	}

	exists, err := c.NamespaceByNameExists(ctx, name)
	if err != nil {
		return nil, err
	}

	clusterName, err := c.GetName()
	if err != nil {
		return nil, err
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Namespace '%s' already exists in cluster '%s'.", name, clusterName)
	} else {
		cmd := []string{"kubectl"}

		if IsInClusterAuthenticationAvailable(ctx) {
			logging.LogInfoByCtxf(ctx, "Kubernetes in cluster authentication is used. cluster context is not used.")
		} else {
			kubectlContext, err := c.GetCachedKubectlContext(ctx)
			if err != nil {
				return nil, err
			}

			cmd = append(cmd, "--context", kubectlContext)
		}

		cmd = append(cmd, "create", "namespace", name)

		_, err = c.RunCommand(
			ctx,
			&parameteroptions.RunCommandOptions{
				Command: cmd,
			},
		)
		if err != nil {
			return nil, err
		}

		logging.LogChangedByCtxf(ctx, "Namespace '%s' in cluster '%s' created.", name, clusterName)
	}

	return c.GetNamespaceByName(name)
}

func (c *CommandExecutorKubernetes) DeleteNamespaceByName(ctx context.Context, name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	exists, err := c.NamespaceByNameExists(ctx, name)
	if err != nil {
		return err
	}

	clusterName, err := c.GetName()
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
					"delete",
					"namespace",
					name,
				},
			},
		)
		if err != nil {
			return err
		}

		logging.LogChangedByCtxf(ctx, "Namespace '%s' in cluster '%s' deleted.", name, clusterName)
	} else {
		logging.LogInfoByCtxf(ctx, "Namespace '%s' already absent in cluster '%s'.", name, clusterName)
	}

	return nil
}

func (c *CommandExecutorKubernetes) GetCachedContextName() (cachedContextName string, err error) {
	if c.cachedContextName == "" {
		return "", tracederrors.TracedErrorf("cachedContextName not set")
	}

	return c.cachedContextName, nil
}

func (c *CommandExecutorKubernetes) GetCachedKubectlContext(ctx context.Context) (context string, err error) {
	if c.cachedContextName == "" {
		return c.GetKubectlContext(ctx)
	}

	context = c.cachedContextName

	clusterName, err := c.GetName()
	if err != nil {
		return "", err
	}

	logging.LogInfof(
		"Kubectl context for cluster '%s' is '%s'.",
		clusterName,
		context,
	)

	return
}

func (c *CommandExecutorKubernetes) GetCommandExecutor() (commandExecutor commandexecutor.CommandExecutor, err error) {
	if c.commandExecutor == nil {
		return nil, tracederrors.TracedError("CommandExecutor not set")
	}

	return c.commandExecutor, nil
}

func (c *CommandExecutorKubernetes) GetKubectlContext(ctx context.Context) (context string, err error) {
	contexts, err := c.GetKubectlContexts()
	if err != nil {
		return "", err
	}

	clusterName, err := c.GetName()
	if err != nil {
		return "", err
	}

	for _, con := range contexts {
		clusterNameToCeck, err := con.GetCluster()
		if err != nil {
			return "", err
		}

		if clusterNameToCeck == clusterName {
			context, err = con.GetName()
			if err != nil {
				return "", err
			}

			logging.LogInfoByCtxf(ctx, "Kubectl context for cluster '%s' is '%s'.", clusterName, context)

			return context, nil
		}
	}

	return "", tracederrors.TracedErrorf(
		"No kubectl context for cluster '%s' found.",
		clusterName,
	)
}

func (c *CommandExecutorKubernetes) GetKubectlContexts() (contexts []KubectlContext, err error) {
	lines, err := c.RunCommandAndGetStdoutAsLines(
		contextutils.ContextSilent(),
		&parameteroptions.RunCommandOptions{
			Command: []string{"kubectl", "config", "get-contexts", "--no-headers"},
		},
	)
	if err != nil {
		return nil, err
	}

	contexts = []KubectlContext{}
	for _, line := range lines {
		line = strings.ReplaceAll(line, "\t", " ")
		line = stringsutils.RepeatReplaceAll(line, "  ", " ")
		line = strings.TrimSpace(line)
		line = strings.TrimPrefix(line, "*")
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		splitted := strings.Split(line, " ")
		if len(splitted) <= 2 {
			return nil, tracederrors.TracedErrorf(
				"Unable to get context from line: '%s'",
				line,
			)
		}

		toAdd := NewKubectlContext()
		err = toAdd.SetName(splitted[0])
		if err != nil {
			return nil, err
		}

		err = toAdd.SetCluster(splitted[1])
		if err != nil {
			return nil, err
		}

		contexts = append(contexts, *toAdd)
	}

	return contexts, nil
}

func (c *CommandExecutorKubernetes) GetNamespaceByName(name string) (namespace Namespace, err error) {
	if name == "" {
		return nil, tracederrors.TracedErrorEmptyString("name")
	}

	toReturn := NewCommandExecutorNamespace()

	err = toReturn.SetName(name)
	if err != nil {
		return nil, err
	}

	err = toReturn.SetKubernetesCluster(c)
	if err != nil {
		return nil, err
	}

	return toReturn, nil
}

func (c *CommandExecutorKubernetes) GetResourceByNames(resourceName string, resourceType string, namespaceName string) (resource Resource, err error) {
	if resourceName == "" {
		return nil, tracederrors.TracedErrorEmptyString("resourceName")
	}

	if resourceType == "" {
		return nil, tracederrors.TracedErrorEmptyString("resourceType")
	}

	if namespaceName == "" {
		return nil, tracederrors.TracedErrorEmptyString("namespaceName")
	}

	namespace, err := c.GetNamespaceByName(namespaceName)
	if err != nil {
		return nil, err
	}

	return namespace.GetResourceByNames(resourceName, resourceType)
}

func (c *CommandExecutorKubernetes) ListNamespaceNames(ctx context.Context) (namespaceNames []string, err error) {
	namespaces, err := c.ListNamespaces(ctx)
	if err != nil {
		return nil, err
	}

	namespaceNames = []string{}
	for _, namespace := range namespaces {
		toAdd, err := namespace.GetName()
		if err != nil {
			return nil, err
		}

		namespaceNames = append(namespaceNames, toAdd)
	}

	return namespaceNames, nil
}

func (c *CommandExecutorKubernetes) ListNamespaces(ctx context.Context) (namespaces []Namespace, err error) {

	cmd := []string{"kubectl"}

	if IsInClusterAuthenticationAvailable(ctx) {
		logging.LogInfoByCtxf(ctx, "Kubernetes in cluster authentication is used. Skip validation of kubectlContext for ListNamespaces.")
	} else {
		context, err := c.GetCachedKubectlContext(ctx)
		if err != nil {
			return nil, err
		}

		cmd = append(cmd, "--context", context)
	}

	cmd = append(cmd, "get", "namespaces", "-o", "name")

	lines, err := c.RunCommandAndGetStdoutAsLines(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: cmd,
		},
	)
	if err != nil {
		return nil, err
	}

	namespaces = []Namespace{}
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		line = strings.TrimPrefix(line, "namespace/")

		toAdd, err := c.GetNamespaceByName(line)
		if err != nil {
			return nil, err
		}

		namespaces = append(namespaces, toAdd)
	}

	return namespaces, nil
}

func (c *CommandExecutorKubernetes) ListResourceNames(options *parameteroptions.ListKubernetesResourcesOptions) (resourceNames []string, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	namespaceName, err := options.GetNamespace()
	if err != nil {
		return nil, err
	}

	context, err := c.GetKubectlContext(contextutils.GetVerbosityContextByBool(options.Verbose))
	if err != nil {
		return nil, err
	}

	resourceType, err := options.GetResourceType()
	if err != nil {
		return nil, err
	}

	output, err := commandExecutor.RunCommandAndGetStdoutAsLines(
		contextutils.GetVerbosityContextByBool(options.Verbose),
		&parameteroptions.RunCommandOptions{
			Command: []string{
				"kubectl",
				"get",
				"--context",
				context,
				"--namespace",
				namespaceName,
				"-o",
				"name",
				resourceType,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	resourceNames = []string{}
	for _, name := range output {
		resourceNames = append(resourceNames, strings.TrimPrefix(name, resourceType+"/"))
	}

	sort.Strings(resourceNames)

	return resourceNames, nil
}

func (c *CommandExecutorKubernetes) ListResources(options *parameteroptions.ListKubernetesResourcesOptions) (resources []Resource, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	resourceNames, err := c.ListResourceNames(options)
	if err != nil {
		return nil, err
	}

	namespaceName, err := options.GetNamespace()
	if err != nil {
		return nil, err
	}

	resourceType, err := options.GetResourceType()
	if err != nil {
		return nil, err
	}

	resources = []Resource{}
	for _, name := range resourceNames {
		toAdd, err := c.GetResourceByNames(name, resourceType, namespaceName)
		if err != nil {
			return nil, err
		}

		resources = append(resources, toAdd)
	}

	return resources, nil
}

func (c *CommandExecutorKubernetes) NamespaceByNameExists(ctx context.Context, name string) (exists bool, err error) {
	if name == "" {
		return false, tracederrors.TracedErrorEmptyString("name")
	}

	namespaces, err := c.ListNamespaceNames(ctx)
	if err != nil {
		return false, err
	}

	exists = slices.Contains(namespaces, name)

	clusterName, err := c.GetName()
	if err != nil {
		return false, err
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Namespace '%s' exists in kubernetes cluster '%s'.", name, clusterName)
	} else {
		logging.LogInfoByCtxf(ctx, "Namespace '%s' does not exist in kubernetes cluster '%s'.", name, clusterName)
	}

	return exists, nil
}

func (c *CommandExecutorKubernetes) RunCommand(ctx context.Context, runCommandOptions *parameteroptions.RunCommandOptions) (commandOutput *commandexecutor.CommandOutput, err error) {
	if runCommandOptions == nil {
		return nil, tracederrors.TracedErrorNil("runCommandOptions")
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	return commandExecutor.RunCommand(ctx, runCommandOptions)
}

func (c *CommandExecutorKubernetes) RunCommandAndGetStdoutAsLines(ctx context.Context, runCommandOptions *parameteroptions.RunCommandOptions) (lines []string, err error) {
	if runCommandOptions == nil {
		return nil, tracederrors.TracedErrorNil("runCommandOptions")
	}

	output, err := c.RunCommand(ctx, runCommandOptions)
	if err != nil {
		return nil, err
	}

	return output.GetStdoutAsLines(false)
}

func (c *CommandExecutorKubernetes) SetCachedContextName(cachedContextName string) (err error) {
	if cachedContextName == "" {
		return tracederrors.TracedErrorf("cachedContextName is empty string")
	}

	c.cachedContextName = cachedContextName

	return nil
}

func (c *CommandExecutorKubernetes) SetCommandExecutor(commandExecutor commandexecutor.CommandExecutor) (err error) {
	c.commandExecutor = commandExecutor

	return nil
}

func (c *CommandExecutorKubernetes) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorf("name is empty string")
	}

	c.name = name

	return nil
}
