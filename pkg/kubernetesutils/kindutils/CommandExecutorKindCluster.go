package kindutils

import (
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/commandexecutorkubernetes"
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type CommandExecutorKindCluster struct {
	commandexecutorkubernetes.CommandExecutorKubernetes
	kind        Kind
	clusterName string
}

func NewCommandExecutorKindCluster() (k *CommandExecutorKindCluster) {
	return new(CommandExecutorKindCluster)
}

func (k *CommandExecutorKindCluster) GetKind() (kind Kind, err error) {
	if k.kind == nil {
		return nil, tracederrors.TracedError("kind not set")
	}
	return k.kind, nil
}

func (k *CommandExecutorKindCluster) MustGetKind() (kind Kind) {
	kind, err := k.GetKind()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return kind
}

func (k *CommandExecutorKindCluster) MustSetKind(kind Kind) {
	err := k.SetKind(kind)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (k *CommandExecutorKindCluster) SetKind(kind Kind) (err error) {
	if kind == nil {
		return tracederrors.TracedErrorNil("kind")
	}

	k.kind = kind

	return nil
}

func (k *CommandExecutorKindCluster) GetClusterName() (string, error) {
	if k.clusterName == "" {
		return "", tracederrors.TracedError("clusterName not set")
	}
	return k.clusterName, nil
}

func (k *CommandExecutorKindCluster) SetClusterName(clusterName string) error {
	if clusterName == "" {
		return tracederrors.TracedErrorEmptyString("clusterName")
	}
	k.clusterName = clusterName
	return nil
}
// CronJob stub implementations
func (c *CommandExecutorKindCluster) CronJobByNameExists(ctx context.Context, namespaceName string, cronJobName string) (exists bool, err error) {
	return false, tracederrors.TracedErrorNotImplemented()
}

func (c *CommandExecutorKindCluster) CreateCronJob(ctx context.Context, namespaceName string, cronJobName string, schedule string, image string, command []string, labels map[string]string) (kubernetesinterfaces.CronJob, error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (c *CommandExecutorKindCluster) DeleteCronJobByName(ctx context.Context, namespaceName string, cronJobName string) error {
	return tracederrors.TracedErrorNotImplemented()
}
