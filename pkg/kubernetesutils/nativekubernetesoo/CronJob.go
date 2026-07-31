package nativekubernetesoo

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type NativeCronJob struct {
	namespace *NativeNamespace
	name      string
}

func (n *NativeCronJob) GetNamespace() (*NativeNamespace, error) {
	if n.namespace == nil {
		return nil, tracederrors.TracedError("namespace not set")
	}

	return n.namespace, nil
}

func (n *NativeCronJob) GetName() (string, error) {
	if n.name == "" {
		return "", tracederrors.TracedError("name not set")
	}

	return n.name, nil
}

func (n *NativeCronJob) Exists(ctx context.Context) (bool, error) {
	cronJobName, err := n.GetName()
	if err != nil {
		return false, err
	}

	namespace, err := n.GetNamespace()
	if err != nil {
		return false, err
	}

	return namespace.CronJobByNameExists(ctx, cronJobName)
}
