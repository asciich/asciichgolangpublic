package parameteroptions

import (
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type ListKubernetesResourcesOptions struct {
	Namespace    string
	ResourceType string
	Verbose      bool
}

func NewListKubernetesResourcesOptions() (k *ListKubernetesResourcesOptions) {
	return new(ListKubernetesResourcesOptions)
}

func (k *ListKubernetesResourcesOptions) GetNamespace() (namespace string, err error) {
	if k.Namespace == "" {
		return "", tracederrors.TracedErrorf("Namespace not set")
	}

	return k.Namespace, nil
}

func (k *ListKubernetesResourcesOptions) GetResourceType() (resourceType string, err error) {
	if k.ResourceType == "" {
		return "", tracederrors.TracedErrorf("ResourceType not set")
	}

	return k.ResourceType, nil
}

func (k *ListKubernetesResourcesOptions) GetVerbose() (verbose bool) {

	return k.Verbose
}

func (k *ListKubernetesResourcesOptions) MustGetNamespace() (namespace string) {
	namespace, err := k.GetNamespace()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return namespace
}

func (k *ListKubernetesResourcesOptions) MustGetResourceType() (resourceType string) {
	resourceType, err := k.GetResourceType()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return resourceType
}

func (k *ListKubernetesResourcesOptions) MustSetNamespace(namespace string) {
	err := k.SetNamespace(namespace)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (k *ListKubernetesResourcesOptions) MustSetResourceType(resourceType string) {
	err := k.SetResourceType(resourceType)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (k *ListKubernetesResourcesOptions) SetNamespace(namespace string) (err error) {
	if namespace == "" {
		return tracederrors.TracedErrorf("namespace is empty string")
	}

	k.Namespace = namespace

	return nil
}

func (k *ListKubernetesResourcesOptions) SetResourceType(resourceType string) (err error) {
	if resourceType == "" {
		return tracederrors.TracedErrorf("resourceType is empty string")
	}

	k.ResourceType = resourceType

	return nil
}

func (k *ListKubernetesResourcesOptions) SetVerbose(verbose bool) {
	k.Verbose = verbose
}
