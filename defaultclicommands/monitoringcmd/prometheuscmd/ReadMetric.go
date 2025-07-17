package prometheuscmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/prometheusutils"
)

func NewReadMetricCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "read-metric-value",
		Short: "Read the value of the --metric-name from the given --metrics-page-url",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			metricName, err := cmd.Flags().GetString("metric-name")
			if err != nil {
				logging.LogGoErrorFatal(err)
			}

			if metricName == "" {
				logging.LogFatal("Please set --metric-name.")
			}

			metricsPageUrl, err := cmd.Flags().GetString("metrics-page-url")
			if err != nil {
				logging.LogGoErrorFatal(err)
			}

			if metricsPageUrl == "" {
				logging.LogFatal("Please set --metrics-page-url.")
			}

			cliReadMetricCmd(ctx, metricsPageUrl, metricName)
		},
	}

	cmd.PersistentFlags().String("metric-name", "", "Name of the metric to read.")
	cmd.PersistentFlags().String("metrics-page-url", "", "Url of the metrics page. Usually something like 'http://hostname:9123/metrics'.")

	return cmd
}

func cliReadMetricCmd(ctx context.Context, metricsPageUrl string, metricName string) {
	fmt.Println(prometheusutils.MustGetMetricValueFromMetricPage(ctx, metricsPageUrl, metricName))
}
