package kubernetes

import (
	"strings"

	"github.com/asciich/asciichgolangpublic"
)

type CommandExecutorKubernetes struct {
	name            string
	commandExecutor asciichgolangpublic.CommandExecutor
}

func GetCommandExecutorKubernetsByName(commandExecutor asciichgolangpublic.CommandExecutor, clusterName string) (kubernetes KubernetesCluster, err error) {
	if commandExecutor == nil {
		return nil, asciichgolangpublic.TracedErrorNil("commandExecutor")
	}

	if clusterName == "" {
		return nil, asciichgolangpublic.TracedErrorEmptyString("clusterName")
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
		return nil, asciichgolangpublic.TracedErrorEmptyString("clusterName")
	}

	return GetCommandExecutorKubernetsByName(asciichgolangpublic.Bash(), clusterName)
}

func MustGetCommandExecutorKubernetsByName(commandExecutor asciichgolangpublic.CommandExecutor, clusterName string) (kubernetes KubernetesCluster) {
	kubernetes, err := GetCommandExecutorKubernetsByName(commandExecutor, clusterName)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return kubernetes
}

func MustGetLocalCommandExecutorKubernetesByName(clusterName string) (kubernetes KubernetesCluster) {
	kubernetes, err := GetLocalCommandExecutorKubernetesByName(clusterName)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return kubernetes
}

func NewCommandExecutorKubernetes() (c *CommandExecutorKubernetes) {
	return new(CommandExecutorKubernetes)
}

// Returns the kubernetes cluster name
func (c *CommandExecutorKubernetes) GetName() (name string, err error) {
	if c.name == "" {
		return "", asciichgolangpublic.TracedErrorf("name not set")
	}

	return c.name, nil
}

func (c *CommandExecutorKubernetes) CreateNamespaceByName(name string, verbose bool) (createdNamespace Namespace, err error) {
	if name == "" {
		return nil, asciichgolangpublic.TracedErrorEmptyString("name")
	}

	exists, err := c.NamespaceByNameExists(name, verbose)
	if err != nil {
		return nil, err
	}

	clusterName, err := c.GetName()
	if err != nil {
		return nil, err
	}

	if exists {
		if verbose {
			asciichgolangpublic.LogInfof(
				"Namespace '%s' already exists in cluster '%s'.",
				name,
				clusterName,
			)
		}
	} else {
		context, err := c.GetKubectlContext(verbose)
		if err != nil {
			return nil, err
		}

		_, err = c.RunCommand(
			&asciichgolangpublic.RunCommandOptions{
				Command: []string{
					"kubectl",
					"--context",
					context,
					"create",
					"namespace",
					name,
				},
			},
		)
		if err != nil {
			return nil, err
		}

		if verbose {
			asciichgolangpublic.LogChangedf(
				"Namespace '%s' in cluster '%s' created.",
				name,
				clusterName,
			)
		}
	}

	return c.GetNamespaceByName(name)
}

func (c *CommandExecutorKubernetes) DeleteNamespaceByName(name string, verbose bool) (err error) {
	if name == "" {
		return asciichgolangpublic.TracedErrorEmptyString("name")
	}

	exists, err := c.NamespaceByNameExists(name, verbose)
	if err != nil {
		return err
	}

	clusterName, err := c.GetName()
	if err != nil {
		return err
	}

	if exists {

		context, err := c.GetKubectlContext(verbose)
		if err != nil {
			return err
		}

		_, err = c.RunCommand(
			&asciichgolangpublic.RunCommandOptions{
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

		if verbose {
			asciichgolangpublic.LogChangedf(
				"Namespace '%s' in cluster '%s' deleted.",
				name,
				clusterName,
			)
		}
	} else {
		if verbose {
			asciichgolangpublic.LogInfof(
				"Namespace '%s' already absent in cluster '%s'.",
				name,
				clusterName,
			)
		}
	}

	return nil
}

func (c *CommandExecutorKubernetes) GetCommandExecutor() (commandExecutor asciichgolangpublic.CommandExecutor, err error) {

	return c.commandExecutor, nil
}

func (c *CommandExecutorKubernetes) GetKubectlContext(verbose bool) (context string, err error) {
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

			if verbose {
				asciichgolangpublic.LogInfof(
					"Kubectl context for cluster '%s' is '%s'.",
					clusterName,
					context,
				)
			}

			return context, nil
		}
	}

	return "", asciichgolangpublic.TracedErrorf(
		"No kubectl context for cluster '%s' found.",
		clusterName,
	)
}

func (c *CommandExecutorKubernetes) GetKubectlContexts() (contexts []KubectlContext, err error) {
	lines, err := c.RunCommandAndGetStdoutAsLines(
		&asciichgolangpublic.RunCommandOptions{
			Command: []string{"kubectl", "config", "get-contexts", "--no-headers"},
		},
	)
	if err != nil {
		return nil, err
	}

	contexts = []KubectlContext{}
	for _, line := range lines {
		line = strings.ReplaceAll(line, "\t", " ")
		line = asciichgolangpublic.Strings().RepeatReplaceAll(line, "  ", " ")
		line = strings.TrimSpace(line)
		line = strings.TrimPrefix(line, "*")
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		splitted := strings.Split(line, " ")
		if len(splitted) <= 2 {
			return nil, asciichgolangpublic.TracedErrorf(
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
		return nil, asciichgolangpublic.TracedErrorEmptyString("name")
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

func (c *CommandExecutorKubernetes) ListNamespaceNames(verbose bool) (namespaceNames []string, err error) {
	namespaces, err := c.ListNamespaces(verbose)
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

func (c *CommandExecutorKubernetes) ListNamespaces(verbose bool) (namespaces []Namespace, err error) {
	context, err := c.GetKubectlContext(verbose)
	if err != nil {
		return nil, err
	}

	lines, err := c.RunCommandAndGetStdoutAsLines(
		&asciichgolangpublic.RunCommandOptions{
			Command: []string{
				"kubectl",
				"--context",
				context,
				"get",
				"namespaces",
				"-o",
				"name",
			},
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

func (c *CommandExecutorKubernetes) MustCreateNamespaceByName(name string, verbose bool) (createdNamespace Namespace) {
	createdNamespace, err := c.CreateNamespaceByName(name, verbose)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return createdNamespace
}

func (c *CommandExecutorKubernetes) MustDeleteNamespaceByName(name string, verbose bool) {
	err := c.DeleteNamespaceByName(name, verbose)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorKubernetes) MustGetCommandExecutor() (commandExecutor asciichgolangpublic.CommandExecutor) {
	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return commandExecutor
}

func (c *CommandExecutorKubernetes) MustGetKubectlContext(verbose bool) (context string) {
	context, err := c.GetKubectlContext(verbose)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return context
}

func (c *CommandExecutorKubernetes) MustGetKubectlContexts() (contexts []KubectlContext) {
	contexts, err := c.GetKubectlContexts()
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return contexts
}

func (c *CommandExecutorKubernetes) MustGetName() (name string) {
	name, err := c.GetName()
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return name
}

func (c *CommandExecutorKubernetes) MustGetNamespaceByName(name string) (namespace Namespace) {
	namespace, err := c.GetNamespaceByName(name)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return namespace
}

func (c *CommandExecutorKubernetes) MustListNamespaceNames(verbose bool) (namespaceNames []string) {
	namespaceNames, err := c.ListNamespaceNames(verbose)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return namespaceNames
}

func (c *CommandExecutorKubernetes) MustListNamespaces(verbose bool) (namespaces []Namespace) {
	namespaces, err := c.ListNamespaces(verbose)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return namespaces
}

func (c *CommandExecutorKubernetes) MustNamespaceByNameExists(name string, verbose bool) (exists bool) {
	exists, err := c.NamespaceByNameExists(name, verbose)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return exists
}

func (c *CommandExecutorKubernetes) MustRunCommand(runCommandOptions *asciichgolangpublic.RunCommandOptions) (commandOutput *asciichgolangpublic.CommandOutput) {
	commandOutput, err := c.RunCommand(runCommandOptions)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return commandOutput
}

func (c *CommandExecutorKubernetes) MustRunCommandAndGetStdoutAsLines(runCommandOptions *asciichgolangpublic.RunCommandOptions) (lines []string) {
	lines, err := c.RunCommandAndGetStdoutAsLines(runCommandOptions)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return lines
}

func (c *CommandExecutorKubernetes) MustSetCommandExecutor(commandExecutor asciichgolangpublic.CommandExecutor) {
	err := c.SetCommandExecutor(commandExecutor)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorKubernetes) MustSetName(name string) {
	err := c.SetName(name)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorKubernetes) NamespaceByNameExists(name string, verbose bool) (exists bool, err error) {
	if name == "" {
		return false, asciichgolangpublic.TracedErrorEmptyString("name")
	}

	namespaces, err := c.ListNamespaceNames(verbose)
	if err != nil {
		return false, err
	}

	exists = asciichgolangpublic.Slices().ContainsString(namespaces, name)

	if verbose {
		clusterName, err := c.GetName()
		if err != nil {
			return false, err
		}

		if exists {
			asciichgolangpublic.LogInfof(
				"Namespace '%s' exists in kubernetes cluster '%s'.",
				name,
				clusterName,
			)
		} else {
			asciichgolangpublic.LogInfof(
				"Namespace does not '%s' exist in kubernetes cluster '%s'.",
				name,
				clusterName,
			)
		}
	}

	return exists, nil
}

func (c *CommandExecutorKubernetes) RunCommand(runCommandOptions *asciichgolangpublic.RunCommandOptions) (commandOutput *asciichgolangpublic.CommandOutput, err error) {
	if runCommandOptions == nil {
		return nil, asciichgolangpublic.TracedErrorNil("runCommandOptions")
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	return commandExecutor.RunCommand(runCommandOptions)
}

func (c *CommandExecutorKubernetes) RunCommandAndGetStdoutAsLines(runCommandOptions *asciichgolangpublic.RunCommandOptions) (lines []string, err error) {
	if runCommandOptions == nil {
		return nil, asciichgolangpublic.TracedErrorNil("runCommandOptions")
	}

	output, err := c.RunCommand(runCommandOptions)
	if err != nil {
		return nil, err
	}

	return output.GetStdoutAsLines(false)
}

func (c *CommandExecutorKubernetes) SetCommandExecutor(commandExecutor asciichgolangpublic.CommandExecutor) (err error) {
	c.commandExecutor = commandExecutor

	return nil
}

func (c *CommandExecutorKubernetes) SetName(name string) (err error) {
	if name == "" {
		return asciichgolangpublic.TracedErrorf("name is empty string")
	}

	c.name = name

	return nil
}