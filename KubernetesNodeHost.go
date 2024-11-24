package asciichgolangpublic

import (
	"strings"

)

type KubernetesNodeHost struct {
	Host
}

func GetKubernetesNodeByHostname(hostname string) (kubernetesNodeHost *KubernetesNodeHost, err error) {
	if len(hostname) <= 0 {
		return nil, TracedError("hostname is empty string")
	}

	kubernetesNodeHost = NewKubernetesNodeHost()

	err = kubernetesNodeHost.SetHostname(hostname)
	if err != nil {
		return nil, err
	}

	return kubernetesNodeHost, err
}

func MustGetKubernetesNodeByHostname(hostname string) (kubernetesNodeHost *KubernetesNodeHost) {
	kubernetesNodeHost, err := GetKubernetesNodeByHostname(hostname)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return kubernetesNodeHost
}

func NewKubernetesNodeHost() (kubernetesNodeHost *KubernetesNodeHost) {
	kubernetesNodeHost = new(KubernetesNodeHost)
	return kubernetesNodeHost
}

func (k *KubernetesNodeHost) MustCheckIsKubernetesNode(verbose bool) (isKubernetesNode bool) {
	isKubernetesNode, err := k.CheckIsKubernetesNode(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return isKubernetesNode
}

func (k *KubernetesNodeHost) MustIsKubernetesNode(verbose bool) (isKubernetesNode bool) {
	isKubernetesNode, err := k.IsKubernetesNode(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return isKubernetesNode
}

func (n *KubernetesNodeHost) CheckIsKubernetesNode(verbose bool) (isKubernetesNode bool, err error) {
	hostname, err := n.GetHostname()
	if err != nil {
		return false, err
	}

	isKubernetesNode, err = n.IsKubernetesNode(verbose)
	if err != nil {
		return false, err
	}

	if !isKubernetesNode {
		return false, TracedErrorf("Host '%s' is not a kubernetes node", hostname)
	}

	return isKubernetesNode, nil
}

func (n *KubernetesNodeHost) IsKubernetesNode(verbose bool) (isKubernetesNode bool, err error) {
	stdout, err := n.RunCommandAndGetStdoutAsString(
		&RunCommandOptions{
			Command: []string{"ctr", "--namespace", "k8s.io", "containers", "list"},
			Verbose: verbose,
		},
	)
	if err != nil {
		return false, err
	}

	hostname, err := n.GetHostname()
	if err != nil {
		return false, err
	}

	isKubernetesNode = true

	if len(Strings().SplitLines(stdout, false)) <= 5 {
		isKubernetesNode = false
	}

	if !strings.Contains(stdout, "registry.k8s.io/kube-proxy") {
		isKubernetesNode = false
	}

	if strings.Contains(stdout, "registry.k8s.io/etcd") {
		if verbose {
			LogInfof("Host '%s' seems to be a kubernetes controlplane since etcd container was found, not a node itself.", hostname)
		}
		isKubernetesNode = false
	}

	if verbose {
		if isKubernetesNode {
			LogInfof("Host '%s' is a kubernetes node.", hostname)
		} else {
			LogInfof("Host '%s' is not a kubernetes node.", hostname)
		}
	}

	return isKubernetesNode, nil
}
