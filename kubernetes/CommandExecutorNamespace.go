package kubernetes

import (
	"fmt"
	"strings"

	"github.com/asciich/asciichgolangpublic"
	"github.com/asciich/asciichgolangpublic/datatypes"
	aslices "github.com/asciich/asciichgolangpublic/datatypes/slices"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type CommandExecutorNamespace struct {
	name              string
	kubernetesCluster KubernetesCluster
}

func NewCommandExecutorNamespace() (c *CommandExecutorNamespace) {
	return new(CommandExecutorNamespace)
}

func (c *CommandExecutorNamespace) CreateRole(createOptions *CreateRoleOptions) (createdRole Role, err error) {
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

	context, err := c.GetCachedKubectlContext(createOptions.Verbose)
	if err != nil {
		return nil, err
	}

	exists, err := c.RoleByNameExists(roleName, createOptions.Verbose)
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
			&asciichgolangpublic.RunCommandOptions{
				Command: command,
			},
		)
		if err != nil {
			return nil, err
		}
	}

	return c.GetRoleByName(roleName)
}

func (c *CommandExecutorNamespace) DeleteRoleByName(name string, verbose bool) (err error) {
	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	namespaceName, err := c.GetName()
	if err != nil {
		return err
	}

	exists, err := c.RoleByNameExists(name, verbose)
	if err != nil {
		return err
	}

	clusterName, err := c.GetClusterName()
	if err != nil {
		return err
	}

	if exists {
		context, err := c.GetCachedKubectlContext(verbose)
		if err != nil {
			return err
		}

		_, err = c.RunCommand(
			&asciichgolangpublic.RunCommandOptions{
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

		if verbose {
			logging.LogChangedf(
				"Role '%s' in namespace '%s' in kubernetes cluster '%s' deleted.",
				name,
				namespaceName,
				clusterName,
			)
		}
	} else {
		if verbose {
			logging.LogChangedf(
				"Role '%s' in namespace '%s' in kubernetes cluster '%s' already absent.",
				name,
				namespaceName,
				clusterName,
			)
		}
	}

	return nil
}

func (c *CommandExecutorNamespace) GetCachedKubectlContext(verbose bool) (context string, err error) {
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

	return commandExecutorKubernetes.GetCachedKubectlContext(verbose)
}

func (c *CommandExecutorNamespace) GetClusterName() (clusterName string, err error) {
	kubernetesCluster, err := c.GetKubernetesCluster()
	if err != nil {
		return "", err
	}

	return kubernetesCluster.GetName()
}

func (c *CommandExecutorNamespace) GetCommandExecutor() (commandExecutor asciichgolangpublic.CommandExecutor, err error) {
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

func (c *CommandExecutorNamespace) GetKubernetesCluster() (kubernetesCluster KubernetesCluster, err error) {

	return c.kubernetesCluster, nil
}

func (c *CommandExecutorNamespace) GetName() (name string, err error) {
	if c.name == "" {
		return "", tracederrors.TracedErrorf("name not set")
	}

	return c.name, nil
}

func (c *CommandExecutorNamespace) GetRoleByName(name string) (role Role, err error) {
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

func (c *CommandExecutorNamespace) ListRoleNames(verbose bool) (roleNames []string, err error) {
	context, err := c.GetCachedKubectlContext(verbose)
	if err != nil {
		return nil, err
	}

	name, err := c.GetName()
	if err != nil {
		return nil, err
	}

	lines, err := c.RunCommandAndGetStdoutAsLines(
		&asciichgolangpublic.RunCommandOptions{
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

func (c *CommandExecutorNamespace) MustCreateRole(createOptions *CreateRoleOptions) (createdRole Role) {
	createdRole, err := c.CreateRole(createOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return createdRole
}

func (c *CommandExecutorNamespace) MustDeleteRoleByName(name string, verbose bool) {
	err := c.DeleteRoleByName(name, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorNamespace) MustGetCachedKubectlContext(verbose bool) (context string) {
	context, err := c.GetCachedKubectlContext(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return context
}

func (c *CommandExecutorNamespace) MustGetClusterName() (clusterName string) {
	clusterName, err := c.GetClusterName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return clusterName
}

func (c *CommandExecutorNamespace) MustGetCommandExecutor() (commandExecutor asciichgolangpublic.CommandExecutor) {
	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return commandExecutor
}

func (c *CommandExecutorNamespace) MustGetKubernetesCluster() (kubernetesCluster KubernetesCluster) {
	kubernetesCluster, err := c.GetKubernetesCluster()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return kubernetesCluster
}

func (c *CommandExecutorNamespace) MustGetName() (name string) {
	name, err := c.GetName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return name
}

func (c *CommandExecutorNamespace) MustGetRoleByName(name string) (role Role) {
	role, err := c.GetRoleByName(name)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return role
}

func (c *CommandExecutorNamespace) MustListRoleNames(verbose bool) (roleNames []string) {
	roleNames, err := c.ListRoleNames(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return roleNames
}

func (c *CommandExecutorNamespace) MustRoleByNameExists(name string, verbose bool) (exists bool) {
	exists, err := c.RoleByNameExists(name, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return exists
}

func (c *CommandExecutorNamespace) MustRunCommand(runCommandOptions *asciichgolangpublic.RunCommandOptions) (commandOutput *asciichgolangpublic.CommandOutput) {
	commandOutput, err := c.RunCommand(runCommandOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return commandOutput
}

func (c *CommandExecutorNamespace) MustRunCommandAndGetStdoutAsLines(runCommandOptions *asciichgolangpublic.RunCommandOptions) (lines []string) {
	lines, err := c.RunCommandAndGetStdoutAsLines(runCommandOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return lines
}

func (c *CommandExecutorNamespace) MustSetKubernetesCluster(kubernetesCluster KubernetesCluster) {
	err := c.SetKubernetesCluster(kubernetesCluster)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorNamespace) MustSetName(name string) {
	err := c.SetName(name)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorNamespace) RoleByNameExists(name string, verbose bool) (exists bool, err error) {
	if name == "" {
		return false, tracederrors.TracedErrorEmptyString("name")
	}

	roleNames, err := c.ListRoleNames(false)
	if err != nil {
		return false, err
	}

	exists = aslices.ContainsString(roleNames, name)

	if verbose {
		clusterName, err := c.GetClusterName()
		if err != nil {
			return false, err
		}

		if exists {
			logging.LogInfof(
				"Role '%s' in kubernetes cluster '%s' exists.",
				name,
				clusterName,
			)
		} else {
			logging.LogInfof(
				"Role '%s' in kubernetes cluster '%s' does not exist.",
				name,
				clusterName,
			)
		}
	}

	return exists, nil
}

func (c *CommandExecutorNamespace) RunCommand(runCommandOptions *asciichgolangpublic.RunCommandOptions) (commandOutput *asciichgolangpublic.CommandOutput, err error) {
	if runCommandOptions == nil {
		return nil, tracederrors.TracedErrorNil("runCommandOptions")
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	return commandExecutor.RunCommand(runCommandOptions)
}

func (c *CommandExecutorNamespace) RunCommandAndGetStdoutAsLines(runCommandOptions *asciichgolangpublic.RunCommandOptions) (lines []string, err error) {
	if runCommandOptions == nil {
		return nil, tracederrors.TracedErrorNil("runCommandOptions")
	}

	commandOutput, err := c.RunCommand(runCommandOptions)
	if err != nil {
		return nil, err
	}

	return commandOutput.GetStdoutAsLines(false)
}

func (c *CommandExecutorNamespace) SetKubernetesCluster(kubernetesCluster KubernetesCluster) (err error) {
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
