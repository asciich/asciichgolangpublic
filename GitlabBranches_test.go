package asciichgolangpublic

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func TestGitlabProjectBranches_pagination(t *testing.T) {
	testutils.SkipIfRunningInGithub(t)

	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				gitlab, err := GetGitlabByFQDN("gitlab.asciich.ch")
				require.NoError(t, err)

				err = gitlab.Authenticate(ctx, &GitlabAuthenticationOptions{AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"}})
				require.NoError(t, err)

				gitlabProject, err := gitlab.GetGitlabProjectByPath(ctx, "test_group/testproject")
				require.NoError(t, err)

				exists, err := gitlabProject.Exists(ctx)
				require.NoError(t, err)
				require.True(t, exists)

				branchName, err := gitlabProject.GetDefaultBranchName(ctx)
				require.NoError(t, err)

				err = gitlabProject.DeleteAllBranchesExceptDefaultBranch(ctx)
				require.NoError(t, err)

				branchList, err := gitlabProject.GetBranchNames(ctx)
				require.NoError(t, err)
				require.EqualValues(t, []string{branchName}, branchList)

				expectedBranchesList := []string{}
				for i := 0; i < 21; i++ {
					expectedBranchesList = append(
						expectedBranchesList,
						fmt.Sprintf("testbranch_%03d", i),
					)
				}

				for _, toCreate := range expectedBranchesList {
					_, err = gitlabProject.CreateBranchFromDefaultBranch(ctx, toCreate)
					require.NoError(t, err)
				}

				branchList, err = gitlabProject.GetBranchNames(ctx)
				require.NoError(t, err)
				require.EqualValues(t, branchList, branchList)
			},
		)
	}
}
