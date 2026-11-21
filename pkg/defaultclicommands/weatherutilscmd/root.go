package weatherutilscmd

import (
	"os"

	"github.com/spf13/cobra"
)

func NewWeatherUtilsCmd() *cobra.Command {
	const short = "Get weather data."

	cmd := &cobra.Command{
		Use:   "weather",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` weather <get-current-temperature> [flags]
`,
	}

	cmd.AddCommand(
		NewGetCurrentTemperatureCmd(),
	)

	return cmd
}
