package kubernetes

import (
	"strings"

	"github.com/asciich/asciichgolangpublic"
	"github.com/asciich/asciichgolangpublic/hosts"
)

type KubernetesControlplaneHost struct {
	hosts.Host
}

func GetKubernetesControlplaneByHost(host hosts.Host) (kubernetesControlplaneHost *KubernetesControlplaneHost, err error) {
	if host == nil {
		return nil, asciichgolangpublic.TracedErrorNil("host")
	}

	kubernetesControlplaneHost = NewKubernetesControlplaneHost()
	kubernetesControlplaneHost.Host = host

	return kubernetesControlplaneHost, nil
}

func GetKubernetesControlplaneByHostname(hostname string) (kubernetesControlplaneHost *KubernetesControlplaneHost, err error) {
	if len(hostname) <= 0 {
		return nil, asciichgolangpublic.TracedError("hostname is empty string")
	}

	host, err := hosts.GetHostByHostname(hostname)
	if err != nil {
		return nil, err
	}

	return GetKubernetesControlplaneByHost(host)
}

func MustGetKubernetesControlplaneByHost(host hosts.Host) (kubernetesControlplaneHost *KubernetesControlplaneHost) {
	kubernetesControlplaneHost, err := GetKubernetesControlplaneByHost(host)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return kubernetesControlplaneHost
}

func MustGetKubernetesControlplaneByHostname(hostname string) (kubernetesControlplaneHost *KubernetesControlplaneHost) {
	kubernetesControlplaneHost, err := GetKubernetesControlplaneByHostname(hostname)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return kubernetesControlplaneHost
}

func NewKubernetesControlplaneHost() (kubernetesControlplaneHost *KubernetesControlplaneHost) {
	kubernetesControlplaneHost = new(KubernetesControlplaneHost)
	return kubernetesControlplaneHost
}

func (k *KubernetesControlplaneHost) CheckIsKubernetesControlplane(verbose bool) (isKubernetesControlplane bool, err error) {
	hostname, err := k.GetHostName()
	if err != nil {
		return false, err
	}

	isKubernetesControlplane, err = k.IsKubernetesControlplane(verbose)
	if err != nil {
		return false, err
	}

	if !isKubernetesControlplane {
		return false, asciichgolangpublic.TracedErrorf("Host '%s' is not a kubernetes controlplane", hostname)
	}

	return isKubernetesControlplane, nil
}

func (k *KubernetesControlplaneHost) GetJoinCommandAsString(verbose bool) (joinCommand string, err error) {
	hostname, err := k.GetHostName()
	if err != nil {
		return "", err
	}

	isControlPlane, err := k.IsKubernetesControlplane(verbose)
	if err != nil {
		return "", err
	}

	if !isControlPlane {
		return "", asciichgolangpublic.TracedErrorf(
			"host '%s' is not a kubernetes control plane and therefore join command can be generated.",
			hostname,
		)
	}

	joinCommand, err = k.RunCommandAndGetStdoutAsString(
		&asciichgolangpublic.RunCommandOptions{
			Command: []string{"kubeadm", "token", "create", "--print-join-command"},
			Verbose: verbose,
		},
	)
	if err != nil {
		return "", err
	}

	joinCommand = strings.TrimSpace(joinCommand)

	if len(joinCommand) <= 0 {
		return "", asciichgolangpublic.TracedError("Unable to get joinCommand. Evaluated joinCommand is empty string")
	}

	if verbose {
		asciichgolangpublic.LogChangedf("Generated join command for a new kubernetes node on control plane host '%s'", hostname)
	}

	return joinCommand, nil
}

func (k *KubernetesControlplaneHost) GetJoinCommandAsStringSlice(verbose bool) (joinCommand []string, err error) {
	joinCommandString, err := k.GetJoinCommandAsString(verbose)
	if err != nil {
		return nil, err
	}

	joinCommand, err = asciichgolangpublic.ShellLineHandler().Split(
		joinCommandString,
	)
	if err != nil {
		return nil, err
	}

	return joinCommand, nil
}

func (k *KubernetesControlplaneHost) IsKubernetesControlplane(verbose bool) (isKubernetesControlplane bool, err error) {
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

	isKubernetesControlplane = true

	if len(asciichgolangpublic.Strings().SplitLines(stdout, false)) <= 5 {
		isKubernetesControlplane = false
	}

	if !strings.Contains(stdout, "registry.k8s.io/kube-proxy") {
		isKubernetesControlplane = false
	}

	if !strings.Contains(stdout, "registry.k8s.io/etcd") {
		if verbose {
			asciichgolangpublic.LogInfof("Host '%s' seems to be a kubernetes node since etcd container was not found, not a controlplane itself.", hostname)
		}
		isKubernetesControlplane = false
	}

	if verbose {
		if isKubernetesControlplane {
			asciichgolangpublic.LogInfof("Host '%s' is a kubernetes controlplane.", hostname)
		} else {
			asciichgolangpublic.LogInfof("Host '%s' is not a kubernetes controlplane.", hostname)
		}
	}

	return isKubernetesControlplane, nil
}

func (k *KubernetesControlplaneHost) MustCheckIsKubernetesControlplane(verbose bool) (isKubernetesControlplane bool) {
	isKubernetesControlplane, err := k.CheckIsKubernetesControlplane(verbose)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return isKubernetesControlplane
}

func (k *KubernetesControlplaneHost) MustGetJoinCommandAsString(verbose bool) (joinCommand string) {
	joinCommand, err := k.GetJoinCommandAsString(verbose)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return joinCommand
}

func (k *KubernetesControlplaneHost) MustGetJoinCommandAsStringSlice(verbose bool) (joinCommand []string) {
	joinCommand, err := k.GetJoinCommandAsStringSlice(verbose)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return joinCommand
}

func (k *KubernetesControlplaneHost) MustIsKubernetesControlplane(verbose bool) (isKubernetesControlplane bool) {
	isKubernetesControlplane, err := k.IsKubernetesControlplane(verbose)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return isKubernetesControlplane
}
