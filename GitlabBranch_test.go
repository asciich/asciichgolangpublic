package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/continuousintegration"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func TestGitlabProjectBranchCreateAndDelete(t *testing.T) {
	if continuousintegration.IsRunningInGithub() {
		logging.LogInfo("Unavailable in Github CI")
		return
	}

	tests := []struct {
		branchName string
	}{
		{"testbranch"},
		{"testbranch2"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				gitlab := MustGetGitlabByFqdn("gitlab.asciich.ch")
				gitlab.MustAuthenticate(&GitlabAuthenticationOptions{
					AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"},
				})

				gitlabProject := gitlab.MustGetGitlabProjectByPath("test_group/testproject", verbose)
				require.True(gitlabProject.MustExists(verbose))

				branch := gitlabProject.MustGetBranchByName(tt.branchName)

				branch.MustDelete(&GitlabDeleteBranchOptions{
					SkipWaitForDeletion: false,
					Verbose:             verbose,
				})
				require.False(branch.MustExists())

				for i := 0; i < 2; i++ {
					branch.CreateFromDefaultBranch(verbose)
					require.True(branch.MustExists())
				}

				for i := 0; i < 2; i++ {
					branch.MustDelete(&GitlabDeleteBranchOptions{
						SkipWaitForDeletion: false,
						Verbose:             verbose,
					})
					require.False(branch.MustExists())
				}
			},
		)
	}
}
