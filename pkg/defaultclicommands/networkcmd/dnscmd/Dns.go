package dnscmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/dnsutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
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

			ctx := contextutils.WithVerbosityContextByBool(cmd.Context(), verbose)

			if len(args) != 1 {
				logging.LogFatal("Please specify exaclty 1 hostname to resolve.")
			}

			hostname := args[0]

			ips := mustutils.Must(dnsutils.DnsLookupIpV4(ctx, hostname))

			for _, ip := range ips {
				fmt.Println(ip)
			}
		},
	}

	dnsReverseLookupCmd := &cobra.Command{
		Use:   "reverse-lookup",
		Short: "Reverse lookup for given IP address.",
		Run: func(cmd *cobra.Command, args []string) {
			verbose, err := cmd.Flags().GetBool("verbose")
			if err != nil {
				logging.LogGoErrorFatal(err)
			}

			ctx := contextutils.WithVerbosityContextByBool(cmd.Context(), verbose)

			if len(args) != 1 {
				logging.LogFatal("Please specify exaclty 1 IP address to resolve.")
			}

			hostname := args[0]

			ips := mustutils.Must(dnsutils.DnsReverseLookup(ctx, hostname))

			for _, ip := range ips {
				fmt.Println(ip)
			}
		},
	}

	dnsCmd.AddCommand(dnsLookupV4Cmd)
	dnsCmd.AddCommand(dnsReverseLookupCmd)

	return dnsCmd
}

func AddDnsCommand(rootCmd *cobra.Command) (err error) {
	if rootCmd == nil {
		return tracederrors.TracedErrorNil("rootCmd")
	}

	rootCmd.AddCommand(NewDnsCommand())

	return nil
}
