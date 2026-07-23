package commandexecutorkubernetes

import (
	"context"
	"encoding/json"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type CommandExecutorConfigMap struct {
	name      string
	namespace kubernetesinterfaces.Namespace
}

func (c *CommandExecutorConfigMap) GetName() (string, error) {
	if c.name == "" {
		return "", tracederrors.TracedError("name not set")
	}

	return c.name, nil
}

func (c *CommandExecutorConfigMap) GetCommandExecutorNamespace() (*CommandExecutorNamespace, error) {
	if c.namespace == nil {
		return nil, tracederrors.TracedError("namespace not set")
	}

	namespace, ok := c.namespace.(*CommandExecutorNamespace)
	if !ok {
		return nil, tracederrors.TracedErrorf("namespace is not of type *CommandExecutorNamespace but '%T'", c.namespace)
	}

	return namespace, nil
}

func (c *CommandExecutorConfigMap) GetNamespace() (kubernetesinterfaces.Namespace, error) {
	if c.namespace == nil {
		return nil, tracederrors.TracedError("namespace not set")
	}

	return c.namespace, nil
}

func (c *CommandExecutorConfigMap) Exists(ctx context.Context) (bool, error) {
	configMapName, err := c.GetName()
	if err != nil {
		return false, err
	}

	namespace, err := c.GetNamespace()
	if err != nil {
		return false, err
	}

	return namespace.ConfigMapByNameExists(ctx, configMapName)
}

func (c *CommandExecutorConfigMap) GetKubernetesCluster() (kubernetesinterfaces.KubernetesCluster, error) {
	namespace, err := c.GetNamespace()
	if err != nil {
		return nil, err
	}

	return namespace.GetKubernetesCluster()
}

func (c *CommandExecutorConfigMap) GetCachedKubectlContext(ctx context.Context) (string, error) {
	namespace, err := c.GetCommandExecutorNamespace()
	if err != nil {
		return "", err
	}

	return namespace.GetCachedKubectlContext(ctx)
}

func (c *CommandExecutorConfigMap) RunCommand(ctx context.Context, options *parameteroptions.RunCommandOptions) (*commandoutput.CommandOutput, error) {
	namespace, err := c.GetCommandExecutorNamespace()
	if err != nil {
		return nil, err
	}

	return namespace.RunCommand(ctx, options)
}

func (c *CommandExecutorConfigMap) GetAllData(ctx context.Context) (map[string]string, error) {
	contextName, err := c.GetCachedKubectlContext(ctx)
	if err != nil {
		return nil, err
	}

	configMapName, err := c.GetName()
	if err != nil {
		return nil, err
	}

	namespaceName, err := c.GetNamespaceName()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Get configmap '%s' in namespace '%s' of kubernetes '%s' started.", configMapName, namespaceName, contextName)

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
				"configmap",
				configMapName,
				"-o",
				"json",
			},
		},
	)
	if err != nil {
		return nil, err
	}

	stdout, err := output.GetStdoutAsBytes()
	if err != nil {
		return nil, err
	}

	var configMap struct {
		Data map[string]string `json:"data"`
	}
	if err := json.Unmarshal(stdout, &configMap); err != nil {
		return nil, tracederrors.TracedErrorf("failed to parse configmap JSON: %w", err)
	}

	logging.LogInfoByCtxf(ctx, "Get configmap '%s' in namespace '%s' of kubernetes '%s' finished.", configMapName, namespaceName, contextName)

	return configMap.Data, nil
}

func (c *CommandExecutorConfigMap) GetAllLabels(ctx context.Context) (map[string]string, error) {
	contextName, err := c.GetCachedKubectlContext(ctx)
	if err != nil {
		return nil, err
	}

	configMapName, err := c.GetName()
	if err != nil {
		return nil, err
	}

	namespaceName, err := c.GetNamespaceName()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Get configmap labels '%s' in namespace '%s' of kubernetes '%s' started.", configMapName, namespaceName, contextName)

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
				"configmap",
				configMapName,
				"-o",
				"json",
			},
		},
	)
	if err != nil {
		return nil, err
	}

	stdout, err := output.GetStdoutAsBytes()
	if err != nil {
		return nil, err
	}

	type MetaData struct {
		Labels map[string]string `json:"labels"`
	}

	var configMap struct {
		Metadata MetaData `json:"metadata"`
	}
	if err := json.Unmarshal(stdout, &configMap); err != nil {
		return nil, tracederrors.TracedErrorf("failed to parse configmap JSON: %w", err)
	}

	if configMap.Metadata.Labels == nil {
		configMap.Metadata.Labels = make(map[string]string)
	}

	logging.LogInfoByCtxf(ctx, "Get configmap labels '%s' in namespace '%s' of kubernetes '%s' finished.", configMapName, namespaceName, contextName)

	return configMap.Metadata.Labels, nil
}

func (c *CommandExecutorConfigMap) GetData(ctx context.Context, fieldName string) (string, error) {
	data, err := c.GetAllData(ctx)
	if err != nil {
		return "", err
	}

	value, ok := data[fieldName]
	if !ok {
		return "", tracederrors.TracedErrorf("ConfigMap '%s' has no field '%s'", c.name, fieldName)
	}

	return value, nil
}

func (c *CommandExecutorConfigMap) GetNamespaceName() (string, error) {
	namespace, err := c.GetNamespace()
	if err != nil {
		return "", err
	}

	return namespace.GetName()
}

func (c *CommandExecutorConfigMap) Delete(ctx context.Context) error {
	namespace, err := c.GetNamespace()
	if err != nil {
		return err
	}

	configMapName, err := c.GetName()
	if err != nil {
		return err
	}

	return namespace.DeleteConfigMapByName(ctx, configMapName)
}
