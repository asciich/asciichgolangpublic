package operateheadscaleapikeyscmd

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/contextutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/defaultclicommands/networkcmd/vpncmd/headscalecmd/operateheadscalecmd/operateheadscaleoptions"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/logging"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/mustutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/netutils/vpnutils/headscaleutils/headscaleinterfaces"
)

func NewListCmd(options *operateheadscaleoptions.OperateOptions) *cobra.Command {
	const short = "List headscale API keys."

	cmd := &cobra.Command{
		Use:   "list",
		Short: short,
		Long: short + `

The prefixes of the API keys are listed alphabetically ordered.
With --wide the API key details are shown as table.
The full API keys are never shown, since headscale only reveals them once at creation.

Usage:
    ` + os.Args[0] + ` network vpn headscale operate apikeys list [--wide]`,

		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			wide := mustutils.Must(cmd.Flags().GetBool("wide"))
			headscale := options.GetHeadScale(ctx, cmd)

			keys := mustutils.Must(headscale.ListApiKeys(ctx))

			if wide {
				printApiKeysWide(keys)
			} else {
				for _, prefix := range getSortedPrefixes(keys) {
					fmt.Println(prefix)
				}
			}

			logging.LogGoodByCtxf(ctx, "List headscale API keys finished")
		},
	}

	cmd.Flags().Bool("wide", false, "Show ID, Prefix, Expiration and Created")

	return cmd
}

func getSortedPrefixes(keys []headscaleinterfaces.ApiKey) []string {
	ret := make([]string, 0, len(keys))
	for _, key := range keys {
		ret = append(ret, mustutils.Must(key.GetPrefix()))
	}

	sort.Strings(ret)

	return ret
}

// printApiKeysWide prints the API keys as table, sorted alphabetically by prefix.
func printApiKeysWide(keys []headscaleinterfaces.ApiKey) {
	type row struct {
		prefix string
		fields []string
	}

	rows := make([]row, 0, len(keys))
	for _, key := range keys {
		id := mustutils.Must(key.GetID())
		prefix := mustutils.Must(key.GetPrefix())
		expiration := mustutils.Must(key.GetExpiration())
		created := mustutils.Must(key.GetCreated())

		expirationString := "never"
		if expiration != nil {
			expirationString = expiration.Local().Format("2006-01-02 15:04:05")
		}

		rows = append(rows, row{
			prefix: prefix,
			fields: []string{
				strconv.FormatUint(id, 10),
				prefix,
				expirationString,
				created.Local().Format("2006-01-02 15:04:05"),
			},
		})
	}

	sort.SliceStable(rows, func(i, j int) bool {
		return rows[i].prefix < rows[j].prefix
	})

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tPrefix\tExpiration\tCreated")
	for _, r := range rows {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", r.fields[0], r.fields[1], r.fields[2], r.fields[3])
	}
	mustutils.Must0(w.Flush())
}
