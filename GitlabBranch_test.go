package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func TestGitlabProjectBranchCreateAndDelete(t *testing.T) {
	testutils.SkipIfRunningInGithub(t)

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
				const verbose bool = true

				gitlab := MustGetGitlabByFqdn("gitlab.asciich.ch")
				gitlab.MustAuthenticate(&GitlabAuthenticationOptions{
					AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"},
				})

				gitlabProject := gitlab.MustGetGitlabProjectByPath("test_group/testproject", verbose)
				require.True(t, gitlabProject.MustExists(verbose))

				branch := gitlabProject.MustGetBranchByName(tt.branchName)

				err := branch.Delete(&GitlabDeleteBranchOptions{
					SkipWaitForDeletion: false,
					Verbose:             verbose,
				})
				require.NoError(t, err)
				exists, err := branch.Exists()
				require.NoError(t, err)
				require.False(t, exists)

				for i := 0; i < 2; i++ {
					branch.CreateFromDefaultBranch(verbose)
					exists, err := branch.Exists()
					require.NoError(t, err)
					require.True(t, exists)
				}

				for i := 0; i < 2; i++ {
					err := branch.Delete(&GitlabDeleteBranchOptions{
						SkipWaitForDeletion: false,
						Verbose:             verbose,
					})
					require.NoError(t, err)
					exists, err := branch.Exists()
					require.NoError(t, err)
					require.False(t, exists)
				}
			},
		)
	}
}
