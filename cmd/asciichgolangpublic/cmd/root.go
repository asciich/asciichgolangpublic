package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/defaultclicommands"
)

var rootCmd = &cobra.Command{
	Use:   "asciichgolangpublic",
	Short: "System admin helper",
}

func init() {
	defaultclicommands.MustAddDefaultCommands(rootCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
