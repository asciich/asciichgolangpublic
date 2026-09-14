package kvmutils

import (
	"github.com/asciich/asciichgolangpublic/pkg/hostsutils"
	"github.com/asciich/asciichgolangpublic/pkg/hostsutils/hostsutilsinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"github.com/asciich/asciichgolangpublic/pkg/vmutils/kvmutils/commandexecutorkvmutils"
	"github.com/asciich/asciichgolangpublic/pkg/vmutils/kvmutils/kvmutilsinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/vmutils/kvmutils/nativekvmutils"
)

func GetKvmHypervisorByHostName(hostname string) (kvmutilsinterfaces.Hypervisor, error) {
	if hostname == "" {
		return nil, tracederrors.TracedErrorEmptyString("hostname")
	}

	if hostname == "localhost" {
		return GetKvmHypervisorOnLocalhost()
	}

	host, err := hostsutils.GetHostByHostname(hostname)
	if err != nil {
		return nil, err
	}

	return GetKvmHypervisorByHost(host)
}

func GetKvmHypervisorByHost(host hostsutilsinterfaces.Host) (kvmutilsinterfaces.Hypervisor, error) {
	if host == nil {
		return nil, tracederrors.TracedError("host is nil")
	}

	kvmHypervisor := commandexecutorkvmutils.NewKVMHypervisor()
	err := kvmHypervisor.SetHost(host)
	if err != nil {
		return nil, err
	}

	return kvmHypervisor, nil
}

func GetKvmHypervisorOnLocalhost() (kvmutilsinterfaces.Hypervisor, error) {
	kvmHypervisor := nativekvmutils.NewKvmHypervisor()
	err := kvmHypervisor.SetUseLocalhost(true)
	if err != nil {
		return nil, err
	}

	return kvmHypervisor, nil
}
