package kvmcmdutils

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/vmutils/kvmutils"
)

func GetCtxAndHostname(cmd *cobra.Command) (context.Context, string) {
	ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

	hostname, err := cmd.Flags().GetString("hostname")
	if err != nil {
		logging.LogGoErrorFatalWithTrace(err)
	}

	if hostname == "" {
		logging.LogFatal("Please specify the --hostname of the KVM hypervisor. Use --hostname=localhost to run against the local machine.")
	}

	return ctx, hostname
}

func GetCtxAndKvmHypervisor(cmd *cobra.Command) (context.Context, *kvmutils.KVMHypervisor) {
	ctx, hostname := GetCtxAndHostname(cmd)

	hypervisor, err := kvmutils.GetKvmHypervisorByHostName(hostname)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return ctx, hypervisor
}
