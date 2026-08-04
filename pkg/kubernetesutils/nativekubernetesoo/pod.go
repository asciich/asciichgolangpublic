package nativekubernetesoo

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"k8s.io/client-go/kubernetes"
)

type Pod struct {
	namespace *NativeNamespace
	name      string
}

func (p *Pod) GetNamespace() (kubernetesinterfaces.Namespace, error) {
	return p.GetNativeNamespace()
}

func (p *Pod) GetNativeNamespace() (*NativeNamespace, error) {
	if p.namespace == nil {
		return nil, tracederrors.TracedErrorNil("namespace")
	}

	return p.namespace, nil
}

func (p *Pod) SetNamespace(namespace *NativeNamespace) error {
	if namespace == nil {
		return tracederrors.TracedErrorNil("namespace")
	}

	p.namespace = namespace

	return nil
}

func (p *Pod) GetName() (string, error) {
	if p.name == "" {
		return "", tracederrors.TracedErrorEmptyString("name")
	}

	return p.name, nil
}

func (p *Pod) SetName(name string) error {
	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	p.name = name

	return nil
}

func (p *Pod) GetNamespaceName() (string, error) {
	namespace, err := p.GetNamespace()
	if err != nil {
		return "", err
	}

	return namespace.GetName()
}

func (p *Pod) GetClientSet() (*kubernetes.Clientset, error) {
	namespace, err := p.GetNativeNamespace()
	if err != nil {
		return nil, err
	}

	return namespace.GetClientSet()
}

func (p *Pod) Delete(ctx context.Context) (err error) {
	podName, err := p.GetName()
	if err != nil {
		return err
	}

	namespaceName, err := p.GetNamespaceName()
	if err != nil {
		return err
	}

	clientSet, err := p.GetClientSet()
	if err != nil {
		return err
	}

	return nativekubernetes.DeletePod(ctx, clientSet, podName, namespaceName)
}

func (p *Pod) Exists(ctx context.Context) (bool, error) {
	podName, err := p.GetName()
	if err != nil {
		return false, err
	}

	namespaceName, err := p.GetNamespaceName()
	if err != nil {
		return false, err
	}

	clientSet, err := p.GetClientSet()
	if err != nil {
		return false, err
	}

	return nativekubernetes.PodExists(ctx, clientSet, podName, namespaceName)
}

func (p *Pod) GetContainerLogs(ctx context.Context, containerName string) (stdout []byte, stderr []byte, err error) {
	podName, err := p.GetName()
	if err != nil {
		return nil, nil, err
	}

	namespaceName, err := p.GetNamespaceName()
	if err != nil {
		return nil, nil, err
	}

	clientSet, err := p.GetClientSet()
	if err != nil {
		return nil, nil, err
	}

	return nativekubernetes.GetContainerLogs(ctx, clientSet, namespaceName, podName, containerName)
}

func (p *Pod) CopyFileToPod(ctx context.Context, localFile string, destPath string, containerName string) error {
	podName, err := p.GetName()
	if err != nil {
		return err
	}

	namespaceName, err := p.GetNamespaceName()
	if err != nil {
		return err
	}

	config, err := p.namespace.GetConfig()
	if err != nil {
		return err
	}

	return nativekubernetes.CopyFileToPod(ctx, config, localFile, destPath, podName, containerName, namespaceName)
}

func (p *Pod) CopyFileFromPod(ctx context.Context, srcPath string, destFile string, containerName string) error {
	podName, err := p.GetName()
	if err != nil {
		return err
	}

	namespaceName, err := p.GetNamespaceName()
	if err != nil {
		return err
	}

	config, err := p.namespace.GetConfig()
	if err != nil {
		return err
	}

	return nativekubernetes.CopyFileFromPod(ctx, config, podName, namespaceName, containerName, srcPath, destFile)
}
