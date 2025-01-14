package kubernetes

import (
	"github.com/asciich/asciichgolangpublic/errors"
	"github.com/asciich/asciichgolangpublic/logging"
)

type KubectlContext struct {
	name    string
	cluster string
}

func NewKubectlContext() (k *KubectlContext) {
	return new(KubectlContext)
}

func (k *KubectlContext) GetCluster() (cluster string, err error) {
	if k.cluster == "" {
		return "", errors.TracedErrorf("cluster not set")
	}

	return k.cluster, nil
}

func (k *KubectlContext) GetName() (name string, err error) {
	if k.name == "" {
		return "", errors.TracedErrorf("name not set")
	}

	return k.name, nil
}

func (k *KubectlContext) MustGetCluster() (cluster string) {
	cluster, err := k.GetCluster()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return cluster
}

func (k *KubectlContext) MustGetName() (name string) {
	name, err := k.GetName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return name
}

func (k *KubectlContext) MustSetCluster(cluster string) {
	err := k.SetCluster(cluster)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (k *KubectlContext) MustSetName(name string) {
	err := k.SetName(name)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (k *KubectlContext) SetCluster(cluster string) (err error) {
	if cluster == "" {
		return errors.TracedErrorf("cluster is empty string")
	}

	k.cluster = cluster

	return nil
}

func (k *KubectlContext) SetName(name string) (err error) {
	if name == "" {
		return errors.TracedErrorf("name is empty string")
	}

	k.name = name

	return nil
}
