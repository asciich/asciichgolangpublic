package commandexecutorkubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kuberneteserrors"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type CommandExecutorSecret struct {
	name      string
	namespace kubernetesinterfaces.Namespace
}

func (c *CommandExecutorSecret) GetName() (string, error) {
	if c.name == "" {
		return "", tracederrors.TracedError("name not set")
	}

	return c.name, nil
}

func (c *CommandExecutorSecret) GetCommandExecutorNamespace() (*CommandExecutorNamespace, error) {
	if c.namespace == nil {
		return nil, tracederrors.TracedError("namespace not set")
	}

	namespace, ok := c.namespace.(*CommandExecutorNamespace)
	if !ok {
		return nil, tracederrors.TracedErrorf("namespace is not of type *CommandExecutorNamespace but '%T'", c.namespace)
	}

	return namespace, nil
}

func (c *CommandExecutorSecret) GetNamespace() (kubernetesinterfaces.Namespace, error) {
	if c.namespace == nil {
		return nil, tracederrors.TracedError("namespace not set")
	}

	return c.namespace, nil
}

func (c *CommandExecutorSecret) Exists(ctx context.Context) (bool, error) {
	secretName, err := c.GetName()
	if err != nil {
		return false, err
	}

	namespace, err := c.GetNamespace()
	if err != nil {
		return false, err
	}

	return namespace.SecretByNameExists(ctx, secretName)
}

func (c *CommandExecutorSecret) GetKubernetesCluster() (kubernetesinterfaces.KubernetesCluster, error) {
	namespace, err := c.GetNamespace()
	if err != nil {
		return nil, err
	}

	return namespace.GetKubernetesCluster()
}

func (c *CommandExecutorSecret) GetCachedKubectlContext(ctx context.Context) (string, error) {
	namespace, err := c.GetCommandExecutorNamespace()
	if err != nil {
		return "", err
	}

	return namespace.GetCachedKubectlContext(ctx)
}

func (c *CommandExecutorSecret) RunCommand(ctx context.Context, options *parameteroptions.RunCommandOptions) (*commandoutput.CommandOutput, error) {
	namespace, err := c.GetCommandExecutorNamespace()
	if err != nil {
		return nil, err
	}

	return namespace.RunCommand(ctx, options)
}
func (c *CommandExecutorSecret) Read(ctx context.Context) (map[string][]byte, error) {
	contextName, err := c.GetCachedKubectlContext(ctx)
	if err != nil {
		return nil, err
	}

	secretName, err := c.GetName()
	if err != nil {
		return nil, err
	}

	namespaceName, err := c.GetNamespaceName()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Get secret '%s' in namespace '%s' of kubernetes '%s' started.", secretName, namespaceName, contextName)

	output, err := c.RunCommand(
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
				"-o",
				"json",
			},
		},
	)
	if err != nil {
		expectedNotFound := fmt.Sprintf("Error from server (NotFound): secrets \"%s\" not found", secretName)
		if strings.Contains(err.Error(), expectedNotFound) {
			return nil, tracederrors.TracedErrorf("%w: secret '%s' in namespace '%s' of kubernetes '%s'", kuberneteserrors.ErrSecretNotFound, secretName, namespaceName, contextName)
		}
		return nil, err
	}

	stdout, err := output.GetStdoutAsBytes()
	if err != nil {
		return nil, err
	}

	var secret struct {
		Data map[string][]byte `json:"data"`
	}
	if err := json.Unmarshal(stdout, &secret); err != nil {
		return nil, fmt.Errorf("failed to parse secret JSON: %w", err)
	}

	logging.LogInfoByCtxf(ctx, "Get secret '%s' in namespace '%s' of kubernetes '%s' finished.", secretName, namespaceName, contextName)

	return secret.Data, nil
}

func (c *CommandExecutorSecret) GetNamespaceName() (string, error) {
	namespace, err := c.GetNamespace()
	if err != nil {
		return "", err
	}

	return namespace.GetName()
}

func (c *CommandExecutorSecret) Delete(ctx context.Context) error {
	namespace, err := c.GetNamespace()
	if err != nil {
		return err
	}

	secretName, err := c.GetName()
	if err != nil {
		return err
	}

	return namespace.DeleteSecretByName(ctx, secretName)
}
