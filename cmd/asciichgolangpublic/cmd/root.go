package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands"
)

func NewRootCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "asciichgolangpublic",
		Short: "System admin helper",
	}

	err := defaultclicommands.AddDefaultCommands(cmd)
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

func Execute() {
	rootCmd, err := NewRootCmd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
