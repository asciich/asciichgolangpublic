package kubernetesutils

import (
	"strings"

	"gitlab.asciich.ch/tools/asciichgolangpublic.git/datatypes/stringsutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/hosts"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/parameteroptions"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/contextutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/logging"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/tracederrors"
)

type KubernetesNodeHost struct {
	hosts.Host
}

func GetKubernetesNodeByHost(host hosts.Host) (kubernetesNodeHost *KubernetesNodeHost, err error) {
	if host == nil {
		return nil, tracederrors.TracedErrorNil("host")
	}

	kubernetesNodeHost = NewKubernetesNodeHost()

	kubernetesNodeHost.Host = host

	return kubernetesNodeHost, nil
}

func GetKubernetesNodeByHostname(hostname string) (kubernetesNodeHost *KubernetesNodeHost, err error) {
	if len(hostname) <= 0 {
		return nil, tracederrors.TracedError("hostname is empty string")
	}

	host, err := hosts.GetHostByHostname("hostname")
	if err != nil {
		return nil, err
	}

	return GetKubernetesNodeByHost(host)
}

func NewKubernetesNodeHost() (kubernetesNodeHost *KubernetesNodeHost) {
	kubernetesNodeHost = new(KubernetesNodeHost)
	return kubernetesNodeHost
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
		return false, tracederrors.TracedErrorf("Host '%s' is not a kubernetes node", hostname)
	}

	return isKubernetesNode, nil
}

func (k *KubernetesNodeHost) IsKubernetesNode(verbose bool) (isKubernetesNode bool, err error) {
	stdout, err := k.RunCommandAndGetStdoutAsString(
		contextutils.GetVerbosityContextByBool(verbose),
		&parameteroptions.RunCommandOptions{
			Command: []string{"ctr", "--namespace", "k8s.io", "containers", "list"},
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

	if len(stringsutils.SplitLines(stdout, false)) <= 5 {
		isKubernetesNode = false
	}

	if !strings.Contains(stdout, "registry.k8s.io/kube-proxy") {
		isKubernetesNode = false
	}

	if strings.Contains(stdout, "registry.k8s.io/etcd") {
		if verbose {
			logging.LogInfof("Host '%s' seems to be a kubernetes controlplane since etcd container was found, not a node itself.", hostname)
		}
		isKubernetesNode = false
	}

	if verbose {
		if isKubernetesNode {
			logging.LogInfof("Host '%s' is a kubernetes node.", hostname)
		} else {
			logging.LogInfof("Host '%s' is not a kubernetes node.", hostname)
		}
	}

	return isKubernetesNode, nil
}
