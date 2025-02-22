package asciichgolangpublic

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/testutils"
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
				require := require.New(t)

				const verbose bool = true

				gitlab := MustGetGitlabByFqdn("gitlab.asciich.ch")
				gitlab.MustAuthenticate(&GitlabAuthenticationOptions{
					AccessTokensFromGopass: []string{"hosts/gitlab.asciich.ch/users/reto/access_token"},
				})

				gitlabProject := gitlab.MustGetGitlabProjectByPath("test_group/testproject", verbose)
				require.True(gitlabProject.MustExists(verbose))

				branchName := gitlabProject.MustGetDefaultBranchName()

				gitlabProject.MustDeleteAllBranchesExceptDefaultBranch(verbose)

				branchList := gitlabProject.MustGetBranchNames(verbose)
				require.EqualValues(
					[]string{branchName},
					branchList,
				)

				expectedBranchesList := []string{}
				for i := 0; i < 21; i++ {
					expectedBranchesList = append(
						expectedBranchesList,
						fmt.Sprintf("testbranch_%03d", i),
					)
				}

				for _, toCreate := range expectedBranchesList {
					gitlabProject.MustCreateBranchFromDefaultBranch(toCreate, verbose)
				}

				branchList = gitlabProject.MustGetBranchNames(verbose)
				require.EqualValues(branchList, branchList)
			},
		)
	}
}
