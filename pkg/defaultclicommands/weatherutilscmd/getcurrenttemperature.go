package weatherutilscmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/weatherutils"
)

func NewGetCurrentTemperatureCmd() *cobra.Command {
	const short = "Get the current temperature in degree celsius at the given coordinates."

	cmd := &cobra.Command{
		Use:   "get-current-temperature",
		Short: short,
		Long: short + `

The weather data is provided by open-meteo.com. It is free for non-commercial use
and licensed under CC BY 4.0, so please credit Open-Meteo.

Usage:
    ` + os.Args[0] + ` weather get-current-temperature --lat <latitude> --lon <longitude>

Usage examples:
    # Bern (Switzerland)
    ` + os.Args[0] + ` weather get-current-temperature --lat 46.948 --lon 7.4474
`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			lat := mustutils.Must(cmd.Flags().GetFloat64("lat"))
			lon := mustutils.Must(cmd.Flags().GetFloat64("lon"))

			temperature := mustutils.Must(weatherutils.GetCurrentTemperature(ctx, lat, lon))

			// Only the value goes to stdout so the output is easy to use in scripts.
			fmt.Fprintf(cmd.OutOrStdout(), "%.1f\n", temperature)
		},
	}

	cmd.Flags().Float64("lat", 0, "Latitude of the position in decimal degrees (-90 to 90)")
	cmd.Flags().Float64("lon", 0, "Longitude of the position in decimal degrees (-180 to 180)")

	mustutils.Must0(cmd.MarkFlagRequired("lat"))
	mustutils.Must0(cmd.MarkFlagRequired("lon"))

	return cmd
}
