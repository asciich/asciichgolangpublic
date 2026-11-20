package gitlabcmdgeneric

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic"
)

type GitlabCmdOptions interface {
	GetGitlabInstance(cmd *cobra.Command, args []string) (*asciichgolangpublic.GitlabInstance, error)
}
