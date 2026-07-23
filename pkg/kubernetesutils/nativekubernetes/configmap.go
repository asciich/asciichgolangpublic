package nativekubernetes

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func ConfigMapExists(ctx context.Context, clientset *kubernetes.Clientset, namespaceName string, configMapName string) (bool, error) {
	if clientset == nil {
		return false, tracederrors.TracedErrorNil("clientset")
	}

	if namespaceName == "" {
		return false, tracederrors.TracedErrorEmptyString("namespaceName")
	}

	if configMapName == "" {
		return false, tracederrors.TracedErrorEmptyString("configMapName")
	}

	logging.LogInfoByCtxf(ctx, "Check if ConfigMap '%s' in namespace '%s' exists.", configMapName, namespaceName)

	_, err := clientset.CoreV1().ConfigMaps(namespaceName).Get(ctx, configMapName, metav1.GetOptions{})
	if err != nil {
		logging.LogInfoByCtxf(ctx, "ConfigMap '%s' in namespace '%s' does not exist.", configMapName, namespaceName)
		return false, nil
	}

	logging.LogInfoByCtxf(ctx, "ConfigMap '%s' in namespace '%s' exists.", configMapName, namespaceName)
	return true, nil
}

func GetConfigMap(ctx context.Context, clientset *kubernetes.Clientset, namespaceName string, configMapName string) (*v1.ConfigMap, error) {
	if clientset == nil {
		return nil, tracederrors.TracedErrorNil("clientset")
	}

	if namespaceName == "" {
		return nil, tracederrors.TracedErrorEmptyString("namespaceName")
	}

	if configMapName == "" {
		return nil, tracederrors.TracedErrorEmptyString("configMapName")
	}

	logging.LogInfoByCtxf(ctx, "Get ConfigMap '%s' in namespace '%s'.", configMapName, namespaceName)

	configMap, err := clientset.CoreV1().ConfigMaps(namespaceName).Get(ctx, configMapName, metav1.GetOptions{})
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to retrieve ConfigMap '%s' in namespace '%s': %w", configMapName, namespaceName, err)
	}

	return configMap, nil
}

func GetConfigMapData(ctx context.Context, clientset *kubernetes.Clientset, namespaceName string, configMapName string) (map[string]string, error) {
	configMap, err := GetConfigMap(ctx, clientset, namespaceName, configMapName)
	if err != nil {
		return nil, err
	}

	return configMap.Data, nil
}

func GetConfigMapLabels(ctx context.Context, clientset *kubernetes.Clientset, namespaceName string, configMapName string) (map[string]string, error) {
	configMap, err := GetConfigMap(ctx, clientset, namespaceName, configMapName)
	if err != nil {
		return nil, err
	}

	return configMap.Labels, nil
}

func GetConfigMapField(ctx context.Context, clientset *kubernetes.Clientset, namespaceName string, configMapName string, fieldName string) (string, error) {
	if fieldName == "" {
		return "", tracederrors.TracedErrorEmptyString("fieldName")
	}

	data, err := GetConfigMapData(ctx, clientset, namespaceName, configMapName)
	if err != nil {
		return "", err
	}

	fieldData, ok := data[fieldName]
	if !ok {
		return "", tracederrors.TracedErrorf("ConfigMap '%s' in namespace '%s' has no field '%s'.", configMapName, namespaceName, fieldName)
	}

	return fieldData, nil
}

func DeleteConfigMap(ctx context.Context, clientset *kubernetes.Clientset, namespaceName string, configMapName string) error {
	if clientset == nil {
		return tracederrors.TracedErrorNil("clientset")
	}

	if namespaceName == "" {
		return tracederrors.TracedErrorEmptyString("namespaceName")
	}

	if configMapName == "" {
		return tracederrors.TracedErrorEmptyString("configMapName")
	}

	logging.LogInfoByCtxf(ctx, "Delete ConfigMap '%s' in namespace '%s' started.", configMapName, namespaceName)

	exists, err := ConfigMapExists(ctx, clientset, namespaceName, configMapName)
	if err != nil {
		return err
	}

	if exists {
		err = clientset.CoreV1().ConfigMaps(namespaceName).Delete(ctx, configMapName, metav1.DeleteOptions{})
		if err != nil {
			return tracederrors.TracedErrorf("Failed to delete ConfigMap '%s' in namespace '%s'.", configMapName, namespaceName)
		}

		logging.LogChangedByCtxf(ctx, "ConfigMap '%s' in namespace '%s' deleted.", configMapName, namespaceName)
	} else {
		logging.LogInfoByCtxf(ctx, "ConfigMap '%s' in namespace '%s' does not exist. Skip delete.", configMapName, namespaceName)
	}

	logging.LogInfoByCtxf(ctx, "Delete ConfigMap '%s' in namespace '%s' finished.", configMapName, namespaceName)

	return nil
}

func CreateConfigMap(ctx context.Context, clientset *kubernetes.Clientset, namespaceName string, configMapName string, data map[string]string, labels map[string]string) error {
	if clientset == nil {
		return tracederrors.TracedErrorNil("clientset")
	}

	if namespaceName == "" {
		return tracederrors.TracedErrorEmptyString("namespaceName")
	}

	if configMapName == "" {
		return tracederrors.TracedErrorEmptyString("configMapName")
	}

	logging.LogInfoByCtxf(ctx, "Create ConfigMap '%s' in namespace '%s' started.", configMapName, namespaceName)

	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:   configMapName,
			Labels: labels,
		},
		Data: data,
	}

	_, err := clientset.CoreV1().ConfigMaps(namespaceName).Create(ctx, configMap, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			logging.LogInfoByCtxf(ctx, "ConfigMap '%s' in namespace '%s' already exists.", configMapName, namespaceName)
			return nil
		}
		return tracederrors.TracedErrorf("Failed to create ConfigMap '%s' in namespace '%s': %w", configMapName, namespaceName, err)
	}

	logging.LogChangedByCtxf(ctx, "ConfigMap '%s' in namespace '%s' created.", configMapName, namespaceName)
	logging.LogInfoByCtxf(ctx, "Create ConfigMap '%s' in namespace '%s' finished.", configMapName, namespaceName)

	return nil
}

func UpdateConfigMap(ctx context.Context, clientset *kubernetes.Clientset, namespaceName string, configMapName string, data map[string]string, labels map[string]string) error {
	if clientset == nil {
		return tracederrors.TracedErrorNil("clientset")
	}

	if namespaceName == "" {
		return tracederrors.TracedErrorEmptyString("namespaceName")
	}

	if configMapName == "" {
		return tracederrors.TracedErrorEmptyString("configMapName")
	}

	logging.LogInfoByCtxf(ctx, "Update ConfigMap '%s' in namespace '%s' started.", configMapName, namespaceName)

	configMap, err := clientset.CoreV1().ConfigMaps(namespaceName).Get(ctx, configMapName, metav1.GetOptions{})
	if err != nil {
		return tracederrors.TracedErrorf("Failed to get ConfigMap '%s' in namespace '%s' for update: %w", configMapName, namespaceName, err)
	}

	configMap.Data = data
	if labels != nil {
		configMap.Labels = labels
	}

	_, err = clientset.CoreV1().ConfigMaps(namespaceName).Update(ctx, configMap, metav1.UpdateOptions{})
	if err != nil {
		return tracederrors.TracedErrorf("Failed to update ConfigMap '%s' in namespace '%s': %w", configMapName, namespaceName, err)
	}

	logging.LogChangedByCtxf(ctx, "ConfigMap '%s' in namespace '%s' updated.", configMapName, namespaceName)
	logging.LogInfoByCtxf(ctx, "Update ConfigMap '%s' in namespace '%s' finished.", configMapName, namespaceName)

	return nil
}

func ListConfigMaps(ctx context.Context, clientset *kubernetes.Clientset, namespaceName string) ([]string, error) {
	if clientset == nil {
		return nil, tracederrors.TracedErrorNil("clientset")
	}

	if namespaceName == "" {
		return nil, tracederrors.TracedErrorEmptyString("namespaceName")
	}

	configMapList, err := clientset.CoreV1().ConfigMaps(namespaceName).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to list ConfigMaps in namespace '%s'.", namespaceName)
	}

	configMapNames := []string{}
	for _, cm := range configMapList.Items {
		configMapNames = append(configMapNames, cm.Name)
	}

	logging.LogInfoByCtxf(ctx, "Found '%d' ConfigMaps in namespace '%s'.", len(configMapNames), namespaceName)

	return configMapNames, nil
}
