package kubernetes

import "github.com/asciich/asciichgolangpublic"

type CommandExecutorNamespace struct {
	name              string
	kubernetesCluster KubernetesCluster
}

func NewCommandExecutorNamespace() (c *CommandExecutorNamespace) {
	return new(CommandExecutorNamespace)
}

func (c *CommandExecutorNamespace) GetKubernetesCluster() (kubernetesCluster KubernetesCluster, err error) {

	return c.kubernetesCluster, nil
}

func (c *CommandExecutorNamespace) GetName() (name string, err error) {
	if c.name == "" {
		return "", asciichgolangpublic.TracedErrorf("name not set")
	}

	return c.name, nil
}

func (c *CommandExecutorNamespace) MustGetKubernetesCluster() (kubernetesCluster KubernetesCluster) {
	kubernetesCluster, err := c.GetKubernetesCluster()
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return kubernetesCluster
}

func (c *CommandExecutorNamespace) MustGetName() (name string) {
	name, err := c.GetName()
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return name
}

func (c *CommandExecutorNamespace) MustSetKubernetesCluster(kubernetesCluster KubernetesCluster) {
	err := c.SetKubernetesCluster(kubernetesCluster)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorNamespace) MustSetName(name string) {
	err := c.SetName(name)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorNamespace) SetKubernetesCluster(kubernetesCluster KubernetesCluster) (err error) {
	c.kubernetesCluster = kubernetesCluster

	return nil
}

func (c *CommandExecutorNamespace) SetName(name string) (err error) {
	if name == "" {
		return asciichgolangpublic.TracedErrorf("name is empty string")
	}

	c.name = name

	return nil
}
