package kubernetesparameteroptions

import (
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/logging"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/tracederrors"
)

type ListKubernetesObjectsOptions struct {
	Namespace  string
	ObjectType string
	Verbose    bool
}

func NewListKubernetesObjectsOptions() (k *ListKubernetesObjectsOptions) {
	return new(ListKubernetesObjectsOptions)
}

func (k *ListKubernetesObjectsOptions) GetNamespace() (namespace string, err error) {
	if k.Namespace == "" {
		return "", tracederrors.TracedErrorf("Namespace not set")
	}

	return k.Namespace, nil
}

func (k *ListKubernetesObjectsOptions) GetObjectType() (resourceType string, err error) {
	if k.ObjectType == "" {
		return "", tracederrors.TracedErrorf("ObjectType not set")
	}

	return k.ObjectType, nil
}

func (k *ListKubernetesObjectsOptions) GetVerbose() (verbose bool) {

	return k.Verbose
}

func (k *ListKubernetesObjectsOptions) MustGetNamespace() (namespace string) {
	namespace, err := k.GetNamespace()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return namespace
}

func (k *ListKubernetesObjectsOptions) MustGetObjectType() (resourceType string) {
	resourceType, err := k.GetObjectType()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return resourceType
}

func (k *ListKubernetesObjectsOptions) MustSetNamespace(namespace string) {
	err := k.SetNamespace(namespace)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (k *ListKubernetesObjectsOptions) MustSetObjectType(resourceType string) {
	err := k.SetObjectType(resourceType)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (k *ListKubernetesObjectsOptions) SetNamespace(namespace string) (err error) {
	if namespace == "" {
		return tracederrors.TracedErrorf("namespace is empty string")
	}

	k.Namespace = namespace

	return nil
}

func (k *ListKubernetesObjectsOptions) SetObjectType(resourceType string) (err error) {
	if resourceType == "" {
		return tracederrors.TracedErrorf("resourceType is empty string")
	}

	k.ObjectType = resourceType

	return nil
}

func (k *ListKubernetesObjectsOptions) SetVerbose(verbose bool) {
	k.Verbose = verbose
}
