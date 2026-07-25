package nativekubernetesoo

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"k8s.io/client-go/kubernetes"
)

type ReplicaSet struct {
	namespace *NativeNamespace
	name      string
}

func (r *ReplicaSet) GetNamespace() (kubernetesinterfaces.Namespace, error) {
	return r.GetNativeNamespace()
}

func (r *ReplicaSet) GetNativeNamespace() (*NativeNamespace, error) {
	if r.namespace == nil {
		return nil, tracederrors.TracedErrorNil("namespace")
	}

	return r.namespace, nil
}

func (r *ReplicaSet) SetNamespace(namespace *NativeNamespace) error {
	if namespace == nil {
		return tracederrors.TracedErrorNil("namespace")
	}

	r.namespace = namespace

	return nil
}

func (r *ReplicaSet) GetName() (string, error) {
	if r.name == "" {
		return "", tracederrors.TracedErrorEmptyString("name")
	}

	return r.name, nil
}

func (r *ReplicaSet) SetName(name string) error {
	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	r.name = name

	return nil
}

func (r *ReplicaSet) GetNamespaceName() (string, error) {
	namespace, err := r.GetNamespace()
	if err != nil {
		return "", err
	}

	return namespace.GetName()
}

func (r *ReplicaSet) GetClientSet() (*kubernetes.Clientset, error) {
	namespace, err := r.GetNativeNamespace()
	if err != nil {
		return nil, err
	}

	return namespace.GetClientSet()
}

func (r *ReplicaSet) Delete(ctx context.Context) (err error) {
	replicaSetName, err := r.GetName()
	if err != nil {
		return err
	}

	namespaceName, err := r.GetNamespaceName()
	if err != nil {
		return err
	}

	clientSet, err := r.GetClientSet()
	if err != nil {
		return err
	}

	return nativekubernetes.DeleteReplicaSet(ctx, clientSet, replicaSetName, namespaceName)
}

func (r *ReplicaSet) Exists(ctx context.Context) (bool, error) {
	replicaSetName, err := r.GetName()
	if err != nil {
		return false, err
	}

	namespaceName, err := r.GetNamespaceName()
	if err != nil {
		return false, err
	}

	clientSet, err := r.GetClientSet()
	if err != nil {
		return false, err
	}

	return nativekubernetes.ReplicaSetExists(ctx, clientSet, replicaSetName, namespaceName)
}
