package nativekubernetes

import (
	"context"

	"github.com/asciich/asciichgolangpublic/tracederrors"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
		return false, nil
	}

	namespace, err := n.GetNamespace()
	if err != nil {
		return false, nil
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

	configMap, err := clientset.CoreV1().ConfigMaps(namespaceName).Get(ctx, configMapName, metav1.GetOptions{})
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to retrieve raw config map '%s' in namespace '%s': %w", configMapName, namespaceName, err)
	}

	return configMap, nil
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
	rawResponse, err := n.GetRawResponse(ctx)
	if err != nil {
		return nil, err
	}

	return rawResponse.Data, nil
}

func (n *NativeConfigMap) GetAllLabels(ctx context.Context) (map[string]string, error) {
	rawResponse, err := n.GetRawResponse(ctx)
	if err != nil {
		return nil, err
	}

	return rawResponse.Labels, nil
}

func (n *NativeConfigMap) GetData(ctx context.Context, fieldName string) (string, error) {
	if fieldName == "" {
		return "", tracederrors.TracedErrorEmptyString("fileName")
	}

	configMapName, err := n.GetName()
	if err != nil {
		return "", err
	}

	namespaceName, err := n.GetNamespaceName()
	if err != nil {
		return "", err
	}

	data, err := n.GetAllData(ctx)
	if err != nil {
		return "", err
	}

	fieldData, ok := data[fieldName]
	if !ok {
		return "", tracederrors.TracedErrorf("ConfigMap '%s' in namespace '%s' has no field '%s'.", configMapName, namespaceName, fieldName)
	}

	return fieldData, nil
}
