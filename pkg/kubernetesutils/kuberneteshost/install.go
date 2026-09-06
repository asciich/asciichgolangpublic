package kuberneteshost

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubeadmutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubectlutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubeletutils"
	"github.com/asciich/asciichgolangpublic/pkg/osutils/linuxutils"
	"github.com/asciich/asciichgolangpublic/pkg/osutils/linuxutils/containerdutils"
	"github.com/asciich/asciichgolangpublic/pkg/osutils/linuxutils/systemdutils"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/packagemanagergeneric"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/packagemanageroptions"
	"github.com/asciich/asciichgolangpublic/pkg/runbook"
)

func NewInstallRunbook(commandExecutor commandexecutorinterfaces.CommandExecutor) (*runbook.RunBook, error) {
	runbook := &runbook.RunBook{
		Name:        "Install kubernetes host",
		Description: "Install everything needed to act as kubernetes host (control-plane or worker-node).",
		Steps: []runbook.Runnable{
			&runbook.Step{
				Name:        "disable-swap",
				Description: "Disable swap now and persistently in /etc/fstab (required by kubelet).",
				Run: func(ctx context.Context) error {
					return linuxutils.TurnSwapOff(ctx, commandExecutor)
				},
			},
			&runbook.Step{
				Name:        "load-kernel-modules",
				Description: "Enable the overlay and br_netfilter kernel modules needed by the container runtime and networking.",
				Run: func(ctx context.Context) error {
					return linuxutils.LoadKernelModules(ctx, commandExecutor, []string{
						"overlay",
						"br_netfilter",
					})
				},
			},
			&runbook.Step{
				Name:        "configure-sysctl",
				Description: "Configure sysctl parameters for bridged traffic and IP forwarding, then apply them.",
				Run: func(ctx context.Context) error {
					return linuxutils.SetSysctlValues(ctx, commandExecutor, map[string]string{
						"net.bridge.bridge-nf-call-iptables":  "1",
						"net.bridge.bridge-nf-call-ip6tables": "1",
						"net.ipv4.ip_forward":                 "1",
					})
				},
			},
			&runbook.Step{
				Name:        "install-containerd",
				Description: "Install and configure the containerd container runtime with the systemd cgroup driver.",
				Run: func(ctx context.Context) error {
					packageManager, err := packagemanagergeneric.NewPackageManagerGeneric(ctx, commandExecutor)
					if err != nil {
						return err
					}

					err = packageManager.InstallPackages(
						ctx,
						[]string{"containerd"},
						&packagemanageroptions.InstallPackageOptions{
							UpdateDatabaseFirst: true,
							UpdateKeyringFirst:  true,
							Force:               true,
						},
					)
					if err != nil {
						return err
					}

					err = systemdutils.EnableAndStartService(ctx, commandExecutor, "containerd")
					if err != nil {
						return err
					}

					err = containerdutils.EnableCGroup(ctx, commandExecutor)
					if err != nil {
						return err
					}

					return nil
				},
			},
			&runbook.Step{
				Name:        "install-kube-packages",
				Description: "Install kubelet, kubeadm and kubectl and pin them to prevent unintended upgrades.",
				Run: func(ctx context.Context) error {
					err := kubectlutils.InstallKubectlUsingCommandExecutor(ctx, commandExecutor, &kubectlutils.InstallKubectlOptions{})
					if err != nil {
						return err
					}

					err = kubeadmutils.InstallKubeadmUsingCommandExecutor(ctx, commandExecutor, &kubeadmutils.InstallKubeadmOptions{})
					if err != nil {
						return err
					}

					err = kubeletutils.InstallKubeletUsingCommandExecutor(ctx, commandExecutor, &kubeletutils.InstallKubeletOptions{})
					if err != nil {
						return err
					}

					return nil
				},
			},
			&runbook.Step{
				Name:        "enable-kubelet",
				Description: "Enable the kubelet service so it starts on boot.",
				Run: func(ctx context.Context) error {
					return nil
				},
			},
		},
	}

	return runbook, nil
}

func Install(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor) error {
	runbook, err := NewInstallRunbook(commandExecutor)
	if err != nil {
		return err
	}

	return runbook.Execute(ctx)
}
