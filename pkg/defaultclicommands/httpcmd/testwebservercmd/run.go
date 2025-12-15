package testwebservercmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/testwebserver"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/userinteraction"
)

func NewRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run the test webserver locally on given --port",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			port, err := cmd.Flags().GetInt("port")
			if err != nil {
				logging.LogGoErrorFatal(err)
			}

			if port <= 0 {
				logging.LogFatalf("port '%d' is smaller or equal 0 and not valid.", port)
			}

			testwebserver := mustutils.Must(testwebserver.GetTestWebServer(port))

			mustutils.Must0(testwebserver.StartInBackground(ctx))

			userinteraction.WaitUserAbortf("Testwebserver running. It is reachable on http://localhost:%d .", port)

			mustutils.Must0(testwebserver.Stop(ctx))

			logging.LogInfoByCtxf(ctx, "local testwebserver on port '%d' stopped.", port)
		},
	}

	cmd.Flags().Int("port", 80, "Port of the test webserver to listen to.")

	return cmd
}
