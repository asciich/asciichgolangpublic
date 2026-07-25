package nativekubernetesoo

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesimplementationindependend"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	sigyaml "sigs.k8s.io/yaml"
)

type NativeObject struct {
	name       string
	kind       string
	apiVersion string
	namespace  *NativeNamespace
}

func (n *NativeObject) GetNativeNamespace() (*NativeNamespace, error) {
	return n.namespace, nil
}

func (n *NativeObject) GetApiVersion(ctx context.Context) (string, error) {
	const defaultVersion = "v1"

	if n.apiVersion == "" {
		if n.kind == "FluxInstance" {
			const fluxControlPlaneIoV1 = "fluxcd.controlplane.io/v1"
			logging.LogInfoByCtxf(ctx, "ApiVersion not set, use '%s' as default API version for kind='%s'.", defaultVersion, n.kind)

			return fluxControlPlaneIoV1, nil
		}

		logging.LogInfoByCtxf(ctx, "ApiVersion not set, use '%s' as default API version.", defaultVersion)
		return defaultVersion, nil
	}

	return n.apiVersion, nil
}

func (n *NativeObject) GetName() (string, error) {
	if n.name == "" {
		return "", tracederrors.TracedError("name not set")
	}

	return n.name, nil
}

func (n *NativeObject) GetKind() (string, error) {
	if n.kind == "" {
		return "", tracederrors.TracedError("kind not set")
	}

	ret, err := kubernetesimplementationindependend.SanitizeKindName(n.kind)
	if err != nil {
		return "", err
	}

	return ret, nil
}

func (n *NativeObject) GetGroupVersionKind(ctx context.Context) (*schema.GroupVersionKind, error) {
	groupVersion, err := n.GetGroupVersion(ctx)
	if err != nil {
		return nil, err
	}

	kind, err := n.GetKind()
	if err != nil {
		return nil, err
	}

	gvk := schema.GroupVersionKind{
		Group:   groupVersion.Group,
		Version: groupVersion.Version,
		Kind:    kind,
	}

	return &gvk, nil
}

func (n *NativeObject) EnsureNamespaceExists(ctx context.Context) error {
	namespace, err := n.GetNamespace()
	if err != nil {
		return err
	}

	return namespace.Create(ctx)
}

func (n *NativeObject) CreateByYamlString(ctx context.Context, options *kubernetesparameteroptions.CreateObjectOptions) (err error) {
	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	yamlString, err := options.GetYamlString()
	if err != nil {
		return err
	}

	namespaceName, err := n.GetNamespaceName()
	if err != nil {
		return err
	}

	kind, err := n.GetKind()
	if err != nil {
		return err
	}

	name, err := n.GetName()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Create kubernetes object by yaml '%s/%s' in namespace '%s' started.", kind, name, namespaceName)

	exists, err := n.Exists(ctx)
	if err != nil {
		return err
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Object '%s' named '%s' in namespace '%s' already exists, skip creation.", kind, name, namespaceName)
	} else {
		err = n.EnsureNamespaceExists(ctx)
		if err != nil {
			return err
		}

		// Parse the YAML string into an unstructured object to extract the GVK from the YAML itself
		var unstructuredObj unstructured.Unstructured
		if err := sigyaml.Unmarshal([]byte(yamlString), &unstructuredObj); err != nil {
			return tracederrors.TracedErrorf("Failed to parse yamlString as unstructuredObj: %w", err)
		}

		// Extract the GroupVersionKind from the YAML content (e.g. apps/v1 Deployment)
		gvk := unstructuredObj.GroupVersionKind()
		if gvk.Kind == "" || gvk.Version == "" {
			return tracederrors.TracedErrorf(
				"YAML string does not contain a valid apiVersion/kind for object '%s' named '%s' in namespace '%s'",
				kind, name, namespaceName,
			)
		}

		unstructuredObj.SetNamespace(namespaceName)

		// Use the GVK from the YAML to resolve the correct dynamic resource interface
		objectInterface, err := n.GetObjectInterfaceByGroupVersionKind(ctx, &gvk)
		if err != nil {
			return err
		}

		_, err = objectInterface.Create(ctx, &unstructuredObj, v1.CreateOptions{})
		if err != nil {
			return tracederrors.TracedErrorf("Failed to create object '%s' named '%s' in namespace '%s': %w", kind, name, namespaceName, err)
		}

		logging.LogChangedByCtxf(ctx, "Created object '%s' named '%s' in namespace '%s'.", kind, name, namespaceName)
	}

	logging.LogInfoByCtxf(ctx, "Create kubernetes object by yaml '%s/%s' in namespace '%s' finished.", kind, name, namespaceName)

	return nil
}

func (n *NativeObject) GetRestConfig() (*rest.Config, error) {
	namespace, err := n.GetNativeNamespace()
	if err != nil {

	}

	return namespace.GetConfig()
}

func (n *NativeObject) GetDiscoveryClient(ctx context.Context) (discovery.DiscoveryInterface, error) {
	restConfig, err := n.GetRestConfig()
	if err != nil {
		return nil, err
	}

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(restConfig)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to create discovery client: %w", err)
	}

	return discoveryClient, nil
}

func (n *NativeObject) GetObjectInterfaceByGroupVersionKind(
	ctx context.Context,
	gvk *schema.GroupVersionKind,
) (dynamic.ResourceInterface, error) {
	dynamicClient, err := n.GetDynamicClient()
	if err != nil {
		return nil, err
	}

	discoveryClient, err := n.GetDiscoveryClient(ctx)
	if err != nil {
		return nil, err
	}

	// Map the GVK to a GVR (GroupVersionResource) using the discovery API
	apiResourceList, err := discoveryClient.ServerResourcesForGroupVersion(gvk.GroupVersion().String())
	if err != nil {
		return nil, tracederrors.TracedErrorf(
			"Failed to discover resources for group version '%s': %w",
			gvk.GroupVersion().String(), err,
		)
	}

	var resource string
	for _, apiResource := range apiResourceList.APIResources {
		if apiResource.Kind == gvk.Kind {
			resource = apiResource.Name
			break
		}
	}

	if resource == "" {
		return nil, tracederrors.TracedErrorf(
			"Could not find resource name for kind '%s' in group version '%s'",
			gvk.Kind, gvk.GroupVersion().String(),
		)
	}

	gvr := schema.GroupVersionResource{
		Group:    gvk.Group,
		Version:  gvk.Version,
		Resource: resource,
	}

	namespaceName, err := n.GetNamespaceName()
	if err != nil {
		return nil, err
	}

	return dynamicClient.Resource(gvr).Namespace(namespaceName), nil
}

func (n *NativeObject) GetNamespace() (*NativeNamespace, error) {
	if n.namespace == nil {
		return nil, tracederrors.TracedErrorNil("namespace")
	}

	return n.namespace, nil
}

func (n *NativeObject) GetNamespaceName() (string, error) {
	namespace, err := n.GetNamespace()
	if err != nil {
		return "", err
	}

	return namespace.GetName()
}

func (n *NativeObject) GetDynamicClient() (*dynamic.DynamicClient, error) {
	namespace, err := n.GetNamespace()
	if err != nil {
		return nil, err
	}

	return namespace.GetDynamicClient()
}

func (n *NativeObject) GetGroupVersion(ctx context.Context) (*schema.GroupVersion, error) {
	apiVersion, err := n.GetApiVersion(ctx)
	if err != nil {
		return nil, err
	}

	groupVersion, err := schema.ParseGroupVersion(apiVersion)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to parse group version '%s': %w", apiVersion, err)
	}

	return &groupVersion, nil
}

func (n *NativeObject) GetObjectPlural() (string, error) {
	kind, err := n.GetKind()
	if err != nil {
		return "", err
	}

	return kubernetesimplementationindependend.GetObjectPlural(kind)
}

func (n *NativeObject) GetGroupVersionObject(ctx context.Context) (*schema.GroupVersionResource, error) {
	groupVersion, err := n.GetGroupVersion(ctx)
	if err != nil {
		return nil, err
	}

	objectPlural, err := n.GetObjectPlural()
	if err != nil {
		return nil, err
	}

	gvr := schema.GroupVersionResource{
		Group:    groupVersion.Group,
		Version:  groupVersion.Version,
		Resource: objectPlural,
	}

	return &gvr, nil
}

func (n *NativeObject) GetObjectInterface(ctx context.Context) (dynamic.ResourceInterface, error) {
	groupVersionObject, err := n.GetGroupVersionObject(ctx)
	if err != nil {
		return nil, err
	}

	namspaceName, err := n.GetNamespaceName()
	if err != nil {
		return nil, err
	}

	dynamicClient, err := n.GetDynamicClient()
	if err != nil {
		return nil, err
	}

	objectInterface := dynamicClient.Resource(*groupVersionObject).Namespace(namspaceName)

	return objectInterface, nil
}

func (n *NativeObject) Exists(ctx context.Context) (bool, error) {
	objectInterface, err := n.GetObjectInterface(ctx)
	if err != nil {
		return false, err
	}

	name, err := n.GetName()
	if err != nil {
		return false, err
	}

	kind, err := n.GetKind()
	if err != nil {
		return false, err
	}

	namespaceName, err := n.GetNamespaceName()
	if err != nil {
		return false, err
	}

	var exists bool
	_, err = objectInterface.Get(ctx, name, v1.GetOptions{})
	if err == nil {
		exists = true
	} else {
		if !errors.IsNotFound(err) {
			return false, tracederrors.TracedErrorf("failed to get object '%s' named '%s' in namespace '%s': %w", kind, name, namespaceName, err)
		}
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Object '%s' named '%s' in namespace '%s' exists.", kind, name, namespaceName)
	} else {
		logging.LogInfoByCtxf(ctx, "Object '%s' named '%s' in namespace '%s' does not exist.", kind, name, namespaceName)
	}

	return exists, nil
}

func (n *NativeObject) Delete(ctx context.Context) error {
	exists, err := n.Exists(ctx)
	if err != nil {
		return err
	}

	kind, err := n.GetKind()
	if err != nil {
		return err
	}

	name, err := n.GetName()
	if err != nil {
		return err
	}

	namespaceName, err := n.GetNamespaceName()
	if err != nil {
		return err
	}

	if exists {
		objectInterface, err := n.GetObjectInterface(ctx)
		if err != nil {
			return err
		}

		err = objectInterface.Delete(ctx, name, v1.DeleteOptions{})
		if err != nil {
			return tracederrors.TracedErrorf("Failed to delete object '%s' named '%s' in namespace '%s': %w", kind, name, namespaceName, err)
		}
		logging.LogChangedByCtxf(ctx, "Object '%s' named '%s' in namespace '%s' deleted.", kind, name, namespaceName)
	} else {
		logging.LogInfoByCtxf(ctx, "Object '%s' named '%s' in namespace '%s' already absent. Skip delete.", kind, name, namespaceName)
	}

	return nil
}

func (n *NativeObject) GetAsYamlString() (yamlString string, err error) {
	return "", tracederrors.TracedErrorNotImplemented()
}

func (n *NativeObject) SetKind(kind string) error {
	if kind == "" {
		return tracederrors.TracedErrorEmptyString("kind")
	}

	n.kind = kind

	return nil
}

func (n *NativeObject) SetApiVersion(apiVersion string) error {
	if apiVersion == "" {
		return tracederrors.TracedErrorEmptyString("apiVersion")
	}

	n.apiVersion = apiVersion

	return nil
}
