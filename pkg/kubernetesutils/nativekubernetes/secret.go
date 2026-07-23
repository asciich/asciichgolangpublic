package nativekubernetes

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kuberneteserrors"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CreateSecret(ctx context.Context, clientset *kubernetes.Clientset, namespaceName string, secretName string, data map[string][]byte) error {
	if clientset == nil {
		return tracederrors.TracedErrorNil("clientset")
	}

	if namespaceName == "" {
		return tracederrors.TracedErrorEmptyString("namespaceName")
	}

	if secretName == "" {
		return tracederrors.TracedErrorEmptyString("secretName")
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespaceName,
		},
		Data: data,
	}

	_, err := clientset.CoreV1().Secrets(namespaceName).Create(ctx, secret, metav1.CreateOptions{})
	if err != nil {
		return tracederrors.TracedErrorf("Failed to create secret '%s' in namespace '%s'.", secretName, namespaceName)
	}

	logging.LogChangedByCtxf(ctx, "Secret '%s' in namespace '%s' created.", secretName, namespaceName)

	return nil
}

func ReadSecret(ctx context.Context, clientset *kubernetes.Clientset, namespaceName string, secretName string) (map[string][]byte, error) {
	if clientset == nil {
		return nil, tracederrors.TracedErrorNil("clientset")
	}

	if namespaceName == "" {
		return nil, tracederrors.TracedErrorEmptyString("namespaceName")
	}

	if secretName == "" {
		return nil, tracederrors.TracedErrorEmptyString("secretName")
	}

	secret, err := clientset.CoreV1().Secrets(namespaceName).Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, tracederrors.TracedErrorf("%w. Secret '%s' does not exist in namespace '%s'", kuberneteserrors.ErrSecretNotFound, secretName, namespaceName)
		}
	}

	return secret.Data, nil
}

func DeleteSecret(ctx context.Context, clientset *kubernetes.Clientset, namespaceName string, secretName string) error {
	if clientset == nil {
		return tracederrors.TracedErrorNil("clientset")
	}

	if namespaceName == "" {
		return tracederrors.TracedErrorEmptyString("namespaceName")
	}

	if secretName == "" {
		return tracederrors.TracedErrorEmptyString("secretName")
	}

	err := clientset.CoreV1().Secrets(namespaceName).Delete(ctx, secretName, metav1.DeleteOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			logging.LogInfoByCtxf(ctx, "Secret '%s' in namespace '%s' is already absent. Skip delete.", secretName, namespaceName)
			return nil
		}
		return tracederrors.TracedErrorf("Failed to delete secret '%s' in namespace '%s'.", secretName, namespaceName)
	}

	logging.LogChangedByCtxf(ctx, "Secret '%s' in namespace '%s' deleted.", secretName, namespaceName)

	return nil
}

func ListSecrets(ctx context.Context, clientset *kubernetes.Clientset, namespaceName string) ([]string, error) {
	if clientset == nil {
		return nil, tracederrors.TracedErrorNil("clientset")
	}

	if namespaceName == "" {
		return nil, tracederrors.TracedErrorEmptyString("namespaceName")
	}

	secretList, err := clientset.CoreV1().Secrets(namespaceName).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to list secrets in namespace '%s'.", namespaceName)
	}

	secretNames := []string{}
	for _, secret := range secretList.Items {
		secretNames = append(secretNames, secret.Name)
	}

	logging.LogInfoByCtxf(ctx, "Found '%d' secrets in namespace '%s'.", len(secretNames), namespaceName)

	return secretNames, nil
}
