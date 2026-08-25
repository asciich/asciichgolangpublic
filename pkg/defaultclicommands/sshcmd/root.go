package sshcmd

import (
	"os"

	"github.com/spf13/cobra"
)

func NewSshCmd() *cobra.Command {
	const short = "SSH related commands"

	cmd := &cobra.Command{
		Use:   "ssh",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` ssh ssh`,
	}

	return cmd
}
