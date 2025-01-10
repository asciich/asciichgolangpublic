package kubernetes

import (
	"strings"

	"github.com/asciich/asciichgolangpublic"
	"github.com/asciich/asciichgolangpublic/hosts"
)

type KubernetesNodeHost struct {
	hosts.Host
}

func GetKubernetesNodeByHostname(hostname string) (kubernetesNodeHost *KubernetesNodeHost, err error) {
	if len(hostname) <= 0 {
		return nil, asciichgolangpublic.TracedError("hostname is empty string")
	}

	kubernetesNodeHost = NewKubernetesNodeHost()

	err = kubernetesNodeHost.SetHostName(hostname)
	if err != nil {
		return nil, err
	}

	return kubernetesNodeHost, err
}

func MustGetKubernetesNodeByHostname(hostname string) (kubernetesNodeHost *KubernetesNodeHost) {
	kubernetesNodeHost, err := GetKubernetesNodeByHostname(hostname)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
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
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return isKubernetesNode
}

func (k *KubernetesNodeHost) MustIsKubernetesNode(verbose bool) (isKubernetesNode bool) {
	isKubernetesNode, err := k.IsKubernetesNode(verbose)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return isKubernetesNode
}

func (k *KubernetesNodeHost) CheckIsKubernetesNode(verbose bool) (isKubernetesNode bool, err error) {
	hostname, err := k.GetHostName()
	if err != nil {
		return false, err
	}

	isKubernetesNode, err = k.IsKubernetesNode(verbose)
	if err != nil {
		return false, err
	}

	if !isKubernetesNode {
		return false, asciichgolangpublic.TracedErrorf("Host '%s' is not a kubernetes node", hostname)
	}

	return isKubernetesNode, nil
}

func (k *KubernetesNodeHost) IsKubernetesNode(verbose bool) (isKubernetesNode bool, err error) {
	stdout, err := k.RunCommandAndGetStdoutAsString(
		&asciichgolangpublic.RunCommandOptions{
			Command: []string{"ctr", "--namespace", "k8s.io", "containers", "list"},
			Verbose: verbose,
		},
	)
	if err != nil {
		return false, err
	}

	hostname, err := k.GetHostName()
	if err != nil {
		return false, err
	}

	isKubernetesNode = true

	if len(asciichgolangpublic.Strings().SplitLines(stdout, false)) <= 5 {
		isKubernetesNode = false
	}

	if !strings.Contains(stdout, "registry.k8s.io/kube-proxy") {
		isKubernetesNode = false
	}

	if strings.Contains(stdout, "registry.k8s.io/etcd") {
		if verbose {
			asciichgolangpublic.LogInfof("Host '%s' seems to be a kubernetes controlplane since etcd container was found, not a node itself.", hostname)
		}
		isKubernetesNode = false
	}

	if verbose {
		if isKubernetesNode {
			asciichgolangpublic.LogInfof("Host '%s' is a kubernetes node.", hostname)
		} else {
			asciichgolangpublic.LogInfof("Host '%s' is not a kubernetes node.", hostname)
		}
	}

	return isKubernetesNode, nil
}
