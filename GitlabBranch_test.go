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
				ctx := getCtx()

				gitlab, err := GetGitlabByFQDN("gitlab.asciich.ch")
				require.NoError(t, err)

				err = gitlab.Authenticate(ctx,
					&GitlabAuthenticationOptions{
						AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"},
					},
				)
				require.NoError(t, err)

				gitlabProject, err := gitlab.GetGitlabProjectByPath(ctx, "test_group/testproject")
				require.NoError(t, err)

				exists, err := gitlabProject.Exists(ctx)
				require.NoError(t, err)
				require.True(t, exists)

				branch, err := gitlabProject.GetBranchByName(tt.branchName)
				require.NoError(t, err)

				err = branch.Delete(
					ctx,
					&GitlabDeleteBranchOptions{
						SkipWaitForDeletion: false,
					})
				require.NoError(t, err)
				exists, err = branch.Exists(ctx)
				require.NoError(t, err)
				require.False(t, exists)

				for i := 0; i < 2; i++ {
					branch.CreateFromDefaultBranch(ctx)
					exists, err := branch.Exists(ctx)
					require.NoError(t, err)
					require.True(t, exists)
				}

				for i := 0; i < 2; i++ {
					err := branch.Delete(
						ctx,
						&GitlabDeleteBranchOptions{
							SkipWaitForDeletion: false,
						})
					require.NoError(t, err)
					exists, err := branch.Exists(ctx)
					require.NoError(t, err)
					require.False(t, exists)
				}
			},
		)
	}
}
