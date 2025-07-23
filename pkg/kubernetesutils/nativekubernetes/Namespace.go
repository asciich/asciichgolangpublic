package nativekubernetes

import (
	"context"
	"reflect"
	"time"

	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/contextutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/kubernetesutils/kubernetesimplementationindependend"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/kubernetesutils/kubernetesinterfaces"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/kubernetesutils/kubernetesparameteroptions"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/logging"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/tracederrors"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

type NativeNamespace struct {
	name              string
	kubernetesCluster *NativeKubernetesCluster
}

func (n *NativeNamespace) GetKubernetesCluster() (*NativeKubernetesCluster, error) {
	if n.kubernetesCluster == nil {
		return nil, tracederrors.TracedError("kubernetesCluster not set")
	}

	return n.kubernetesCluster, nil
}

func (n *NativeNamespace) GetClientSet() (*kubernetes.Clientset, error) {
	cluster, err := n.GetKubernetesCluster()
	if err != nil {
		return nil, err
	}

	return cluster.GetClientSet()
}

func (n *NativeNamespace) GetDynamicClient() (*dynamic.DynamicClient, error) {
	cluster, err := n.GetKubernetesCluster()
	if err != nil {
		return nil, err
	}

	return cluster.GetDynamicClient()
}

func (n *NativeNamespace) Create(ctx context.Context) (err error) {
	cluster, err := n.GetKubernetesCluster()
	if err != nil {
		return err
	}

	namespaceName, err := n.GetName()
	if err != nil {
		return err
	}

	_, err = cluster.CreateNamespaceByName(ctx, namespaceName)
	if err != nil {
		return err
	}

	return err
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
	return nil, tracederrors.TracedErrorNotImplemented()
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
		currentData, err := ReadSecret(ctx, clientset, namespaceName, secretName)
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

	var exists bool
	_, err = clientset.CoreV1().ConfigMaps(namespaceName).Get(ctx, configmapName, metav1.GetOptions{})
	if err == nil {
		exists = true
	} else {
		if !errors.IsNotFound(err) {
			return false, tracederrors.TracedErrorf("failed to get configmap '%s' in namespace '%s': %w", configmapName, namespaceName, err)
		}
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "ConfigMap '%s' in namespace '%s' exists.", configmapName, namespaceName)
	} else {
		logging.LogInfoByCtxf(ctx, "ConfigMap '%s' in namespace '%s' does not exist.", configmapName, namespaceName)
	}

	return exists, nil
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

	exists, err := n.ConfigMapByNameExists(ctx, configmapName)
	if err != nil {
		return err
	}

	if exists {
		clientset, err := n.GetClientSet()
		if err != nil {
			return err
		}

		err = clientset.CoreV1().ConfigMaps(namespaceName).Delete(ctx, configmapName, metav1.DeleteOptions{})
		if err != nil {
			return tracederrors.TracedErrorf("Failed to delete configmap '%s' in namespace '%s'.", configmapName, namespaceName)
		}

		logging.LogChangedByCtxf(ctx, "ConfigMap '%s' in namespace '%s' deleted.", configmapName, namespaceName)
	} else {
		logging.LogInfoByCtxf(ctx, "ConfigMap '%s' in namespace '%s' does not exist. Skip delete.", configmapName, namespaceName)
	}

	return nil
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
		select {
		case <-ctx.Done():
			if verbose {
				logging.LogInfof("Watch ConfigMap '%s' in namespace '%s' canceled.", configMapName, namespaceName)
			}
		}
	}()

	if !cache.WaitForCacheSync(ctx.Done(), informer.HasSynced) {
		return tracederrors.TracedErrorf("Failed to sync cache for watching ConfigMap '%s' in namespace '%s'.", configMapName, namespaceName)
	}

	logging.LogInfoByCtxf(ctx, "Watch ConfigMap '%s' in namespace '%s' set up. Create, update and delete are now watched.", configMapName, namespaceName)

	return nil
}

func (n *NativeNamespace) GetDiscoveryClient() (discovery.DiscoveryInterface, error) {
	cluster, err := n.GetKubernetesCluster()
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
