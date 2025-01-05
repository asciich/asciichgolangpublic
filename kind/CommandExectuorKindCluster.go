package kind

import "github.com/asciich/asciichgolangpublic"

import (
	"github.com/asciich/asciichgolangpublic/kubernetes"
)

type KindCluster struct {
	kubernetes.CommandExecutorKubernetes
	kind Kind
}

func NewKindCluster() (k *KindCluster) {
	return new(KindCluster)
}

func (k *KindCluster) GetKind() (kind Kind, err error) {

	return k.kind, nil
}

func (k *KindCluster) MustGetKind() (kind Kind) {
	kind, err := k.GetKind()
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return kind
}

func (k *KindCluster) MustSetKind(kind Kind) {
	err := k.SetKind(kind)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}
}

func (k *KindCluster) SetKind(kind Kind) (err error) {
	k.kind = kind

	return nil
}
