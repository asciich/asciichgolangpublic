package asciichgolangpublic

import (
	"strings"
)

type KubernetesControlplaneHost struct {
	Host
}

func GetKubernetesControlplaneByHostname(hostname string) (kubernetesControlplaneHost *KubernetesControlplaneHost, err error) {
	if len(hostname) <= 0 {
		return nil, TracedError("hostname is empty string")
	}

	kubernetesControlplaneHost = NewKubernetesControlplaneHost()

	err = kubernetesControlplaneHost.SetHostname(hostname)
	if err != nil {
		return nil, err
	}

	return kubernetesControlplaneHost, nil
}

func MustGetKubernetesControlplaneByHostname(hostname string) (kubernetesControlplaneHost *KubernetesControlplaneHost) {
	kubernetesControlplaneHost, err := GetKubernetesControlplaneByHostname(hostname)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return kubernetesControlplaneHost
}

func NewKubernetesControlplaneHost() (kubernetesControlplaneHost *KubernetesControlplaneHost) {
	kubernetesControlplaneHost = new(KubernetesControlplaneHost)
	return kubernetesControlplaneHost
}

func (k *KubernetesControlplaneHost) MustCheckIsKubernetesControlplane(verbose bool) (isKubernetesControlplane bool) {
	isKubernetesControlplane, err := k.CheckIsKubernetesControlplane(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return isKubernetesControlplane
}

func (k *KubernetesControlplaneHost) MustGetJoinCommandAsString(verbose bool) (joinCommand string) {
	joinCommand, err := k.GetJoinCommandAsString(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return joinCommand
}

func (k *KubernetesControlplaneHost) MustGetJoinCommandAsStringSlice(verbose bool) (joinCommand []string) {
	joinCommand, err := k.GetJoinCommandAsStringSlice(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return joinCommand
}

func (k *KubernetesControlplaneHost) MustIsKubernetesControlplane(verbose bool) (isKubernetesControlplane bool) {
	isKubernetesControlplane, err := k.IsKubernetesControlplane(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return isKubernetesControlplane
}

func (n *KubernetesControlplaneHost) CheckIsKubernetesControlplane(verbose bool) (isKubernetesControlplane bool, err error) {
	hostname, err := n.GetHostname()
	if err != nil {
		return false, err
	}

	isKubernetesControlplane, err = n.IsKubernetesControlplane(verbose)
	if err != nil {
		return false, err
	}

	if !isKubernetesControlplane {
		return false, TracedErrorf("Host '%s' is not a kubernetes controlplane", hostname)
	}

	return isKubernetesControlplane, nil
}

func (n *KubernetesControlplaneHost) GetJoinCommandAsString(verbose bool) (joinCommand string, err error) {
	hostname, err := n.GetHostname()
	if err != nil {
		return "", err
	}

	isControlPlane, err := n.IsKubernetesControlplane(verbose)
	if err != nil {
		return "", err
	}

	if !isControlPlane {
		return "", TracedErrorf(
			"host '%s' is not a kubernetes control plane and therefore join command can be generated.",
			hostname,
		)
	}

	joinCommand, err = n.RunCommandAndGetStdoutAsString(
		&RunCommandOptions{
			Command: []string{"kubeadm", "token", "create", "--print-join-command"},
			Verbose: verbose,
		},
	)
	if err != nil {
		return "", err
	}

	joinCommand = strings.TrimSpace(joinCommand)

	if len(joinCommand) <= 0 {
		return "", TracedError("Unable to get joinCommand. Evaluated joinCommand is empty string")
	}

	if verbose {
		LogChangedf("Generated join command for a new kubernetes node on control plane host '%s'", hostname)
	}

	return joinCommand, nil
}

func (n *KubernetesControlplaneHost) GetJoinCommandAsStringSlice(verbose bool) (joinCommand []string, err error) {
	joinCommandString, err := n.GetJoinCommandAsString(verbose)
	if err != nil {
		return nil, err
	}

	joinCommand, err = ShellLineHandler().Split(
		joinCommandString,
	)
	if err != nil {
		return nil, err
	}

	return joinCommand, nil
}

func (n *KubernetesControlplaneHost) IsKubernetesControlplane(verbose bool) (isKubernetesControlplane bool, err error) {
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

	isKubernetesControlplane = true

	if len(Strings().SplitLines(stdout, false)) <= 5 {
		isKubernetesControlplane = false
	}

	if !strings.Contains(stdout, "registry.k8s.io/kube-proxy") {
		isKubernetesControlplane = false
	}

	if !strings.Contains(stdout, "registry.k8s.io/etcd") {
		if verbose {
			LogInfof("Host '%s' seems to be a kubernetes node since etcd container was not found, not a controlplane itself.", hostname)
		}
		isKubernetesControlplane = false
	}

	if verbose {
		if isKubernetesControlplane {
			LogInfof("Host '%s' is a kubernetes controlplane.", hostname)
		} else {
			LogInfof("Host '%s' is not a kubernetes controlplane.", hostname)
		}
	}

	return isKubernetesControlplane, nil
}
