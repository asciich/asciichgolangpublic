package nativekubernetesoo

import (
	"context"
	"reflect"

	"github.com/asciich/asciichgolangpublic/pkg/fileformats/yamlutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	v1 "k8s.io/api/core/v1"
)

type NativeConfigMap struct {
	namespace *NativeNamespace
	name      string
}

func (n *NativeConfigMap) GetNamespace() (*NativeNamespace, error) {
	if n.namespace == nil {
		return nil, tracederrors.TracedError("namespace not set")
	}

	return n.namespace, nil
}

func (n *NativeConfigMap) GetName() (string, error) {
	if n.name == "" {
		return "", tracederrors.TracedError("name not set")
	}

	return n.name, nil
}

func (n *NativeConfigMap) Exists(ctx context.Context) (bool, error) {
	configMapName, err := n.GetName()
	if err != nil {
		return false, err
	}

	namespace, err := n.GetNamespace()
	if err != nil {
		return false, err
	}

	return namespace.ConfigMapByNameExists(ctx, configMapName)
}

func (n *NativeConfigMap) GetRawResponse(ctx context.Context) (*v1.ConfigMap, error) {
	namespace, err := n.GetNamespace()
	if err != nil {
		return nil, err
	}

	clientset, err := namespace.GetClientSet()
	if err != nil {
		return nil, err
	}

	namespaceName, err := namespace.GetName()
	if err != nil {
		return nil, err
	}

	configMapName, err := n.GetName()
	if err != nil {
		return nil, err
	}

	return nativekubernetes.GetConfigMap(ctx, clientset, namespaceName, configMapName)
}

func (n *NativeConfigMap) GetNamespaceName() (string, error) {
	namespace, err := n.GetNamespace()
	if err != nil {
		return "", err
	}

	namespaceName, err := namespace.GetName()
	if err != nil {
		return "", err
	}

	return namespaceName, nil
}

func (n *NativeConfigMap) GetAllData(ctx context.Context) (map[string]string, error) {
	namespace, err := n.GetNamespace()
	if err != nil {
		return nil, err
	}

	clientset, err := namespace.GetClientSet()
	if err != nil {
		return nil, err
	}

	namespaceName, err := namespace.GetName()
	if err != nil {
		return nil, err
	}

	configMapName, err := n.GetName()
	if err != nil {
		return nil, err
	}

	return nativekubernetes.GetConfigMapData(ctx, clientset, namespaceName, configMapName)
}

func (n *NativeConfigMap) GetAllLabels(ctx context.Context) (map[string]string, error) {
	namespace, err := n.GetNamespace()
	if err != nil {
		return nil, err
	}

	clientset, err := namespace.GetClientSet()
	if err != nil {
		return nil, err
	}

	namespaceName, err := namespace.GetName()
	if err != nil {
		return nil, err
	}

	configMapName, err := n.GetName()
	if err != nil {
		return nil, err
	}

	return nativekubernetes.GetConfigMapLabels(ctx, clientset, namespaceName, configMapName)
}

func (n *NativeConfigMap) GetData(ctx context.Context, fieldName string) (string, error) {
	namespace, err := n.GetNamespace()
	if err != nil {
		return "", err
	}

	clientset, err := namespace.GetClientSet()
	if err != nil {
		return "", err
	}

	namespaceName, err := namespace.GetName()
	if err != nil {
		return "", err
	}

	configMapName, err := n.GetName()
	if err != nil {
		return "", err
	}

	return nativekubernetes.GetConfigMapField(ctx, clientset, namespaceName, configMapName, fieldName)
}

func (n *NativeConfigMap) ReadYamlValueAsString(ctx context.Context, yamlPath string) (string, error) {
	if yamlPath == "" {
		return "", tracederrors.TracedErrorEmptyString("yamlPath")
	}

	namespace, err := n.GetNamespace()
	if err != nil {
		return "", err
	}

	namespaceName, err := namespace.GetName()
	if err != nil {
		return "", err
	}

	configMapName, err := n.GetName()
	if err != nil {
		return "", err
	}

	clusterName, err := namespace.GetClusterName()
	if err != nil {
		return "", err
	}

	logging.LogInfoByCtxf(ctx, "Read yaml value '%s' from configmap '%s' in namespace '%s' of kubernetes cluster '%s' started.", yamlPath, configMapName, namespaceName, clusterName)

	data, err := n.GetAllData(ctx)
	if err != nil {
		return "", err
	}

	yamlString, err := kubernetesgeneric.ConfigMapDataToYamlString(configMapName, data)
	if err != nil {
		return "", err
	}

	value, err := yamlutils.RunYqQueryAginstYamlStringAsString(yamlString, yamlPath)
	if err != nil {
		return "", err
	}

	if value == "" || value == "null" {
		return "", tracederrors.TracedErrorf("ConfigMap '%s' in namespace '%s' of kubernetes cluster '%s' has no value at yaml path '%s'.", configMapName, namespaceName, clusterName, yamlPath)
	}

	logging.LogInfoByCtxf(ctx, "Read yaml value '%s' from configmap '%s' in namespace '%s' of kubernetes cluster '%s' finished.", yamlPath, configMapName, namespaceName, clusterName)

	return value, nil
}

func (n *NativeConfigMap) IsYamlValueAsStringEqual(ctx context.Context, yamlPath string, expectedValue string) (bool, error) {
	currentValue, err := n.ReadYamlValueAsString(ctx, yamlPath)
	if err != nil {
		return false, err
	}

	return currentValue == expectedValue, nil
}

func IsConfigMapContentEqual(configMap1 map[string]string, configMap2 map[string]string) bool {
	if len(configMap1) <= 0 && len(configMap2) <= 0 {
		return true
	}

	return reflect.DeepEqual(configMap1, configMap2)
}

func IsConfigMapLabelsEqual(labels1 map[string]string, labels2 map[string]string) bool {
	if len(labels1) <= 0 && len(labels2) <= 0 {
		return true
	}

	return reflect.DeepEqual(labels1, labels2)
}
