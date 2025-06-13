package nativekubernetes

import (
	"context"

	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesimplementationindependend"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/tracederrors"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	sigyaml "sigs.k8s.io/yaml"
)

type NativeResource struct {
	name       string
	kind       string
	apiVersion string
	namespace  *NativeNamespace
}

func (n *NativeResource) GetApiVersion(ctx context.Context) (string, error) {
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

func (n *NativeResource) GetName() (string, error) {
	if n.name == "" {
		return "", tracederrors.TracedError("name not set")
	}

	return n.name, nil
}

func (n *NativeResource) GetKind() (string, error) {
	if n.kind == "" {
		return "", tracederrors.TracedError("kind not set")
	}

	ret, err := kubernetesimplementationindependend.SanitizeKindName(n.kind)
	if err != nil {
		return "", err
	}

	return ret, nil
}

func (n *NativeResource) GetGroupVersionKind(ctx context.Context) (*schema.GroupVersionKind, error) {
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

func (n *NativeResource) EnsureNamespaceExists(ctx context.Context) error {
	namespace, err := n.GetNamespace()
	if err != nil {
		return err
	}

	return namespace.Create(ctx)
}

func (n *NativeResource) CreateByYamlString(ctx context.Context, options *kubernetesparameteroptions.CreateResourceOptions) (err error) {
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

	logging.LogInfoByCtxf(ctx, "Create kubernetes resource by yaml '%s/%s' in namespace '%s' started.", kind, name, namespaceName)

	exists, err := n.Exists(ctx)
	if err != nil {
		return err
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Resource '%s' named '%s' in namespace '%s' already exists, skip creation.", kind, name, namespaceName)
	} else {
		err = n.EnsureNamespaceExists(ctx)
		if err != nil {
			return err
		}

		gvk, err := n.GetGroupVersionKind(ctx)
		if err != nil {
			return err
		}

		var unstructuredObj unstructured.Unstructured
		if err := sigyaml.Unmarshal([]byte(yamlString), &unstructuredObj); err != nil {
			return tracederrors.TracedErrorf("Failed to parse yamlString as unstructuredObj: %w", err)
		}

		unstructuredObj.SetGroupVersionKind(*gvk)
		unstructuredObj.SetNamespace(namespaceName)

		resourceInterface, err := n.GetResourceInterface(ctx)
		if err != nil {
			return err
		}

		_, err = resourceInterface.Create(ctx, &unstructuredObj, v1.CreateOptions{})
		if err != nil {
			return tracederrors.TracedErrorf("Failed to create resource '%s' names '%s' in namespace '%s': %w", kind, name, namespaceName, err)
		}

		logging.LogChangedByCtxf(ctx, "Created resource '%s' named '%s' in namespace '%s'.", kind, name, namespaceName)
	}

	logging.LogInfoByCtxf(ctx, "Create kubernetes resource by yaml '%s/%s' in namespace '%s' finished.", kind, name, namespaceName)

	return nil
}

func (n *NativeResource) GetNamespace() (*NativeNamespace, error) {
	if n.namespace == nil {
		return nil, tracederrors.TracedErrorNil("namespace")
	}

	return n.namespace, nil
}

func (n *NativeResource) GetNamespaceName() (string, error) {
	namespace, err := n.GetNamespace()
	if err != nil {
		return "", err
	}

	return namespace.GetName()
}

func (n *NativeResource) GetDynamicClient() (*dynamic.DynamicClient, error) {
	namespace, err := n.GetNamespace()
	if err != nil {
		return nil, err
	}

	return namespace.GetDynamicClient()
}

func (n *NativeResource) GetGroupVersion(ctx context.Context) (*schema.GroupVersion, error) {
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

func (n *NativeResource) GetResourcePlural() (string, error) {
	kind, err := n.GetKind()
	if err != nil {
		return "", err
	}

	return kubernetesimplementationindependend.GetResourcePlural(kind)
}

func (n *NativeResource) GetGroupVersionResource(ctx context.Context) (*schema.GroupVersionResource, error) {
	groupVersion, err := n.GetGroupVersion(ctx)
	if err != nil {
		return nil, err
	}

	resourcePlural, err := n.GetResourcePlural()
	if err != nil {
		return nil, err
	}

	gvr := schema.GroupVersionResource{
		Group:    groupVersion.Group,
		Version:  groupVersion.Version,
		Resource: resourcePlural,
	}

	return &gvr, nil
}

func (n *NativeResource) GetResourceInterface(ctx context.Context) (dynamic.ResourceInterface, error) {
	groupVersionResource, err := n.GetGroupVersionResource(ctx)
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

	resourceInterface := dynamicClient.Resource(*groupVersionResource).Namespace(namspaceName)

	return resourceInterface, nil
}

func (n *NativeResource) Exists(ctx context.Context) (bool, error) {
	resourceInterface, err := n.GetResourceInterface(ctx)
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
	_, err = resourceInterface.Get(ctx, name, v1.GetOptions{})
	if err == nil {
		exists = true
	} else {
		if !errors.IsNotFound(err) {
			return false, tracederrors.TracedErrorf("failed to get resource '%s' named '%s' in namespace '%s': %w", kind, name, namespaceName, err)
		}
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Resource '%s' named '%s' in namespace '%s' exists.", kind, name, namespaceName)
	} else {
		logging.LogInfoByCtxf(ctx, "Resource '%s' named '%s' in namespace '%s' does not exist.", kind, name, namespaceName)
	}

	return exists, nil
}

func (n *NativeResource) Delete(ctx context.Context) error {
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
		resourceInterface, err := n.GetResourceInterface(ctx)
		if err != nil {
			return err
		}

		err = resourceInterface.Delete(ctx, name, v1.DeleteOptions{})
		if err != nil {
			return tracederrors.TracedErrorf("Failed to delete resource '%s' named '%s' in namespace '%s': %w", kind, name, namespaceName, err)
		}
		logging.LogChangedByCtxf(ctx, "Resource '%s' named '%s' in namespace '%s' deleted.", kind, name, namespaceName)
	} else {
		logging.LogInfoByCtxf(ctx, "Resource '%s' named '%s' in namespace '%s' already absent. Skip delete.", kind, name, namespaceName)
	}

	return nil
}

func (n *NativeResource) GetAsYamlString() (yamlString string, err error) {
	return "", tracederrors.TracedErrorNotImplemented()
}

func (n *NativeResource) SetKind(kind string) error {
	if kind == "" {
		return tracederrors.TracedErrorEmptyString("kind")
	}

	n.kind = kind

	return nil
}

func (n *NativeResource) SetAPIVersion(apiVersion string) error {
	if apiVersion == "" {
		return tracederrors.TracedErrorEmptyString("apiVersion")
	}

	n.apiVersion = apiVersion

	return nil
}
