package operateheadscaleuserscmd

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/networkcmd/vpncmd/headscalecmd/operateheadscalecmd/operateheadscaleoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/headscaleutils/headscaleinterfaces"
)

func NewListCmd(options *operateheadscaleoptions.OperateOptions) *cobra.Command {
	const short = "List headscale users."

	cmd := &cobra.Command{
		Use:   "list",
		Short: short,
		Long: short + `

The usernames are listed alphabetically ordered.
With --wide the user details are shown as table.

Usage:
    ` + os.Args[0] + ` network vpn headscale operate users list [--wide]`,

		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			wide := mustutils.Must(cmd.Flags().GetBool("wide"))
			headscale := options.GetHeadScale(ctx, cmd)

			if wide {
				users := mustutils.Must(headscale.ListUsers(ctx))
				printUsersWide(users)
			} else {
				for _, name := range mustutils.Must(headscale.ListUserNames(ctx)) {
					fmt.Println(name)
				}
			}

			logging.LogGoodByCtxf(ctx, "List headscale users finished")
		},
	}

	cmd.Flags().Bool("wide", false, "Show ID, Name, Username, Email and Created")

	return cmd
}

// printUsersWide prints the users as table, sorted alphabetically by username.
func printUsersWide(users []headscaleinterfaces.User) {
	type row struct {
		username string
		fields   []string
	}

	rows := make([]row, 0, len(users))
	for _, user := range users {
		id := mustutils.Must(user.GetID())
		name := mustutils.Must(user.GetName())
		username := mustutils.Must(user.GetUsername())
		email := mustutils.Must(user.GetEmail())
		created := mustutils.Must(user.GetCreated())

		if name == "" {
			name = "-"
		}

		if email == "" {
			email = "-"
		}

		rows = append(rows, row{
			username: username,
			fields: []string{
				strconv.FormatUint(id, 10),
				name,
				username,
				email,
				created.Local().Format("2006-01-02 15:04:05"),
			},
		})
	}

	sort.SliceStable(rows, func(i, j int) bool {
		return rows[i].username < rows[j].username
	})

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tName\tUsername\tEmail\tCreated")
	for _, r := range rows {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", r.fields[0], r.fields[1], r.fields[2], r.fields[3], r.fields[4])
	}
	mustutils.Must0(w.Flush())
}
