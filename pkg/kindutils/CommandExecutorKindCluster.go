package kindutils

import (
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/commandexecutorkubernetes"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type CommandExecutorKindCluster struct {
	commandexecutorkubernetes.CommandExecutorKubernetes
	kind Kind
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
