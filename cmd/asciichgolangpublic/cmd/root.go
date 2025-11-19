package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "asciichgolangpublic",
		Short: "System admin helper",
	}

	err := defaultclicommands.AddDefaultCommands(cmd)
	if err != nil {
		return nil
	}

	return cmd
}

func Execute() {
	rootCmd := NewRootCmd()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
