package nativekubernetesoo

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"k8s.io/client-go/kubernetes"
)

type DaemonSet struct {
	namespace *NativeNamespace
	name      string
}

func (d *DaemonSet) GetNamespace() (kubernetesinterfaces.Namespace, error) {
	return d.GetNativeNamespace()
}

func (d *DaemonSet) GetNativeNamespace() (*NativeNamespace, error) {
	if d.namespace == nil {
		return nil, tracederrors.TracedErrorNil("namespace")
	}

	return d.namespace, nil
}

func (d *DaemonSet) SetNamespace(namespace *NativeNamespace) error {
	if namespace == nil {
		return tracederrors.TracedErrorNil("namespace")
	}

	d.namespace = namespace

	return nil
}

func (d *DaemonSet) GetName() (string, error) {
	if d.name == "" {
		return "", tracederrors.TracedErrorEmptyString("name")
	}

	return d.name, nil
}

func (d *DaemonSet) SetName(name string) error {
	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	d.name = name

	return nil
}

func (d *DaemonSet) GetNamespaceName() (string, error) {
	namespace, err := d.GetNamespace()
	if err != nil {
		return "", err
	}

	return namespace.GetName()
}

func (d *DaemonSet) GetClientSet() (*kubernetes.Clientset, error) {
	namespace, err := d.GetNativeNamespace()
	if err != nil {
		return nil, err
	}

	return namespace.GetClientSet()
}

func (d *DaemonSet) Delete(ctx context.Context) (err error) {
	daemonSetName, err := d.GetName()
	if err != nil {
		return err
	}

	namespaceName, err := d.GetNamespaceName()
	if err != nil {
		return err
	}

	clientSet, err := d.GetClientSet()
	if err != nil {
		return err
	}

	return nativekubernetes.DeleteDaemonSet(ctx, clientSet, daemonSetName, namespaceName)
}

func (d *DaemonSet) Exists(ctx context.Context) (bool, error) {
	daemonSetName, err := d.GetName()
	if err != nil {
		return false, err
	}

	namespaceName, err := d.GetNamespaceName()
	if err != nil {
		return false, err
	}

	clientSet, err := d.GetClientSet()
	if err != nil {
		return false, err
	}

	return nativekubernetes.DaemonSetExists(ctx, clientSet, daemonSetName, namespaceName)
}
