package kubernetesinterfaces

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
)

type Pod interface {
	commandexecutorinterfaces.CommandExecutor
	Delete(ctx context.Context) (err error)
	Exists(ctx context.Context) (bool, error)
	GetName() (name string, err error)
	GetNamespace() (namespace Namespace, err error)
	GetContainerLogs(ctx context.Context, containerName string) (stdout []byte, stderr []byte, err error)
	CopyFileToPod(ctx context.Context, localFile string, destPath string, containerName string) error
	CopyFileFromPod(ctx context.Context, srcPath string, destFile string, containerName string) error
	RunCommandInContainer(ctx context.Context, options *kubernetesparameteroptions.KubernetesRunCommandOptions) (*commandoutput.CommandOutput, error)
}
