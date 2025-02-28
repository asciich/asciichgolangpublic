package dns

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/contextutils"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/netutils"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

func NewDnsCommand() (dnsCmd *cobra.Command) {
	dnsCmd = &cobra.Command{
		Use:   "dns",
		Short: "DNS related commands.",
	}

	dnsLookupV4Cmd := &cobra.Command{
		Use:   "lookup-v4",
		Short: "Lookup IPv4 addresses for given hostname.",
		Run: func(cmd *cobra.Command, args []string) {
			verbose, err := cmd.Flags().GetBool("verbose")
			if err != nil {
				logging.LogGoErrorFatal(err)
			}

			if len(args) != 1 {
				logging.LogFatal("Please specify exaclty 1 hostname to resolve.")
			}

			hostname := args[0]

			ips := netutils.MustLookupIpV4(
				contextutils.GetVerbosityContextByBool(verbose),
				hostname,
			)

			for _, ip := range ips {
				fmt.Println(ip)
			}
		},
	}

	dnsCmd.AddCommand(dnsLookupV4Cmd)

	return dnsCmd
}

func AddDnsCommand(rootCmd *cobra.Command) (err error) {
	if rootCmd == nil {
		return tracederrors.TracedErrorNil("rootCmd")
	}

	rootCmd.AddCommand(NewDnsCommand())

	return nil
}
