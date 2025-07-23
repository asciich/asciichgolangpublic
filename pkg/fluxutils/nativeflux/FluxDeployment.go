package nativeflux

import (
	"context"
	"time"

	"gitlab.asciich.ch/tools/asciichgolangpublic.git/datatypes"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/kubernetesutils/kubernetesinterfaces"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/kubernetesutils/nativekubernetes"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/logging"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/tracederrors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

type FluxDeployment struct {
	cluster *nativekubernetes.NativeKubernetesCluster

	namespace string
}

func GetFluxDeployment(cluster kubernetesinterfaces.KubernetesCluster, namespaceName string) (*FluxDeployment, error) {
	if cluster == nil {
		return nil, tracederrors.TracedErrorNil("cluster")
	}

	nativeCluster, ok := cluster.(*nativekubernetes.NativeKubernetesCluster)
	if !ok {
		return nil, tracederrors.TracedErrorf("cluster is not a native kubernetes cluster, it is of type: '%s'", datatypes.MustGetTypeName(cluster))
	}

	if namespaceName == "" {
		return nil, tracederrors.TracedErrorEmptyString("namespaceName")
	}

	return &FluxDeployment{
		namespace: namespaceName,
		cluster:   nativeCluster,
	}, nil
}

func (f *FluxDeployment) GetKubernetesCluster() (*nativekubernetes.NativeKubernetesCluster, error) {
	if f.cluster == nil {
		return nil, tracederrors.TracedError("kubernetes cluster not set")
	}

	return f.cluster, nil
}

func (f *FluxDeployment) GetNamespaceName() (string, error) {
	if f.namespace == "" {
		return "", tracederrors.TracedError("Namespace not set")
	}

	return f.namespace, nil
}

func (f *FluxDeployment) GetDynamicClient() (*dynamic.DynamicClient, error) {
	cluster, err := f.GetKubernetesCluster()
	if err != nil {
		return nil, err
	}

	return cluster.GetDynamicClient()
}

func GetGitRepositoryGVR() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "source.toolkit.fluxcd.io",
		Version:  "v1",
		Resource: "gitrepositories",
	}
}

func GetKustomizationGVR() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "kustomize.toolkit.fluxcd.io",
		Version:  "v1",
		Resource: "kustomizations",
	}
}

func GetHelmReleaseGVR() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "helm.toolkit.fluxcd.io",
		Version:  "v2",
		Resource: "helmreleases",
	}
}

func (f *FluxDeployment) WaitUntilGitRepositoryDeleted(ctx context.Context, name string, namespaceName string) error {
	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	if namespaceName == "" {
		return tracederrors.TracedErrorEmptyString("namespace")
	}

	logging.LogInfoByCtxf(ctx, "Wait for Flux GitRepository '%s' in namespace '%s' deleted started.", name, namespaceName)

	for {
		err := ctx.Err()
		if err != nil {
			return tracederrors.TracedErrorf("context errored while waiting for Flux GitRepository '%s' in namespace '%s' to be deleted: %w", name, namespaceName, err)
		}

		exists, err := f.GitRepositoryExists(ctx, name, namespaceName)
		if err != nil {
			return err
		}

		if exists {
			logging.LogInfoByCtxf(ctx, "Wait for Flux GitRepository '%s' in namespace '%s' to be deleted.", name, namespaceName)
			time.Sleep(time.Second * 1)
			continue
		}

		break
	}

	logging.LogInfoByCtxf(ctx, "Flux GitRepository '%s' in namespace '%s' is deleted.", name, namespaceName)

	return nil
}

func (f *FluxDeployment) WaitUntilHelmReleaseDeleted(ctx context.Context, name string, namespaceName string) error {
	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	if namespaceName == "" {
		return tracederrors.TracedErrorEmptyString("namespace")
	}

	logging.LogInfoByCtxf(ctx, "Wait for Flux HelmRelease '%s' in namespace '%s' deleted started.", name, namespaceName)

	for {
		err := ctx.Err()
		if err != nil {
			return tracederrors.TracedErrorf("context errored while waiting for Flux HelmRelease '%s' in namespace '%s' to be deleted: %w", name, namespaceName, err)
		}

		exists, err := f.HelmReleaseExists(ctx, name, namespaceName)
		if err != nil {
			return err
		}

		if exists {
			logging.LogInfoByCtxf(ctx, "Wait for Flux HelmRelease '%s' in namespace '%s' to be deleted.", name, namespaceName)
			time.Sleep(time.Second * 1)
			continue
		}

		break
	}

	logging.LogInfoByCtxf(ctx, "Flux HelmRelease '%s' in namespace '%s' is deleted.", name, namespaceName)

	return nil
}

func (f *FluxDeployment) WaitUntilKustomizationDeleted(ctx context.Context, name string, namespaceName string) error {
	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	if namespaceName == "" {
		return tracederrors.TracedErrorEmptyString("namespace")
	}

	logging.LogInfoByCtxf(ctx, "Wait for Flux Kustomization '%s' in namespace '%s' deleted started.", name, namespaceName)

	for {
		err := ctx.Err()
		if err != nil {
			return tracederrors.TracedErrorf("context errored while waiting for Flux Kustomization '%s' in namespace '%s' to be deleted: %w", name, namespaceName, err)
		}

		exists, err := f.KustomizationExists(ctx, name, namespaceName)
		if err != nil {
			return err
		}

		if exists {
			logging.LogInfoByCtxf(ctx, "Wait for Flux Kustomization '%s' in namespace '%s' to be deleted.", name, namespaceName)
			time.Sleep(time.Second * 1)
			continue
		}

		break
	}

	logging.LogInfoByCtxf(ctx, "Flux Kustomization '%s' in namespace '%s' is deleted.", name, namespaceName)

	return nil
}

func (f *FluxDeployment) DeleteGitRepository(ctx context.Context, name string, namespace string) error {
	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	if namespace == "" {
		return tracederrors.TracedErrorEmptyString("namespace")
	}

	exists, err := f.GitRepositoryExists(ctx, name, namespace)
	if err != nil {
		return err
	}

	if exists {
		dynamicClient, err := f.GetDynamicClient()
		if err != nil {
			return err
		}

		deleteOptions := metav1.DeleteOptions{
			GracePeriodSeconds: func() *int64 {
				grace := int64(0)
				return &grace
			}(),
		}

		err = dynamicClient.Resource(GetGitRepositoryGVR()).Namespace(namespace).Delete(ctx, name, deleteOptions)

		timeoutCtx, _ := context.WithTimeout(ctx, time.Second*10)
		err = f.WaitUntilGitRepositoryDeleted(timeoutCtx, name, namespace)
		if err != nil {
			return err
		}

		logging.LogChangedByCtxf(ctx, "Flux GitRepository '%s' in namespace '%s' deleted.", name, namespace)
	} else {
		logging.LogInfoByCtxf(ctx, "Flux GitRepository '%s' is already absent in namespace '%s'. Skip deletion.", name, namespace)
	}

	return nil
}

func (f *FluxDeployment) DeleteKustomization(ctx context.Context, name string, namespace string) error {
	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	if namespace == "" {
		return tracederrors.TracedErrorEmptyString("namespace")
	}

	exists, err := f.KustomizationExists(ctx, name, namespace)
	if err != nil {
		return err
	}

	if exists {
		dynamicClient, err := f.GetDynamicClient()
		if err != nil {
			return err
		}

		deleteOptions := metav1.DeleteOptions{
			GracePeriodSeconds: func() *int64 {
				grace := int64(0)
				return &grace
			}(),
		}

		err = dynamicClient.Resource(GetKustomizationGVR()).Namespace(namespace).Delete(ctx, name, deleteOptions)

		timeoutCtx, _ := context.WithTimeout(ctx, time.Second*10)
		err = f.WaitUntilKustomizationDeleted(timeoutCtx, name, namespace)
		if err != nil {
			return err
		}

		logging.LogChangedByCtxf(ctx, "Flux Kustomization '%s' in namespace '%s' deleted.", name, namespace)
	} else {
		logging.LogInfoByCtxf(ctx, "Flux Kustomization '%s' is already absent in namespace '%s'. Skip deletion.", name, namespace)
	}

	return nil
}

func (f *FluxDeployment) DeleteHelmRelease(ctx context.Context, name string, namespace string) error {
	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	if namespace == "" {
		return tracederrors.TracedErrorEmptyString("namespace")
	}

	exists, err := f.HelmReleaseExists(ctx, name, namespace)
	if err != nil {
		return err
	}

	if exists {
		dynamicClient, err := f.GetDynamicClient()
		if err != nil {
			return err
		}

		deleteOptions := metav1.DeleteOptions{
			GracePeriodSeconds: func() *int64 {
				grace := int64(0)
				return &grace
			}(),
		}

		err = dynamicClient.Resource(GetHelmReleaseGVR()).Namespace(namespace).Delete(ctx, name, deleteOptions)

		timeoutCtx, _ := context.WithTimeout(ctx, time.Second*10)
		err = f.WaitUntilHelmReleaseDeleted(timeoutCtx, name, namespace)
		if err != nil {
			return err
		}

		logging.LogChangedByCtxf(ctx, "Flux HelmRelease '%s' in namespace '%s' deleted.", name, namespace)
	} else {
		logging.LogInfoByCtxf(ctx, "Flux HelmRelease '%s' is already absent in namespace '%s'. Skip deletion.", name, namespace)
	}

	return nil
}

func (f *FluxDeployment) GitRepositoryExists(ctx context.Context, name string, namespace string) (bool, error) {
	if name == "" {
		return false, tracederrors.TracedErrorEmptyString("name")
	}

	if namespace == "" {
		return false, tracederrors.TracedErrorEmptyString("namespace")
	}

	dynamicClient, err := f.GetDynamicClient()
	if err != nil {
		return false, err
	}

	_, err = dynamicClient.Resource(GetGitRepositoryGVR()).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})

	var exists bool
	if err == nil {
		exists = true
	} else {
		if apierrors.IsNotFound(err) {
			exists = false
		} else {
			return false, tracederrors.TracedErrorf("Failed to request if flux gitrepository exists: %w", err)
		}
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Flux GitRepository '%s' in namespace '%s' exists.", name, namespace)
	} else {
		logging.LogInfoByCtxf(ctx, "Flux GitRepository '%s' in namespace '%s' does not exist.", name, namespace)
	}

	return exists, nil
}

func (f *FluxDeployment) HelmReleaseExists(ctx context.Context, name string, namespace string) (bool, error) {
	if name == "" {
		return false, tracederrors.TracedErrorEmptyString("name")
	}

	if namespace == "" {
		return false, tracederrors.TracedErrorEmptyString("namespace")
	}

	dynamicClient, err := f.GetDynamicClient()
	if err != nil {
		return false, err
	}

	_, err = dynamicClient.Resource(GetHelmReleaseGVR()).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})

	var exists bool
	if err == nil {
		exists = true
	} else {
		if apierrors.IsNotFound(err) {
			exists = false
		} else {
			return false, tracederrors.TracedErrorf("Failed to request if flux HelmRelease exists: %w", err)
		}
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Flux HelmRelease '%s' in namespace '%s' exists.", name, namespace)
	} else {
		logging.LogInfoByCtxf(ctx, "Flux HelmRelease '%s' in namespace '%s' does not exist.", name, namespace)
	}

	return exists, nil
}

func (f *FluxDeployment) KustomizationExists(ctx context.Context, name string, namespace string) (bool, error) {
	if name == "" {
		return false, tracederrors.TracedErrorEmptyString("name")
	}

	if namespace == "" {
		return false, tracederrors.TracedErrorEmptyString("namespace")
	}

	dynamicClient, err := f.GetDynamicClient()
	if err != nil {
		return false, err
	}

	_, err = dynamicClient.Resource(GetKustomizationGVR()).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})

	var exists bool
	if err == nil {
		exists = true
	} else {
		if apierrors.IsNotFound(err) {
			exists = false
		} else {
			return false, tracederrors.TracedErrorf("Failed to request if flux kustomization exists: %w", err)
		}
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Flux Kustomization '%s' in namespace '%s' exists.", name, namespace)
	} else {
		logging.LogInfoByCtxf(ctx, "Flux Kustomization '%s' in namespace '%s' does not exist.", name, namespace)
	}

	return exists, nil
}

func (f *FluxDeployment) GetGitRepositoryStatusMessage(ctx context.Context, name string, namespaceName string) (string, error) {
	if name == "" {
		return "", tracederrors.TracedErrorEmptyString("name")
	}

	if namespaceName == "" {
		return "", tracederrors.TracedErrorEmptyString("namespace")
	}

	dynamicClient, err := f.GetDynamicClient()
	if err != nil {
		return "", err
	}

	unstructuredObj, err := dynamicClient.Resource(GetGitRepositoryGVR()).Namespace(namespaceName).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return "", tracederrors.TracedErrorf("Unable to get Flux GitRepository status of '%s' in namespace '%s': GitRepository does not exist.", name, namespaceName)
		} else {
			return "", tracederrors.TracedErrorf("Failed to get Flux GitRepository '%s' in namespace '%s': %w", name, namespaceName, err)
		}
	}

	status, found := unstructuredObj.Object["status"]
	if !found {
		return "", tracederrors.TracedErrorf("Flux GitRepository '%s' in namespace '%s' does not have a 'status' field", name, namespaceName)
	}

	statusMap, ok := status.(map[string]interface{})
	if !ok {
		return "", tracederrors.TracedErrorf("GitRepository '%s' in namespace '%s' 'status' field is not a map.", name, namespaceName)
	}

	conditions, ok := statusMap["conditions"]
	if !ok {
		return "", tracederrors.TracedErrorf("GitRepository '%s' in namespace '%s' 'statusMap' has no conditions.", name, namespaceName)
	}

	conditionsList, ok := conditions.([]interface{})
	if !ok {
		return "", tracederrors.TracedErrorf("GitRepository '%s' in namespace '%s' 'conditions' is not a list.", name, namespaceName)
	}

	var statusMessage string
	for _, cond := range conditionsList {
		conditionMap, ok := cond.(map[string]interface{})
		if !ok {
			continue
		}

		message, ok := conditionMap["message"]
		if !ok {
			continue
		}

		messageString, ok := message.(string)
		if !ok {
			continue
		}

		statusMessage = messageString
	}

	if statusMessage == "" {
		return "", tracederrors.TracedErrorf("No status message found for Flux GitRepository '%s' in namespace '%s'.", name, namespaceName)
	}

	logging.LogInfoByCtxf(ctx, "Status message of Flux GitRepository '%s' is '%s'.", name, namespaceName)

	return statusMessage, nil
}

func (f *FluxDeployment) GetKustomizationStatusMessage(ctx context.Context, name string, namespaceName string) (string, error) {
	if name == "" {
		return "", tracederrors.TracedErrorEmptyString("name")
	}

	if namespaceName == "" {
		return "", tracederrors.TracedErrorEmptyString("namespace")
	}

	dynamicClient, err := f.GetDynamicClient()
	if err != nil {
		return "", err
	}

	unstructuredObj, err := dynamicClient.Resource(GetKustomizationGVR()).Namespace(namespaceName).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return "", tracederrors.TracedErrorf("Unable to get Flux Kustomization status of '%s' in namespace '%s': Kustomization does not exist.", name, namespaceName)
		} else {
			return "", tracederrors.TracedErrorf("Failed to get Flux Kustomization '%s' in namespace '%s': %w", name, namespaceName, err)
		}
	}

	status, found := unstructuredObj.Object["status"]
	if !found {
		return "", tracederrors.TracedErrorf("Flux Kustomization '%s' in namespace '%s' does not have a 'status' field", name, namespaceName)
	}

	statusMap, ok := status.(map[string]interface{})
	if !ok {
		return "", tracederrors.TracedErrorf("Kustomization '%s' in namespace '%s' 'status' field is not a map.", name, namespaceName)
	}

	conditions, ok := statusMap["conditions"]
	if !ok {
		return "", tracederrors.TracedErrorf("Kustomization '%s' in namespace '%s' 'statusMap' has no conditions.", name, namespaceName)
	}

	conditionsList, ok := conditions.([]interface{})
	if !ok {
		return "", tracederrors.TracedErrorf("Kustomization '%s' in namespace '%s' 'conditions' is not a list.", name, namespaceName)
	}

	var statusMessage string
	for _, cond := range conditionsList {
		conditionMap, ok := cond.(map[string]interface{})
		if !ok {
			continue
		}

		message, ok := conditionMap["message"]
		if !ok {
			continue
		}

		messageString, ok := message.(string)
		if !ok {
			continue
		}

		statusMessage = messageString
	}

	if statusMessage == "" {
		return "", tracederrors.TracedErrorf("No status message found for Flux Kustomization '%s' in namespace '%s'.", name, namespaceName)
	}

	logging.LogInfoByCtxf(ctx, "Status message of Flux Kustomization '%s' is '%s'.", name, namespaceName)

	return statusMessage, nil
}

func (f *FluxDeployment) GetHelmReleaseStatusMessage(ctx context.Context, name string, namespaceName string) (string, error) {
	if name == "" {
		return "", tracederrors.TracedErrorEmptyString("name")
	}

	if namespaceName == "" {
		return "", tracederrors.TracedErrorEmptyString("namespace")
	}

	dynamicClient, err := f.GetDynamicClient()
	if err != nil {
		return "", err
	}

	unstructuredObj, err := dynamicClient.Resource(GetHelmReleaseGVR()).Namespace(namespaceName).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return "", tracederrors.TracedErrorf("Unable to get Flux HelmRelease status of '%s' in namespace '%s': Kustomization does not exist.", name, namespaceName)
		} else {
			return "", tracederrors.TracedErrorf("Failed to get Flux HelmRelease '%s' in namespace '%s': %w", name, namespaceName, err)
		}
	}

	status, found := unstructuredObj.Object["status"]
	if !found {
		return "", tracederrors.TracedErrorf("Flux HelmRelease '%s' in namespace '%s' does not have a 'status' field", name, namespaceName)
	}

	statusMap, ok := status.(map[string]interface{})
	if !ok {
		return "", tracederrors.TracedErrorf("HelmRelease '%s' in namespace '%s' 'status' field is not a map.", name, namespaceName)
	}

	conditions, ok := statusMap["conditions"]
	if !ok {
		return "", tracederrors.TracedErrorf("HelmRelease '%s' in namespace '%s' 'statusMap' has no conditions.", name, namespaceName)
	}

	conditionsList, ok := conditions.([]interface{})
	if !ok {
		return "", tracederrors.TracedErrorf("HelmRelease '%s' in namespace '%s' 'conditions' is not a list.", name, namespaceName)
	}

	var statusMessage string
	for _, cond := range conditionsList {
		conditionMap, ok := cond.(map[string]interface{})
		if !ok {
			continue
		}

		message, ok := conditionMap["message"]
		if !ok {
			continue
		}

		messageString, ok := message.(string)
		if !ok {
			continue
		}

		statusMessage = messageString
	}

	if statusMessage == "" {
		return "", tracederrors.TracedErrorf("No status message found for Flux HelmRelease '%s' in namespace '%s'.", name, namespaceName)
	}

	logging.LogInfoByCtxf(ctx, "Status message of Flux HelmRelease '%s' is '%s'.", name, namespaceName)

	return statusMessage, nil
}

func (f *FluxDeployment) WatchGitRepository(ctx context.Context, name string, namespaceName string, create func(*unstructured.Unstructured), update func(*unstructured.Unstructured), delete func(*unstructured.Unstructured)) error {
	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	if namespaceName == "" {
		return tracederrors.TracedErrorEmptyString("namespace")
	}

	dynamicClient, err := f.GetDynamicClient()
	if err != nil {
		return err
	}

	watcher, err := dynamicClient.Resource(GetGitRepositoryGVR()).Namespace(namespaceName).Watch(ctx, metav1.ListOptions{
		FieldSelector: "metadata.name=" + name,
	})
	if err != nil {
		return err
	}

	go func() {
		for event := range watcher.ResultChan() {
			obj, ok := event.Object.(*unstructured.Unstructured)
			if !ok {
				logging.LogWarnByCtxf(ctx, "Unexpected type for watch event object: %T", event.Object)
				continue
			}

			switch event.Type {
			case watch.Added:
				if create != nil {
					create(obj)
				}
			case watch.Modified:
				if update != nil {
					update(obj)
				}
			case watch.Deleted:
				if delete != nil {
					delete(obj)
				}
			case watch.Error:
				logging.LogErrorByCtxf(ctx, "Watcher for Flux GitRepository '%s' in namespace '%s' failed: %v", name, namespaceName, obj)
			}
		}
	}()

	logging.LogInfoByCtxf(ctx, "Watcher for Flux GitRepositoy '%s' in namespace '%s' registered.", name, namespaceName)

	return nil
}

func (f *FluxDeployment) WatchHelmRelease(ctx context.Context, name string, namespaceName string, create func(*unstructured.Unstructured), update func(*unstructured.Unstructured), delete func(*unstructured.Unstructured)) error {
	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	if namespaceName == "" {
		return tracederrors.TracedErrorEmptyString("namespace")
	}

	dynamicClient, err := f.GetDynamicClient()
	if err != nil {
		return err
	}

	watcher, err := dynamicClient.Resource(GetHelmReleaseGVR()).Namespace(namespaceName).Watch(ctx, metav1.ListOptions{
		FieldSelector: "metadata.name=" + name,
	})
	if err != nil {
		return err
	}

	go func() {
		for event := range watcher.ResultChan() {
			obj, ok := event.Object.(*unstructured.Unstructured)
			if !ok {
				logging.LogWarnByCtxf(ctx, "Unexpected type for watch event object: %T", event.Object)
				continue
			}

			switch event.Type {
			case watch.Added:
				if create != nil {
					create(obj)
				}
			case watch.Modified:
				if update != nil {
					update(obj)
				}
			case watch.Deleted:
				if delete != nil {
					delete(obj)
				}
			case watch.Error:
				logging.LogErrorByCtxf(ctx, "Watcher for Flux HelmRelease '%s' in namespace '%s' failed: %v", name, namespaceName, obj)
			}
		}
	}()

	logging.LogInfoByCtxf(ctx, "Watcher for Flux HelmRelease '%s' in namespace '%s' registered.", name, namespaceName)

	return nil
}

func (f *FluxDeployment) WatchKustomization(ctx context.Context, name string, namespaceName string, create func(*unstructured.Unstructured), update func(*unstructured.Unstructured), delete func(*unstructured.Unstructured)) error {
	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	if namespaceName == "" {
		return tracederrors.TracedErrorEmptyString("namespace")
	}

	dynamicClient, err := f.GetDynamicClient()
	if err != nil {
		return err
	}

	watcher, err := dynamicClient.Resource(GetKustomizationGVR()).Namespace(namespaceName).Watch(ctx, metav1.ListOptions{
		FieldSelector: "metadata.name=" + name,
	})
	if err != nil {
		return err
	}

	go func() {
		for event := range watcher.ResultChan() {
			obj, ok := event.Object.(*unstructured.Unstructured)
			if !ok {
				logging.LogWarnByCtxf(ctx, "Unexpected type for watch event object: %T", event.Object)
				continue
			}

			switch event.Type {
			case watch.Added:
				if create != nil {
					create(obj)
				}
			case watch.Modified:
				if update != nil {
					update(obj)
				}
			case watch.Deleted:
				if delete != nil {
					delete(obj)
				}
			case watch.Error:
				logging.LogErrorByCtxf(ctx, "Watcher for Flux Kustomization '%s' in namespace '%s' failed: %v", name, namespaceName, obj)
			}
		}
	}()

	logging.LogInfoByCtxf(ctx, "Watcher for Flux Kustomization '%s' in namespace '%s' registered.", name, namespaceName)

	return nil
}
