package operateheadscalenodescmd

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/networkcmd/vpncmd/headscalecmd/operateheadscalecmd/operateheadscaleoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/headscaleutils/headscaleinterfaces"
)

func NewListCmd(options *operateheadscaleoptions.OperateOptions) *cobra.Command {
	const short = "List headscale nodes."

	cmd := &cobra.Command{
		Use:   "list",
		Short: short,
		Long: short + `

The node names are listed alphabetically ordered.
With --wide the node details are shown as table.

Usage:
    ` + os.Args[0] + ` network vpn headscale operate nodes list [--wide]`,

		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			wide := mustutils.Must(cmd.Flags().GetBool("wide"))
			headscale := options.GetHeadScale(ctx, cmd)

			if wide {
				nodes := mustutils.Must(headscale.ListNodes(ctx))
				printNodesWide(nodes)
			} else {
				for _, name := range mustutils.Must(headscale.ListNodeNames(ctx)) {
					fmt.Println(name)
				}
			}

			logging.LogGoodByCtxf(ctx, "List headscale nodes finished")
		},
	}

	cmd.Flags().Bool("wide", false, "Show ID, Hostname, Name, User, Tags, IP addresses, Last seen and Connected")

	return cmd
}

// printNodesWide prints the nodes as table, sorted alphabetically by name.
func printNodesWide(nodes []headscaleinterfaces.Node) {
	type row struct {
		name   string
		fields []string
	}

	rows := make([]row, 0, len(nodes))
	for _, node := range nodes {
		id := mustutils.Must(node.GetID())
		hostname := mustutils.Must(node.GetHostname())
		name := mustutils.Must(node.GetName())
		user := mustutils.Must(node.GetUser())
		tags := mustutils.Must(node.GetTags())
		ips := mustutils.Must(node.GetIPAddresses())
		lastSeen := mustutils.Must(node.GetLastSeen())
		online := mustutils.Must(node.IsOnline())

		lastSeenString := "-"
		if lastSeen != nil {
			lastSeenString = lastSeen.Local().Format("2006-01-02 15:04:05")
		}

		tagsString := "-"
		if len(tags) > 0 {
			tagsString = strings.Join(tags, ",")
		}

		connected := "offline"
		if online {
			connected = "online"
		}

		rows = append(rows, row{
			name: name,
			fields: []string{
				strconv.FormatUint(id, 10),
				hostname,
				name,
				user,
				tagsString,
				strings.Join(ips, ", "),
				lastSeenString,
				connected,
			},
		})
	}

	sort.SliceStable(rows, func(i, j int) bool {
		return rows[i].name < rows[j].name
	})

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, strings.Join([]string{"ID", "Hostname", "Name", "User", "Tags", "IP addresses", "Last seen", "Connected"}, "\t"))
	for _, r := range rows {
		fmt.Fprintln(w, strings.Join(r.fields, "\t"))
	}
	mustutils.Must0(w.Flush())
}
