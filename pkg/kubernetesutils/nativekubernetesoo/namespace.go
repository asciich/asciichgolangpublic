package nativekubernetesoo

import (
	"context"
	"reflect"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesimplementationindependend"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

type NativeNamespace struct {
	name              string
	kubernetesCluster *NativeKubernetesCluster
}

func (n *NativeNamespace) GetNativeKubernetesCluster() (*NativeKubernetesCluster, error) {
	if n.kubernetesCluster == nil {
		return nil, tracederrors.TracedError("kubernetesCluster not set")
	}

	return n.kubernetesCluster, nil
}

func (n *NativeNamespace) GetKubernetesCluster() (kubernetesinterfaces.KubernetesCluster, error) {
	if n.kubernetesCluster == nil {
		return nil, tracederrors.TracedError("kubernetesCluster not set")
	}

	return n.kubernetesCluster, nil
}

func (n *NativeNamespace) GetConfig() (*rest.Config, error) {
	cluster, err := n.GetNativeKubernetesCluster()
	if err != nil {

	}

	return cluster.GetConfig()
}

func (n *NativeNamespace) DeletePodByName(ctx context.Context, podName string) error {
	pod, err := n.GetPodByName(podName)
	if err != nil {
		return err
	}

	return pod.Delete(ctx)
}

func (n *NativeNamespace) DeleteReplicaSetByName(ctx context.Context, replicaSetName string) error {
	replicaSet, err := n.GetReplicaSetByName(replicaSetName)
	if err != nil {
		return err
	}

	return replicaSet.Delete(ctx)
}

func (n *NativeNamespace) DeleteDeploymentByName(ctx context.Context, deploymentName string) error {
	deployment, err := n.GetDeploymentByName(deploymentName)
	if err != nil {
		return err
	}

	return deployment.Delete(ctx)
}

func (n *NativeNamespace) CreatePod(ctx context.Context, options *kubernetesparameteroptions.RunCommandOptions) (kubernetesinterfaces.Pod, error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	podName, err := options.GetPodName()
	if err != nil {
		return nil, err
	}

	clientSet, err := n.GetClientSet()
	if err != nil {
		return nil, err
	}

	namespaceName, err := n.GetName()
	if err != nil {
		return nil, err
	}

	err = nativekubernetes.CreatePod(ctx, clientSet, namespaceName, options)
	if err != nil {
		return nil, err
	}

	return n.GetPodByName(podName)
}

func (n *NativeNamespace) CreateReplicaSet(ctx context.Context, options *kubernetesparameteroptions.RunCommandOptions) (kubernetesinterfaces.ReplicaSet, error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	replicaSetName, err := options.GetReplicaSetName()
	if err != nil {
		return nil, err
	}

	clientSet, err := n.GetClientSet()
	if err != nil {
		return nil, err
	}

	namespaceName, err := n.GetName()
	if err != nil {
		return nil, err
	}

	err = nativekubernetes.CreateReplicaSet(ctx, clientSet, namespaceName, options)
	if err != nil {
		return nil, err
	}

	return n.GetReplicaSetByName(replicaSetName)
}

func (n *NativeNamespace) CreateDeployment(ctx context.Context, options *kubernetesparameteroptions.RunCommandOptions) (kubernetesinterfaces.Deployment, error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	deploymentName, err := options.GetDeploymentName()
	if err != nil {
		return nil, err
	}

	clientSet, err := n.GetClientSet()
	if err != nil {
		return nil, err
	}

	namespaceName, err := n.GetName()
	if err != nil {
		return nil, err
	}

	err = nativekubernetes.CreateDeployment(ctx, clientSet, namespaceName, options)
	if err != nil {
		return nil, err
	}

	return n.GetDeploymentByName(deploymentName)
}

func (n *NativeNamespace) GetPodByName(podName string) (kubernetesinterfaces.Pod, error) {
	if podName == "" {
		return nil, tracederrors.TracedErrorEmptyString("podName")
	}

	pod := &Pod{}

	err := pod.SetName(podName)
	if err != nil {
		return nil, err
	}

	err = pod.SetNamespace(n)
	if err != nil {
		return nil, err
	}

	return pod, nil
}

func (n *NativeNamespace) GetReplicaSetByName(replicaSetName string) (kubernetesinterfaces.ReplicaSet, error) {
	if replicaSetName == "" {
		return nil, tracederrors.TracedErrorEmptyString("replicaSetName")
	}

	replicaSet := &ReplicaSet{}

	err := replicaSet.SetName(replicaSetName)
	if err != nil {
		return nil, err
	}

	err = replicaSet.SetNamespace(n)
	if err != nil {
		return nil, err
	}

	return replicaSet, nil
}

func (n *NativeNamespace) GetDeploymentByName(deploymentName string) (kubernetesinterfaces.Deployment, error) {
	if deploymentName == "" {
		return nil, tracederrors.TracedErrorEmptyString("deploymentName")
	}

	deployment := &Deployment{}

	err := deployment.SetName(deploymentName)
	if err != nil {
		return nil, err
	}

	err = deployment.SetNamespace(n)
	if err != nil {
		return nil, err
	}

	return deployment, nil
}

func (n *NativeNamespace) GetClientSet() (*kubernetes.Clientset, error) {
	cluster, err := n.GetNativeKubernetesCluster()
	if err != nil {
		return nil, err
	}

	return cluster.GetClientSet()
}

func (n *NativeNamespace) GetDynamicClient() (*dynamic.DynamicClient, error) {
	cluster, err := n.GetNativeKubernetesCluster()
	if err != nil {
		return nil, err
	}

	return cluster.GetDynamicClient()
}

func (n *NativeNamespace) Create(ctx context.Context) (err error) {
	namespaceName, err := n.GetName()
	if err != nil {
		return err
	}

	clientSet, err := n.GetClientSet()
	if err != nil {
		return err
	}

	return nativekubernetes.CreateNamespace(ctx, clientSet, namespaceName)
}

func (n *NativeNamespace) CreateRole(ctx context.Context, createOptions *kubernetesparameteroptions.CreateRoleOptions) (createdRole kubernetesinterfaces.Role, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (n *NativeNamespace) DeleteRoleByName(ctx context.Context, name string) (err error) {
	return tracederrors.TracedErrorNotImplemented()
}

func (n *NativeNamespace) GetClusterName() (clusterName string, err error) {
	return "", tracederrors.TracedErrorNotImplemented()
}

func (n *NativeNamespace) GetKubectlContext(ctx context.Context) (contextName string, err error) {
	return "", tracederrors.TracedErrorNotImplemented()
}

func (n *NativeNamespace) GetName() (name string, err error) {
	if n.name == "" {
		return "", tracederrors.TracedError("name not set")
	}

	return n.name, nil
}

func (n *NativeNamespace) GetObjectByNames(objectName string, objectKind string) (object kubernetesinterfaces.Object, err error) {
	if objectName == "" {
		return nil, tracederrors.TracedErrorEmptyString("objectName")
	}

	if objectKind == "" {
		return nil, tracederrors.TracedErrorEmptyString("objectType")
	}

	return &NativeObject{
		name:      objectName,
		kind:      objectKind,
		namespace: n,
	}, nil
}

func (n *NativeNamespace) GetRoleByName(name string) (role kubernetesinterfaces.Role, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (n *NativeNamespace) ListRoleNames(ctx context.Context) (roleNames []string, err error) {
	clientSet, err := n.GetClientSet()
	if err != nil {
		return nil, err
	}

	namespaceName, err := n.GetName()
	if err != nil {
		return nil, err
	}

	return nativekubernetes.ListRoleNames(ctx, clientSet, namespaceName)
}

func (n *NativeNamespace) RoleByNameExists(ctx context.Context, name string) (exists bool, err error) {
	return false, tracederrors.TracedErrorNotImplemented()
}

func (n *NativeNamespace) SecretByNameExists(ctx context.Context, secretName string) (bool, error) {
	if secretName == "" {
		return false, tracederrors.TracedErrorEmptyString("name")
	}

	clientset, err := n.GetClientSet()
	if err != nil {
		return false, err
	}

	namespaceName, err := n.GetName()
	if err != nil {
		return false, err
	}

	var exists bool
	_, err = clientset.CoreV1().Secrets(namespaceName).Get(ctx, secretName, metav1.GetOptions{})
	if err == nil {
		exists = true
	} else {
		if !errors.IsNotFound(err) {
			return false, tracederrors.TracedErrorf("failed to get secret '%s' in namespace '%s': %w", secretName, namespaceName, err)
		}
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Secret '%s' in namespace '%s' exists.", secretName, namespaceName)
	} else {
		logging.LogInfoByCtxf(ctx, "Secret '%s' in namespace '%s' does not exist.", secretName, namespaceName)
	}

	return exists, nil
}

// CheckSecretByNameExists checks if a secret exists by name.
// Returns nil if it exists, error if it does not exist.
func (n *NativeNamespace) CheckSecretByNameExists(ctx context.Context, secretName string) error {
	exists, err := n.SecretByNameExists(ctx, secretName)
	if err != nil {
		return err
	}
	if !exists {
		namespaceName, _ := n.GetName()
		return tracederrors.TracedErrorf("Secret '%s' does not exist in namespace '%s'", secretName, namespaceName)
	}
	return nil
}

// CheckNamespaceByNameExists checks if a namespace exists by name.
// Returns nil if it exists, error if it does not exist.
func (n *NativeNamespace) CheckNamespaceByNameExists(ctx context.Context) error {
	exists, err := n.Exists(ctx)
	if err != nil {
		return err
	}
	if !exists {
		namespaceName, _ := n.GetName()
		return tracederrors.TracedErrorf("Namespace '%s' does not exist", namespaceName)
	}
	return nil
}

// CheckPodByNameExists checks if a pod exists by name.
// Returns nil if it exists, error if it does not exist.
func (n *NativeNamespace) CheckPodByNameExists(ctx context.Context, podName string) error {
	exists, err := n.PodByNameExists(ctx, podName)
	if err != nil {
		return err
	}
	if !exists {
		namespaceName, _ := n.GetName()
		return tracederrors.TracedErrorf("Pod '%s' does not exist in namespace '%s'", podName, namespaceName)
	}
	return nil
}

// CheckReplicaSetByNameExists checks if a replicaSet exists by name.
// Returns nil if it exists, error if it does not exist.
func (n *NativeNamespace) CheckReplicaSetByNameExists(ctx context.Context, replicaSetName string) error {
	exists, err := n.ReplicaSetByNameExists(ctx, replicaSetName)
	if err != nil {
		return err
	}
	if !exists {
		namespaceName, _ := n.GetName()
		return tracederrors.TracedErrorf("ReplicaSet '%s' does not exist in namespace '%s'", replicaSetName, namespaceName)
	}
	return nil
}

// CheckDeploymentByNameExists checks if a deployment exists by name.
// Returns nil if it exists, error if it does not exist.
func (n *NativeNamespace) CheckDeploymentByNameExists(ctx context.Context, deploymentName string) error {
	exists, err := n.DeploymentByNameExists(ctx, deploymentName)
	if err != nil {
		return err
	}
	if !exists {
		namespaceName, _ := n.GetName()
		return tracederrors.TracedErrorf("Deployment '%s' does not exist in namespace '%s'", deploymentName, namespaceName)
	}
	return nil
}

// CheckCronJobByNameExists checks if a cronJob exists by name.
// Returns nil if it exists, error if it does not exist.
func (n *NativeNamespace) CheckCronJobByNameExists(ctx context.Context, cronJobName string) error {
	exists, err := n.CronJobByNameExists(ctx, cronJobName)
	if err != nil {
		return err
	}
	if !exists {
		namespaceName, _ := n.GetName()
		return tracederrors.TracedErrorf("CronJob '%s' does not exist in namespace '%s'", cronJobName, namespaceName)
	}
	return nil
}

func (n *NativeNamespace) DeleteSecretByName(ctx context.Context, secretName string) (err error) {
	if secretName == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	namespaceName, err := n.GetName()
	if err != nil {
		return err
	}

	exists, err := n.SecretByNameExists(ctx, secretName)
	if err != nil {
		return err
	}

	if exists {
		clientset, err := n.GetClientSet()
		if err != nil {
			return err
		}

		err = clientset.CoreV1().Secrets(namespaceName).Delete(ctx, secretName, metav1.DeleteOptions{})
		if err != nil {
			return tracederrors.TracedErrorf("Failed to delete secret '%s' in namespace '%s'.", secretName, namespaceName)
		}

		logging.LogChangedByCtxf(ctx, "Secret '%s' in namespace '%s' deleted.", secretName, namespaceName)
	} else {
		logging.LogInfoByCtxf(ctx, "Secret '%s' in namespace '%s' does not exist. Skip delete.", secretName, namespaceName)
	}

	return nil
}

func (n *NativeNamespace) GetSecretByName(name string) (secret kubernetesinterfaces.Secret, err error) {
	if name == "" {
		return nil, tracederrors.TracedErrorEmptyString("name")
	}

	return &NativeSecret{
		namespace: n,
		name:      name,
	}, nil
}

func (n *NativeNamespace) CreateSecret(ctx context.Context, secretName string, options *kubernetesparameteroptions.CreateSecretOptions) (createdSecret kubernetesinterfaces.Secret, err error) {
	if secretName == "" {
		return nil, tracederrors.TracedErrorEmptyString("secret")
	}

	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	exists, err := n.SecretByNameExists(ctx, secretName)
	if err != nil {
		return nil, err
	}

	namespaceName, err := n.GetName()
	if err != nil {
		return nil, err
	}

	clientset, err := n.GetClientSet()
	if err != nil {
		return nil, err
	}

	if exists {
		currentData, err := nativekubernetes.ReadSecret(ctx, clientset, namespaceName, secretName)
		if err != nil {
			return nil, err
		}

		if reflect.DeepEqual(currentData, options.SecretData) {
			logging.LogInfoByCtxf(ctx, "Secret '%s' in namespace '%s' is already up to date. Skip creation and update.", secretName, namespaceName)
		} else {
			secret, err := clientset.CoreV1().Secrets(namespaceName).Get(ctx, secretName, metav1.GetOptions{})
			if err != nil {
				return nil, tracederrors.TracedErrorf("Failed to read secret '%s' in namespace '%s' to update it: %w", secretName, namespaceName, err)
			}

			secret.Data = options.SecretData

			_, err = clientset.CoreV1().Secrets(namespaceName).Update(ctx, secret, metav1.UpdateOptions{})
			if err != nil {
				return nil, tracederrors.TracedErrorf("Failed to read update '%s' in namespace '%s': %w", secretName, namespaceName, err)
			}

			logging.LogChangedByCtxf(ctx, "Secret '%s' in namespace '%s' updated.", secretName, namespaceName)
		}
	} else {
		secretData, err := options.GetSecretData()
		if err != nil {
			return nil, err
		}

		secret := &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:   secretName,
				Labels: map[string]string{},
			},
			Data: secretData,
			Type: v1.SecretTypeOpaque,
		}

		_, err = clientset.CoreV1().Secrets(namespaceName).Create(ctx, secret, metav1.CreateOptions{})
		if err != nil {
			return nil, tracederrors.TracedErrorf("failed to create secret '%s' in namespace '%s': %w", secretName, namespaceName, err)
		}

		logging.LogChangedByCtxf(ctx, "Created secret '%s' in kubernetes namespace '%s'.", secretName, namespaceName)
	}

	return n.GetSecretByName(secretName)
}

func (n *NativeNamespace) ConfigMapByNameExists(ctx context.Context, configmapName string) (bool, error) {
	if configmapName == "" {
		return false, tracederrors.TracedErrorEmptyString("name")
	}

	clientset, err := n.GetClientSet()
	if err != nil {
		return false, err
	}

	namespaceName, err := n.GetName()
	if err != nil {
		return false, err
	}

	return nativekubernetes.ConfigMapExists(ctx, clientset, namespaceName, configmapName)
}

func (n *NativeNamespace) CreateConfigMap(ctx context.Context, configMapName string, options *kubernetesparameteroptions.CreateConfigMapOptions) (createdConfigMap kubernetesinterfaces.ConfigMap, err error) {
	if configMapName == "" {
		return nil, tracederrors.TracedErrorEmptyString("configmap")
	}

	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	exists, err := n.ConfigMapByNameExists(ctx, configMapName)
	if err != nil {
		return nil, err
	}

	namespaceName, err := n.GetName()
	if err != nil {
		return nil, err
	}

	configmapData, err := options.GetConfigMapData()
	if err != nil {
		return nil, err
	}

	labels := options.GetLabels()

	clientset, err := n.GetClientSet()
	if err != nil {
		return nil, err
	}

	configmap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:   configMapName,
			Labels: labels,
		},
		Data: configmapData,
	}

	if exists {
		configMap, err := n.GetConfigMapByName(configMapName)
		if err != nil {
			return nil, err
		}

		nativeConfigMap, ok := configMap.(*NativeConfigMap)
		if !ok {
			return nil, tracederrors.TracedError("Returned config map is not a nativeConfigMap")
		}

		rawResponse, err := nativeConfigMap.GetRawResponse(ctx)
		if err != nil {
			return nil, err
		}

		if IsConfigMapContentEqual(rawResponse.Data, configmapData) && IsConfigMapLabelsEqual(rawResponse.Labels, labels) {
			logging.LogInfoByCtxf(ctx, "ConfigMap '%s' already exists in namespace '%s' and is up to date.", configMapName, namespaceName)
		} else {
			_, err := clientset.CoreV1().ConfigMaps(namespaceName).Update(ctx, configmap, metav1.UpdateOptions{})
			if err != nil {
				return nil, tracederrors.TracedErrorf("failed to create ConfigMap '%s' in namespace '%s': %w", configMapName, namespaceName, err)
			}

			logging.LogChangedByCtxf(ctx, "Updated ConfigMap '%s' in kubernetes namespace '%s'.", configMapName, namespaceName)
		}
	} else {
		_, err = clientset.CoreV1().ConfigMaps(namespaceName).Create(ctx, configmap, metav1.CreateOptions{})
		if err != nil {
			return nil, tracederrors.TracedErrorf("failed to create configmap '%s' in namespace '%s': %w", configMapName, namespaceName, err)
		}

		logging.LogChangedByCtxf(ctx, "Created ConfigMap '%s' in kubernetes namespace '%s'.", configMapName, namespaceName)
	}

	return n.GetConfigMapByName(configMapName)
}

func (n *NativeNamespace) GetConfigMapByName(name string) (configMap kubernetesinterfaces.ConfigMap, err error) {
	if name == "" {
		return nil, tracederrors.TracedErrorEmptyString("name")
	}

	return &NativeConfigMap{
		namespace: n,
		name:      name,
	}, nil
}

func (n *NativeNamespace) DeleteConfigMapByName(ctx context.Context, configmapName string) (err error) {
	if configmapName == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	namespaceName, err := n.GetName()
	if err != nil {
		return err
	}

	clientset, err := n.GetClientSet()
	if err != nil {
		return err
	}

	return nativekubernetes.DeleteConfigMap(ctx, clientset, namespaceName, configmapName)
}

func (n *NativeNamespace) WatchConfigMap(ctx context.Context, configMapName string, onCreate func(kubernetesinterfaces.ConfigMap), onUpdate func(kubernetesinterfaces.ConfigMap), onDelete func(kubernetesinterfaces.ConfigMap)) error {
	if configMapName == "" {
		return tracederrors.TracedErrorEmptyString("configMapName")
	}

	namespaceName, err := n.GetName()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Watch ConfigMap '%s' in namespace '%s' started.", configMapName, namespaceName)

	clientset, err := n.GetClientSet()
	if err != nil {
		return err
	}

	fieldSelector := fields.OneTermEqualSelector("metadata.name", configMapName)

	listWatcher := cache.NewListWatchFromClient(
		clientset.CoreV1().RESTClient(),
		"configmaps",
		v1.NamespaceAll,
		fieldSelector,
	)

	informer := cache.NewSharedIndexInformer(
		listWatcher,
		&v1.ConfigMap{},
		5*time.Minute,
		cache.Indexers{},
	)

	_, err = informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			nativeConfigMap, ok := obj.(*v1.ConfigMap)
			if !ok {
				tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
				if !ok {
					return
				}
				nativeConfigMap, ok = tombstone.Obj.(*v1.ConfigMap)
				if !ok {
					return
				}
			}
			cm, err := n.GetConfigMapByName(nativeConfigMap.Name)
			if err != nil {
				return
			}
			onCreate(cm)
		},
		UpdateFunc: func(oldObj interface{}, newObj interface{}) {
			nativeConfigMap, ok := newObj.(*v1.ConfigMap)
			if !ok {
				tombstone, ok := newObj.(cache.DeletedFinalStateUnknown)
				if !ok {
					return
				}
				nativeConfigMap, ok = tombstone.Obj.(*v1.ConfigMap)
				if !ok {
					return
				}
			}
			cm, err := n.GetConfigMapByName(nativeConfigMap.Name)
			if err != nil {
				return
			}
			onUpdate(cm)
		},
		DeleteFunc: func(obj interface{}) {
			nativeConfigMap, ok := obj.(*v1.ConfigMap)
			if !ok {
				tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
				if !ok {
					return
				}
				nativeConfigMap, ok = tombstone.Obj.(*v1.ConfigMap)
				if !ok {
					return
				}
			}
			cm, err := n.GetConfigMapByName(nativeConfigMap.Name)
			if err != nil {
				return
			}
			onDelete(cm)
		},
	})
	if err != nil {
		return err
	}

	go informer.Run(ctx.Done())

	go func() {
		verbose := contextutils.GetVerboseFromContext(ctx)
		<-ctx.Done()
		if verbose {
			logging.LogInfof("Watch ConfigMap '%s' in namespace '%s' canceled.", configMapName, namespaceName)
		}
	}()

	if !cache.WaitForCacheSync(ctx.Done(), informer.HasSynced) {
		return tracederrors.TracedErrorf("Failed to sync cache for watching ConfigMap '%s' in namespace '%s'.", configMapName, namespaceName)
	}

	logging.LogInfoByCtxf(ctx, "Watch ConfigMap '%s' in namespace '%s' set up. Create, update and delete are now watched.", configMapName, namespaceName)

	return nil
}

func (n *NativeNamespace) GetDiscoveryClient() (discovery.DiscoveryInterface, error) {
	cluster, err := n.GetNativeKubernetesCluster()
	if err != nil {
		return nil, err
	}

	return cluster.GetDiscoveryClient()
}

func (n *NativeNamespace) WaitUntilAllPodsInNamespaceAreRunning(ctx context.Context, options *kubernetesparameteroptions.WaitForPodsOptions) error {
	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	namspaceName, err := n.GetName()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Wait until all pods in namespace '%s' are running started.", namspaceName)

	clientset, err := n.GetClientSet()
	if err != nil {
		return err
	}

	var nPods int
	for {
		err := ctx.Err()
		if err != nil {
			return err
		}

		pods, err := clientset.CoreV1().Pods(namspaceName).List(ctx, metav1.ListOptions{})
		if err != nil {
			return tracederrors.TracedErrorf("Failed to list pods to wait for: %w", err)
		}

		allRunning := true
		nPods = 0
		for _, pod := range pods.Items {
			nPods++
			if pod.Status.Phase != v1.PodRunning {
				allRunning = false
				logging.LogInfoByCtxf(ctx, "Pod %s is in phase '%s' and not 'running' yet.", pod.Name, pod.Status.Phase)
				break
			}
		}

		if options.MinNumberOfPods > 0 {
			minPods := options.MinNumberOfPods
			if nPods < minPods {
				allRunning = false
				logging.LogInfoByCtxf(ctx, "Only %d pods present in namespace '%s'. Waiting until minimum required pods of %d are present.", nPods, namspaceName, minPods)
			}
		}

		if allRunning {
			break
		}

		delay := time.Second * 3
		logging.LogInfoByCtxf(ctx, "Wait '%s' before checking again if all pods in namespace '%s are running.'", delay, namspaceName)
		time.Sleep(delay)
	}

	logging.LogInfoByCtxf(ctx, "Wait until all pods in namespace '%s' are running finished. There are now '%d' pods running.", namspaceName, nPods)

	return nil
}

func (n *NativeNamespace) GetObjectByYamlString(yaml string) (kubernetesinterfaces.Object, error) {
	if yaml == "" {
		return nil, tracederrors.TracedErrorEmptyString("yaml")
	}

	objectYamls, err := kubernetesimplementationindependend.UnmarshalObjectYaml(yaml)
	if err != nil {
		return nil, err
	}

	nObjects := len(objectYamls)
	if nObjects != 1 {
		return nil, tracederrors.TracedErrorf("Exepected one yaml document to get resouce by yaml string but got '%d'.", nObjects)
	}

	ret, err := n.GetObjectByNames(objectYamls[0].Name(), objectYamls[0].Kind())
	if err != nil {
		return nil, err
	}

	err = ret.SetApiVersion(objectYamls[0].ApiVersion())
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (n *NativeNamespace) Exists(ctx context.Context) (bool, error) {
	namespaceName, err := n.GetName()
	if err != nil {
		return false, err
	}

	cluster, err := n.GetKubernetesCluster()
	if err != nil {
		return false, err
	}

	return cluster.NamespaceByNameExists(ctx, namespaceName)
}

func (n *NativeNamespace) CreateObject(ctx context.Context, options *kubernetesparameteroptions.CreateObjectOptions) (kubernetesinterfaces.Object, error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	object, err := n.GetObjectByYamlString(options.YamlString)
	if err != nil {
		return nil, err
	}

	err = object.CreateByYamlString(ctx, options)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func (n *NativeNamespace) ListSecretNames(ctx context.Context) ([]string, error) {
	clientSet, err := n.GetClientSet()
	if err != nil {
		return nil, err
	}

	namespaceName, err := n.GetName()
	if err != nil {
		return nil, err
	}

	return nativekubernetes.ListSecretNames(ctx, clientSet, namespaceName)
}

func (n *NativeNamespace) ListSecrets(ctx context.Context) ([]kubernetesinterfaces.Secret, error) {
	clientSet, err := n.GetClientSet()
	if err != nil {
		return nil, err
	}

	namespaceName, err := n.GetName()
	if err != nil {
		return nil, err
	}

	secretNames, err := nativekubernetes.ListSecretNames(ctx, clientSet, namespaceName)
	if err != nil {
		return nil, err
	}

	secrets := []kubernetesinterfaces.Secret{}
	for _, n := range secretNames {
		toAdd := &NativeSecret{}

		err := toAdd.SetName(n)
		if err != nil {
			return nil, err
		}

		secrets = append(secrets, toAdd)
	}

	return secrets, nil
}

func (n *NativeNamespace) ListPodNames(ctx context.Context) ([]string, error) {
	clientSet, err := n.GetClientSet()
	if err != nil {
		return nil, err
	}

	namespaceName, err := n.GetName()
	if err != nil {
		return nil, err
	}

	return nativekubernetes.ListPodNames(ctx, clientSet, namespaceName)
}

func (n *NativeNamespace) ListDeploymentNames(ctx context.Context) ([]string, error) {
	clientSet, err := n.GetClientSet()
	if err != nil {
		return nil, err
	}

	namespaceName, err := n.GetName()
	if err != nil {
		return nil, err
	}

	return nativekubernetes.ListDeploymentNames(ctx, clientSet, namespaceName)
}

func (n *NativeNamespace) ListReplicaSetNames(ctx context.Context) ([]string, error) {
	clientSet, err := n.GetClientSet()
	if err != nil {
		return nil, err
	}

	namespaceName, err := n.GetName()
	if err != nil {
		return nil, err
	}

	return nativekubernetes.ListReplicaSetNames(ctx, clientSet, namespaceName)
}

func (n *NativeNamespace) ListConfigMapNames(ctx context.Context) ([]string, error) {
	clientSet, err := n.GetClientSet()
	if err != nil {
		return nil, err
	}

	namespaceName, err := n.GetName()
	if err != nil {
		return nil, err
	}

	return nativekubernetes.ListConfigMapNames(ctx, clientSet, namespaceName)
}

func (n *NativeNamespace) ListObjectNames(options *kubernetesparameteroptions.ListKubernetesObjectsOptions) ([]string, error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	objectType, err := options.GetObjectType()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	if options.Verbose {
		ctx = contextutils.WithVerbose(ctx)
	}

	clientSet, err := n.GetClientSet()
	if err != nil {
		return nil, err
	}

	namespaceName, err := n.GetName()
	if err != nil {
		return nil, err
	}

	switch objectType {
	case "pod", "pods":
		return nativekubernetes.ListPodNames(ctx, clientSet, namespaceName)
	case "deployment", "deployments":
		return nativekubernetes.ListDeploymentNames(ctx, clientSet, namespaceName)
	case "replicaset", "replicasets":
		return nativekubernetes.ListReplicaSetNames(ctx, clientSet, namespaceName)
	case "configmap", "configmaps":
		return nativekubernetes.ListConfigMapNames(ctx, clientSet, namespaceName)
	case "secret", "secrets":
		return nativekubernetes.ListSecretNames(ctx, clientSet, namespaceName)
	default:
		return nil, tracederrors.TracedErrorf("Unsupported object type '%s'", objectType)
	}
}

func (n *NativeNamespace) PodByNameExists(ctx context.Context, podName string) (bool, error) {
	pod, err := n.GetPodByName(podName)
	if err != nil {
		return false, err
	}

	return pod.Exists(ctx)
}

func (n *NativeNamespace) ReplicaSetByNameExists(ctx context.Context, replicaSetName string) (bool, error) {
	replicaSet, err := n.GetReplicaSetByName(replicaSetName)
	if err != nil {
		return false, err
	}

	return replicaSet.Exists(ctx)
}

func (n *NativeNamespace) DeploymentByNameExists(ctx context.Context, deploymentName string) (bool, error) {
	deployment, err := n.GetDeploymentByName(deploymentName)
	if err != nil {
		return false, err
	}

	return deployment.Exists(ctx)
}

func (n *NativeNamespace) CronJobByNameExists(ctx context.Context, cronJobName string) (bool, error) {
	if cronJobName == "" {
		return false, tracederrors.TracedErrorEmptyString("cronJobName")
	}

	namespaceName, err := n.GetName()
	if err != nil {
		return false, err
	}

	clientset, err := n.GetClientSet()
	if err != nil {
		return false, err
	}

	return nativekubernetes.CronJobExists(ctx, clientset, namespaceName, cronJobName)
}

func (n *NativeNamespace) CreateCronJob(ctx context.Context, cronJobName string, schedule string, image string, command []string, labels map[string]string) (createdCronJob kubernetesinterfaces.CronJob, err error) {
	if cronJobName == "" {
		return nil, tracederrors.TracedErrorEmptyString("cronJobName")
	}

	namespaceName, err := n.GetName()
	if err != nil {
		return nil, err
	}

	clientset, err := n.GetClientSet()
	if err != nil {
		return nil, err
	}

	err = nativekubernetes.CreateCronJob(ctx, clientset, namespaceName, cronJobName, schedule, image, command, labels)
	if err != nil {
		return nil, err
	}

	return n.GetCronJobByName(cronJobName)
}

func (n *NativeNamespace) GetCronJobByName(name string) (cronJob kubernetesinterfaces.CronJob, err error) {
	if name == "" {
		return nil, tracederrors.TracedErrorEmptyString("name")
	}

	return &NativeCronJob{
		name:      name,
		namespace: n,
	}, nil
}

func (n *NativeNamespace) DeleteCronJobByName(ctx context.Context, cronJobName string) (err error) {
	if cronJobName == "" {
		return tracederrors.TracedErrorEmptyString("cronJobName")
	}

	namespaceName, err := n.GetName()
	if err != nil {
		return err
	}

	clientset, err := n.GetClientSet()
	if err != nil {
		return err
	}

	return nativekubernetes.DeleteCronJob(ctx, clientset, namespaceName, cronJobName)
}

func (n *NativeNamespace) ListCronJobNames(ctx context.Context) ([]string, error) {
	namespaceName, err := n.GetName()
	if err != nil {
		return nil, err
	}

	clientSet, err := n.GetClientSet()
	if err != nil {
		return nil, err
	}

	return nativekubernetes.ListCronJobNames(ctx, clientSet, namespaceName)
}
