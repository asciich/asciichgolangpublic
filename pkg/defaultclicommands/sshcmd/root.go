package sshcmd

import "github.com/spf13/cobra"

func NewSshCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "ssh",
		Short: "SSH related commands",
	}

	return cmd
}