package nativekubernetesoo

import (
	"context"
	"fmt"
	"io"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Pod struct {
	commandexecutorgeneric.CommandExecutorBase
	namespace *NativeNamespace
	name      string
}

func NewPod() *Pod {
	ret := new(Pod)
	ret.SetParentCommandExecutorForBaseClass(ret)
	return ret
}

func (p *Pod) GetDeepCopyAsCommandExecutor() commandexecutorinterfaces.CommandExecutor {
	ret := &Pod{}
	*ret = *p

	return p
}

func (p *Pod) RunCommandAndGetStdoutAsIoReadCloser(ctx context.Context, options *parameteroptions.RunCommandOptions) (io.ReadCloser, error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (p *Pod) RunCommandAndGetStdinAsIoWriteCloser(ctx context.Context, options *parameteroptions.RunCommandOptions) (io.WriteCloser, error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

// Returns the default Container name to exec commands.
// - Takes the first running container name.
// - Init containers are not taken into account.
func (p *Pod) GetDefaultContainerName(ctx context.Context) (string, error) {
	podName, err := p.GetName()
	if err != nil {
		return "", err
	}

	namespaceName, err := p.GetNamespaceName()
	if err != nil {
		return "", err
	}

	clientSet, err := p.GetClientSet()
	if err != nil {
		return "", err
	}

	pod, err := clientSet.CoreV1().Pods(namespaceName).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return "", tracederrors.TracedErrorf("Failed to get pod '%s' in namespace '%s': %w", podName, namespaceName, err)
	}

	// Find the first running container (excluding init containers)
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if containerStatus.State.Running != nil {
			logging.LogInfoByCtxf(ctx, "Default container name is '%s' in pod '%s'", containerStatus.Name, podName)
			return containerStatus.Name, nil
		}
	}

	return "", tracederrors.TracedErrorf("No running container found in pod '%s' in namespace '%s'", podName, namespaceName)
}

func (p *Pod) RunCommand(ctx context.Context, options *parameteroptions.RunCommandOptions) (*commandoutput.CommandOutput, error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	podName, err := p.GetName()
	if err != nil {
		return nil, err
	}

	containerName, err := p.GetDefaultContainerName(ctx)
	if err != nil {
		return nil, err
	}

	return p.RunCommandInContainer(ctx, &kubernetesparameteroptions.KubernetesRunCommandOptions{
		PodName:           podName,
		ContainerName:     containerName,
		RunCommandOptions: options,
	})
}

func (p *Pod) IsRunningOnLocalhost() (bool, error) {
	return false, nil
}

func (p *Pod) GetClusterName() (string, error) {
	namespace, err := p.GetNamespace()
	if err != nil {
		return "", err
	}

	return namespace.GetClusterName()
}

func (p *Pod) GetHostDescription() (string, error) {
	podName, err := p.GetName()
	if err != nil {
		return "", err
	}

	namespaceName, err := p.GetNamespaceName()
	if err != nil {
		return "", err
	}

	clusterName, err := p.GetClusterName()
	if err != nil {
		return "", err
	}

	got := fmt.Sprintf("Pod '%s' in namespace '%s' of kubernetes cluster '%s'.", podName, namespaceName, clusterName)

	return got, nil
}

func (p *Pod) GetCPUArchitecture(ctx context.Context) (string, error) {
	return "", tracederrors.TracedErrorNotImplemented()
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

func (p *Pod) RunCommandInContainer(ctx context.Context, options *kubernetesparameteroptions.KubernetesRunCommandOptions) (*commandoutput.CommandOutput, error) {
	podName, err := p.GetName()
	if err != nil {
		return nil, err
	}

	namespaceName, err := p.GetNamespaceName()
	if err != nil {
		return nil, err
	}

	config, err := p.namespace.GetConfig()
	if err != nil {
		return nil, err
	}

	// Set pod name in options if not already set
	if options.PodName == "" {
		options.PodName = podName
	}

	return nativekubernetes.RunCommand(ctx, config, namespaceName, options)
}
